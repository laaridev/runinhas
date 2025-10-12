package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
	throttleDuration time.Duration
}

// NewHeroConsumer creates a new hero consumer with handlers
func NewHeroConsumer(eventBus *events.EventBus, logger *logrus.Entry, handlerList []handlers.Handler) *HeroConsumer {
	throttleSeconds := 10 // Default throttle

	return &HeroConsumer{
		logger:           logger,
		eventChan:        eventBus.Subscribe(),
		stopChan:         make(chan struct{}),
		handlers:         handlerList,
		eventThrottle:    make(map[string]time.Time),
		throttleDuration: time.Duration(throttleSeconds) * time.Second,
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
	// Parse JSON once for better performance
	jsonData := gjson.ParseBytes(event.RawJSON)

	// Extract hero data from parsed result
	deaths := jsonData.Get("player.deaths").Int()
	health := jsonData.Get("hero.health_percent").Int()
	mana := jsonData.Get("hero.mana_percent").Int()
	level := jsonData.Get("hero.level").Int()

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
	lastTime, exists := hc.eventThrottle[eventType]
	if !exists || time.Since(lastTime) > hc.throttleDuration {
		hc.eventThrottle[eventType] = time.Now()
		return true
	}
	return false
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
