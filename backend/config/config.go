package config

import (
	"fmt"
	"sync"
)

// ============================================================================
// Application Configuration
// ============================================================================
// This file manages the main application configuration including:
// - Server settings (port)
// - ElevenLabs voice configuration (API key, voice ID)
// - Voice cache path
// - Game configuration

// Singleton instance and mutex for thread-safe access
var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port int

	// App mode: "free" or "pro"
	Mode string
	
	// License key for PRO version
	LicenseKey string

	// Voice configuration (Pro only)
	ElevenLabsAPIKey  string
	ElevenLabsVoiceID string
	VoiceCachePath    string
	
	// Game configuration
	Game *GameConfig
}

// Load loads configuration from config.json only (no .env)
func Load() (*Config, error) {
	var loadErr error
	
	once.Do(func() {
		// Get system voice cache path
		voiceCachePath, err := GetVoiceCachePath()
		if err != nil {
			voiceCachePath = "./voice-cache" // Fallback
		}

		// Load game configuration from system config path
		gameConfig, err := LoadOrCreateConfig()
		if err != nil {
			loadErr = err
			return
		}

		// Extract mode from game config (default: free)
		mode := "free"
		if gameConfig.Mode != "" {
			mode = gameConfig.Mode
		}

		// Extract voice config from game config only (no env vars)
		apiKey := ""
		voiceID := DefaultVoiceID
		
		if gameConfig.Voice != nil {
			if key, ok := gameConfig.Voice["apiKey"].(string); ok {
				apiKey = key
			}
			if id, ok := gameConfig.Voice["voiceId"].(string); ok && id != "" {
				voiceID = id
			}
		}

		instance = &Config{
			Port:              3001, // Fixed port
			Mode:              mode,
			ElevenLabsAPIKey:  apiKey,
			ElevenLabsVoiceID: voiceID,
			VoiceCachePath:    voiceCachePath,
			Game:              gameConfig,
		}
	})
	
	if loadErr != nil {
		return nil, loadErr
	}
	
	return instance, nil
}

// GetConfig returns the singleton config instance (thread-safe)
func GetConfig() (*Config, error) {
	return Load()
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	// Voice is optional, but if API key is provided, voice ID should also be provided
	if c.ElevenLabsAPIKey != "" && c.ElevenLabsVoiceID == "" {
		return fmt.Errorf("ELEVENLABS_VOICE_ID is required when ELEVENLABS_API_KEY is provided")
	}

	return nil
}

// HasVoiceConfig returns true if voice configuration is complete
func (c *Config) HasVoiceConfig() bool {
	return c.ElevenLabsAPIKey != "" && c.ElevenLabsVoiceID != ""
}

// UpdateVoiceConfig updates the voice configuration and saves to disk
func (c *Config) UpdateVoiceConfig(apiKey, voiceID string) error {
	mu.Lock()
	defer mu.Unlock()
	
	c.ElevenLabsAPIKey = apiKey
	c.ElevenLabsVoiceID = voiceID
	
	// Update in game config
	if c.Game.Voice == nil {
		c.Game.Voice = make(map[string]interface{})
	}
	c.Game.Voice["apiKey"] = apiKey
	c.Game.Voice["voiceId"] = voiceID
	
	// Save to disk
	configPath, _ := GetConfigPath()
	return SaveGameConfig(configPath, c.Game)
}
