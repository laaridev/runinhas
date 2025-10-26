package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetAppDataDir returns the appropriate app data directory for the current OS
func GetAppDataDir() (string, error) {
	var baseDir string
	
	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\Runinhas
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	case "darwin":
		// macOS: ~/Library/Application Support/Runinhas
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, "Library", "Application Support")
	default:
		// Linux: ~/.config/runinhas
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(home, ".config")
		}
		baseDir = configDir
	}
	
	appDir := filepath.Join(baseDir, "runinhas")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory %s: %w", appDir, err)
	}
	
	return appDir, nil
}

// GetCacheDir returns the appropriate cache directory for the current OS
func GetCacheDir() (string, error) {
	var baseDir string
	
	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\Runinhas\Cache
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		baseDir = filepath.Join(baseDir, "runinhas", "cache")
	case "darwin":
		// macOS: ~/Library/Caches/Runinhas
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, "Library", "Caches", "runinhas")
	default:
		// Linux: ~/.cache/runinhas
		cacheDir := os.Getenv("XDG_CACHE_HOME")
		if cacheDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			cacheDir = filepath.Join(home, ".cache")
		}
		baseDir = filepath.Join(cacheDir, "runinhas")
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}
	
	return baseDir, nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	appDir, err := GetAppDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

// GetVoiceCachePath returns the full path to the voice cache directory
func GetVoiceCachePath() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	voiceCacheDir := filepath.Join(cacheDir, "voice")
	
	// Create voice cache directory if it doesn't exist
	if err := os.MkdirAll(voiceCacheDir, 0755); err != nil {
		return "", err
	}
	
	return voiceCacheDir, nil
}

// GetLogPath returns the full path to the log file
func GetLogPath() (string, error) {
	appDir, err := GetAppDataDir()
	if err != nil {
		return "", err
	}
	logDir := filepath.Join(appDir, "logs")
	
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(logDir, "runinhas.log"), nil
}
