package config

import (
	"fmt"
	"os"
)

// DefaultGameConfig returns the default game configuration
func DefaultGameConfig() *GameConfig {
	return &GameConfig{
		Timings: map[string]map[string]interface{}{
			"bounty_rune": {
				"enabled":         true,
				"warning_seconds": 30,
			},
			"power_rune": {
				"enabled":         true,
				"warning_seconds": 30,
			},
			"wisdom_rune": {
				"enabled":         true,
				"warning_seconds": 30,
			},
			"water_rune": {
				"enabled":         true,
				"warning_seconds": 30,
			},
			"stack_timing": {
				"enabled":         true,
				"warning_seconds": 20,
			},
			"catapult_timing": {
				"enabled":         true,
				"warning_seconds": 15,
			},
			"day_night_cycle": {
				"enabled":         true,
				"warning_seconds": 20,
			},
			"day_night_transition": {
				"enabled":         true,
				"warning_seconds": 0,
			},
			"tormentor": {
				"enabled": true,
				"time":    20,
			},
			"roshan": {
				"enabled": true,
				"minimum": 30,
				"maximum": 30,
			},
			"glyph": {
				"enabled": true,
				"time":    30,
			},
			"buyback": {
				"enabled": true,
				"time":    30,
			},
			"lotus": {
				"enabled": true,
				"time":    20,
			},
			"outpost": {
				"enabled": true,
				"time":    20,
			},
		},
		Audio: AudioConfig{
			VoiceSpeed: 1.0,
		},
		Messages: map[string]string{
			"bounty_rune":    "Runa de Recompensa em {seconds} segundos",
			"power_rune":     "Runa de Poder em {seconds} segundos",
			"wisdom_rune":    "Runa de Sabedoria em {seconds} segundos",
			"water_rune":     "Runa de Água em {seconds} segundos",
			"stack_timing":   "Stacks em {seconds} segundos",
			"catapult_timing": "Catapulta em {seconds} segundos",
			"day_night_cycle": "Vai {cycle_type} em {seconds} segundos",
			"day_night_transition": "{cycle_type}",
			"tormentor":      "Tormentor em {seconds} segundos",
			"roshan":         "Roshan em {seconds} segundos",
			"glyph":          "Glyph em {seconds} segundos",
			"buyback":        "Buyback em {seconds} segundos",
			"lotus":          "Lotus em {seconds} segundos",
			"outpost":        "Outpost em {seconds} segundos",
		},
		System: &SystemConfig{
			FirstRun:     true,
			GSIInstalled: false,
		},
		Voice: map[string]interface{}{
			"apiKey":       "",
			"voiceId":      "eVXYtPVYB9wDoz9NVTIy",
			"stability":    0.5,
			"similarity":   0.75,
			"style":        0.0,
			"speakerBoost": true,
		},
	}
}

// EnsureConfigExists creates the config file with defaults if it doesn't exist
func EnsureConfigExists() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := DefaultGameConfig()
		
		// Set cache path to system cache directory
		cachePath, err := GetVoiceCachePath()
		if err != nil {
			return fmt.Errorf("failed to get voice cache path: %w", err)
		}
		defaultConfig.Audio.CachePath = cachePath
		
		// Save default config
		if err := SaveGameConfig(configPath, defaultConfig); err != nil {
			return fmt.Errorf("failed to save default config to %s: %w", configPath, err)
		}
		
		// Log config creation
		fmt.Printf("Created default config at: %s\n", configPath)
	}
	
	return nil
}

// LoadOrCreateConfig loads the config file or creates it with defaults
func LoadOrCreateConfig() (*GameConfig, error) {
	// Ensure config exists
	if err := EnsureConfigExists(); err != nil {
		return nil, err
	}
	
	// Load config
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	
	return LoadGameConfig(configPath)
}
