package config

import (
	"dota-gsi/backend/i18n"
	"fmt"
	"os"
)

// DefaultGameConfig returns the default game configuration
func DefaultGameConfig() *GameConfig {
	return &GameConfig{
		Mode:     "free",  // Default mode (free version)
		Language: "pt-BR", // Default language
		Timings: map[string]map[string]interface{}{
			"bounty_rune": {
				"enabled":         true,
				"warning_seconds": DefaultRuneWarning,
			},
			"power_rune": {
				"enabled":         true,
				"warning_seconds": DefaultRuneWarning,
			},
			"wisdom_rune": {
				"enabled":         true,
				"warning_seconds": DefaultRuneWarning,
			},
			"water_rune": {
				"enabled":         true,
				"warning_seconds": DefaultRuneWarning,
			},
			"stack_timing": {
				"enabled":         true,
				"warning_seconds": DefaultStackWarning,
			},
			"catapult_timing": {
				"enabled":         true,
				"warning_seconds": DefaultCatapultWarning,
			},
			"day_night_cycle": {
				"enabled":         true,
				"warning_seconds": DefaultDayNightWarning,
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
			VoiceSpeed: DefaultVoiceSpeed,
		},
		Messages: map[string]string{
			"bounty_rune":     i18n.T("messages.bounty_rune", map[string]interface{}{"seconds": "{seconds}"}),
			"power_rune":      i18n.T("messages.power_rune", map[string]interface{}{"seconds": "{seconds}"}),
			"wisdom_rune":     i18n.T("messages.wisdom_rune", map[string]interface{}{"seconds": "{seconds}"}),
			"water_rune":      i18n.T("messages.water_rune", map[string]interface{}{"seconds": "{seconds}"}),
			"stack_timing":    i18n.T("messages.stack_timing", map[string]interface{}{"seconds": "{seconds}"}),
			"catapult_timing": i18n.T("messages.catapult_timing", map[string]interface{}{"seconds": "{seconds}"}),
			"day_night_cycle": i18n.T("messages.day_night_cycle", map[string]interface{}{"seconds": "{seconds}"}),
		},
		System: &SystemConfig{
			FirstRun:     DefaultFirstRun,
			GSIInstalled: DefaultGSIInstalled,
		},
		Voice: map[string]interface{}{
			"apiKey":       "",
			"voiceId":      DefaultVoiceID,
			"stability":    DefaultStability,
			"similarity":   DefaultSimilarity,
			"style":        DefaultStyle,
			"speakerBoost": DefaultSpeakerBoost,
		},
		Events: map[string]TimingEvent{
			"bounty_rune": {
				Enabled:        true,
				WarningSeconds: DefaultRuneWarning,
				Min:            10,
				Max:            90,
				Step:           5,
				Name:           "Runa de Recompensa",
				Description:    "Spawns de ouro para todo o time (0:00, depois a cada 3min)",
				Category:       "rune",
			},
			"power_rune": {
				Enabled:        true,
				WarningSeconds: DefaultRuneWarning,
				Min:            10,
				Max:            90,
				Step:           5,
				Name:           "Runa de Poder",
				Description:    "Runas de utilidade ou dano no rio (a cada 2min)",
				Category:       "rune",
			},
			"wisdom_rune": {
				Enabled:        true,
				WarningSeconds: DefaultRuneWarning,
				Min:            10,
				Max:            90,
				Step:           5,
				Name:           "Runa de Sabedoria",
				Description:    "XP bônus para suporte e offlane (7:00, depois a cada 7min)",
				Category:       "rune",
			},
			"water_rune": {
				Enabled:        true,
				WarningSeconds: DefaultRuneWarning,
				Min:            10,
				Max:            50,
				Step:           5,
				Name:           "Runa de Água",
				Description:    "Regeneração instantânea de HP/Mana (2:00 e 4:00)",
				Category:       "rune",
			},
			"stack_timing": {
				Enabled:        true,
				WarningSeconds: DefaultStackWarning,
				Min:            5,
				Max:            15,
				Step:           1,
				Name:           "Stack Timing",
				Description:    "Aviso para stackar camps de neutrals (sempre aos :53)",
				Category:       "timing",
			},
			"catapult_timing": {
				Enabled:        true,
				WarningSeconds: DefaultCatapultWarning,
				Min:            10,
				Max:            30,
				Step:           5,
				Name:           "Catapulta",
				Description:    "Spawn de catapultas para pressão em lanes e push",
				Category:       "timing",
			},
			"day_night_cycle": {
				Enabled:        true,
				WarningSeconds: DefaultDayNightWarning,
				Min:            10,
				Max:            30,
				Step:           5,
				Name:           "Ciclo Dia/Noite",
				Description:    "Alertas de mudança dia/noite para timing estratégico",
				Category:       "timing",
			},
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
