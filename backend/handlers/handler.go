package handlers

// Handler interface for processing domain events
type Handler interface {
	Handle(eventType string, data interface{})
}

// EventData represents processed domain data
type EventData struct {
	Type    string      `json:"type"`    // "game_state_change", "hero_health_low", etc.
	Payload interface{} `json:"payload"` // Specific data for the event
}

// Common event types
const (
	// Map events
	EventGameStateChange = "game_state_change"
	EventDayNightChange  = "day_night_change"
	EventScoreChange     = "score_change"

	// Hero events
	EventHeroHealthLow = "hero_health_low"
	EventHeroManaLow   = "hero_mana_low"
	EventHeroDeath     = "hero_death"
	EventHeroLevelUp   = "hero_level_up"

	// Abilities events
	EventUltimateReady   = "ultimate_ready"
	EventAbilityCooldown = "ability_cooldown"

	// Items events
	EventTPReady       = "tp_ready"
	EventItemPurchased = "item_purchased"
)
