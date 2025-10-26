package server

import (
	"dota-gsi/backend/config"
	"dota-gsi/backend/handlers"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// ============================================================================
// ElevenLabs Voice Configuration Endpoints
// ============================================================================
// This file handles ElevenLabs API configuration and voice testing
// Uses VoiceHandler for actual voice generation to avoid duplication

// ElevenLabsConfig represents the ElevenLabs configuration
type ElevenLabsConfig struct {
	APIKey       string  `json:"apiKey"`
	VoiceID      string  `json:"voiceId"`
	Stability    float64 `json:"stability"`
	Similarity   float64 `json:"similarity"`
	Style        float64 `json:"style"`
	SpeakerBoost bool    `json:"speakerBoost"`
}

// VoiceSettings for ElevenLabs API requests
type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

// AddElevenLabsEndpoints registers all ElevenLabs-related HTTP endpoints
func (s *GSIServer) AddElevenLabsEndpoints(router *mux.Router) {
	router.HandleFunc("/api/elevenlabs/config", s.handleGetElevenLabsConfig).Methods("GET")
	router.HandleFunc("/api/elevenlabs/config", s.handleSaveElevenLabsConfig).Methods("POST")
	router.HandleFunc("/api/elevenlabs/test", s.handleTestVoice).Methods("POST")
	router.HandleFunc("/api/elevenlabs/voices", s.handleGetVoices).Methods("GET")
}

// handleGetElevenLabsConfig returns the current ElevenLabs configuration
func (s *GSIServer) handleGetElevenLabsConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		// Return default config if error
		defaultConfig := ElevenLabsConfig{
			APIKey:       "",
			VoiceID:      "EXAVITQu4vr4xnSDxMaL", // Default Sarah voice
			Stability:    0.5,
			Similarity:   0.75,
			Style:        0,
			SpeakerBoost: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultConfig)
		return
	}

	// Extract ElevenLabs config from game config
	elevenLabsConfig := ElevenLabsConfig{
		APIKey:       cfg.ElevenLabsAPIKey,
		VoiceID:      "EXAVITQu4vr4xnSDxMaL", // Default for now
		Stability:    0.5,
		Similarity:   0.75,
		Style:        0,
		SpeakerBoost: true,
	}

	// Check if we have voice settings in the config
	if cfg.Game != nil && cfg.Game.Voice != nil {
		if voiceID, ok := cfg.Game.Voice["voice_id"].(string); ok {
			elevenLabsConfig.VoiceID = voiceID
		}
		if stability, ok := cfg.Game.Voice["stability"].(float64); ok {
			elevenLabsConfig.Stability = stability
		}
		if similarity, ok := cfg.Game.Voice["similarity"].(float64); ok {
			elevenLabsConfig.Similarity = similarity
		}
		if style, ok := cfg.Game.Voice["style"].(float64); ok {
			elevenLabsConfig.Style = style
		}
		if speakerBoost, ok := cfg.Game.Voice["speaker_boost"].(bool); ok {
			elevenLabsConfig.SpeakerBoost = speakerBoost
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(elevenLabsConfig)
}

// handleSaveElevenLabsConfig saves the ElevenLabs configuration
func (s *GSIServer) handleSaveElevenLabsConfig(w http.ResponseWriter, r *http.Request) {
	var elevenLabsConfig ElevenLabsConfig
	if err := json.NewDecoder(r.Body).Decode(&elevenLabsConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update ElevenLabs API key
	cfg.ElevenLabsAPIKey = elevenLabsConfig.APIKey
	cfg.ElevenLabsVoiceID = elevenLabsConfig.VoiceID

	// Initialize Voice config if needed
	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}
	if cfg.Game.Voice == nil {
		cfg.Game.Voice = make(map[string]interface{})
	}

	// Update voice settings
	cfg.Game.Voice["voice_id"] = elevenLabsConfig.VoiceID
	cfg.Game.Voice["stability"] = elevenLabsConfig.Stability
	cfg.Game.Voice["similarity"] = elevenLabsConfig.Similarity
	cfg.Game.Voice["style"] = elevenLabsConfig.Style
	cfg.Game.Voice["speaker_boost"] = elevenLabsConfig.SpeakerBoost

	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update VoiceHandler if it exists
	if s.voiceHandler != nil {
		if vh, ok := s.voiceHandler.(*handlers.VoiceHandler); ok {
			vh.UpdateSettings(
				elevenLabsConfig.APIKey,
				elevenLabsConfig.VoiceID,
				elevenLabsConfig.Stability,
				elevenLabsConfig.Similarity,
				elevenLabsConfig.Style,
				elevenLabsConfig.SpeakerBoost,
			)
		}
	}
	
	s.logger.WithField("voice_id", elevenLabsConfig.VoiceID).Info("ElevenLabs configuration saved")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

// handleTestVoice tests the voice synthesis with current settings
func (s *GSIServer) handleTestVoice(w http.ResponseWriter, r *http.Request) {
	var testRequest struct {
		Text     string        `json:"text"`
		VoiceID  string        `json:"voiceId"`
		Settings VoiceSettings `json:"settings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&testRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load current config to get API key
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	apiKey := cfg.ElevenLabsAPIKey
	if apiKey == "" && cfg.Game != nil && cfg.Game.Voice != nil {
		if key, ok := cfg.Game.Voice["apiKey"].(string); ok {
			apiKey = key
		}
	}
	
	if apiKey == "" {
		http.Error(w, "API Key not configured", http.StatusBadRequest)
		return
	}

	// Use VoiceHandler if available
	if s.voiceHandler != nil {
		if vh, ok := s.voiceHandler.(*handlers.VoiceHandler); ok {
			// Update settings for this test (including fresh API key)
			vh.UpdateSettings(
				apiKey,
				testRequest.VoiceID,
				testRequest.Settings.Stability,
				testRequest.Settings.SimilarityBoost,
				testRequest.Settings.Style,
				testRequest.Settings.UseSpeakerBoost,
			)
			
			// Generate voice
			audioData, err := vh.GenerateVoice(testRequest.Text)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			// Save to cache for future use
			cachePath, _ := config.GetVoiceCachePath()
			tempFile := filepath.Join(cachePath, "elevenlabs_test.mp3")
			
			// Ensure cache directory exists
			os.MkdirAll(cachePath, 0755)
			
			if err := os.WriteFile(tempFile, audioData, 0644); err != nil {
				s.logger.WithError(err).Error("Failed to save audio file")
				http.Error(w, "Failed to save audio", http.StatusInternalServerError)
				return
			}
			
			// Return filename for frontend to play
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":   "generated",
				"filename": "elevenlabs_test.mp3",
			})
			return
		}
	}

	// VoiceHandler not available - should not happen in normal operation
	s.logger.Error("VoiceHandler not available for voice test")
	http.Error(w, "Voice handler not initialized", http.StatusServiceUnavailable)
}

// handleGetVoices fetches available voices from ElevenLabs
func (s *GSIServer) handleGetVoices(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("xi-api-key")
	if apiKey == "" {
		// Try to get from config
		cfg, err := config.Load()
		if err != nil || cfg.ElevenLabsAPIKey == "" {
			http.Error(w, "API key not provided", http.StatusBadRequest)
			return
		}
		apiKey = cfg.ElevenLabsAPIKey
	}

	url := "https://api.elevenlabs.io/v1/voices"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("xi-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch voices", resp.StatusCode)
		return
	}

	// Parse and forward the response
	var voicesResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&voicesResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voicesResponse)
}
