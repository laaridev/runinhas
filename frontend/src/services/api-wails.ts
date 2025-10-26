// API Service - Versão que usa o proxy do Wails para evitar CORS
// Use esta versão se a api.ts direta não funcionar por causa de CORS

import { ProxyToBackend } from '../../wailsjs/go/main/App';
import type {
  VoiceConfig,
  Voice,
  SystemStatus,
  AudioCheckResponse,
  ConfigResponse,
  VoicesResponse,
  VoiceSettings,
  EventsMetadataResponse,
  TimingEventMetadata,
} from '@/types/api';
import { safeJsonParse } from '@/utils/json';
import { DEFAULT_VOICE_CONFIG } from '@/constants/defaults';

// Timing API via Wails proxy
export const timingAPI = {
  async getValue(key: string, field: string): Promise<number> {
    try {
      const response = await ProxyToBackend('GET', `/api/timing/${key}/${field}`, '');
      const data = JSON.parse(response);
      return data.value;
    } catch (error) {
      console.error('Failed to get timing value:', error);
      return 0;
    }
  },

  async setValue(key: string, field: string, value: number): Promise<void> {
    try {
      await ProxyToBackend('POST', `/api/timing/${key}/${field}`, JSON.stringify({ value }));
    } catch (error) {
      console.error('Failed to set timing value:', error);
      throw error;
    }
  },

  async getEnabled(key: string): Promise<boolean> {
    try {
      const response = await ProxyToBackend('GET', `/api/timing/${key}/enabled`, '');
      const data = JSON.parse(response);
      return data.enabled;
    } catch (error) {
      console.error('Failed to get timing enabled:', error);
      return false;
    }
  },

  async setEnabled(key: string, enabled: boolean): Promise<void> {
    try {
      await ProxyToBackend('POST', `/api/timing/${key}/enabled`, JSON.stringify({ enabled }));
    } catch (error) {
      console.error('Failed to set timing enabled:', error);
      throw error;
    }
  }
};

// Message API via Wails proxy
export const messageAPI = {
  async get(key: string): Promise<string> {
    try {
      const response = await ProxyToBackend('GET', `/api/message/${key}`, '');
      const data = JSON.parse(response);
      return data.message || '';
    } catch (error) {
      console.error(`[API-Wails] Failed to get message:`, error);
      return '';
    }
  },

  async set(key: string, message: string): Promise<void> {
    try {
      await ProxyToBackend('POST', `/api/message/${key}`, JSON.stringify({ message }));
    } catch (error) {
      console.error(`[API-Wails] Failed to set message:`, error);
      throw error;
    }
  }
};

// System API via Wails proxy
export const systemAPI = {
  async getStatus(): Promise<SystemStatus> {
    try {
      const response = await ProxyToBackend('GET', '/api/system/status', '');
      return safeJsonParse<SystemStatus>(response, { 
        first_run: true, 
        gsi_installed: false 
      });
    } catch (error) {
      console.error('Failed to get system status:', error);
      return { first_run: true, gsi_installed: false };
    }
  }
};

// Audio API via Wails proxy
export const audioAPI = {
  async generate(eventType: string): Promise<string> {
    try {
      const response = await ProxyToBackend('POST', `/api/audio/generate/${eventType}`, '');
      return response;
    } catch (error) {
      console.error(`[API-Wails] Failed to generate audio:`, error);
      throw error;
    }
  },

  async preview(eventType: string): Promise<void> {
    try {
      await ProxyToBackend('POST', `/api/audio/preview/${eventType}`, '');
    } catch (error) {
      console.error(`[API-Wails] Failed to preview audio:`, error);
      throw error;
    }
  },

  // Get audio file URL through proxy
  getAudioUrl(filename: string): string {
    // Return a data URL or blob URL after fetching through proxy
    return `http://localhost:3001/api/audio/file/${filename}`;
  },

  // Fetch audio file as blob directly (faster than base64)
  async getAudioBlob(filename: string): Promise<Blob> {
    try {
      // Direct fetch is faster and works in Wails
      const response = await fetch(`http://localhost:3001/api/audio/file/${filename}`);
      if (!response.ok) {
        throw new Error(`Failed to fetch audio: ${response.status}`);
      }
      return await response.blob();
    } catch (error) {
      console.error(`[API-Wails] Failed to get audio blob:`, error);
      throw error;
    }
  },

  async check(eventType: string): Promise<boolean> {
    try {
      const response = await ProxyToBackend('GET', `/api/audio/check/${eventType}`, '');
      const data = safeJsonParse<AudioCheckResponse>(response, { exists: false });
      return data.exists || false;
    } catch (error) {
      console.error(`[API-Wails] Failed to check audio:`, error);
      return false;
    }
  }
};

// Events Metadata API via Wails proxy
export const eventsAPI = {
  async getAll(): Promise<EventsMetadataResponse> {
    try {
      const response = await ProxyToBackend('GET', '/api/events', '');
      return safeJsonParse<EventsMetadataResponse>(response, {});
    } catch (error) {
      console.error('[API-Wails] Failed to get events metadata:', error);
      return {};
    }
  },

  async getEvent(key: string): Promise<TimingEventMetadata | null> {
    try {
      const response = await ProxyToBackend('GET', `/api/events/${key}`, '');
      const parsed = JSON.parse(response);
      return parsed as TimingEventMetadata;
    } catch (error) {
      console.error(`[API-Wails] Failed to get event ${key}:`, error);
      return null;
    }
  }
};

// ElevenLabs Voice API via Wails proxy
export const voiceAPI = {
  async getConfig(): Promise<VoiceConfig> {
    try {
      const response = await ProxyToBackend('GET', '/api/config', '');
      const data = safeJsonParse<ConfigResponse>(response, { voice: undefined });
      return data.voice || { apiKey: '', ...DEFAULT_VOICE_CONFIG };
    } catch (error) {
      console.error('[API-Wails] Failed to get voice config:', error);
      return { apiKey: '', ...DEFAULT_VOICE_CONFIG };
    }
  },

  async saveConfig(config: VoiceConfig): Promise<void> {
    try {
      // Get current config first
      const response = await ProxyToBackend('GET', '/api/config', '');
      const currentData = JSON.parse(response);
      
      // Update only the voice part
      currentData.voice = config;
      
      // Save back
      await ProxyToBackend('POST', '/api/config', JSON.stringify(currentData));
    } catch (error) {
      console.error('[API-Wails] Failed to save voice config:', error);
      throw error;
    }
  },

  async testVoice(text: string, voiceId: string, settings: VoiceSettings): Promise<string> {
    try {
      const payload = {
        text,
        voiceId,
        settings
      };
      const response = await ProxyToBackend('POST', '/api/elevenlabs/test', JSON.stringify(payload));
      const data = JSON.parse(response);
      return data.filename; // Return filename for frontend to play
    } catch (error) {
      console.error('[API-Wails] Failed to test voice:', error);
      throw error;
    }
  },

  async getVoices(): Promise<Voice[]> {
    try {
      const response = await ProxyToBackend('GET', '/api/elevenlabs/voices', '');
      const data = safeJsonParse<VoicesResponse>(response, { voices: [] });
      return data.voices || [];
    } catch (error) {
      console.error('[API-Wails] Failed to get voices:', error);
      return [];
    }
  }
};
