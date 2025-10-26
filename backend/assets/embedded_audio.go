package assets

import (
	"embed"
	"io/fs"
)

// ============================================================================
// Embedded Audio Files - Free Version
// ============================================================================
// This package embeds default audio files into the binary for the free version.
// These are generic messages like "Runa de Poder em alguns segundos"

//go:embed audio
var AudioFiles embed.FS

// GetAudioFile returns the content of an embedded audio file
func GetAudioFile(filename string) ([]byte, error) {
	return fs.ReadFile(AudioFiles, "audio/"+filename)
}

// ListAudioFiles returns all available audio files
func ListAudioFiles() ([]string, error) {
	entries, err := fs.ReadDir(AudioFiles, "audio")
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// HasAudioFile checks if an audio file exists in the embedded files
func HasAudioFile(filename string) bool {
	_, err := fs.Stat(AudioFiles, "audio/"+filename)
	return err == nil
}
