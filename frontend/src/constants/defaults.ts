/**
 * Default configuration values for the application
 */

// Voice configuration defaults
export const DEFAULT_VOICE_CONFIG = {
  voiceId: 'eVXYtPVYB9wDoz9NVTIy',
  stability: 0.5,
  similarity: 0.75,
  style: 0,
  speakerBoost: true,
} as const;

// Timing constants
export const TIMING = {
  DEBOUNCE_DELAY: 800, // ms to wait after slider stops moving
  API_RETRY_DELAY: 100, // ms to wait after API call before next action
  TOAST_DURATION: 3000, // ms for toast auto-dismiss
} as const;

// Retry configuration
export const RETRY_CONFIG = {
  MAX_RETRIES_SAVE: 3,
  MAX_RETRIES_AUDIO: 2,
  DELAY_MS: 500,
  BACKOFF_DELAY: 1000,
} as const;

// Animation durations
export const ANIMATION = {
  TAB_TRANSITION: 300, // ms for tab transitions
  BUTTON_TRANSITION: 200, // ms for button state changes
} as const;
