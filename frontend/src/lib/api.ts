// API client for direct backend communication
// Frontend calls backend DIRECTLY, not through Wails app.go!

const BACKEND_URL = 'http://localhost:3001';

// Configuration API
export const configAPI = {
  async get() {
    const res = await fetch(`${BACKEND_URL}/api/config`);
    return res.json();
  },

  async save(config: any) {
    const res = await fetch(`${BACKEND_URL}/api/config`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config)
    });
    return res.ok;
  },

  async isFirstRun() {
    const res = await fetch(`${BACKEND_URL}/api/config/first-run`);
    const data = await res.json();
    return data.firstRun;
  },

  async isGSIInstalled() {
    const res = await fetch(`${BACKEND_URL}/api/config/gsi-installed`);
    const data = await res.json();
    return data.installed;
  },

  async setGSIInstalled() {
    const res = await fetch(`${BACKEND_URL}/api/config/gsi-installed`, {
      method: 'POST'
    });
    return res.ok;
  }
};

// Timing API
export const timingAPI = {
  async getValue(key: string, field: string) {
    const res = await fetch(`${BACKEND_URL}/api/timing/${key}/${field}`);
    const data = await res.json();
    return data.value;
  },

  async setValue(key: string, field: string, value: number) {
    const res = await fetch(`${BACKEND_URL}/api/timing/${key}/${field}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ value })
    });
    return res.ok;
  },

  async getEnabled(key: string) {
    const res = await fetch(`${BACKEND_URL}/api/timing/${key}/enabled`);
    const data = await res.json();
    return data.enabled;
  },

  async setEnabled(key: string, enabled: boolean) {
    const res = await fetch(`${BACKEND_URL}/api/timing/${key}/enabled`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled })
    });
    return res.ok;
  }
};

// Message API
export const messageAPI = {
  async get(key: string) {
    const res = await fetch(`${BACKEND_URL}/api/message/${key}`);
    const data = await res.json();
    return data.message;
  },

  async set(key: string, message: string) {
    const res = await fetch(`${BACKEND_URL}/api/message/${key}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message })
    });
    return res.ok;
  }
};

// Audio API
export const audioAPI = {
  async check(eventType: string) {
    const res = await fetch(`${BACKEND_URL}/api/audio/check/${eventType}`);
    const data = await res.json();
    return data.exists;
  },

  async generate(eventType: string) {
    const res = await fetch(`${BACKEND_URL}/api/audio/generate/${eventType}`, {
      method: 'POST'
    });
    return res.ok;
  },

  async preview(eventType: string) {
    const res = await fetch(`${BACKEND_URL}/api/audio/preview/${eventType}`, {
      method: 'POST'
    });
    return res.ok;
  }
};

// Health check
export async function checkBackendHealth() {
  try {
    const res = await fetch(`${BACKEND_URL}/health`);
    return res.ok;
  } catch {
    return false;
  }
}
