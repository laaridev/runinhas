export interface RuneConfig {
  key: string;
  name: string;
  description: string;
  icon: string;
  min: number;
  max: number;
  step: number;
  enabled?: boolean;
  value?: number;
  customMessage?: string;
}

export interface TimingConfig {
  key: string;
  name: string;
  description: string;
  icon: string;
  min: number;
  max: number;
  step: number;
  enabled?: boolean;
  value?: number;
  customMessage?: string;
}

export interface HandlerConfig {
  voice: {
    enabled: boolean;
    apiKey?: string;
    voiceId?: string;
  };
  discord: {
    enabled: boolean;
    webhookUrl?: string;
  };
  notify: {
    enabled: boolean;
  };
  overlay: {
    enabled: boolean;
    port?: number;
  };
}

export interface GameConfig {
  runes: Record<string, RuneConfig>;
  timings: Record<string, TimingConfig>;
  handlers: HandlerConfig;
}
