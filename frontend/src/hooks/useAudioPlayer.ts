import { useEffect, useRef, useState } from 'react';
import { audioAPI } from '@/services/api-wails';

interface AudioEvent {
  filename: string;
  eventType: string;
  data: Record<string, any>;
  timestamp: number;
}

interface AudioPlayerState {
  isPlaying: boolean;
  currentFile: string | null;
  queueLength: number;
}

/**
 * Hook to manage audio playback from backend events
 * Handles sequential playback of audio files
 */
export function useAudioPlayer() {
  const [state, setState] = useState<AudioPlayerState>({
    isPlaying: false,
    currentFile: null,
    queueLength: 0,
  });

  const audioRef = useRef<HTMLAudioElement | null>(null);
  const queueRef = useRef<AudioEvent[]>([]);
  const isProcessingRef = useRef(false);

  // Process audio queue
  const processQueue = async () => {
    if (isProcessingRef.current || queueRef.current.length === 0) {
      return;
    }

    isProcessingRef.current = true;
    const event = queueRef.current.shift()!;

    setState({
      isPlaying: true,
      currentFile: event.filename,
      queueLength: queueRef.current.length,
    });

    try {
      await playAudio(event.filename);
    } catch (error) {
      console.error('Failed to play audio:', error);
    }

    setState({
      isPlaying: false,
      currentFile: null,
      queueLength: queueRef.current.length,
    });

    isProcessingRef.current = false;

    // Process next item in queue
    if (queueRef.current.length > 0) {
      setTimeout(processQueue, 100); // Small delay between tracks
    }
  };

  // Play a single audio file
  const playAudio = async (filename: string): Promise<void> => {
    return new Promise(async (resolve, reject) => {
      if (!audioRef.current) {
        reject(new Error('Audio element not initialized'));
        return;
      }

      try {
        // Fetch audio as blob through Wails proxy
        const blob = await audioAPI.getAudioBlob(filename);
        const blobUrl = URL.createObjectURL(blob);

        const audio = audioRef.current;
        audio.src = blobUrl;

        audio.onended = () => {
          URL.revokeObjectURL(blobUrl); // Clean up
          resolve();
        };

        audio.onerror = () => {
          URL.revokeObjectURL(blobUrl); // Clean up
          reject(new Error(`Failed to load audio: ${filename}`));
        };

        await audio.play();
      } catch (error) {
        reject(error);
      }
    });
  };

  // Add audio event to queue
  const enqueueAudio = (event: AudioEvent) => {
    queueRef.current.push(event);
    setState((prev) => ({
      ...prev,
      queueLength: queueRef.current.length,
    }));
    processQueue();
  };

  // Initialize audio element and SSE connection
  useEffect(() => {
    audioRef.current = new Audio();
    audioRef.current.preload = 'auto';

    // Connect to SSE stream for audio events
    const eventSource = new EventSource('http://localhost:3001/api/audio/events');

    eventSource.onmessage = (event) => {
      try {
        const audioEvent: AudioEvent = JSON.parse(event.data);
        console.log('🎵 Received audio event:', audioEvent);
        enqueueAudio(audioEvent);
      } catch (error) {
        console.error('Failed to parse audio event:', error);
      }
    };

    eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
      // Will automatically reconnect
    };

    // Cleanup on unmount
    return () => {
      eventSource.close();
      if (audioRef.current) {
        audioRef.current.pause();
        audioRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // enqueueAudio is stable, no need to add to deps

  // Play audio immediately (for testing)
  const playNow = async (filename: string) => {
    try {
      await playAudio(filename);
    } catch (error) {
      console.error('Failed to play audio:', error);
    }
  };

  // Clear queue
  const clearQueue = () => {
    queueRef.current = [];
    setState((prev) => ({
      ...prev,
      queueLength: 0,
    }));
    if (audioRef.current) {
      audioRef.current.pause();
    }
  };

  return {
    state,
    enqueueAudio,
    playNow,
    clearQueue,
  };
}
