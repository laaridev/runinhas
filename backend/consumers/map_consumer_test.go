package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHandler implements handlers.Handler interface for testing
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(eventType string, data interface{}) {
	m.Called(eventType, data)
}

// Helper function to create test JSON data
func createTestJSON(gameState string, daytime bool, radiantScore, direScore int64) []byte {
	data := map[string]interface{}{
		"map": map[string]interface{}{
			"game_state":    gameState,
			"daytime":       daytime,
			"radiant_score": radiantScore,
			"dire_score":    direScore,
		},
	}
	jsonBytes, _ := json.Marshal(data)
	return jsonBytes
}

// Helper function to create a test logger that doesn't output during tests
func createTestLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel) // Suppress all output during tests
	return logger.WithField("test", true)
}

func TestNewMapConsumer(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	handlerList := []handlers.Handler{mockHandler}

	// Act
	consumer := NewMapConsumer(eventBus, logger, handlerList)

	// Assert
	assert.NotNil(t, consumer)
	assert.Equal(t, logger, consumer.logger)
	assert.Equal(t, handlerList, consumer.handlers)
	assert.NotNil(t, consumer.eventChan)
	assert.NotNil(t, consumer.stopChan)
	assert.Equal(t, "", consumer.lastGameState)
	assert.Equal(t, false, consumer.lastDaytime)
	assert.Equal(t, int64(0), consumer.lastRadScore)
	assert.Equal(t, int64(0), consumer.lastDireScore)
}

func TestMapConsumer_ProcessMapChanges_FirstTick_NoEvents(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// First tick should not trigger any events (lastGameState is empty)
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_HERO_SELECTION", true, 0, 0)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertNotCalled(t, "Handle")
	assert.Equal(t, "DOTA_GAMERULES_STATE_HERO_SELECTION", consumer.lastGameState)
	assert.Equal(t, true, consumer.lastDaytime)
	assert.Equal(t, int64(0), consumer.lastRadScore)
	assert.Equal(t, int64(0), consumer.lastDireScore)
}

func TestMapConsumer_ProcessMapChanges_GameStateChange(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state
	consumer.lastGameState = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	consumer.lastDaytime = true
	consumer.lastRadScore = 0
	consumer.lastDireScore = 0

	// Prepare mock expectation
	expectedData := map[string]interface{}{
		"from": "DOTA_GAMERULES_STATE_HERO_SELECTION",
		"to":   "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS",
	}
	mockHandler.On("Handle", handlers.EventGameStateChange, expectedData).Once()

	// Create tick with game state change
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", true, 0, 0)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertExpectations(t)
	assert.Equal(t, "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", consumer.lastGameState)
}

func TestMapConsumer_ProcessMapChanges_DayNightChange(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state (game already started)
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = true
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	// Prepare mock expectation
	expectedData := map[string]interface{}{
		"daytime": false,
	}
	mockHandler.On("Handle", handlers.EventDayNightChange, expectedData).Once()

	// Create tick with day/night change
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 5, 3)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertExpectations(t)
	assert.Equal(t, false, consumer.lastDaytime)
}

func TestMapConsumer_ProcessMapChanges_ScoreChange_RadiantScores(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state (game already started)
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = false
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	// Prepare mock expectation
	expectedData := map[string]interface{}{
		"radiant_score": int64(6),
		"dire_score":    int64(3),
		"radiant_diff":  int64(1),
		"dire_diff":     int64(0),
	}
	mockHandler.On("Handle", handlers.EventScoreChange, expectedData).Once()

	// Create tick with score change
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 6, 3)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertExpectations(t)
	assert.Equal(t, int64(6), consumer.lastRadScore)
	assert.Equal(t, int64(3), consumer.lastDireScore)
}

func TestMapConsumer_ProcessMapChanges_ScoreChange_DireScores(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state (game already started)
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = false
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	// Prepare mock expectation
	expectedData := map[string]interface{}{
		"radiant_score": int64(5),
		"dire_score":    int64(5),
		"radiant_diff":  int64(0),
		"dire_diff":     int64(2),
	}
	mockHandler.On("Handle", handlers.EventScoreChange, expectedData).Once()

	// Create tick with dire score change
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 5, 5)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertExpectations(t)
	assert.Equal(t, int64(5), consumer.lastRadScore)
	assert.Equal(t, int64(5), consumer.lastDireScore)
}

func TestMapConsumer_ProcessMapChanges_MultipleChanges(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state (game already started)
	consumer.lastGameState = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	consumer.lastDaytime = true
	consumer.lastRadScore = 0
	consumer.lastDireScore = 0

	// Prepare mock expectations for multiple events
	gameStateData := map[string]interface{}{
		"from": "DOTA_GAMERULES_STATE_HERO_SELECTION",
		"to":   "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS",
	}
	dayNightData := map[string]interface{}{
		"daytime": false,
	}
	scoreData := map[string]interface{}{
		"radiant_score": int64(1),
		"dire_score":    int64(1),
		"radiant_diff":  int64(1),
		"dire_diff":     int64(1),
	}

	mockHandler.On("Handle", handlers.EventGameStateChange, gameStateData).Once()
	mockHandler.On("Handle", handlers.EventDayNightChange, dayNightData).Once()
	mockHandler.On("Handle", handlers.EventScoreChange, scoreData).Once()

	// Create tick with multiple changes
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 1, 1)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler.AssertExpectations(t)
	assert.Equal(t, "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", consumer.lastGameState)
	assert.Equal(t, false, consumer.lastDaytime)
	assert.Equal(t, int64(1), consumer.lastRadScore)
	assert.Equal(t, int64(1), consumer.lastDireScore)
}

func TestMapConsumer_ProcessMapChanges_NoChanges(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state (game already started)
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = false
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	// Create tick with no changes
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 5, 3)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert - no handlers should be called
	mockHandler.AssertNotCalled(t, "Handle")
}

func TestMapConsumer_ProcessMapChanges_EmptyGameState(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state - scores match the test JSON to avoid score change events
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	// Create tick with empty game state but same scores
	testJSON := createTestJSON("", false, 5, 3)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert - no events should be triggered because game state becomes empty
	// and scores didn't change
	mockHandler.AssertNotCalled(t, "Handle")
	assert.Equal(t, "", consumer.lastGameState) // State should be updated to empty
}

func TestMapConsumer_ProcessMapChanges_MultipleHandlers(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler1 := &MockHandler{}
	mockHandler2 := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler1, mockHandler2})

	// Set initial state
	consumer.lastGameState = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	consumer.lastDaytime = true

	// Prepare mock expectations for both handlers
	expectedData := map[string]interface{}{
		"from": "DOTA_GAMERULES_STATE_HERO_SELECTION",
		"to":   "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS",
	}
	mockHandler1.On("Handle", handlers.EventGameStateChange, expectedData).Once()
	mockHandler2.On("Handle", handlers.EventGameStateChange, expectedData).Once()

	// Create tick with game state change
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", true, 0, 0)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Act
	consumer.processMapChanges(tickEvent)

	// Assert
	mockHandler1.AssertExpectations(t)
	mockHandler2.AssertExpectations(t)
}

func TestMapConsumer_StartStop_Integration(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Prepare mock expectation
	expectedData := map[string]interface{}{
		"daytime": true,
	}
	mockHandler.On("Handle", handlers.EventDayNightChange, expectedData).Once()

	// Act - Start consumer
	consumer.Start()

	// Set initial state to trigger event
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = false

	// Publish an event
	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", true, 0, 0)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}
	eventBus.Publish(tickEvent)

	// Give some time for event processing
	time.Sleep(10 * time.Millisecond)

	// Stop consumer
	consumer.Stop()

	// Give some time for graceful shutdown
	time.Sleep(10 * time.Millisecond)

	// Assert
	mockHandler.AssertExpectations(t)
}

func TestMapConsumer_HandleEvent(t *testing.T) {
	// Arrange
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	eventType := handlers.EventGameStateChange
	eventData := map[string]interface{}{
		"from": "OLD_STATE",
		"to":   "NEW_STATE",
	}

	// Prepare mock expectation
	mockHandler.On("Handle", eventType, eventData).Once()

	// Act
	consumer.handleEvent(eventType, eventData)

	// Assert
	mockHandler.AssertExpectations(t)
}

// Benchmark test for performance
func BenchmarkMapConsumer_ProcessMapChanges(b *testing.B) {
	// Setup
	eventBus := events.NewEventBus()
	logger := createTestLogger()
	mockHandler := &MockHandler{}
	consumer := NewMapConsumer(eventBus, logger, []handlers.Handler{mockHandler})

	// Set initial state
	consumer.lastGameState = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	consumer.lastDaytime = true
	consumer.lastRadScore = 5
	consumer.lastDireScore = 3

	testJSON := createTestJSON("DOTA_GAMERULES_STATE_GAME_IN_PROGRESS", false, 6, 3)
	tickEvent := events.TickEvent{
		RawJSON: testJSON,
		Time:    time.Now(),
	}

	// Mock handler setup
	mockHandler.On("Handle", mock.Anything, mock.Anything).Maybe()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consumer.processMapChanges(tickEvent)
	}
}
