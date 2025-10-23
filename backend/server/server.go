package server

import (
	"context"
	"dota-gsi/backend/config"
	"dota-gsi/backend/consumers"
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"dota-gsi/backend/metrics"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GSIServer handles GSI POST requests and publishes events
type GSIServer struct {
	eventBus        *events.EventBus
	eventEmitter    func(eventName string, data interface{}) // Wails event emitter
	logger          *logrus.Entry
	port            int
	server          *http.Server
	voiceHandler    interface{} // Will be set if voice is enabled
	consumerManager *consumers.ConsumerManager
	startTime       time.Time
}

// NewGSIServer creates a new GSI server with event streaming
func NewGSIServer(port int, logger *logrus.Entry, eventBus *events.EventBus) *GSIServer {
	return &GSIServer{
		eventBus: eventBus,
		logger:   logger,
		port:     port,
	}
}

// New creates a new GSI server with default configuration
func New() (*GSIServer, error) {
	// Set up logging
	logrus.SetLevel(logrus.DebugLevel) // Enable debug logging
	logEntry := logrus.WithField("component", "server")

	// Create logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create event bus
	eventBus := events.NewEventBus()

	// Create server
	server := NewGSIServer(3001, logEntry, eventBus)
	server.startTime = time.Now()

	// Load configuration to get voice settings
	cfg, err := config.Load()
	if err == nil && cfg.ElevenLabsAPIKey != "" {
		// Create VoiceHandler with config from file
		voiceHandler, err := handlers.NewVoiceHandler(
			cfg.ElevenLabsAPIKey,
			cfg.ElevenLabsVoiceID,
			cfg.VoiceCachePath,
			logEntry.WithField("component", "voice"),
		)
		if err == nil {
			// Set game config for voice handler
			voiceHandler.SetGameConfig(cfg.Game)
			// Set voice handler in the server
			server.SetVoiceHandler(voiceHandler)
			logEntry.Info("Voice handler initialized with API key from config")

			// Create and start consumers with voice handler
			server.consumerManager = consumers.NewConsumerManager(logEntry.WithField("component", "consumers"))
			handlerList := []handlers.Handler{voiceHandler}

			// Add rune and timing consumers
			server.consumerManager.AddRuneConsumer(eventBus, handlerList, cfg.Game)
			server.consumerManager.AddTimingConsumer(eventBus, handlerList, cfg.Game)

			// Start all consumers
			server.consumerManager.StartAll()
		} else {
			logEntry.WithError(err).Warn("Failed to create voice handler")
		}
	} else {
		logEntry.Info("Voice handler disabled (no API key configured)")
	}

	return server, nil
}

// SetVoiceHandler sets the voice handler for audio endpoints
func (s *GSIServer) SetVoiceHandler(handler interface{}) {
	s.voiceHandler = handler
	
	// If we already have an event emitter, set up direct emitter ONLY (no channel listener)
	if s.eventEmitter != nil {
		if vh, ok := handler.(*handlers.VoiceHandler); ok {
			s.logger.Info("🎵 Voice handler set, configuring direct emitter")
			// Set the direct emitter on the voice handler
			vh.SetDirectEmitter(s.eventEmitter)
			// NOTE: Removed channel listener to avoid duplicate events
		}
	}
}

// SetEventEmitter sets the event emitter callback for Wails
func (s *GSIServer) SetEventEmitter(emitter func(eventName string, data interface{})) {
	s.eventEmitter = emitter
	s.logger.Info("📡 Event emitter configured for GSI server")
	
	// Also set it on the voice handler if available
	if s.voiceHandler != nil {
		if vh, ok := s.voiceHandler.(*handlers.VoiceHandler); ok {
			s.logger.Info("🔗 Linking event emitter to voice handler")
			vh.SetDirectEmitter(emitter)
			// NOTE: Using direct emitter only to avoid duplicate events
		} else {
			s.logger.Warn("⚠️ Voice handler exists but is not the correct type")
		}
	} else {
		s.logger.Warn("⚠️ No voice handler available when setting event emitter")
	}
}

// Start starts the HTTP server on the specified address
func (s *GSIServer) Start(addr string) error {
	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/gsi", s.handleGSI).Methods("POST")
	router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Add config endpoints
	s.AddConfigEndpoints(router)

	// Add ElevenLabs endpoints
	s.AddElevenLabsEndpoints(router)

	// Add audio endpoints
	s.AddAudioEndpoints(router)
	router.Use(s.corsMiddleware)

	// Create HTTP server
	s.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	s.logger.WithField("addr", addr).Info("🚀 GSI Server starting - Event Publisher Only")

	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server
func (s *GSIServer) Stop() error {
	s.logger.Info("🛑 Shutting down GSI Server...")

	// Stop consumers first
	if s.consumerManager != nil {
		s.consumerManager.StopAll()
	}

	if s.server == nil {
		return nil
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// handleGSI processes GSI POST requests and publishes TickEvents
func (s *GSIServer) handleGSI(w http.ResponseWriter, r *http.Request) {
	// Read raw JSON body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.WithError(err).Error("Failed to read GSI body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log that we received GSI data
	s.logger.WithField("size", len(body)).Info("📡 GSI DATA RECEIVED!")

	// Create TickEvent with raw JSON
	tickEvent := events.TickEvent{
		RawJSON: body,
		Time:    time.Now(),
	}

	// Publish to event bus - all consumers will receive it
	s.eventBus.Publish(tickEvent)

	s.logger.Info("✅ Event published to bus")

	w.WriteHeader(http.StatusOK)
}

// handleHealth returns server health status with metrics
func (s *GSIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Build health response with metrics
	stats := map[string]interface{}{
		"status":      "healthy",
		"architecture": "event-streaming",
		"uptime":      time.Since(s.startTime).Seconds(),
		"metrics":     metrics.Instance.GetStats(),
	}
	
	// Add consumer count if available
	if s.consumerManager != nil {
		stats["consumers"] = s.consumerManager.Count()
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// corsMiddleware adds CORS headers
func (s *GSIServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permite qualquer origem (incluindo wails://)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Audio API Handlers moved to audio_endpoints.go
