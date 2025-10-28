package server

import (
	"encoding/json"
	"net/http"

	"dota-gsi/backend/config"
	"github.com/gorilla/mux"
)

// AddVirtualMicEndpoints adds virtual mic configuration endpoints
func (s *GSIServer) AddVirtualMicEndpoints(router *mux.Router) {
	router.HandleFunc("/api/audio/virtualmicEnabled", s.handleGetVirtualMicEnabled).Methods("GET")
	router.HandleFunc("/api/audio/virtualmicEnabled", s.handleSetVirtualMicEnabled).Methods("POST")
	router.HandleFunc("/api/audio/virtualmicDevice", s.handleDetectVirtualMic).Methods("GET")
	router.HandleFunc("/api/audio/virtualmicDevice", s.handleSetVirtualMicDevice).Methods("POST")
}

// handleGetVirtualMicEnabled returns current virtual mic status
func (s *GSIServer) handleGetVirtualMicEnabled(w http.ResponseWriter, r *http.Request) {
	enabled := s.audioPlayer.GetVirtualMicEnabled()
	device := s.audioPlayer.GetVirtualMicDevice()

	response := map[string]interface{}{
		"enabled": enabled,
		"device":  device,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSetVirtualMicEnabled enables/disables virtual mic output
func (s *GSIServer) handleSetVirtualMicEnabled(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Enabled bool `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update audio player
	s.audioPlayer.SetVirtualMicEnabled(request.Enabled)

	// Save to config
	cfg, err := config.Load()
	if err != nil {
		s.logger.WithError(err).Error("Failed to load config")
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}

	cfg.Game.Audio.VirtualMicEnabled = request.Enabled

	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		s.logger.WithError(err).Error("Failed to save virtual mic config")
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	s.logger.WithField("enabled", request.Enabled).Info("Virtual mic status updated")

	response := map[string]interface{}{
		"success": true,
		"enabled": request.Enabled,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDetectVirtualMic detects available virtual microphone
func (s *GSIServer) handleDetectVirtualMic(w http.ResponseWriter, r *http.Request) {
	device, found := s.audioPlayer.DetectVirtualMic()

	response := map[string]interface{}{
		"found":  found,
		"device": device,
	}

	if !found {
		response["message"] = "No virtual microphone detected on this system"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSetVirtualMicDevice sets the virtual mic device name
func (s *GSIServer) handleSetVirtualMicDevice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Device string `json:"device"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Device == "" {
		http.Error(w, "Device name cannot be empty", http.StatusBadRequest)
		return
	}

	// Update audio player
	s.audioPlayer.SetVirtualMicDevice(request.Device)

	// Save to config
	cfg, err := config.Load()
	if err != nil {
		s.logger.WithError(err).Error("Failed to load config")
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}

	cfg.Game.Audio.VirtualMicDevice = request.Device

	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		s.logger.WithError(err).Error("Failed to save virtual mic device")
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	s.logger.WithField("device", request.Device).Info("Virtual mic device updated")

	response := map[string]interface{}{
		"success": true,
		"device":  request.Device,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
