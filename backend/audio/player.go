package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/sirupsen/logrus"
)

// Player handles dual audio output (speakers + virtual mic)
type Player struct {
	logger            *logrus.Entry
	virtualMicEnabled bool
	virtualMicDevice  string
	speakerVolume     float64
	virtualMicVolume  float64
	mu                sync.RWMutex
	initialized       bool
	sampleRate        beep.SampleRate
}

// NewPlayer creates a new audio player
func NewPlayer(logger *logrus.Logger) *Player {
	return &Player{
		logger:            logger.WithField("component", "audio"),
		virtualMicEnabled: false,
		virtualMicDevice:  "virt_mic",
		speakerVolume:     1.0,
		virtualMicVolume:  0.8,
		initialized:       false,
		sampleRate:        beep.SampleRate(44100),
	}
}

// Initialize sets up the audio speaker
func (p *Player) Initialize() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.initialized {
		return nil
	}

	// Initialize beep speaker for default output
	err := speaker.Init(p.sampleRate, p.sampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	p.initialized = true
	p.logger.Info("🔊 Audio player initialized")
	
	return nil
}

// Play plays audio on speakers and optionally on virtual mic
func (p *Player) Play(audioData []byte, filename string) error {
	p.mu.RLock()
	virtualMicEnabled := p.virtualMicEnabled
	p.mu.RUnlock()

	// Always play on default speakers
	if err := p.playOnSpeakers(audioData); err != nil {
		p.logger.WithError(err).Warn("Failed to play audio on speakers")
	}

	// If virtual mic is enabled, also play there
	if virtualMicEnabled {
		if err := p.playOnVirtualMic(audioData, filename); err != nil {
			p.logger.WithError(err).Warn("Failed to play audio on virtual mic")
		}
	}

	return nil
}

// playOnSpeakers plays audio on default output using beep
func (p *Player) playOnSpeakers(audioData []byte) error {
	if !p.initialized {
		if err := p.Initialize(); err != nil {
			return err
		}
	}

	// Decode MP3
	streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(audioData)))
	if err != nil {
		return fmt.Errorf("failed to decode mp3: %w", err)
	}
	defer streamer.Close()

	// Resample if necessary
	var finalStreamer beep.Streamer = streamer
	if format.SampleRate != p.sampleRate {
		finalStreamer = beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
	}

	// Play and wait for completion
	done := make(chan bool)
	speaker.Play(beep.Seq(finalStreamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

// playOnVirtualMic plays audio on virtual microphone using paplay
func (p *Player) playOnVirtualMic(audioData []byte, filename string) error {
	// Only works on Linux with PulseAudio
	if runtime.GOOS != "linux" {
		p.logger.Debug("Virtual mic output only supported on Linux")
		return nil
	}

	p.mu.RLock()
	device := p.virtualMicDevice
	p.mu.RUnlock()

	// Create temp file for audio data
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("runinhas_%d_%s", time.Now().Unix(), filename))
	
	if err := os.WriteFile(tmpFile, audioData, 0644); err != nil {
		return fmt.Errorf("failed to write temp audio file: %w", err)
	}
	defer os.Remove(tmpFile)

	// Use paplay to output to virtual mic
	cmd := exec.Command("paplay", "--device="+device, tmpFile)
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start paplay: %w", err)
	}

	// Don't wait for completion (async playback on virtual mic)
	go func() {
		if err := cmd.Wait(); err != nil {
			p.logger.WithError(err).Debug("paplay finished with error")
		}
	}()

	p.logger.WithFields(logrus.Fields{
		"device": device,
		"file":   filename,
	}).Debug("🎤 Playing audio on virtual mic")

	return nil
}

// SetVirtualMicEnabled enables or disables virtual mic output
func (p *Player) SetVirtualMicEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.virtualMicEnabled = enabled
	
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	
	p.logger.WithField("enabled", enabled).Infof("🎤 Virtual mic output %s", status)
}

// SetVirtualMicDevice sets the virtual mic device name
func (p *Player) SetVirtualMicDevice(device string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.virtualMicDevice = device
	p.logger.WithField("device", device).Info("🎤 Virtual mic device set")
}

// GetVirtualMicEnabled returns whether virtual mic is enabled
func (p *Player) GetVirtualMicEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.virtualMicEnabled
}

// GetVirtualMicDevice returns the current virtual mic device
func (p *Player) GetVirtualMicDevice() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.virtualMicDevice
}

// DetectVirtualMic tries to detect virtual microphone on the system
func (p *Player) DetectVirtualMic() (string, bool) {
	if runtime.GOOS != "linux" {
		return "", false
	}

	// Try to list PulseAudio sources
	cmd := exec.Command("pactl", "list", "sources", "short")
	output, err := cmd.Output()
	if err != nil {
		p.logger.WithError(err).Debug("Failed to list PulseAudio sources")
		return "", false
	}

	// Look for common virtual mic names
	virtualMicNames := []string{"virt_mic", "virtual_mic", "Virtual_Microphone"}
	
	for _, name := range virtualMicNames {
		if bytes.Contains(output, []byte(name)) {
			p.logger.WithField("device", name).Info("✅ Virtual microphone detected")
			return name, true
		}
	}

	p.logger.Debug("No virtual microphone detected")
	return "", false
}

// Close cleans up resources
func (p *Player) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.initialized {
		// Speaker cleanup is automatic in beep
		p.initialized = false
		p.logger.Info("🔊 Audio player closed")
	}
}
