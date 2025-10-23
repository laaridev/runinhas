package metrics

import (
	"sync/atomic"
	"time"
)

// Metrics holds application metrics
type Metrics struct {
	EventsProcessed uint64
	EventsDropped   uint64
	CacheHits       uint64
	CacheMisses     uint64
	StartTime       time.Time
	ParseCount      uint64
	ParseTime       uint64 // nanoseconds
}

// Instance is the global metrics singleton
var Instance = &Metrics{
	StartTime: time.Now(),
}

// IncrementProcessed increments processed events counter
func (m *Metrics) IncrementProcessed() {
	atomic.AddUint64(&m.EventsProcessed, 1)
}

// IncrementDropped increments dropped events counter
func (m *Metrics) IncrementDropped() {
	atomic.AddUint64(&m.EventsDropped, 1)
}

// IncrementCacheHit increments cache hits counter
func (m *Metrics) IncrementCacheHit() {
	atomic.AddUint64(&m.CacheHits, 1)
}

// IncrementCacheMiss increments cache misses counter
func (m *Metrics) IncrementCacheMiss() {
	atomic.AddUint64(&m.CacheMisses, 1)
}

// AddParseTime adds parse time to metrics
func (m *Metrics) AddParseTime(duration time.Duration) {
	atomic.AddUint64(&m.ParseCount, 1)
	atomic.AddUint64(&m.ParseTime, uint64(duration.Nanoseconds()))
}

// GetStats returns current metrics as map
func (m *Metrics) GetStats() map[string]interface{} {
	parseCount := atomic.LoadUint64(&m.ParseCount)
	parseTime := atomic.LoadUint64(&m.ParseTime)
	
	avgParseTime := float64(0)
	if parseCount > 0 {
		avgParseTime = float64(parseTime) / float64(parseCount) / 1000000 // Convert to ms
	}
	
	return map[string]interface{}{
		"events_processed":    atomic.LoadUint64(&m.EventsProcessed),
		"events_dropped":      atomic.LoadUint64(&m.EventsDropped),
		"cache_hits":          atomic.LoadUint64(&m.CacheHits),
		"cache_misses":        atomic.LoadUint64(&m.CacheMisses),
		"uptime_seconds":      time.Since(m.StartTime).Seconds(),
		"parse_count":         parseCount,
		"avg_parse_time_ms":   avgParseTime,
		"drop_rate":           m.GetDropRate(),
	}
}

// GetDropRate returns the percentage of events dropped
func (m *Metrics) GetDropRate() float64 {
	processed := float64(atomic.LoadUint64(&m.EventsProcessed))
	dropped := float64(atomic.LoadUint64(&m.EventsDropped))
	total := processed + dropped
	
	if total == 0 {
		return 0
	}
	
	return (dropped / total) * 100
}
