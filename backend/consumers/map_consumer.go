package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"

	"github.com/sirupsen/logrus"
)

// MapConsumer processes map-related events (game state, day/night, score)
type MapConsumer struct {
	logger        *logrus.Entry
	lastGameState string
	lastDaytime   bool
	lastRadScore  int64
	lastDireScore int64
	eventChan     <-chan events.TickEvent
	stopChan      chan struct{}
	handlers      []handlers.Handler
}

// NewMapConsumer creates a new map consumer with handlers
func NewMapConsumer(eventBus *events.EventBus, logger *logrus.Entry, handlerList []handlers.Handler) *MapConsumer {
	return &MapConsumer{
		logger:    logger,
		eventChan: eventBus.Subscribe(),
		stopChan:  make(chan struct{}),
		handlers:  handlerList,
	}
}

// Start begins consuming events
func (mc *MapConsumer) Start() {
	go mc.consume()
	mc.logger.Info("🗺️ MapConsumer started")
}

// Stop stops the consumer
func (mc *MapConsumer) Stop() {
	close(mc.stopChan)
	mc.logger.Info("🗺️ MapConsumer stopped")
}

// consume processes TickEvents and detects map changes
func (mc *MapConsumer) consume() {
	for {
		select {
		case event := <-mc.eventChan:
			mc.processMapChanges(event)
		case <-mc.stopChan:
			return
		}
	}
}

// processMapChanges extracts map data and detects changes
func (mc *MapConsumer) processMapChanges(event events.TickEvent) {
	// Use parsed event for efficient JSON access
	parsed := events.NewParsedTickEvent(event)
	
	// Extract map data using cached parse
	gameState := parsed.GetString("map.game_state")
	daytime := parsed.GetBool("map.daytime")
	radiantScore := parsed.GetInt64("map.radiant_score")
	direScore := parsed.GetInt64("map.dire_score")

	// Check for game state changes
	if gameState != "" && gameState != mc.lastGameState && mc.lastGameState != "" {
		if mc.isEventEnabled("game_state_change") {
			mc.handleEvent("game_state_change", map[string]interface{}{
				"from": mc.lastGameState,
				"to":   gameState,
			})
		}
	}

	// Check for day/night changes
	if daytime != mc.lastDaytime && mc.lastGameState != "" {
		if mc.isEventEnabled("day_night_change") {
			mc.handleEvent("day_night_change", map[string]interface{}{
				"daytime": daytime,
			})
		}
	}

	// Check for score changes
	if (radiantScore != mc.lastRadScore || direScore != mc.lastDireScore) && mc.lastGameState != "" {
		if mc.isEventEnabled("score_change") {
			mc.handleEvent("score_change", map[string]interface{}{
				"radiant_score": radiantScore,
				"dire_score":    direScore,
				"radiant_diff":  radiantScore - mc.lastRadScore,
				"dire_diff":     direScore - mc.lastDireScore,
			})
		}
	}

	// Update last known values
	mc.lastGameState = gameState
	mc.lastDaytime = daytime
	mc.lastRadScore = radiantScore
	mc.lastDireScore = direScore
}

// isEventEnabled checks if an event is enabled (always true for now)
func (mc *MapConsumer) isEventEnabled(eventType string) bool {
	return true // All events enabled by default
}

// handleEvent sends event to all handlers
func (mc *MapConsumer) handleEvent(eventType string, data interface{}) {
	mc.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"data":       data,
	}).Debug("🗺️ Map event detected")

	// Send to all handlers (VoiceHandler, etc.)
	for _, handler := range mc.handlers {
		handler.Handle(eventType, data)
	}
}
