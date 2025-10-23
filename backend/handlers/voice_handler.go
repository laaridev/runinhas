package handlers

import (
	"bytes"
	"dota-gsi/backend/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// AudioEvent represents an audio event to be played by the frontend
type AudioEvent struct {
	Filename  string                 `json:"filename"`
	EventType string                 `json:"eventType"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// VoiceHandler handles voice announcements using ElevenLabs
type VoiceHandler struct {
	apiKey      string
	voiceID     string
	cachePath   string
	logger      *logrus.Entry
	enabled     bool
	initialized bool
	gameConfig  interface{} // Game configuration for messages
	initMutex   sync.Mutex  // Mutex to prevent double initialization
	// Audio event notification channel (for frontend playback)
	audioEventChan chan AudioEvent
	queueMutex     sync.Mutex // Mutex for queue operations
	// Direct emitter for Wails events
	directEmitter func(eventName string, data interface{})
	// Voice settings (can be updated dynamically)
	stability    float64
	similarity   float64
	style        float64
	speakerBoost bool
}

// UpdateSettings updates the voice handler settings dynamically
func (vh *VoiceHandler) UpdateSettings(apiKey, voiceID string, stability, similarity, style float64, speakerBoost bool) {
	vh.apiKey = apiKey
	vh.voiceID = voiceID
	vh.stability = stability
	vh.similarity = similarity
	vh.style = style
	vh.speakerBoost = speakerBoost
	vh.enabled = apiKey != ""

	vh.logger.WithFields(logrus.Fields{
		"voice_id": voiceID,
		"enabled":  vh.enabled,
	}).Info("Voice settings updated")
}

// GetSettings returns current voice settings
func (vh *VoiceHandler) GetSettings() map[string]interface{} {
	return map[string]interface{}{
		"apiKey":       vh.apiKey,
		"voiceId":      vh.voiceID,
		"stability":    vh.stability,
		"similarity":   vh.similarity,
		"style":        vh.style,
		"speakerBoost": vh.speakerBoost,
		"enabled":      vh.enabled,
	}
}

// NewVoiceHandler creates a new voice handler
func NewVoiceHandler(apiKey, voiceID, cachePath string, logger *logrus.Entry) (*VoiceHandler, error) {
	// Validate configuration
	if apiKey != "" && voiceID == "" {
		return nil, fmt.Errorf("ELEVENLABS_VOICE_ID is required when ELEVENLABS_API_KEY is provided")
	}

	// Use provided cache path
	finalCachePath := cachePath

	// Create cache directory
	if err := os.MkdirAll(finalCachePath, 0755); err != nil {
	}

	enabled := apiKey != ""

	vh := &VoiceHandler{
		cachePath:      cachePath,
		apiKey:         apiKey,
		voiceID:        voiceID,
		logger:         logger.WithField("component", "voice"),
		gameConfig:     nil, // Will be set later by SetGameConfig
		enabled:        enabled,
		audioEventChan: make(chan AudioEvent, 10), // Buffer up to 10 audio events
			// Default voice settings (will be overridden by config)
		stability:    config.DefaultStability,
		similarity:   config.DefaultSimilarity,
		style:        config.DefaultStyle,
		speakerBoost: config.DefaultSpeakerBoost,
	}

	// Log configuration status
	if enabled {
		logger.Info("🎵 Voice handler initialized with ElevenLabs")
		// Start cache cleanup goroutine
		go vh.startCacheCleanup()
	}

	return vh, nil
}

// SetGameConfig sets the game configuration for the voice handler
func (vh *VoiceHandler) SetGameConfig(gameConfig interface{}) {
	vh.gameConfig = gameConfig
}

// SetDirectEmitter sets the direct emitter callback for Wails events
func (vh *VoiceHandler) SetDirectEmitter(emitter func(eventName string, data interface{})) {
	vh.directEmitter = emitter
	vh.logger.Info("✅ Direct emitter configured for voice handler")
}

// Handle processes voice events
func (vh *VoiceHandler) Handle(eventType string, data interface{}) {
	if !vh.enabled {
		return
	}

	// Initialize speaker if needed
	if !vh.initialized {
		vh.initSpeaker()
	}

	message := vh.getMessageForEvent(eventType, data)
	if message != "" {
		vh.speakWithData(message, eventType, data)
	}
}

// initSpeaker is deprecated - audio playback now handled by frontend
func (vh *VoiceHandler) initSpeaker() {
	vh.initMutex.Lock()
	defer vh.initMutex.Unlock()

	if vh.initialized {
		return
	}

	vh.initialized = true
	vh.logger.Info("✅ Voice handler initialized (frontend playback mode)")
}

// speak generates and plays voice audio with semantic caching
func (vh *VoiceHandler) speakWithData(text string, eventType string, data interface{}) {
	go func() {
		// Get semantic cache path
		cacheFile := vh.getCacheFilePathSemantic(eventType, data)

		// Check if cache file exists
		if _, err := os.Stat(cacheFile); err == nil {
			vh.logger.WithField("cache_file", filepath.Base(cacheFile)).Debug("Using cached audio")
			vh.emitAudioEvent(filepath.Base(cacheFile), eventType, data)
			return
		}

		// Generate with ElevenLabs
		audioData, err := vh.GenerateVoice(text)
		if err != nil {
			vh.logger.WithError(err).Error("Failed to generate voice")
			return
		}

		// Save to cache
		if err := os.WriteFile(cacheFile, audioData, 0644); err != nil {
			vh.logger.WithError(err).Error("Failed to cache audio")
			return
		}

		vh.logger.WithField("cache_file", filepath.Base(cacheFile)).Debug("Audio generated and cached")
		vh.emitAudioEvent(filepath.Base(cacheFile), eventType, data)
	}()
}

// speak generates and plays voice audio (legacy method for compatibility)
func (vh *VoiceHandler) speak(text string) {
	vh.speakWithData(text, "", nil)
}

// GenerateVoice calls ElevenLabs API (exported for API usage)
func (vh *VoiceHandler) GenerateVoice(text string) ([]byte, error) {
	// Check if API key and voice ID are configured
	if vh.apiKey == "" {
		return nil, fmt.Errorf("ElevenLabs API key not configured")
	}
	if vh.voiceID == "" {
		return nil, fmt.Errorf("ElevenLabs voice ID not configured")
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", vh.voiceID)

	// Log the request for debugging
	vh.logger.WithFields(logrus.Fields{
		"voice_id": vh.voiceID,
		"text_len": len(text),
		"url":      url,
	}).Debug("Generating voice with ElevenLabs")

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":         vh.stability,
			"similarity_boost":  vh.similarity,
			"style":             vh.style,
			"use_speaker_boost": vh.speakerBoost,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", vh.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		vh.logger.WithFields(logrus.Fields{
			"status":   resp.StatusCode,
			"body":     string(body),
			"voice_id": vh.voiceID,
		}).Error("ElevenLabs API error")
		return nil, fmt.Errorf("ElevenLabs API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// emitAudioEvent emits an audio event for frontend to play
func (vh *VoiceHandler) emitAudioEvent(filename, eventType string, data interface{}) {
	// Convert data to map
	dataMap := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		dataMap = m
	}

	event := AudioEvent{
		Filename:  filename,
		EventType: eventType,
		Data:      dataMap,
		Timestamp: time.Now().Unix(),
	}

	// Use direct emitter (Wails events)
	if vh.directEmitter != nil {
		vh.directEmitter("audio:play", event)
		vh.logger.WithField("filename", filename).Debug("✅ Audio event emitted via Wails")
	} else {
		vh.logger.Warn("⚠️ No direct emitter configured, audio event not sent")
	}
}

// GetAudioEventChannel returns the audio event channel (deprecated, kept for compatibility)
func (vh *VoiceHandler) GetAudioEventChannel() <-chan AudioEvent {
	return vh.audioEventChan
}

// PlayAudioFile emits event for frontend to play audio
func (vh *VoiceHandler) PlayAudioFile(filePath string) {
	vh.logger.WithField("file", filepath.Base(filePath)).Debug("🎵 PlayAudioFile called")
	vh.emitAudioEvent(filepath.Base(filePath), "manual_play", nil)
}

// getCacheFilePathSemantic generates semantic cache file path based on event type and parameters
func (vh *VoiceHandler) getCacheFilePathSemantic(eventType string, data interface{}) string {
	cachePath := vh.cachePath

	// Convert data to map for easier access
	dataMap := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		dataMap = m
	}

	var filename string

	// Generate semantic filenames based on event type and parameters
	switch eventType {
	// Rune warnings - use generic filename (reuse same audio)
	case "bounty_rune_warning", "power_rune_warning", "water_rune_warning", "wisdom_rune_warning":
		filename = fmt.Sprintf("%s.mp3", eventType)

	// Hero events - use generic filename (reuse same audio)
	case "hero_health_low", "hero_health_critical", "hero_mana_low":
		filename = fmt.Sprintf("%s.mp3", eventType)

	case "hero_level_up", "hero_ultimate_ready", "hero_death":
		filename = fmt.Sprintf("%s.mp3", eventType)

	// Timing events - use generic filenames (reuse same audio)
	case "catapult_timing":
		filename = "catapult_timing.mp3"

	// Day/night: 1 arquivo apenas para cycle (warning), transition não precisa de cache
	case "day_night_cycle":
		filename = "day_night_cycle.mp3"
		
	case "day_night_transition":
		// Transition não usa cache, fala diretamente "Amanheceu" ou "Anoiteceu"
		filename = "day_night_transition.mp3"

	case "day_night_change":
		if isDaytime, exists := dataMap["daytime"]; exists {
			if isDaytime.(bool) {
				filename = "day_night_change_day.mp3"
			} else {
				filename = "day_night_change_night.mp3"
			}
		} else {
			filename = fmt.Sprintf("%s.mp3", eventType)
		}

	case "game_state_change":
		if to, exists := dataMap["to"]; exists {
			state := strings.ToLower(strings.ReplaceAll(to.(string), "_", ""))
			filename = fmt.Sprintf("%s_%s.mp3", eventType, state)
		} else {
			filename = fmt.Sprintf("%s.mp3", eventType)
		}

	case "score_change":
		if radiantScore, rExists := dataMap["radiant_score"]; rExists {
			if direScore, dExists := dataMap["dire_score"]; dExists {
				if radiantDiff, exists := dataMap["radiant_diff"]; exists {
					if radiantDiff.(int64) > 0 {
						filename = fmt.Sprintf("score_radiant_%dv%d.mp3", radiantScore, direScore)
					} else {
						filename = fmt.Sprintf("score_dire_%dv%d.mp3", direScore, radiantScore)
					}
				} else {
					filename = fmt.Sprintf("score_%dv%d.mp3", radiantScore, direScore)
				}
			} else {
				filename = fmt.Sprintf("%s.mp3", eventType)
			}
		} else {
			filename = fmt.Sprintf("%s.mp3", eventType)
		}

	case "stack_timing":
		// Use generic filename (reuse same audio for all minutes)
		filename = fmt.Sprintf("%s.mp3", eventType)

	default:
		// Fallback to simple event type naming
		filename = fmt.Sprintf("%s.mp3", strings.ReplaceAll(eventType, " ", "_"))
	}

	// Sanitize filename
	filename = strings.ToLower(filename)
	filename = strings.ReplaceAll(filename, " ", "_")

	return filepath.Join(cachePath, filename)
}

// getMessageForEvent returns appropriate message for the event
func (vh *VoiceHandler) getMessageForEvent(eventType string, data interface{}) string {
	// Convert data to map if needed
	dataMap := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		dataMap = m
	}

	// ALWAYS reload config from disk to get latest messages
	// This ensures we get updated messages after user edits them
	type GameConfigInterface interface {
		GetMessage(string) string
	}

	if gc, ok := vh.gameConfig.(GameConfigInterface); ok {
		// Handle day_night_transition - automatic "Amanheceu" or "Anoiteceu"
		if eventType == "day_night_transition" {
			if cycleType, exists := dataMap["cycle_type"]; exists {
				cycleTypeStr := cycleType.(string)
				if cycleTypeStr == "day" {
					return "Amanheceu"
				} else {
					return "Anoiteceu"
				}
			}
		}

		// Try to get message from config
		if msg := gc.GetMessage(eventType); msg != "" {
			return vh.replaceParameters(msg, dataMap)
		}
	}

	// Fallback to hardcoded messages
	// Special handling for game state changes
	if eventType == "game_state_change" {
		if to, exists := dataMap["to"]; exists {
			return vh.getGameStateMessage(to.(string))
		}
	}

	// Special handling for day/night changes
	if eventType == "day_night_change" {
		if isDaytime, exists := dataMap["daytime"]; exists {
			if isDaytime.(bool) {
				return "Amanheceu"
			}
			return "Anoiteceu"
		}
	}

	// Special handling for score changes
	if eventType == "score_change" {
		if radiantDiff, exists := dataMap["radiant_diff"]; exists {
			if radiantDiff.(int64) > 0 {
				return "Radiant marcou!"
			}
			return "Dire marcou!"
		}
	}

	// Special handling for hero death
	if eventType == "hero_death" {
		if deaths, exists := dataMap["deaths"]; exists {
			return fmt.Sprintf("Você morreu %d vez", deaths)
		}
		return "Você morreu!"
	}

	// Return static message
	return vh.getStaticMessage(eventType, data)
}

// replaceParameters replaces {param} placeholders in message with actual values
func (vh *VoiceHandler) replaceParameters(message string, data map[string]interface{}) string {
	result := message

	// Replace each parameter in the data map
	for key, value := range data {
		placeholder := fmt.Sprintf("{%s}", key)
		var replacement string

		switch v := value.(type) {
		case int64:
			replacement = fmt.Sprintf("%d", v)
		case int:
			replacement = fmt.Sprintf("%d", v)
		case float64:
			replacement = fmt.Sprintf("%.0f", v)
		case string:
			replacement = v
		case bool:
			if v {
				replacement = "sim"
			} else {
				replacement = "não"
			}
		default:
			replacement = fmt.Sprintf("%v", v)
		}

		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return result
}

// getStaticMessage returns hardcoded fallback messages
func (vh *VoiceHandler) getStaticMessage(eventType string, data interface{}) string {
	switch eventType {
	case "hero_health_low":
		return "Vida baixa!"
	case "hero_mana_low":
		return "Mana baixa!"
	case "hero_death":
		return "Você morreu!"
	case "hero_level_up":
		return "Level up!"
	case "hero_ultimate_ready":
		return "Ultimate pronto!"
	case "bounty_rune_warning":
		return "Runa de recompensa em breve!"
	case "power_rune_warning":
		return "Runa de poder em breve!"
	case "water_rune_warning":
		return "Runa de água em breve!"
	case "wisdom_rune_warning":
		return "Runa de sabedoria em breve!"
	case "catapult_timing":
		return "Catapulta chegando!"
	case "stack_timing":
		return "Hora de stackar!"
	case "day_night_cycle":
		return "Atenção: mudança de ciclo em breve!"
	default:
		return ""
	}
}

// getGameStateMessage returns Portuguese message for game state
func (vh *VoiceHandler) getGameStateMessage(gameState string) string {
	switch gameState {
	case "DOTA_GAMERULES_STATE_HERO_SELECTION":
		return "Seleção de heróis iniciada"
	case "DOTA_GAMERULES_STATE_STRATEGY_TIME":
		return "Tempo de estratégia"
	case "DOTA_GAMERULES_STATE_PRE_GAME":
		return "Pré-jogo iniciado"
	case "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS":
		return "Partida em andamento"
	case "DOTA_GAMERULES_STATE_POST_GAME":
		return "Partida finalizada"
	default:
		return ""
	}
}

// SetEnabled enables/disables voice
func (vh *VoiceHandler) SetEnabled(enabled bool) {
	vh.enabled = enabled && vh.apiKey != ""
	vh.logger.WithField("enabled", vh.enabled).Info("Voice handler state changed")
}

// startCacheCleanup runs a goroutine to clean old cache files
func (vh *VoiceHandler) startCacheCleanup() {
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	defer ticker.Stop()

	vh.logger.Info("🗑️ Cache cleanup goroutine started (24h interval)")
	
	// Run initial cleanup
	vh.cleanupOldCache()
	
	for range ticker.C {
		vh.cleanupOldCache()
	}
}

// cleanupOldCache removes cache files older than 7 days
func (vh *VoiceHandler) cleanupOldCache() {
	if vh.cachePath == "" {
		return
	}

	now := time.Now()
	removed := 0
	totalSize := int64(0)

	// Read all files in cache directory
	entries, err := os.ReadDir(vh.cachePath)
	if err != nil {
		vh.logger.WithError(err).Warn("Failed to read cache directory")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Remove files older than 7 days
		if now.Sub(info.ModTime()) > 7*24*time.Hour {
			filePath := filepath.Join(vh.cachePath, entry.Name())
			totalSize += info.Size()
			
			if err := os.Remove(filePath); err != nil {
				vh.logger.WithError(err).Warn("Failed to remove cache file")
			} else {
				removed++
			}
		}
	}

	if removed > 0 {
		vh.logger.WithFields(logrus.Fields{
			"files_removed": removed,
			"size_freed_mb": float64(totalSize) / (1024 * 1024),
		}).Info("🧹 Cache cleanup completed")
	}
}
