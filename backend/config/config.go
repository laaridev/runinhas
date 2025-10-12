package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ============================================================================
// Application Configuration
// ============================================================================
// This file manages the main application configuration including:
// - Server settings (port)
// - ElevenLabs voice configuration (API key, voice ID)
// - Voice cache path
// - Game configuration

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port int

	// Voice configuration
	ElevenLabsAPIKey  string
	ElevenLabsVoiceID string
	VoiceCachePath    string

	// Paths
	ConfigPath string
	
	// Game configuration
	Game *GameConfig
}

// Load loads configuration from environment variables and config.json
func Load() (*Config, error) {
	// Try to load .env file from root directory (optional)
	_ = godotenv.Load(".env")

	// Get system voice cache path
	voiceCachePath, err := GetVoiceCachePath()
	if err != nil {
		voiceCachePath = "./voice-cache" // Fallback
	}

	// Load game configuration from system config path first
	gameConfig, err := LoadOrCreateConfig()
	if err != nil {
		return nil, err
	}

	// Extract voice config from game config if available
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	voiceID := os.Getenv("ELEVENLABS_VOICE_ID")
	
	// Prefer config.json values over environment variables
	if gameConfig.Voice != nil {
		if key, ok := gameConfig.Voice["apiKey"].(string); ok && key != "" {
			apiKey = key
		}
		if id, ok := gameConfig.Voice["voiceId"].(string); ok && id != "" {
			voiceID = id
		}
	}

	cfg := &Config{
		Port:              getEnvAsInt("SERVER_PORT", 3001),
		ElevenLabsAPIKey:  apiKey,
		ElevenLabsVoiceID: voiceID,
		VoiceCachePath:    getEnv("VOICE_CACHE_PATH", voiceCachePath),
		ConfigPath:        "internal/config/.env",
		Game:              gameConfig,
	}

	return cfg, nil
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

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
