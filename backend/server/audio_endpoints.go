package server

import (
	"dota-gsi/backend/assets"
	"dota-gsi/backend/config"
	"dota-gsi/backend/handlers"
	"encoding/base64"
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
	// Serve MP3 files from user cache (PRO mode)
	router.HandleFunc("/api/audio/file/{filename}", s.handleServeAudioFile).Methods("GET")
	// Serve embedded MP3 files (FREE mode)
	router.HandleFunc("/api/audio/embedded/{filename}", s.handleServeEmbeddedAudio).Methods("GET")
	// Get audio file as base64 (for Wails proxy)
	router.HandleFunc("/api/audio/base64/{filename}", s.handleGetAudioBase64).Methods("GET")
	// Stream audio events to frontend
	router.HandleFunc("/api/audio/events", s.handleAudioEvents).Methods("GET")
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

	// Only log if audio doesn't exist (to avoid spam)
	if !exists {
		s.logger.WithFields(logrus.Fields{
			"eventType": eventType,
			"filename":  filename,
		}).Trace("Audio file not found (will be generated on first use)")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}
// handleGenerateAudio generates audio for a specific event - ALWAYS regenerates
func (s *GSIServer) handleGenerateAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["eventType"]

	s.logger.WithField("eventType", eventType).Info("🎯 handleGenerateAudio called")

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
	
	// Debug: Log entire config
	s.logger.WithFields(logrus.Fields{
		"hasGame": cfg.Game != nil,
		"hasTimings": cfg.Game != nil && cfg.Game.Timings != nil,
		"timingsCount": len(cfg.Game.Timings),
	}).Debug("🔍 Config loaded")

	// Get warning seconds from config
	warningSeconds := 30 // Default fallback
	if cfg.Game != nil && cfg.Game.Timings != nil {
		if timing, ok := cfg.Game.Timings[eventType]; ok {
			s.logger.WithFields(logrus.Fields{
				"eventType": eventType,
				"timing": timing,
			}).Debug("🔍 Found timing config")
			
			if ws, ok := timing["warning_seconds"].(float64); ok {
				warningSeconds = int(ws)
				s.logger.WithField("warningSeconds", warningSeconds).Debug("✅ Got warning_seconds as float64")
			} else if ws, ok := timing["warning_seconds"].(int); ok {
				warningSeconds = ws
				s.logger.WithField("warningSeconds", warningSeconds).Debug("✅ Got warning_seconds as int")
			} else {
				s.logger.WithField("timing", timing).Warn("⚠️ warning_seconds not found or wrong type")
			}
		} else {
			s.logger.WithField("eventType", eventType).Warn("⚠️ No timing config found for event")
		}
	} else {
		s.logger.Warn("⚠️ No Game or Timings config found")
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

	// Return success (frontend will play the audio)
	response := map[string]string{
		"status":   "generated",
		"filename": filename,
		"message":  message,
	}
	
	s.logger.WithFields(logrus.Fields{
		"response": response,
	}).Info("📤 Sending response to frontend")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

// handleServeAudioFile serves MP3 files from cache directory
func (s *GSIServer) handleServeAudioFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	
	// Remove query params if present (used for cache busting)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Get cache path
	cachePath, err := config.GetVoiceCachePath()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get cache path")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build full path
	audioPath := filepath.Join(cachePath, filename)

	// Check if file exists
	if _, err := os.Stat(audioPath); err != nil {
		s.logger.WithField("file", filename).Debug("Audio file not found")
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	// Add timestamp to log
	s.logger.WithFields(logrus.Fields{
		"file": filename,
		"path": audioPath,
		"size": getFileSize(audioPath),
	}).Info("📁 Serving audio file")
	
	// Serve the file with NO CACHE to always get fresh audio
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, audioPath)
}

// handleServeEmbeddedAudio serves embedded audio files (Free mode)
func (s *GSIServer) handleServeEmbeddedAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Get embedded audio file
	audioData, err := assets.GetAudioFile(filename)
	if err != nil {
		s.logger.WithField("file", filename).Debug("Embedded audio file not found")
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"file": filename,
		"size": len(audioData),
		"mode": "free",
	}).Debug("📁 Serving embedded audio file")

	// Serve the embedded file
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache embedded files for 1 year
	w.Write(audioData)
}

// handleGetAudioBase64 returns audio file as base64 JSON (for Wails proxy compatibility)
func (s *GSIServer) handleGetAudioBase64(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Get cache path
	cachePath, err := config.GetVoiceCachePath()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get cache path")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build full path
	audioPath := filepath.Join(cachePath, filename)

	// Read file
	audioData, err := os.ReadFile(audioPath)
	if err != nil {
		s.logger.WithField("file", filename).WithError(err).Debug("Audio file not found")
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	// Encode to base64
	base64Str := base64.StdEncoding.EncodeToString(audioData)

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"data": base64Str,
	})
}

// handleAudioEvents streams audio events to frontend via Server-Sent Events
func (s *GSIServer) handleAudioEvents(w http.ResponseWriter, r *http.Request) {
	// Get voice handler
	vh, ok := s.voiceHandler.(*handlers.VoiceHandler)
	if !ok {
		http.Error(w, "Voice handler not available", http.StatusServiceUnavailable)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get event channel
	eventChan := vh.GetAudioEventChannel()
	
	// Create a done channel for cleanup
	done := r.Context().Done()

	s.logger.Info("🎵 Frontend connected to audio event stream")

	for {
		select {
		case event := <-eventChan:
			// Send event to frontend
			eventJSON, err := json.Marshal(event)
			if err != nil {
				s.logger.WithError(err).Error("Failed to marshal audio event")
				continue
			}

			// SSE format: data: {json}\n\n
			fmt.Fprintf(w, "data: %s\n\n", eventJSON)
			
			// Flush immediately
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

		case <-done:
			s.logger.Info("🎵 Frontend disconnected from audio event stream")
			return
		}
	}
}

// getFileSize returns the file size in bytes
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// getSemanticFilename returns the semantic filename for an event type
func getSemanticFilename(eventType string) string {
	// Add _warning suffix if not present (for compatibility)
	if !strings.HasSuffix(eventType, "_warning") &&
		(strings.Contains(eventType, "rune") ||
			eventType == "catapult" ||
			eventType == "catapult_timing" ||
			eventType == "stack_timing" ||
			eventType == "day_night" ||
			eventType == "day_night_cycle") {
		eventType = eventType + "_warning"
	}

	// Sanitize and return
	filename := strings.ToLower(eventType)
	filename = strings.ReplaceAll(filename, " ", "_")
	return filename + ".mp3"
}
