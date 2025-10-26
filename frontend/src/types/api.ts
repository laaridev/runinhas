/**
 * Type definitions for API responses and requests
 */

export interface VoiceConfig {
  apiKey: string;
  voiceId: string;
  stability: number;
  similarity: number;
  style: number;
  speakerBoost: boolean;
}

export interface VoiceSettings {
  stability: number;
  similarity_boost: number;
  style: number;
  use_speaker_boost: boolean;
}

export interface Voice {
  voice_id: string;
  name: string;
  preview_url?: string;
  category?: string;
}

export interface AudioGenerateResponse {
  filename: string;
  success: boolean;
  message?: string;
}

export interface SystemStatus {
  first_run: boolean;
  gsi_installed: boolean;
}

export interface TimingValue {
  value: number;
}

export interface TimingEnabled {
  enabled: boolean;
}

export interface MessageResponse {
  message: string;
}

export interface AudioCheckResponse {
  exists: boolean;
}

export interface ConfigResponse {
  voice?: VoiceConfig;
  [key: string]: any;
}

export interface VoicesResponse {
  voices: Voice[];
}

export interface TestVoiceRequest {
  text: string;
  voiceId: string;
  settings: VoiceSettings;
}

export interface TestVoiceResponse {
  filename: string;
  success: boolean;
}

export interface TimingEventMetadata {
  enabled: boolean;
  warning_seconds: number;
  min: number;
  max: number;
  step: number;
  name: string;
  description: string;
  category: 'rune' | 'timing';
}

export interface EventsMetadataResponse {
  [key: string]: TimingEventMetadata;
}
