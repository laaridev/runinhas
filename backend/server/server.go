package server

import (
	"context"
	"dota-gsi/backend/config"
	"dota-gsi/backend/consumers"
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GSIServer handles HTTP requests and publishes TickEvents
type GSIServer struct {
	eventBus        *events.EventBus
	logger          *logrus.Entry
	port            int
	server          *http.Server
	voiceHandler    interface{} // Will be set if voice is enabled
	consumerManager *consumers.ConsumerManager
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
	// Create logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logEntry := logger.WithField("component", "server")

	// Create event bus
	eventBus := events.NewEventBus()

	// Create server
	server := NewGSIServer(3001, logEntry, eventBus)

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
}

// Start starts the HTTP server on the specified address
func (s *GSIServer) Start(addr string) error {
	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/gsi", s.handleGSI).Methods("POST")
	router.HandleFunc("/health", s.handleHealth).Methods("GET")
	// Add configuration endpoints
	s.AddConfigEndpoints(router)
	// Add ElevenLabs endpoints
	s.AddElevenLabsEndpoints(router)

	// Add audio endpoints
	s.AddAudioEndpoints(router)

	// Add CORS middleware
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

// Shutdown is an alias for Stop to maintain compatibility
func (s *GSIServer) Shutdown() error {
	return s.Stop()
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

// handleHealth returns server health status
func (s *GSIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "architecture": "event-streaming"}`))
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
