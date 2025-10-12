// API Service - Versão que usa o proxy do Wails para evitar CORS
// Use esta versão se a api.ts direta não funcionar por causa de CORS

import { ProxyToBackend } from '../../wailsjs/go/main/App';

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
  async getStatus(): Promise<{ first_run: boolean; gsi_installed: boolean }> {
    try {
      const response = await ProxyToBackend('GET', '/api/system/status', '');
      const data = JSON.parse(response);
      return data;
    } catch (error) {
      console.error('Failed to get system status:', error);
      return { first_run: true, gsi_installed: false };
    }
  }
};

// Audio API via Wails proxy
export const audioAPI = {
  async generate(eventType: string): Promise<void> {
    try {
      await ProxyToBackend('POST', `/api/audio/generate/${eventType}`, '');
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

  async check(eventType: string): Promise<boolean> {
    try {
      const response = await ProxyToBackend('GET', `/api/audio/check/${eventType}`, '');
      const data = JSON.parse(response);
      return data.exists || false;
    } catch (error) {
      console.error(`[API-Wails] Failed to check audio:`, error);
      return false;
    }
  }
};

// ElevenLabs Voice API via Wails proxy
export const voiceAPI = {
  async getConfig(): Promise<any> {
    try {
      const response = await ProxyToBackend('GET', '/api/config', '');
      const data = JSON.parse(response);
      return data.voice || {
        apiKey: '',
        voiceId: 'eVXYtPVYB9wDoz9NVTIy',
        stability: 0.5,
        similarity: 0.75,
        style: 0,
        speakerBoost: true
      };
    } catch (error) {
      console.error('[API-Wails] Failed to get voice config:', error);
      return {
        apiKey: '',
        voiceId: 'eVXYtPVYB9wDoz9NVTIy',
        stability: 0.5,
        similarity: 0.75,
        style: 0,
        speakerBoost: true
      };
    }
  },

  async saveConfig(config: any): Promise<void> {
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

  async testVoice(text: string, voiceId: string, settings: any): Promise<void> {
    try {
      const payload = {
        text,
        voiceId,
        settings
      };
      await ProxyToBackend('POST', '/api/elevenlabs/test', JSON.stringify(payload));
    } catch (error) {
      console.error('[API-Wails] Failed to test voice:', error);
      throw error;
    }
  },

  async getVoices(): Promise<any[]> {
    try {
      const response = await ProxyToBackend('GET', '/api/elevenlabs/voices', '');
      const data = JSON.parse(response);
      return data.voices || [];
    } catch (error) {
      console.error('[API-Wails] Failed to get voices:', error);
      return [];
    }
  }
};
