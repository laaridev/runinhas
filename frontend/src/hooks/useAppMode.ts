import { useState, useEffect } from 'react';
import { GetMode } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export interface AppMode {
  mode: 'free' | 'pro';
  features: {
    customMessages: boolean;
    customVoice: boolean;
    elevenLabs: boolean;
    audioGeneration: boolean;
    embeddedAudioOnly: boolean;
  };
}

/**
 * Hook to get the current app mode (free or pro)
 * and feature flags
 */
export function useAppMode() {
  const [appMode, setAppMode] = useState<AppMode>({
    mode: 'free',
    features: {
      customMessages: false,
      customVoice: false,
      elevenLabs: false,
      audioGeneration: false,
      embeddedAudioOnly: true,
    },
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAppMode();
    
    // Listen for mode changes from backend
    EventsOn("mode:changed", (newMode: string) => {
      console.log(`🔄 Mode changed to: ${newMode}`);
      loadAppMode();
    });
  }, []);

  const loadAppMode = async () => {
    try {
      setLoading(true);
      const response = await GetMode();
      const data = JSON.parse(response);
      setAppMode(data);
    } catch (err) {
      console.error('Failed to load app mode:', err);
      setError('Failed to load app mode');
      // Default to free mode on error
      setAppMode({
        mode: 'free',
        features: {
          customMessages: false,
          customVoice: false,
          elevenLabs: false,
          audioGeneration: false,
          embeddedAudioOnly: true,
        },
      });
    } finally {
      setLoading(false);
    }
  };

  return { appMode, loading, error, reload: loadAppMode };
}
