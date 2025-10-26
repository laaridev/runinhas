package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"time"

	"github.com/sirupsen/logrus"
)

// HeroConsumer processes hero-related events (deaths, health, mana, level)
type HeroConsumer struct {
	logger           *logrus.Entry
	lastDeaths       int64
	lastHealth       int64
	lastMana         int64
	lastLevel        int64
	eventChan        <-chan events.TickEvent
	stopChan         chan struct{}
	handlers         []handlers.Handler
	eventThrottle    map[string]time.Time // Throttle events to avoid spam
	throttleConfig   map[string]time.Duration // Configurable throttle per event type
}

// NewHeroConsumer creates a new hero consumer with handlers
func NewHeroConsumer(eventBus *events.EventBus, logger *logrus.Entry, handlerList []handlers.Handler) *HeroConsumer {
	// Configurable throttle per event type for better control
	throttleConfig := map[string]time.Duration{
		"hero_health_low":     5 * time.Second,  // More responsive for critical events
		"hero_health_critical": 3 * time.Second,  // Even faster for critical health
		"hero_mana_low":       3 * time.Second,  // Quick mana warnings
		"hero_death":          0,                // No throttle for death events
		"hero_level_up":       0,                // No throttle for level up
		"hero_ultimate_ready": 0,                // No throttle for ultimate
	}

	return &HeroConsumer{
		logger:         logger,
		eventChan:      eventBus.Subscribe(),
		stopChan:       make(chan struct{}),
		handlers:       handlerList,
		eventThrottle:  make(map[string]time.Time),
		throttleConfig: throttleConfig,
	}
}

// Start begins consuming events
func (hc *HeroConsumer) Start() {
	go hc.consume()
	hc.logger.Info("🦸 HeroConsumer started")
}

// Stop stops the consumer
func (hc *HeroConsumer) Stop() {
	close(hc.stopChan)
	hc.logger.Info("🦸 HeroConsumer stopped")
}

// consume processes TickEvents and detects hero changes
func (hc *HeroConsumer) consume() {
	for {
		select {
		case event := <-hc.eventChan:
			hc.processHeroChanges(event)
		case <-hc.stopChan:
			return
		}
	}
}

// processHeroChanges extracts hero data and detects changes
func (hc *HeroConsumer) processHeroChanges(event events.TickEvent) {
	// Use parsed event for efficient JSON access
	parsed := events.NewParsedTickEvent(event)

	// Extract hero data using cached parse
	deaths := parsed.GetInt64("player.deaths")
	health := parsed.GetInt64("hero.health_percent")
	mana := parsed.GetInt64("hero.mana_percent")
	level := parsed.GetInt64("hero.level")

	// Check for death events (deaths increased)
	if deaths > hc.lastDeaths && hc.lastDeaths >= 0 {
		if hc.isEventEnabled("hero_death") {
			hc.handleEvent("hero_death", map[string]interface{}{
				"deaths":      deaths,
				"prev_deaths": hc.lastDeaths,
				"deaths_diff": deaths - hc.lastDeaths,
			})
		}
	}

	// Use default thresholds
	healthThreshold := int64(25)
	manaThreshold := int64(15)

	// Check for low health
	if health > 0 && health <= healthThreshold && hc.lastHealth > healthThreshold {
		if hc.isEventEnabled("hero_health_low") && hc.canTriggerEvent("hero_health_low") {
			hc.handleEvent("hero_health_low", map[string]interface{}{
				"health":      health,
				"prev_health": hc.lastHealth,
				"value":       healthThreshold,
			})
		}
	}

	// Check for low mana
	if mana > 0 && mana <= manaThreshold && hc.lastMana > manaThreshold {
		if hc.isEventEnabled("hero_mana_low") && hc.canTriggerEvent("hero_mana_low") {
			hc.handleEvent("hero_mana_low", map[string]interface{}{
				"mana":      mana,
				"prev_mana": hc.lastMana,
				"value":     manaThreshold,
			})
		}
	}

	// Check for level up
	if level > hc.lastLevel && hc.lastLevel > 0 {
		if hc.isEventEnabled("hero_level_up") {
			hc.handleEvent("hero_level_up", map[string]interface{}{
				"level":      level,
				"prev_level": hc.lastLevel,
				"level_diff": level - hc.lastLevel,
			})
		}

		// Check for ultimate ready at level 6
		if level == 6 && hc.lastLevel < 6 && hc.isEventEnabled("hero_ultimate_ready") {
			hc.handleEvent("hero_ultimate_ready", map[string]interface{}{
				"level": level,
			})
		}
	}

	// Update last known values (only if we have valid data)
	if deaths >= 0 {
		hc.lastDeaths = deaths
	}
	if health > 0 {
		hc.lastHealth = health
	}
	if mana > 0 {
		hc.lastMana = mana
	}
	if level > 0 {
		hc.lastLevel = level
	}
}

// canTriggerEvent checks if enough time has passed since last event
func (hc *HeroConsumer) canTriggerEvent(eventType string) bool {
	// Get throttle duration for this event type
	throttleDuration := hc.getThrottleDuration(eventType)
	
	// No throttle if duration is 0
	if throttleDuration == 0 {
		return true
	}
	
	lastTime, exists := hc.eventThrottle[eventType]
	if !exists || time.Since(lastTime) > throttleDuration {
		hc.eventThrottle[eventType] = time.Now()
		return true
	}
	return false
}

// getThrottleDuration returns the throttle duration for an event type
func (hc *HeroConsumer) getThrottleDuration(eventType string) time.Duration {
	if duration, exists := hc.throttleConfig[eventType]; exists {
		return duration
	}
	return 10 * time.Second // Default fallback
}

// isEventEnabled checks if an event is enabled (always true for now)
func (hc *HeroConsumer) isEventEnabled(eventType string) bool {
	return true // All events enabled by default
}

// handleEvent sends event to all handlers
func (hc *HeroConsumer) handleEvent(eventType string, data interface{}) {
	hc.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"data":       data,
	}).Debug("🦸 Hero event detected")

	// Send to all handlers (VoiceHandler, etc.)
	for _, handler := range hc.handlers {
		handler.Handle(eventType, data)
	}
}
