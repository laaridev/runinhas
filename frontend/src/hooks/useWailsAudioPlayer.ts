import { useEffect, useRef, useState } from 'react';
import { EventsOn, EventsOff } from '../../wailsjs/runtime';

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

export function useWailsAudioPlayer() {
  const [state, setState] = useState<AudioPlayerState>({
    isPlaying: false,
    currentFile: null,
    queueLength: 0,
  });

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
    
    // Process next in queue
    if (queueRef.current.length > 0) {
      setTimeout(processQueue, 500);
    }
  };

  // Play audio file
  const playAudio = async (filename: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      const audioUrl = `http://localhost:3001/api/audio/file/${filename}`;
      
      const audio = new Audio(audioUrl);
      
      audio.onended = () => resolve();
      audio.onerror = () => reject(new Error(`Failed to play ${filename}`));
      
      audio.play().catch(reject);
    });
  };

  // Add audio event to queue
  const enqueueAudio = (event: AudioEvent) => {
    // Check if we should play this audio
    if (!shouldPlayAudio(event)) {
      return;
    }

    queueRef.current.push(event);
    setState(prev => ({
      ...prev,
      queueLength: queueRef.current.length,
    }));

    // Start processing if not already processing
    if (!isProcessingRef.current) {
      processQueue();
    }
  };

  // Filter logic for audio events
  const shouldPlayAudio = (_event: AudioEvent): boolean => {
    // You can add logic here to filter which events should trigger audio
    // For now, play all audio events
    return true;
  };

  // Initialize Wails event listener
  useEffect(() => {
    console.log('🎧 Setting up Wails audio event listener...');
    
    // Listen for audio events from backend
    EventsOn('audio:play', (audioEvent: any) => {
      console.log('🎵 Received Wails audio event:', audioEvent);
      console.log('🔍 Event structure:', {
        hasFilename: 'filename' in audioEvent,
        hasEventType: 'eventType' in audioEvent,
        keys: Object.keys(audioEvent || {}),
        raw: JSON.stringify(audioEvent)
      });
      
      // Try to handle the event regardless of structure
      const event: AudioEvent = {
        filename: audioEvent.filename || audioEvent.Filename || 'unknown.mp3',
        eventType: audioEvent.eventType || audioEvent.EventType || 'unknown',
        data: audioEvent.data || audioEvent.Data || {},
        timestamp: audioEvent.timestamp || audioEvent.Timestamp || Date.now()
      };
      
      enqueueAudio(event);
    });
    
    // Also listen for server started event to confirm connection
    EventsOn('server:started', () => {
      console.log('✅ Server started, audio events should work now');
    });

    // Cleanup on unmount
    return () => {
      EventsOff('audio:play');
      EventsOff('server:started');
    };
  }, []);

  return state;
}
