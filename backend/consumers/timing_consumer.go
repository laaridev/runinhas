package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// TimingConsumer handles all timing-based alerts (catapults, stack, day/night)
type TimingConsumer struct {
	logger         *logrus.Entry
	eventChan      <-chan events.TickEvent
	stopChan       chan struct{}
	handlers       []handlers.Handler
	lastAlertTime  map[string]int64 // Track last alert game time for each event type
	lastGameTime   int64
	gameInProgress bool
	isDaytime      bool
	gameConfig     interface{} // Game configuration (can be *config.GameConfig)
}

// NewTimingConsumer creates a new timing consumer
func NewTimingConsumer(eventBus *events.EventBus, logger *logrus.Entry, handlerList []handlers.Handler, gameConfig interface{}) *TimingConsumer {
	return &TimingConsumer{
		logger:        logger,
		eventChan:     eventBus.Subscribe(),
		stopChan:      make(chan struct{}),
		handlers:      handlerList,
		lastAlertTime: make(map[string]int64),
		gameConfig:    gameConfig,
	}
}

// Start begins consuming events
func (tc *TimingConsumer) Start() {
	go tc.consume()
	tc.logger.Info("⏰ TimingConsumer started")
}

// Stop stops the consumer
func (tc *TimingConsumer) Stop() {
	close(tc.stopChan)
	tc.logger.Info("⏰ TimingConsumer stopped")
}

// consume processes TickEvents
func (tc *TimingConsumer) consume() {
	tc.logger.Info("🔄 TimingConsumer loop started, waiting for events...")
	for {
		select {
		case event := <-tc.eventChan:
			tc.logger.Info("⏰ TimingConsumer received event!")
			tc.processTimingEvents(event)
		case <-tc.stopChan:
			tc.logger.Info("🛑 TimingConsumer stopped")
			return
		}
	}
}

// processTimingEvents checks for all timing-based events
func (tc *TimingConsumer) processTimingEvents(event events.TickEvent) {
	jsonData := gjson.ParseBytes(event.RawJSON)
	// Use clock_time instead of game_time (game_time includes pre-game time)
	clockTime := jsonData.Get("map.clock_time").Int()
	gameState := jsonData.Get("map.game_state").String()
	daytime := jsonData.Get("map.daytime").Bool()

	// Only process if game is in progress
	tc.gameInProgress = gameState == "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	if !tc.gameInProgress || clockTime < 0 {
		return
	}

	// Track day/night for warnings
	tc.isDaytime = daytime

	// Skip if no time change
	if clockTime == tc.lastGameTime {
		return
	}

	// Check all timing events
	tc.checkCatapultWarning(clockTime)
	tc.checkDayNightWarning(clockTime, daytime)
	tc.checkStackTiming(clockTime)

	tc.lastGameTime = clockTime
}

// checkCatapultWarning checks for upcoming catapult waves
func (tc *TimingConsumer) checkCatapultWarning(gameTime int64) {
	if !tc.isEventEnabled("catapult_timing") {
		return
	}

	cfg := tc.getTimingConfig("catapult_timing")
	if cfg == nil {
		return
	}

	// Get warning seconds from config (only field that exists)
	warningSeconds := int64(15) // Default from defaults.go
	if val, ok := cfg["warning_seconds"].(float64); ok {
		warningSeconds = int64(val)
	}

	// Catapults spawn every 5 minutes (300s) - hardcoded game rule
	interval := int64(300)

	// Catapults spawn every 5 minutes (300s)
	// Calculate time until next catapult spawn
	timeUntilNextCatapult := interval - (gameTime % interval)

	// If we're within warning time and haven't alerted yet
	if timeUntilNextCatapult <= warningSeconds {
		// Calculate the actual next spawn time for tracking
		nextSpawn := gameTime + timeUntilNextCatapult
		alertKey := fmt.Sprintf("catapult_warning_%d", nextSpawn)

		if tc.lastAlertTime[alertKey] != nextSpawn {
			tc.handleEvent("catapult_timing", map[string]interface{}{
				"seconds":      timeUntilNextCatapult,
				"spawn_time":   nextSpawn,
				"current_time": gameTime,
			})
			tc.lastAlertTime[alertKey] = nextSpawn
		}
	}
}

// checkDayNightWarning checks for day/night transitions and warns before they happen
func (tc *TimingConsumer) checkDayNightWarning(gameTime int64, daytime bool) {
	if !tc.isEventEnabled("day_night_cycle") {
		return
	}

	cfg := tc.getTimingConfig("day_night_cycle")
	if cfg == nil {
		return
	}

	// Get warning seconds from config (only field that exists)
	warningSeconds := int64(20) // Default from defaults.go
	if val, ok := cfg["warning_seconds"].(float64); ok {
		warningSeconds = int64(val)
	}

	// Day/Night cycle is 5 minutes (300s) - hardcoded game rule
	cycleDuration := int64(300)

	// Day/Night cycle: 0-300 (day), 300-600 (night), 600-900 (day), etc.
	// Calculate time until next transition
	timeInCycle := gameTime % cycleDuration
	timeUntilTransition := cycleDuration - timeInCycle

	// Determine what's coming next based on current state
	var nextTransitionType string
	if daytime {
		nextTransitionType = "night" // Currently day, night is coming
	} else {
		nextTransitionType = "day" // Currently night, day is coming
	}

	// Check if transition just happened (within 2 seconds)
	if timeInCycle <= 2 {
		alertKey := fmt.Sprintf("day_night_transition_%d", gameTime/cycleDuration)
		if tc.lastAlertTime[alertKey] == 0 {
			// Announce the transition that just happened
			var transitionType string
			if daytime {
				transitionType = "day" // Just became day
			} else {
				transitionType = "night" // Just became night
			}
			
			eventData := map[string]interface{}{
				"current_time": gameTime,
				"cycle_type":   transitionType,
				"transition":   true, // Flag to indicate this is the transition itself
			}
			
			tc.handleEvent("day_night_transition", eventData)
			tc.lastAlertTime[alertKey] = 1
		}
	}

	// If we're within warning time and haven't alerted yet
	if timeUntilTransition <= warningSeconds && timeUntilTransition > 2 {
		// Calculate the actual next transition time for tracking
		nextTransition := gameTime + timeUntilTransition
		alertKey := fmt.Sprintf("day_night_warning_%d", nextTransition)

		if tc.lastAlertTime[alertKey] != nextTransition {
			eventData := map[string]interface{}{
				"current_time": gameTime,
				"cycle_type":   nextTransitionType,
				"seconds":      timeUntilTransition,
			}

			tc.handleEvent("day_night_cycle", eventData)
			tc.lastAlertTime[alertKey] = nextTransition
		}
	}
}

// checkStackTiming checks for neutral stack timing windows
func (tc *TimingConsumer) checkStackTiming(gameTime int64) {
	if !tc.isEventEnabled("stack_timing") {
		return
	}

	cfg := tc.getTimingConfig("stack_timing")
	if cfg == nil {
		return
	}

	// Get warning seconds from config (how many seconds BEFORE X:53 to warn)
	warningSeconds := int64(20) // Default from defaults.go
	if val, ok := cfg["warning_seconds"].(float64); ok {
		warningSeconds = int64(val)
	}

	// Stack timing: Players need to pull at X:53 to stack at X:00
	// warningSeconds from config = how many seconds BEFORE X:53 to warn
	startMinute := int64(4)  // Start warning from minute 4
	stackPullTime := int64(53) // Players must pull at X:53

	// Get current minute and second
	currentMinute := gameTime / 60
	currentSecond := gameTime % 60

	// Calculate when to warn: X:53 minus warningSeconds
	// Example: warningSeconds=7 → warn at X:46 (7 seconds before X:53)
	warnAtSecond := stackPullTime - warningSeconds
	if warnAtSecond < 0 {
		warnAtSecond = 0 // Don't go negative
	}

	// Only alert from minute 4 onwards, at the calculated warning time
	if currentMinute >= startMinute && currentSecond == warnAtSecond {
		// Check throttle (only alert once per minute)
		lastAlert, exists := tc.lastAlertTime["stack_timing"]
		if !exists || gameTime-lastAlert >= 60 {
			// Calculate seconds until stack pull time (X:53)
			secondsUntilStackPull := stackPullTime - currentSecond
			tc.handleEvent("stack_timing", map[string]interface{}{
				"seconds":      secondsUntilStackPull,
				"minute":       currentMinute,
				"current_time": gameTime,
			})
			tc.lastAlertTime["stack_timing"] = gameTime
		}
	}
}

// Helper methods

// isEventEnabled checks if event is enabled in config
func (tc *TimingConsumer) isEventEnabled(eventType string) bool {
	if tc.gameConfig == nil {
		return true // Default to enabled if no config
	}

	// Type assertion to access GameConfig methods
	type GameConfigInterface interface {
		IsTimingEnabled(string) bool
	}

	if gc, ok := tc.gameConfig.(GameConfigInterface); ok {
		return gc.IsTimingEnabled(eventType)
	}

	return true // Default to enabled
}

// getTimingConfig returns timing configuration from GameConfig
func (tc *TimingConsumer) getTimingConfig(eventType string) map[string]interface{} {
	if tc.gameConfig == nil {
		return nil
	}

	// Type assertion to access GameConfig methods
	type GameConfigInterface interface {
		GetTimingConfig(string) map[string]interface{}
	}

	if gc, ok := tc.gameConfig.(GameConfigInterface); ok {
		return gc.GetTimingConfig(eventType)
	}

	return nil
}

// hasAlerted checks if we already alerted for this spawn time
func (tc *TimingConsumer) hasAlerted(eventType string, spawnTime int64) bool {
	key := fmt.Sprintf("%s_%d", eventType, spawnTime)
	lastAlert, exists := tc.lastAlertTime[key]
	return exists && lastAlert == spawnTime
}

// handleEvent sends event to all handlers
func (tc *TimingConsumer) handleEvent(eventType string, data interface{}) {
	tc.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"data":       data,
	}).Debug("⏰ Timing event detected")

	// Send to all handlers
	for _, handler := range tc.handlers {
		handler.Handle(eventType, data)
	}
}
