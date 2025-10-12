package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/sirupsen/logrus"
)

// audioQueueItem represents an audio file to be played
type audioQueueItem struct {
	filePath  string
	eventType string
	data      interface{}
}

// VoiceHandler handles voice announcements using ElevenLabs + oto
type VoiceHandler struct {
	apiKey       string
	voiceID      string
	cachePath    string
	logger       *logrus.Entry
	enabled      bool
	initialized  bool
	gameConfig   interface{}  // Game configuration for messages
	initMutex    sync.Mutex   // Mutex to prevent double initialization
	otoContext   *oto.Context // Audio context
	audioQueue   chan audioQueueItem // Queue for sequential audio playback
	queueMutex   sync.Mutex   // Mutex for queue operations
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
		"enabled": vh.enabled,
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
		return nil, fmt.Errorf("failed to create voice cache directory: %w", err)
	}

	enabled := apiKey != ""

	vh := &VoiceHandler{
		apiKey:       apiKey,
		voiceID:      voiceID,
		cachePath:    finalCachePath,
		logger:       logger,
		enabled:      enabled,
		audioQueue:   make(chan audioQueueItem, 10), // Buffer up to 10 audio items
		// Default voice settings
		stability:    0.5,
		similarity:   0.75,
		style:        0,
		speakerBoost: true,
	}

	// Start audio queue processor
	if enabled {
		go vh.processAudioQueue()
	}

	// Log configuration status
	if enabled {
		logger.Info("🎵 Voice handler initialized with ElevenLabs")
	} else {
		logger.Warn("🔇 Voice handler disabled - no API key provided")
	}

	return vh, nil
}

// SetGameConfig sets the game configuration for message templates
func (vh *VoiceHandler) SetGameConfig(gameConfig interface{}) {
	vh.gameConfig = gameConfig
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

// initSpeaker initializes the audio context (thread-safe)
func (vh *VoiceHandler) initSpeaker() {
	vh.initMutex.Lock()
	defer vh.initMutex.Unlock()
	
	// Check if already initialized
	if vh.initialized {
		return
	}
	
	vh.logger.Info("🔊 Initializing oto audio context...")
	
	// Initialize oto context
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	
	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		vh.logger.WithError(err).Error("❌ Failed to initialize oto context")
		vh.enabled = false
		return
	}
	
	// Wait for context to be ready
	<-readyChan
	
	vh.otoContext = ctx
	vh.initialized = true
	vh.logger.Info("✅ Audio context initialized successfully")
}

// speak generates and plays voice audio with semantic caching
func (vh *VoiceHandler) speakWithData(text string, eventType string, data interface{}) {
	go func() {
		// Try semantic cache path first
		cacheFile := vh.getCacheFilePathSemantic(eventType, data)

		// Check if semantic cache file exists
		if _, err := os.Stat(cacheFile); err == nil {
			vh.logger.WithField("cache_file", filepath.Base(cacheFile)).Debug("Using semantic cache")
			// Add to queue instead of playing immediately
			vh.enqueueAudio(cacheFile, eventType, data)
			return
		}

		// Fallback to hash-based cache
		hashCacheFile := vh.getCacheFilePath(text)
		if _, err := os.Stat(hashCacheFile); err == nil {
			vh.logger.WithField("cache_file", filepath.Base(hashCacheFile)).Debug("Using hash cache")
			// Add to queue instead of playing immediately
			vh.enqueueAudio(hashCacheFile, eventType, data)
			return
		}

		// Generate with ElevenLabs
		audioData, err := vh.GenerateVoice(text)
		if err != nil {
			vh.logger.WithError(err).Error("Failed to generate voice")
			return
		}

		// Save to semantic cache path
		if err := os.WriteFile(cacheFile, audioData, 0644); err != nil {
			vh.logger.WithError(err).Error("Failed to cache audio")
			// Try hash-based cache as fallback
			if err := os.WriteFile(hashCacheFile, audioData, 0644); err != nil {
				vh.logger.WithError(err).Error("Failed to cache audio (fallback)")
			}
		}

		// Add to queue instead of playing immediately
			vh.enqueueAudio(cacheFile, eventType, data)
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
		"url": url,
	}).Debug("Generating voice with ElevenLabs")

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":        vh.stability,
			"similarity_boost": vh.similarity,
			"style":            vh.style,
			"use_speaker_boost": vh.speakerBoost,
		},
	}

	jsonData, _ := json.Marshal(payload)

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
			"status": resp.StatusCode,
			"body": string(body),
			"voice_id": vh.voiceID,
		}).Error("ElevenLabs API error")
		return nil, fmt.Errorf("ElevenLabs API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// enqueueAudio adds an audio file to the playback queue
func (vh *VoiceHandler) enqueueAudio(filePath string, eventType string, data interface{}) {
	item := audioQueueItem{
		filePath:  filePath,
		eventType: eventType,
		data:      data,
	}
	
	select {
	case vh.audioQueue <- item:
		vh.logger.WithField("file", filepath.Base(filePath)).Debug("🎵 Audio enqueued")
	default:
		vh.logger.Warn("⚠️ Audio queue full, dropping audio")
	}
}

// processAudioQueue processes the audio queue sequentially
func (vh *VoiceHandler) processAudioQueue() {
	vh.logger.Info("🎵 Audio queue processor started")
	
	for item := range vh.audioQueue {
		vh.PlayAudioFile(item.filePath)
		// Small delay between audio files
		time.Sleep(100 * time.Millisecond)
	}
	
	vh.logger.Info("🎵 Audio queue processor stopped")
}

// PlayAudioFile plays audio using system command (exported for API usage)
func (vh *VoiceHandler) PlayAudioFile(filePath string) {
	vh.logger.WithField("file", filepath.Base(filePath)).Debug("🎵 Playing audio file")
	
	// Try platform-specific players first
	switch runtime.GOOS {
	case "linux":
		// Linux: Try paplay (PulseAudio/PipeWire) first
		if vh.tryPlayWithCommand("paplay", filePath) {
			return
		}
		// Fallback to mpg123
		if vh.tryPlayWithCommand("mpg123", "-q", filePath) {
			return
		}
		
	case "windows":
		// Windows: Use oto directly (works well with default device)
		// No system command needed, will use oto fallback below
		
	case "darwin":
		// macOS: Use afplay (native player)
		if vh.tryPlayWithCommand("afplay", filePath) {
			return
		}
	}
	
	// Fallback to oto if no system player available
	vh.playAudioFileWithOto(filePath)
}

// tryPlayWithCommand tries to play audio with a system command
func (vh *VoiceHandler) tryPlayWithCommand(command string, args ...string) bool {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return false
	}
	vh.logger.Debug("✅ Audio playback completed")
	return true
}

// playAudioFileWithOto plays audio using oto (fallback)
func (vh *VoiceHandler) playAudioFileWithOto(filePath string) {
	vh.logger.WithField("file", filepath.Base(filePath)).Debug("🎵 Playing audio file with oto")
	
	if vh.otoContext == nil {
		vh.logger.Error("❌ Audio context not initialized")
		return
	}
	
	// Open MP3 file
	file, err := os.Open(filePath)
	if err != nil {
		vh.logger.WithError(err).Error("❌ Failed to open audio file")
		return
	}
	defer file.Close()
	
	// Decode MP3
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		vh.logger.WithError(err).Error("❌ Failed to decode MP3")
		return
	}
	
	vh.logger.Debug("🎵 MP3 decoded, creating player...")
	
	// Create player
	player := vh.otoContext.NewPlayer(decoder)
	defer player.Close()
	
	// Play audio
	player.Play()
	
	// Wait for playback to finish
	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}
	
	vh.logger.Debug("✅ Audio playback completed")
}

// getCacheFilePath generates cache file path based on text hash
func (vh *VoiceHandler) getCacheFilePath(text string) string {
	cachePath := vh.cachePath
	// Create filename based on text content for parameter-aware caching
	// This ensures different parameters generate different cache files
	// e.g., "Runa em 30 segundos" -> different file than "Runa em 20 segundos"
	hash := fmt.Sprintf("%x", []byte(text))

	// Limit hash length for filesystem compatibility
	if len(hash) > 16 {
		hash = hash[:16]
	}

	filename := fmt.Sprintf("voice_%s.mp3", hash)
	return filepath.Join(cachePath, filename)
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
	case "catapult_warning":
		filename = fmt.Sprintf("%s.mp3", eventType)

	case "day_night_warning", "day_night_transition":
		if cycleType, exists := dataMap["cycle_type"]; exists {
			filename = fmt.Sprintf("%s_%s.mp3", eventType, cycleType)
		} else {
			filename = fmt.Sprintf("%s.mp3", eventType)
		}

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
		// Handle day_night_warning and day_night_transition with cycle_type
		if eventType == "day_night_warning" || eventType == "day_night_transition" {
			if cycleType, exists := dataMap["cycle_type"]; exists {
				// For transition, cycle_type is "day" or "night" (what it became)
				// For warning, cycle_type is "day" or "night" (what's coming)
				var message string
				if eventType == "day_night_transition" {
					// "Amanheceu" or "Anoiteceu"
					if cycleType == "day" {
						message = "Amanheceu"
					} else {
						message = "Anoiteceu"
					}
				} else {
					// "Vai amanhecer em X segundos" or "Vai anoitecer em X segundos"
					if cycleType == "day" {
						message = "Vai amanhecer em {seconds} segundos"
					} else {
						message = "Vai anoitecer em {seconds} segundos"
					}
				}
				
				// Check if there's a custom message in config
				msgKey := fmt.Sprintf("%s_%s", eventType, cycleType)
				if customMsg := gc.GetMessage(msgKey); customMsg != "" {
					message = customMsg
				}
				
				return vh.replaceParameters(message, dataMap)
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
	case "catapult_warning":
		return "Catapulta chegando!"
	case "stack_timing":
		return "Hora de stackar!"
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
