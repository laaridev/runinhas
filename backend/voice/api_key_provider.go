package voice

import (
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

// APIKeyProvider fetches and caches ElevenLabs API key from remote sources
type APIKeyProvider struct {
	sources    []string
	cacheFile  string
	cacheTTL   time.Duration
	mu         sync.RWMutex
	cachedKey  string
	cacheTime  time.Time
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewAPIKeyProvider creates a new API key provider
func NewAPIKeyProvider(cacheDir string, logger *logrus.Logger) *APIKeyProvider {
	if logger == nil {
		logger = logrus.New()
	}

	// Ensure cache directory exists
	os.MkdirAll(cacheDir, 0755)

	return &APIKeyProvider{
		sources: []string{
			// Primary: GitHub Gist (secret - not indexed but accessible via link)
			"https://gist.githubusercontent.com/laaridev/3efb7853465a59ce9edd7e87b19909e3/raw/eleven.txt",
			// Fallback: Direct gist URL (alternative endpoint)
			"https://gist.github.com/laaridev/3efb7853465a59ce9edd7e87b19909e3/raw",
		},
		cacheFile: filepath.Join(cacheDir, ".apikey_cache"),
		cacheTTL:  24 * time.Hour, // Cache for 24 hours
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetAPIKey returns the ElevenLabs API key (from cache or remote)
func (p *APIKeyProvider) GetAPIKey() (string, error) {
	p.mu.RLock()
	// Check if we have a valid cached key in memory
	if p.cachedKey != "" && time.Since(p.cacheTime) < p.cacheTTL {
		p.logger.Debug("🔑 Using in-memory cached API key")
		key := p.cachedKey
		p.mu.RUnlock()
		return key, nil
	}
	p.mu.RUnlock()

	// Try to load from file cache
	if key, err := p.loadFromFileCache(); err == nil {
		p.mu.Lock()
		p.cachedKey = key
		p.cacheTime = time.Now()
		p.mu.Unlock()
		p.logger.Info("🔑 Loaded API key from file cache")
		return key, nil
	}

	// File cache miss or expired - fetch from remote sources
	p.logger.Info("🌐 Fetching API key from remote sources...")
	key, err := p.fetchFromRemote()
	if err != nil {
		p.logger.WithError(err).Error("❌ Failed to fetch API key from all sources")
		return "", fmt.Errorf("failed to fetch API key: %w", err)
	}

	// Cache the key
	p.mu.Lock()
	p.cachedKey = key
	p.cacheTime = time.Now()
	p.mu.Unlock()

	// Save to file cache
	if err := p.saveToFileCache(key); err != nil {
		p.logger.WithError(err).Warn("⚠️ Failed to save API key to file cache")
	}

	p.logger.Info("✅ Successfully fetched and cached API key")
	return key, nil
}

// fetchFromRemote tries to fetch the API key from remote sources (with fallback)
func (p *APIKeyProvider) fetchFromRemote() (string, error) {
	var lastErr error

	for i, source := range p.sources {
		p.logger.WithField("source", source).WithField("attempt", i+1).Debug("🔄 Trying source...")

		key, err := p.fetchFromURL(source)
		if err != nil {
			p.logger.WithError(err).WithField("source", source).Warn("⚠️ Source failed, trying next...")
			lastErr = err
			continue
		}

		// Validate key format (ElevenLabs keys start with "sk_")
		if !strings.HasPrefix(key, "sk_") {
			p.logger.WithField("source", source).Warn("⚠️ Invalid key format, trying next...")
			lastErr = fmt.Errorf("invalid key format from %s", source)
			continue
		}

		p.logger.WithField("source", source).Info("✅ Successfully fetched API key")
		return key, nil
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("no sources available")
}

// fetchFromURL fetches the API key from a single URL
func (p *APIKeyProvider) fetchFromURL(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add headers to avoid caching issues
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Runinhas/1.0")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Clean up the key (remove whitespace, newlines)
	key := strings.TrimSpace(string(body))
	if key == "" {
		return "", fmt.Errorf("empty response")
	}

	return key, nil
}

// loadFromFileCache loads the API key from file cache
func (p *APIKeyProvider) loadFromFileCache() (string, error) {
	info, err := os.Stat(p.cacheFile)
	if err != nil {
		return "", err
	}

	// Check if cache is expired
	if time.Since(info.ModTime()) > p.cacheTTL {
		return "", fmt.Errorf("cache expired")
	}

	data, err := os.ReadFile(p.cacheFile)
	if err != nil {
		return "", err
	}

	key := strings.TrimSpace(string(data))
	if key == "" || !strings.HasPrefix(key, "sk_") {
		return "", fmt.Errorf("invalid cached key")
	}

	return key, nil
}

// saveToFileCache saves the API key to file cache
func (p *APIKeyProvider) saveToFileCache(key string) error {
	// Create a temporary file first
	tmpFile := p.cacheFile + ".tmp"
	if err := os.WriteFile(tmpFile, []byte(key), 0600); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpFile, p.cacheFile)
}

// InvalidateCache clears the cached API key (useful for testing or key rotation)
func (p *APIKeyProvider) InvalidateCache() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cachedKey = ""
	p.cacheTime = time.Time{}

	// Remove file cache
	os.Remove(p.cacheFile)
	p.logger.Info("🗑️ API key cache invalidated")
}

// AddSource adds a new fallback source URL
func (p *APIKeyProvider) AddSource(url string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sources = append(p.sources, url)
	p.logger.WithField("url", url).Info("➕ Added new API key source")
}
