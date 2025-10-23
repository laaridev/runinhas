package events

import (
	"dota-gsi/backend/metrics"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TickEvent represents a GSI tick with raw JSON for selective parsing
type TickEvent struct {
	RawJSON []byte    // Raw JSON from GSI - each consumer extracts what it needs
	Time    time.Time // When the tick was received
}

// EventBus broadcasts TickEvents to multiple consumers
type EventBus struct {
	subscribers []chan TickEvent
	mutex       sync.RWMutex
	logger      *logrus.Entry
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	logger := logrus.WithField("component", "event-bus")
	return &EventBus{
		subscribers: make([]chan TickEvent, 0),
		logger:      logger,
	}
}

// Subscribe returns a channel to receive TickEvents
func (eb *EventBus) Subscribe() <-chan TickEvent {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	// Increased buffer size to reduce drops
	ch := make(chan TickEvent, 100)
	eb.subscribers = append(eb.subscribers, ch)
	
	return ch
}

// Publish broadcasts a TickEvent to all subscribers
func (eb *EventBus) Publish(event TickEvent) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	// Broadcast with metrics tracking
	for _, subscriber := range eb.subscribers {
		select {
		case subscriber <- event:
			// Event delivered successfully
			metrics.Instance.IncrementProcessed()
		default:
			// Channel full, log and track
			eb.logger.Warn("Event dropped - channel buffer full")
			metrics.Instance.IncrementDropped()
		}
	}
}

// Close shuts down the event bus
func (eb *EventBus) Close() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	for _, subscriber := range eb.subscribers {
		close(subscriber)
	}
	eb.subscribers = nil
}
