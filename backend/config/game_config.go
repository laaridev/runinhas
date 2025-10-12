package config

import (
	"encoding/json"
	"os"
)

// ============================================================================
// Game Configuration
// ============================================================================
// This file manages game-specific configuration including:
// - Timing configurations (runes, stacks, etc.)
// - Audio settings (cache path, voice speed)
// - Custom messages for events
// - System settings (first run, GSI installed)

// GameConfig holds the game configuration
type GameConfig struct {
	Timings  map[string]map[string]interface{} `json:"timings"`
	Audio    AudioConfig                       `json:"audio"`
	Messages map[string]string                 `json:"messages"`
	System   *SystemConfig                     `json:"system,omitempty"`
	Voice    map[string]interface{}            `json:"voice,omitempty"`
}

// SystemConfig holds system configuration
type SystemConfig struct {
	FirstRun     bool `json:"first_run"`
	GSIInstalled bool `json:"gsi_installed"`
}

// AudioConfig holds audio configuration
type AudioConfig struct {
	CachePath  string  `json:"cache_path"`
	VoiceSpeed float64 `json:"voice_speed"`
}

// LoadGameConfig loads game configuration from config.json
func LoadGameConfig(path string) (*GameConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg GameConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetTimingConfig returns timing configuration for a specific event
func (gc *GameConfig) GetTimingConfig(eventType string) map[string]interface{} {
	if gc.Timings == nil {
		return nil
	}
	return gc.Timings[eventType]
}

// SaveGameConfig saves game configuration to file
func SaveGameConfig(path string, cfg *GameConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// IsTimingEnabled checks if a timing event is enabled
func (gc *GameConfig) IsTimingEnabled(eventType string) bool {
	cfg := gc.GetTimingConfig(eventType)
	if cfg == nil {
		return false
	}
	
	// Check if "enabled" field exists
	if enabled, exists := cfg["enabled"]; exists {
		if enabledBool, ok := enabled.(bool); ok {
			return enabledBool
		}
	}
	
	// If no "enabled" field, assume enabled
	return true
}

// GetMessage returns the message template for an event
func (gc *GameConfig) GetMessage(eventType string) string {
	if msg, exists := gc.Messages[eventType]; exists {
		return msg
	}
	return ""
}

