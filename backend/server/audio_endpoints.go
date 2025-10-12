package server

import (
	"dota-gsi/backend/config"
	"dota-gsi/backend/handlers"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AddAudioEndpoints adds audio-related endpoints to the router
func (s *GSIServer) AddAudioEndpoints(router *mux.Router) {
	router.HandleFunc("/api/audio/check/{eventType}", s.handleCheckAudio).Methods("GET")
	router.HandleFunc("/api/audio/generate/{eventType}", s.handleGenerateAudio).Methods("POST")
	router.HandleFunc("/api/audio/preview/{eventType}", s.handlePreviewAudio).Methods("POST")
}

// handleCheckAudio checks if audio file exists for an event
func (s *GSIServer) handleCheckAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["eventType"]

	// Get voice handler (just to check if it's available)
	_, ok := s.voiceHandler.(*handlers.VoiceHandler)
	if !ok {
		http.Error(w, "Voice handler not available", http.StatusServiceUnavailable)
		return
	}

	// Check if audio exists using semantic cache path
	cachePath, _ := config.GetVoiceCachePath()
	filename := getSemanticFilename(eventType)
	audioPath := filepath.Join(cachePath, filename)

	_, err := os.Stat(audioPath)
	exists := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}

// handleGenerateAudio generates audio for a specific event - ALWAYS regenerates
func (s *GSIServer) handleGenerateAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["eventType"]

	s.logger.WithField("eventType", eventType).Info("🎤 Starting audio generation")

	// Get voice handler
	vh, ok := s.voiceHandler.(*handlers.VoiceHandler)
	if !ok {
		s.logger.Error("Voice handler not available")
		http.Error(w, "Voice handler not available", http.StatusServiceUnavailable)
		return
	}

	// ALWAYS reload config to get latest messages and timings
	cfg, err := config.Load()
	if err != nil {
		s.logger.WithError(err).Error("Failed to load config")
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	// Get warning seconds from config
	var warningSeconds int
	if cfg.Game != nil && cfg.Game.Timings != nil {
		if timing, ok := cfg.Game.Timings[eventType]; ok {
			if ws, ok := timing["warning_seconds"].(float64); ok {
				warningSeconds = int(ws)
			} else if ws, ok := timing["warning_seconds"].(int); ok {
				warningSeconds = ws
			}
		}
	}

	// Get message from config
	var message string
	if cfg.Game != nil && cfg.Game.Messages != nil {
		if msg, ok := cfg.Game.Messages[eventType]; ok {
			message = msg
		}
	}

	// Replace placeholders
	if message != "" {
		message = strings.ReplaceAll(message, "{seconds}", fmt.Sprintf("%d", warningSeconds))
		message = strings.ReplaceAll(message, "{time}", fmt.Sprintf("%d segundos", warningSeconds))
	} else {
		// Fallback message
		message = fmt.Sprintf("%s em %d segundos", eventType, warningSeconds)
	}

	s.logger.WithFields(logrus.Fields{
		"eventType": eventType,
		"message":   message,
		"seconds":   warningSeconds,
	}).Info("📝 Message prepared for generation")

	// Generate audio directly with ElevenLabs
	audioData, err := vh.GenerateVoice(message)
	if err != nil {
		s.logger.WithError(err).Error("❌ Failed to generate voice")
		http.Error(w, fmt.Sprintf("Failed to generate voice: %v", err), http.StatusInternalServerError)
		return
	}

	// Get cache path and save
	cachePath, _ := config.GetVoiceCachePath()
	filename := getSemanticFilename(eventType)
	audioPath := filepath.Join(cachePath, filename)

	// Ensure directory exists
	os.MkdirAll(cachePath, 0755)

	// Save the new audio (overwrite if exists)
	if err := os.WriteFile(audioPath, audioData, 0644); err != nil {
		s.logger.WithError(err).Error("❌ Failed to save audio file")
		http.Error(w, "Failed to save audio", http.StatusInternalServerError)
		return
	}

	s.logger.WithField("audioPath", audioPath).Info("✅ Audio generated and saved successfully")

	// Play the audio
	go vh.PlayAudioFile(audioPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "generated",
		"path":    audioPath,
		"message": message,
	})
}

// handlePreviewAudio plays audio for an event
func (s *GSIServer) handlePreviewAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["eventType"]

	// Get voice handler
	vh, ok := s.voiceHandler.(*handlers.VoiceHandler)
	if !ok {
		http.Error(w, "Voice handler not available", http.StatusServiceUnavailable)
		return
	}

	// Play the audio file
	cachePath, _ := config.GetVoiceCachePath()
	filename := getSemanticFilename(eventType)
	audioPath := filepath.Join(cachePath, filename)

	// Check if file exists
	if _, err := os.Stat(audioPath); err != nil {
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	// Play using the handler's method
	vh.PlayAudioFile(audioPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "playing"})
}

// getSemanticFilename returns the semantic filename for an event type
func getSemanticFilename(eventType string) string {
	// Add _warning suffix if not present (for compatibility)
	if !strings.HasSuffix(eventType, "_warning") &&
		(strings.Contains(eventType, "rune") ||
			eventType == "catapult" ||
			eventType == "stack_timing" ||
			eventType == "day_night") {
		eventType = eventType + "_warning"
	}

	// Sanitize and return
	filename := strings.ToLower(eventType)
	filename = strings.ReplaceAll(filename, " ", "_")
	return filename + ".mp3"
}
