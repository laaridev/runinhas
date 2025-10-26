package main

import (
	"context"
	"dota-gsi/backend/installer"
	"dota-gsi/backend/server"
	"dota-gsi/backend/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct - Single binary with embedded backend
type App struct {
	ctx          context.Context
	gsiServer    *server.GSIServer
	gsiInstaller *installer.GSIInstaller
	serverStop   chan struct{}
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := utils.CreateLogger("installer")
	return &App{
		gsiInstaller: installer.NewGSIInstaller(logger),
		serverStop:   make(chan struct{}),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start the embedded backend server
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure context is ready
		if err := a.StartEmbeddedServer(); err != nil {
			wailsruntime.LogError(a.ctx, fmt.Sprintf("Failed to start embedded server: %v", err))
		}
	}()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.serverStop != nil {
		close(a.serverStop)
	}
	if a.gsiServer != nil {
		a.gsiServer.Stop()
	}
}

// === Embedded Server Management ===

// StartEmbeddedServer starts the backend server in the same process
func (a *App) StartEmbeddedServer() error {
	if a.gsiServer != nil {
		return fmt.Errorf("server already running")
	}

	// Create and start the GSI server
	var err error
	a.gsiServer, err = server.New()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Configure event emitter for Wails integration (audio playback, etc)
	if a.ctx != nil {
		a.gsiServer.SetEventEmitter(func(eventName string, data interface{}) {
			wailsruntime.EventsEmit(a.ctx, eventName, data)
		})
	}

	// Start server in goroutine
	go func() {
		wailsruntime.LogInfo(a.ctx, "Starting embedded backend server on :3001")
		if err := a.gsiServer.Start(":3001"); err != nil && err != http.ErrServerClosed {
			wailsruntime.LogError(a.ctx, fmt.Sprintf("Server error: %v", err))
		}
	}()

	// Wait for server to be ready
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://localhost:3001/health")
		if err == nil {
			resp.Body.Close()
			wailsruntime.LogInfo(a.ctx, "Embedded backend server started successfully")
			if a.ctx != nil {
				wailsruntime.EventsEmit(a.ctx, "server:started")
			}
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("server failed to start within timeout")
}

// StopEmbeddedServer stops the embedded backend server
func (a *App) StopEmbeddedServer() error {
	if a.gsiServer == nil {
		return fmt.Errorf("server not running")
	}

	err := a.gsiServer.Stop()
	a.gsiServer = nil

	if a.ctx != nil {
		wailsruntime.EventsEmit(a.ctx, "server:stopped")
	}

	return err
}

// IsServerRunning checks if server is running
func (a *App) IsServerRunning() bool {
	return a.gsiServer != nil
}

// GetBackendURL returns the backend URL for frontend to use
func (a *App) GetBackendURL() string {
	return "http://localhost:3001"
}

// GetMode returns the current app mode (free or pro)
func (a *App) GetMode() (string, error) {
	resp, err := a.ProxyToBackend("GET", "/api/mode", "")
	return resp, err
}

// SetMode sets the app mode (free or pro) with optional license key
func (a *App) SetMode(mode string, licenseKey string) error {
	// Prepare request body
	requestBody := map[string]string{
		"mode":        mode,
		"license_key": licenseKey,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call backend endpoint
	resp, err := a.ProxyToBackend("POST", "/api/mode", string(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to update mode: %w", err)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if successful
	if success, ok := response["success"].(bool); !ok || !success {
		if errMsg, ok := response["error"].(string); ok {
			return fmt.Errorf("mode update failed: %s", errMsg)
		}
		return fmt.Errorf("mode update failed")
	}

	// Emit event to frontend
	wailsruntime.EventsEmit(a.ctx, "mode:changed", mode)

	return nil
}

// ProxyToBackend proxies requests from frontend to backend
func (a *App) ProxyToBackend(method, path, body string) (string, error) {
	url := fmt.Sprintf("http://localhost:3001%s", path)

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("backend returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	return string(responseBody), nil
}

// === GSI Installation (System-level, must stay in Wails) ===

// InstallGSI installs the GSI configuration
func (a *App) InstallGSI() map[string]interface{} {
	result := a.gsiInstaller.Install()
	return map[string]interface{}{
		"success": result.Success,
		"message": result.Message,
	}
}

// IsGSIInstalled checks if GSI is installed
func (a *App) IsGSIInstalled() bool {
	installed, _ := a.gsiInstaller.CheckInstallation()
	return installed
}

// IsDotaInstalled checks if Dota 2 is installed on the system
func (a *App) IsDotaInstalled() map[string]interface{} {
	installed, path := a.gsiInstaller.CheckDotaInstalled()
	if !installed {
		return map[string]interface{}{
			"installed": false,
			"message":   "Dota 2 não encontrado. Por favor, instale o Dota 2 via Steam.",
		}
	}
	return map[string]interface{}{
		"installed": true,
		"path":      path,
	}
}

// GetDotaPath returns the Dota 2 installation path
func (a *App) GetDotaPath() string {
	_, path := a.gsiInstaller.CheckInstallation()
	return path
}

// === Language Management ===

// GetLanguage returns the current language setting
func (a *App) GetLanguage() string {
	resp, err := a.ProxyToBackend("GET", "/api/config/language", "")
	if err != nil {
		return "pt-BR" // Default fallback
	}
	return resp
}

// SetLanguage changes the application language
func (a *App) SetLanguage(language string) error {
	_, err := a.ProxyToBackend("POST", "/api/config/language", fmt.Sprintf(`{"language":"%s"}`, language))
	return err
}

// === Window Management (Wails-specific) ===

// MinimizeWindow minimizes the application window
func (a *App) MinimizeWindow() {
	if a.ctx != nil {
		wailsruntime.WindowMinimise(a.ctx)
	}
}

// CloseWindow closes the application window
func (a *App) CloseWindow() {
	if a.ctx != nil {
		wailsruntime.Quit(a.ctx)
	}
}
