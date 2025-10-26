package config

// ============================================================================
// Application Constants
// ============================================================================
// This file contains all default values and magic numbers used throughout
// the application. These are the single source of truth for default values.

const (
	// Server defaults
	DefaultPort = 3001

	// Voice defaults
	DefaultVoiceID      = "eVXYtPVYB9wDoz9NVTIy"
	DefaultStability    = 0.5
	DefaultSimilarity   = 0.75
	DefaultStyle        = 0.0
	DefaultSpeakerBoost = true
	DefaultVoiceSpeed   = 1.0

	// Timing defaults (used as fallbacks if config is not set)
	DefaultCatapultWarning  = 15
	DefaultDayNightWarning  = 20
	DefaultStackWarning     = 20
	DefaultRuneWarning      = 30

	// System defaults
	DefaultFirstRun     = true
	DefaultGSIInstalled = false
)
