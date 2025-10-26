package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"fmt"

	"github.com/sirupsen/logrus"
)

// toInt64 safely converts interface{} to int64, handling multiple numeric types
func toInt64(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

// RuneConsumer monitors game time and alerts before rune spawns
type RuneConsumer struct {
	logger           *logrus.Entry
	lastGameTime     int64
	eventChan        <-chan events.TickEvent
	stopChan         chan struct{}
	handlers         []handlers.Handler
	lastAlertedRunes map[string]int64 // Track last alerted time for each rune type
	gameConfig       interface{}      // Game configuration
}

// NewRuneConsumer creates a new rune consumer
func NewRuneConsumer(eventBus *events.EventBus, logger *logrus.Entry, handlerList []handlers.Handler, gameConfig interface{}) *RuneConsumer {
	return &RuneConsumer{
		logger:           logger,
		eventChan:        eventBus.Subscribe(),
		stopChan:         make(chan struct{}),
		handlers:         handlerList,
		lastAlertedRunes: make(map[string]int64),
		gameConfig:       gameConfig,
	}
}

// Start begins consuming events
func (rc *RuneConsumer) Start() {
	go rc.consume()
	rc.logger.Info("💎 RuneConsumer started")
}

// Stop stops the consumer
func (rc *RuneConsumer) Stop() {
	close(rc.stopChan)
	rc.logger.Info("💎 RuneConsumer stopped")
}

// consume processes TickEvents
func (rc *RuneConsumer) consume() {
	rc.logger.Info("🔄 RuneConsumer loop started, waiting for events...")
	for {
		select {
		case event := <-rc.eventChan:
			rc.logger.Info("💎 RuneConsumer received event!")
			rc.processRuneTimings(event)
		case <-rc.stopChan:
			rc.logger.Info("🛑 RuneConsumer stopped")
			return
		}
	}
}

// processRuneTimings checks game time and alerts for upcoming runes
func (rc *RuneConsumer) processRuneTimings(event events.TickEvent) {
	// Use parsed event for efficient JSON access
	parsed := events.NewParsedTickEvent(event)
	
	// Use clock_time instead of game_time (game_time includes pre-game time)
	clockTime := parsed.GetInt64("map.clock_time")

	// Only process if game is in progress
	gameState := parsed.GetString("map.game_state")
	if gameState != "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS" {
		return
	}

	// Skip if no time change
	if clockTime == rc.lastGameTime {
		return
	}

	// Get warning seconds from config for each rune type
	bountyWarning := rc.getWarningSeconds("bounty_rune")
	powerWarning := rc.getWarningSeconds("power_rune")
	waterWarning := rc.getWarningSeconds("water_rune")
	wisdomWarning := rc.getWarningSeconds("wisdom_rune")

	// Check bounty runes (first at 0:00, then every 3 minutes)
	rc.checkBountyRunes(clockTime, bountyWarning)

	// Check power runes (first at 6:00, then every 2 minutes)
	rc.checkPowerRunes(clockTime, powerWarning)

	// Check water runes (only at 2:00 and 4:00)
	rc.checkWaterRunes(clockTime, waterWarning)

	// Check wisdom runes (first at 7:00, then every 7 minutes)
	rc.checkWisdomRunes(clockTime, wisdomWarning)

	rc.lastGameTime = clockTime
}

// checkBountyRunes checks for bounty rune spawns
func (rc *RuneConsumer) checkBountyRunes(gameTime, warningSeconds int64) {
	// Bounty runes spawn at 0:00 and every 3 minutes (180s)
	interval := int64(180)

	// Calculate time until next rune spawn
	timeUntilNextRune := interval - (gameTime % interval)

	// If we're within warning time and haven't alerted yet
	if timeUntilNextRune <= warningSeconds {
		// Calculate the actual next spawn time for tracking
		nextSpawn := gameTime + timeUntilNextRune

		if rc.lastAlertedRunes["bounty"] != nextSpawn {
			rc.handleEvent("bounty_rune", map[string]interface{}{
				"seconds":    timeUntilNextRune,
				"spawn_time": nextSpawn,
				"rune_type":  "bounty",
			})
			rc.lastAlertedRunes["bounty"] = nextSpawn
		}
	} else {
		// Reset alert flag when we're far from spawn time
		if timeUntilNextRune > warningSeconds {
			rc.lastAlertedRunes["bounty"] = 0
		}
	}
}

// checkPowerRunes checks for power rune spawns
func (rc *RuneConsumer) checkPowerRunes(gameTime, warningSeconds int64) {
	// Power runes spawn at 6:00 and every 2 minutes
	firstSpawn := int64(360)
	interval := int64(120)

	// Power runes start at 6:00, don't check before that
	if gameTime < firstSpawn-warningSeconds {
		return
	}

	// Calculate time until next rune spawn
	var timeUntilNextRune int64
	if gameTime < firstSpawn {
		// Before first spawn
		timeUntilNextRune = firstSpawn - gameTime
	} else {
		// After first spawn, calculate based on interval
		timeSinceFirst := gameTime - firstSpawn
		timeUntilNextRune = interval - (timeSinceFirst % interval)
	}

	// If we're within warning time and haven't alerted yet
	if timeUntilNextRune <= warningSeconds {
		// Calculate the actual next spawn time for tracking
		nextSpawn := gameTime + timeUntilNextRune

		if rc.lastAlertedRunes["power"] != nextSpawn {
			rc.handleEvent("power_rune", map[string]interface{}{
				"seconds":    timeUntilNextRune,
				"spawn_time": nextSpawn,
				"rune_type":  "power",
			})
			rc.lastAlertedRunes["power"] = nextSpawn
		}
	} else {
		// Reset alert flag when we're far from spawn time
		if timeUntilNextRune > warningSeconds {
			rc.lastAlertedRunes["power"] = 0
		}
	}
}

// checkWaterRunes checks for water rune spawns (only at 2:00 and 4:00)
func (rc *RuneConsumer) checkWaterRunes(gameTime, warningSeconds int64) {
	// Water runes spawn only at 2:00 (120s) and 4:00 (240s)
	spawnTimes := []int64{120, 240}

	for _, spawnTime := range spawnTimes {
		timeUntilSpawn := spawnTime - gameTime

		// If we're within warning time and haven't alerted yet
		if timeUntilSpawn > 0 && timeUntilSpawn <= warningSeconds {
			alertKey := fmt.Sprintf("water_%d", spawnTime)
			if rc.lastAlertedRunes[alertKey] != spawnTime {
				rc.handleEvent("water_rune", map[string]interface{}{
					"seconds":    timeUntilSpawn,
					"spawn_time": spawnTime,
					"rune_type":  "water",
				})
				rc.lastAlertedRunes[alertKey] = spawnTime
			}
		}

		// Reset alert flag after spawn time
		if gameTime > spawnTime {
			alertKey := fmt.Sprintf("water_%d", spawnTime)
			rc.lastAlertedRunes[alertKey] = 0
		}
	}
}

// checkWisdomRunes checks for wisdom rune spawns
func (rc *RuneConsumer) checkWisdomRunes(gameTime, warningSeconds int64) {
	// Wisdom runes spawn at 7:00 and every 7 minutes (420s)
	firstSpawn := int64(420)
	interval := int64(420)

	// Don't check before first spawn
	if gameTime < firstSpawn-warningSeconds {
		return
	}

	// Calculate time until next rune spawn
	var timeUntilNextRune int64
	if gameTime < firstSpawn {
		// Before first spawn
		timeUntilNextRune = firstSpawn - gameTime
	} else {
		// After first spawn, calculate based on interval
		timeSinceFirst := gameTime - firstSpawn
		timeUntilNextRune = interval - (timeSinceFirst % interval)
	}

	// If we're within warning time and haven't alerted yet
	if timeUntilNextRune <= warningSeconds {
		// Calculate the actual next spawn time for tracking
		nextSpawn := gameTime + timeUntilNextRune

		if rc.lastAlertedRunes["wisdom"] != nextSpawn {
			rc.handleEvent("wisdom_rune", map[string]interface{}{
				"seconds":    timeUntilNextRune,
				"spawn_time": nextSpawn,
				"rune_type":  "wisdom",
			})
			rc.lastAlertedRunes["wisdom"] = nextSpawn
		}
	} else {
		// Reset alert flag when we're far from spawn time
		if timeUntilNextRune > warningSeconds {
			rc.lastAlertedRunes["wisdom"] = 0
		}
	}
}

// getWarningSeconds returns warning seconds from config
func (rc *RuneConsumer) getWarningSeconds(runeType string) int64 {
	if rc.gameConfig == nil {
		return int64(30) // Default if no config
	}

	// Type assertion to access GameConfig methods
	type GameConfigInterface interface {
		GetTimingConfig(string) map[string]interface{}
	}

	if gc, ok := rc.gameConfig.(GameConfigInterface); ok {
		if cfg := gc.GetTimingConfig(runeType); cfg != nil {
			if ws, exists := cfg["warning_seconds"]; exists {
				if val, ok := toInt64(ws); ok {
					return val
				}
			}
		}
	}

	return int64(30) // Default fallback
}

// isEventEnabled always returns true for now
func (rc *RuneConsumer) isEventEnabled(eventType string) bool {
	return true // All events enabled by default
}

// handleEvent sends event to all handlers
func (rc *RuneConsumer) handleEvent(eventType string, data interface{}) {
	rc.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"data":       data,
	}).Debug("💎 Rune event detected")

	// Send to all handlers
	for _, handler := range rc.handlers {
		handler.Handle(eventType, data)
	}
}
