package events

import (
	"dota-gsi/backend/metrics"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

// ParsedTickEvent wraps TickEvent with cached parsed JSON
type ParsedTickEvent struct {
	*TickEvent
	parsed gjson.Result
	once   sync.Once
}

// NewParsedTickEvent creates a new ParsedTickEvent from a TickEvent
func NewParsedTickEvent(event TickEvent) *ParsedTickEvent {
	return &ParsedTickEvent{
		TickEvent: &event,
	}
}

// Parse returns the parsed JSON result (cached after first call)
func (pe *ParsedTickEvent) Parse() gjson.Result {
	pe.once.Do(func() {
		start := time.Now()
		pe.parsed = gjson.ParseBytes(pe.RawJSON)
		
		// Track parse time for metrics
		metrics.Instance.AddParseTime(time.Since(start))
	})
	return pe.parsed
}

// Get retrieves a value from the parsed JSON
func (pe *ParsedTickEvent) Get(path string) gjson.Result {
	return pe.Parse().Get(path)
}

// GetInt64 retrieves an int64 value from the parsed JSON
func (pe *ParsedTickEvent) GetInt64(path string) int64 {
	return pe.Get(path).Int()
}

// GetString retrieves a string value from the parsed JSON
func (pe *ParsedTickEvent) GetString(path string) string {
	return pe.Get(path).String()
}

// GetBool retrieves a bool value from the parsed JSON
func (pe *ParsedTickEvent) GetBool(path string) bool {
	return pe.Get(path).Bool()
}
