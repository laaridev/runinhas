import { useState, useEffect } from 'react';
import { ProxyToBackend } from '../../wailsjs/go/main/App';

export function useVirtualMic() {
  const [enabled, setEnabled] = useState(false);
  const [device, setDevice] = useState('virt_mic');
  const [detected, setDetected] = useState(false);
  const [loading, setLoading] = useState(true);

  // Load current status
  const loadStatus = async () => {
    try {
      const response = await ProxyToBackend('GET', '/api/audio/virtualmicEnabled', '');
      const data = JSON.parse(response);
      setEnabled(data.enabled || false);
      setDevice(data.device || 'virt_mic');
    } catch (error) {
      console.error('Failed to load virtual mic status:', error);
    } finally {
      setLoading(false);
    }
  };

  // Detect virtual mic
  const detectVirtualMic = async () => {
    try {
      const response = await ProxyToBackend('GET', '/api/audio/virtualmicDevice', '');
      const data = JSON.parse(response);
      setDetected(data.found || false);
      if (data.found && data.device) {
        setDevice(data.device);
      }
      return data.found;
    } catch (error) {
      console.error('Failed to detect virtual mic:', error);
      return false;
    }
  };

  // Toggle virtual mic
  const toggleVirtualMic = async (newEnabled: boolean) => {
    try {
      const response = await ProxyToBackend(
        'POST',
        '/api/audio/virtualmicEnabled',
        JSON.stringify({ enabled: newEnabled })
      );
      const data = JSON.parse(response);
      if (data.success) {
        setEnabled(newEnabled);
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to toggle virtual mic:', error);
      return false;
    }
  };

  // Set device
  const setVirtualMicDevice = async (deviceName: string) => {
    try {
      const response = await ProxyToBackend(
        'POST',
        '/api/audio/virtualmicDevice',
        JSON.stringify({ device: deviceName })
      );
      const data = JSON.parse(response);
      if (data.success) {
        setDevice(deviceName);
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to set virtual mic device:', error);
      return false;
    }
  };

  useEffect(() => {
    loadStatus();
    detectVirtualMic();
  }, []);

  return {
    enabled,
    device,
    detected,
    loading,
    toggleVirtualMic,
    detectVirtualMic,
    setVirtualMicDevice,
  };
}
