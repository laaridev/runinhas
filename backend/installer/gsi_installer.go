package installer

import (
	"dota-gsi/backend/i18n"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// GSIInstaller handles installation of GSI config file
type GSIInstaller struct {
	detector *DotaDetector
	logger   *logrus.Entry
}

// NewGSIInstaller creates a new GSI installer
func NewGSIInstaller(logger *logrus.Entry) *GSIInstaller {
	return &GSIInstaller{
		detector: NewDotaDetector(),
		logger:   logger,
	}
}

// InstallResult contains the result of installation attempt
type InstallResult struct {
	Success     bool
	Message     string
	InstalledAt string
}

// Install attempts to install GSI config file
func (gi *GSIInstaller) Install() InstallResult {
	// Find Dota 2 installation
	dotaPath, err := gi.detector.FindDota2Path()
	if err != nil {
		gi.logger.WithError(err).Error("Failed to find Dota 2 installation")
		return InstallResult{
			Success: false,
			Message: i18n.T("installer.error_not_found", nil),
		}
	}

	gi.logger.WithField("path", dotaPath).Info("Found Dota 2 installation")

	// Get GSI config directory
	gsiDir := gi.detector.GetGSIConfigPath(dotaPath)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(gsiDir, 0755); err != nil {
		gi.logger.WithError(err).Error("Failed to create GSI config directory")
		return InstallResult{
			Success: false,
			Message: i18n.T("installer.error_create_dir", map[string]interface{}{"error": err.Error()}),
		}
	}

	// Generate config file path
	configPath := filepath.Join(gsiDir, "gamestate_integration_dota-gsi.cfg")

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		gi.logger.Info("GSI config file already exists")
		return InstallResult{
			Success: true,
			Message: i18n.T("installer.already_installed", nil),
			InstalledAt: configPath,
		}
	}

	// Write config file
	if err := gi.writeConfigFile(configPath); err != nil {
		gi.logger.WithError(err).Error("Failed to write GSI config file")
		return InstallResult{
			Success: false,
			Message: i18n.T("installer.error_write_file", map[string]interface{}{"error": err.Error()}),
		}
	}

	gi.logger.WithField("path", configPath).Info("GSI config file installed successfully")
	return InstallResult{
		Success: true,
		Message: i18n.T("installer.success", nil),
		InstalledAt: configPath,
	}
}

// writeConfigFile writes the GSI configuration file
func (gi *GSIInstaller) writeConfigFile(path string) error {
	config := `"dota2-gsi Configuration"
{
    "uri"               "http://localhost:3001/gsi"
    "timeout"           "5.0"
    "buffer"            "0.1"
    "throttle"          "0.1"
    "heartbeat"         "30.0"
    "data"
    {
        "provider"      "1"
        "map"           "1"
        "player"        "1"
        "hero"          "1"
        "abilities"     "1"
        "items"         "1"
    }
}
`
	return os.WriteFile(path, []byte(config), 0644)
}

// CheckInstallation checks if GSI is already installed
func (gi *GSIInstaller) CheckInstallation() (bool, string) {
	dotaPath, err := gi.detector.FindDota2Path()
	if err != nil {
		return false, ""
	}

	gsiDir := gi.detector.GetGSIConfigPath(dotaPath)
	configPath := filepath.Join(gsiDir, "gamestate_integration_dota-gsi.cfg")

	if _, err := os.Stat(configPath); err == nil {
		return true, configPath
	}

	return false, ""
}

// CheckDotaInstalled checks if Dota 2 is installed
func (gi *GSIInstaller) CheckDotaInstalled() (bool, string) {
	dotaPath, err := gi.detector.FindDota2Path()
	if err != nil {
		return false, ""
	}
	return true, dotaPath
}
