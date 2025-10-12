package events

import (
	"sync"
	"time"
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
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]chan TickEvent, 0),
	}
}

// Subscribe returns a channel to receive TickEvents
func (eb *EventBus) Subscribe() <-chan TickEvent {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	// Buffered channel prevents blocking publishers
	ch := make(chan TickEvent, 50)
	eb.subscribers = append(eb.subscribers, ch)
	
	return ch
}

// Publish broadcasts a TickEvent to all subscribers
func (eb *EventBus) Publish(event TickEvent) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	// Non-blocking broadcast
	for _, subscriber := range eb.subscribers {
		select {
		case subscriber <- event:
			// Event delivered
		default:
			// Channel full, skip to prevent blocking
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
