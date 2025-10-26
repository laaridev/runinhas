package installer

import (
	"os"
	"path/filepath"
	"runtime"
)

// DotaDetector finds Dota 2 installation directory
type DotaDetector struct {
	steamPaths []string
}

// NewDotaDetector creates a new Dota 2 detector
func NewDotaDetector() *DotaDetector {
	return &DotaDetector{
		steamPaths: getSteamPaths(),
	}
}

// getSteamPaths returns possible Steam installation paths based on OS
func getSteamPaths() []string {
	switch runtime.GOOS {
	case "linux":
		homeDir, _ := os.UserHomeDir()
		return []string{
			filepath.Join(homeDir, ".steam", "steam"),
			filepath.Join(homeDir, ".local", "share", "Steam"),
			"/usr/share/steam",
		}
	case "windows":
		return []string{
			"C:\\Program Files (x86)\\Steam",
			"C:\\Program Files\\Steam",
			"D:\\Steam",
			"E:\\Steam",
		}
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		return []string{
			filepath.Join(homeDir, "Library", "Application Support", "Steam"),
		}
	default:
		return []string{}
	}
}

// FindDota2Path searches for Dota 2 installation directory
func (dd *DotaDetector) FindDota2Path() (string, error) {
	for _, steamPath := range dd.steamPaths {
		// Check if Steam directory exists
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			continue
		}

		// Try common Dota 2 path
		dotaPath := filepath.Join(steamPath, "steamapps", "common", "dota 2 beta")
		if _, err := os.Stat(dotaPath); err == nil {
			return dotaPath, nil
		}

		// Try alternative library folders
		libraryFolders := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
		if additionalPaths := dd.parseLibraryFolders(libraryFolders); len(additionalPaths) > 0 {
			for _, libPath := range additionalPaths {
				dotaPath := filepath.Join(libPath, "steamapps", "common", "dota 2 beta")
				if _, err := os.Stat(dotaPath); err == nil {
					return dotaPath, nil
				}
			}
		}
	}

	return "", os.ErrNotExist
}

// parseLibraryFolders parses Steam's libraryfolders.vdf to find additional library paths
func (dd *DotaDetector) parseLibraryFolders(vdfPath string) []string {
	// Simple VDF parser - looks for "path" entries
	data, err := os.ReadFile(vdfPath)
	if err != nil {
		return nil
	}

	var paths []string
	lines := string(data)
	
	// Very basic parsing - looks for path entries
	// Format: "path"		"C:\\SteamLibrary"
	for i := 0; i < len(lines); i++ {
		if i+6 < len(lines) && lines[i:i+6] == "\"path\"" {
			// Find the value after "path"
			start := i + 6
			for start < len(lines) && (lines[start] == ' ' || lines[start] == '\t') {
				start++
			}
			if start < len(lines) && lines[start] == '"' {
				start++
				end := start
				for end < len(lines) && lines[end] != '"' {
					end++
				}
				if end < len(lines) {
					path := lines[start:end]
					// Convert Windows path separators
					path = filepath.FromSlash(path)
					paths = append(paths, path)
				}
			}
		}
	}

	return paths
}

// GetGSIConfigPath returns the full path to GSI config directory
func (dd *DotaDetector) GetGSIConfigPath(dotaPath string) string {
	return filepath.Join(dotaPath, "game", "dota", "cfg", "gamestate_integration")
}
