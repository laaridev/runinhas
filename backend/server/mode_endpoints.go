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
