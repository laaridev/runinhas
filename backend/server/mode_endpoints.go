package server

import (
	"dota-gsi/backend/config"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// AddModeEndpoints adds mode-related endpoints to the router
func (s *GSIServer) AddModeEndpoints(router *mux.Router) {
	router.HandleFunc("/api/mode", s.handleGetMode).Methods("GET")
	router.HandleFunc("/api/mode", s.handleSetMode).Methods("POST")
}

// handleGetMode returns the current app mode (free or pro)
func (s *GSIServer) handleGetMode(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	mode := cfg.Mode
	if mode == "" {
		mode = "free" // Default
	}

	response := map[string]interface{}{
		"mode": mode,
		"features": map[string]bool{
			"customMessages":    mode == "pro",
			"customVoice":       mode == "pro",
			"elevenLabs":        mode == "pro",
			"audioGeneration":   mode == "pro",
			"embeddedAudioOnly": mode == "free",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSetMode sets the app mode (free or pro) with optional license key
func (s *GSIServer) handleSetMode(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Mode       string `json:"mode"`
		LicenseKey string `json:"license_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate mode
	if request.Mode != "free" && request.Mode != "pro" {
		http.Error(w, "Invalid mode: must be 'free' or 'pro'", http.StatusBadRequest)
		return
	}

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	// Update mode and license key
	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}
	cfg.Game.Mode = request.Mode
	cfg.Game.LicenseKey = request.LicenseKey
	cfg.Mode = request.Mode
	cfg.LicenseKey = request.LicenseKey

	// Save to config.json
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		s.logger.WithError(err).Error("Failed to save mode to config")
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	s.logger.WithFields(map[string]interface{}{
		"mode":        request.Mode,
		"license_key": request.LicenseKey != "",
	}).Info("App mode updated successfully")

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"mode":    request.Mode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
