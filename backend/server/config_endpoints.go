package server

import (
	"dota-gsi/backend/config"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AddConfigEndpoints adds all configuration-related HTTP endpoints
// This includes: config, timings, messages, and system status
func (s *GSIServer) AddConfigEndpoints(router *mux.Router) {
	// Configuration endpoints
	router.HandleFunc("/api/config", s.handleGetConfig).Methods("GET")
	router.HandleFunc("/api/config", s.handleSaveConfig).Methods("POST")
	
	// Timing endpoints
	router.HandleFunc("/api/timing/{key}/enabled", s.handleGetTimingEnabled).Methods("GET")
	router.HandleFunc("/api/timing/{key}/{field}", s.handleGetTimingValue).Methods("GET")
	router.HandleFunc("/api/timing/{key}/{field}", s.handleSetTimingValue).Methods("POST")
	
	// Message endpoints
	router.HandleFunc("/api/message/{key}", s.handleGetMessage).Methods("GET")
	router.HandleFunc("/api/message/{key}", s.handleSetMessage).Methods("POST")
	
	// System status endpoints
	router.HandleFunc("/api/system/status", s.handleSystemStatus).Methods("GET")
	router.HandleFunc("/api/system/first-run", s.handleSetFirstRun).Methods("POST")
	router.HandleFunc("/api/system/gsi-installed", s.handleSetGSIInstalled).Methods("POST")
}

// ============================================================================
// Configuration Endpoints
// ============================================================================
// handleGetConfig returns the current configuration
func (s *GSIServer) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		// Se falhar ao carregar, cria uma nova configuração
		gameConfig, _ := config.LoadOrCreateConfig()
		cfg = &config.Config{
			Game: gameConfig,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg.Game)
}

// handleSaveConfig saves the entire configuration
func (s *GSIServer) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update voice configuration if present
	if voiceData, ok := updates["voice"].(map[string]interface{}); ok {
		if cfg.Game.Voice == nil {
			cfg.Game.Voice = make(map[string]interface{})
		}
		
		// Update each voice field
		for key, value := range voiceData {
			cfg.Game.Voice[key] = value
		}
		
		s.logger.WithFields(logrus.Fields{
			"voice": voiceData,
		}).Info("Updating voice configuration")
	}
	
	// Update timings if present
	if timingsData, ok := updates["timings"].(map[string]interface{}); ok {
		if cfg.Game.Timings == nil {
			cfg.Game.Timings = make(map[string]map[string]interface{})
		}
		
		for key, value := range timingsData {
			if timingMap, ok := value.(map[string]interface{}); ok {
				cfg.Game.Timings[key] = timingMap
			}
		}
	}
	
	// Update messages if present
	if messagesData, ok := updates["messages"].(map[string]interface{}); ok {
		if cfg.Game.Messages == nil {
			cfg.Game.Messages = make(map[string]string)
		}
		
		for key, value := range messagesData {
			if msg, ok := value.(string); ok {
				cfg.Game.Messages[key] = msg
			}
		}
	}
	
	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	s.logger.Info("Configuration saved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

// ============================================================================
// Timing Configuration Endpoints
// ============================================================================

// handleGetTimingEnabled returns if a timing is enabled
func (s *GSIServer) handleGetTimingEnabled(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	enabled := false
	if timing, exists := cfg.Game.Timings[key]; exists {
		if val, ok := timing["enabled"].(bool); ok {
			enabled = val
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"enabled": enabled})
}

// handleGetTimingValue returns a timing configuration value
func (s *GSIServer) handleGetTimingValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	field := vars["field"]
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	value := 0
	if timing, exists := cfg.Game.Timings[key]; exists {
		if val, ok := timing[field].(float64); ok {
			value = int(val)
		} else if val, ok := timing[field].(int); ok {
			value = val
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"value": value})
}

// handleSetTimingValue sets a timing configuration value
func (s *GSIServer) handleSetTimingValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	field := vars["field"]
	
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Initialize if needed
	if cfg.Game.Timings == nil {
		cfg.Game.Timings = make(map[string]map[string]interface{})
	}
	if cfg.Game.Timings[key] == nil {
		cfg.Game.Timings[key] = make(map[string]interface{})
	}
	
	// Handle enabled field (boolean)
	if field == "enabled" {
		if val, ok := body["enabled"].(bool); ok {
			cfg.Game.Timings[key][field] = val
			s.logger.WithFields(logrus.Fields{
				"key":   key,
				"field": field,
				"value": val,
			}).Info("Setting enabled status for " + key)
		}
	} else {
		// Handle numeric fields
		if val, ok := body["value"].(float64); ok {
			cfg.Game.Timings[key][field] = int(val)
		} else if val, ok := body["value"].(int); ok {
			cfg.Game.Timings[key][field] = val
		}
	}
	
	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	s.logger.WithFields(logrus.Fields{
		"key":   key,
		"field": field,
	}).Info("Successfully updated " + key + "." + field)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// ============================================================================
// Message Configuration Endpoints
// ============================================================================

// handleGetMessage returns a custom message
func (s *GSIServer) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	message := ""
	if cfg.Game.Messages != nil {
		if msg, ok := cfg.Game.Messages[key]; ok {
			message = msg
		}
	}
	
	s.logger.WithFields(logrus.Fields{
		"key": key,
		"message": message,
	}).Debug("Getting message for key")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// handleSetMessage sets a custom message
func (s *GSIServer) handleSetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if cfg.Game.Messages == nil {
		cfg.Game.Messages = make(map[string]string)
	}
	
	if msg, ok := body["message"]; ok {
		cfg.Game.Messages[key] = msg
	}
	
	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// ============================================================================
// System Configuration Endpoints
// ============================================================================

// handleSystemStatus returns system status (first run, GSI installed)
func (s *GSIServer) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Default values
	firstRun := true
	gsiInstalled := false
	
	// Get actual values if available
	if cfg.Game != nil && cfg.Game.System != nil {
		firstRun = cfg.Game.System.FirstRun
		gsiInstalled = cfg.Game.System.GSIInstalled
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"first_run":     firstRun,
		"gsi_installed": gsiInstalled,
	})
}

// handleSetFirstRun updates the first run status
func (s *GSIServer) handleSetFirstRun(w http.ResponseWriter, r *http.Request) {
	var body map[string]bool
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	value, ok := body["value"]
	if !ok {
		http.Error(w, "Missing 'value' field", http.StatusBadRequest)
		return
	}
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Initialize if needed
	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}
	if cfg.Game.System == nil {
		cfg.Game.System = &config.SystemConfig{}
	}
	
	cfg.Game.System.FirstRun = value
	
	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// handleSetGSIInstalled updates the GSI installed status
func (s *GSIServer) handleSetGSIInstalled(w http.ResponseWriter, r *http.Request) {
	var body map[string]bool
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	value, ok := body["value"]
	if !ok {
		http.Error(w, "Missing 'value' field", http.StatusBadRequest)
		return
	}
	
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Initialize if needed
	if cfg.Game == nil {
		cfg.Game = &config.GameConfig{}
	}
	if cfg.Game.System == nil {
		cfg.Game.System = &config.SystemConfig{}
	}
	
	cfg.Game.System.GSIInstalled = value
	
	// Save configuration
	configPath, _ := config.GetConfigPath()
	if err := config.SaveGameConfig(configPath, cfg.Game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
