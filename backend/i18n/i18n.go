package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed translations/*.json
var translationsFS embed.FS

// Translator manages translations for different locales
type Translator struct {
	translations  map[string]map[string]interface{}
	currentLocale string
	mu            sync.RWMutex
}

var (
	defaultTranslator *Translator
	once              sync.Once
)

// Init initializes the translator with a default locale
func Init(locale string) error {
	var initErr error
	
	once.Do(func() {
		t := &Translator{
			translations:  make(map[string]map[string]interface{}),
			currentLocale: locale,
		}
		
		// Load all translation files
		locales := []string{"pt-BR", "en"}
		for _, loc := range locales {
			data, err := translationsFS.ReadFile(fmt.Sprintf("translations/%s.json", loc))
			if err != nil {
				initErr = fmt.Errorf("failed to load translations for %s: %w", loc, err)
				return
			}
			
			var trans map[string]interface{}
			if err := json.Unmarshal(data, &trans); err != nil {
				initErr = fmt.Errorf("failed to parse translations for %s: %w", loc, err)
				return
			}
			
			t.translations[loc] = trans
		}
		
		defaultTranslator = t
	})
	
	return initErr
}

// SetLocale changes the current locale
func SetLocale(locale string) {
	if defaultTranslator != nil {
		defaultTranslator.SetLocale(locale)
	}
}

// GetLocale returns the current locale
func GetLocale() string {
	if defaultTranslator != nil {
		return defaultTranslator.GetLocale()
	}
	return "en"
}

// T translates a key with optional parameters
func T(key string, params map[string]interface{}) string {
	if defaultTranslator == nil {
		return key
	}
	return defaultTranslator.Translate(key, params)
}

// SetLocale changes the current locale (instance method)
func (t *Translator) SetLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.currentLocale = locale
}

// GetLocale returns the current locale (instance method)
func (t *Translator) GetLocale() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentLocale
}

// Translate translates a key with parameters
func (t *Translator) Translate(key string, params map[string]interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	// Split key by dots: "events.bounty_rune.name"
	parts := strings.Split(key, ".")
	
	// Try current locale
	value := t.getValue(t.currentLocale, parts)
	if value == "" {
		// Fallback to English
		value = t.getValue("en", parts)
	}
	
	if value == "" {
		return key // Return key if not found
	}
	
	// Replace parameters like {seconds}, {error}
	return t.interpolate(value, params)
}

// getValue retrieves a nested value from translations
func (t *Translator) getValue(locale string, keys []string) string {
	trans, ok := t.translations[locale]
	if !ok {
		return ""
	}
	
	var current interface{} = trans
	for _, key := range keys {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[key]
		} else {
			return ""
		}
	}
	
	if str, ok := current.(string); ok {
		return str
	}
	return ""
}

// interpolate replaces placeholders like {seconds} with actual values
func (t *Translator) interpolate(text string, params map[string]interface{}) string {
	if params == nil {
		return text
	}
	
	result := text
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}
