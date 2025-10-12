package validation

import (
	"fmt"
	"strings"
)

// Validator provides data validation
type Validator struct {
	errors []string
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: []string{},
	}
}

// ValidateTimingKey validates a timing key
func (v *Validator) ValidateTimingKey(key string) *Validator {
	validKeys := map[string]bool{
		"bounty_rune":      true,
		"power_rune":       true,
		"water_rune":       true,
		"wisdom_rune":      true,
		"stack_timing":     true,
		"day_night_cycle":  true,
		"catapult_timing":  true,
	}
	
	if !validKeys[key] {
		v.errors = append(v.errors, fmt.Sprintf("invalid timing key: %s", key))
	}
	
	return v
}

// ValidateTimingField validates a timing field
func (v *Validator) ValidateTimingField(field string) *Validator {
	validFields := map[string]bool{
		"enabled":         true,
		"warning_seconds": true,
		"first_spawn":     true,
		"interval":        true,
	}
	
	if !validFields[field] {
		v.errors = append(v.errors, fmt.Sprintf("invalid timing field: %s", field))
	}
	
	return v
}

// ValidateTimingValue validates a timing value
func (v *Validator) ValidateTimingValue(value int) *Validator {
	if value < 0 || value > 300 {
		v.errors = append(v.errors, fmt.Sprintf("timing value out of range (0-300): %d", value))
	}
	
	return v
}

// ValidateMessage validates a custom message
func (v *Validator) ValidateMessage(message string) *Validator {
	if len(message) > 500 {
		v.errors = append(v.errors, "message too long (max 500 characters)")
	}
	
	// Check for dangerous content
	dangerous := []string{"<script", "javascript:", "onclick", "onerror"}
	lower := strings.ToLower(message)
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			v.errors = append(v.errors, "message contains potentially dangerous content")
			break
		}
	}
	
	return v
}

// ValidateEventType validates an event type for audio
func (v *Validator) ValidateEventType(eventType string) *Validator {
	validTypes := map[string]bool{
		"bounty_rune":     true,
		"power_rune":      true,
		"water_rune":      true,
		"wisdom_rune":     true,
		"stack_timing":    true,
		"day_night_cycle": true,
		"catapult_timing": true,
	}
	
	if !validTypes[eventType] {
		v.errors = append(v.errors, fmt.Sprintf("invalid event type: %s", eventType))
	}
	
	return v
}

// IsValid returns whether validation passed
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors returns validation errors
func (v *Validator) Errors() []string {
	return v.errors
}

// Error returns a single error string
func (v *Validator) Error() string {
	if len(v.errors) == 0 {
		return ""
	}
	return strings.Join(v.errors, "; ")
}

// Reset clears validation errors
func (v *Validator) Reset() *Validator {
	v.errors = []string{}
	return v
}
