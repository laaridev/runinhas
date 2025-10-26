import React, { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Volume2, RefreshCw, Lock, Settings, Key, AlertCircle, Check, Sliders } from 'lucide-react';
import { voiceAPI } from '@/services/api-wails';
import { useToast } from '@/hooks/useToast';
import { useAppMode } from '@/hooks/useAppMode';
import { useDebouncedCallback } from '@/hooks/useDebounce';

interface VoiceOption {
  voice_id: string;
  name: string;
  preview_url?: string;
}
interface ElevenLabsConfig {
  apiKey: string;
  voiceId: string;
  stability: number;
  similarity: number;
  style: number;
  speakerBoost: boolean;
}

const ElevenLabsSettings: React.FC = () => {
  const { t } = useTranslation('settings');
  const toast = useToast();
  const { appMode } = useAppMode();
  const isPro = appMode.mode === 'pro';
  
  const [config, setConfig] = useState<ElevenLabsConfig>({
    apiKey: '',
    voiceId: 'eVXYtPVYB9wDoz9NVTIy', // Default voice
    stability: 0.5,
    similarity: 0.75,
    style: 0,
    speakerBoost: true
  });
  
  const [voices, setVoices] = useState<VoiceOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [testingVoice, setTestingVoice] = useState(false);
  const [saving, setSaving] = useState(false);

  // Vozes padrão do ElevenLabs
  const defaultVoices: VoiceOption[] = [
    { voice_id: 'EXAVITQu4vr4xnSDxMaL', name: 'Sarah - Feminina suave' },
    { voice_id: 'ErXwobaYiN019PkySvjV', name: 'Antoni - Masculina profunda' },
    { voice_id: 'MF3mGyEYCl7XYWbV9V6O', name: 'Elli - Feminina jovem' },
    { voice_id: 'TxGEqnHWrfWFTfGW9XjX', name: 'Josh - Masculina jovem' },
    { voice_id: 'VR6AewLTigWG4xSOukaG', name: 'Arnold - Masculina grave' },
    { voice_id: 'pNInz6obpgDQGcFmaJgB', name: 'Adam - Masculina neutra' },
    { voice_id: 'Yko7PKHZNXotIFUBG7I9', name: 'Sam - Masculina casual' },
  ];

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      // Carregar configuração do backend usando voiceAPI
      const voiceConfig = await voiceAPI.getConfig();
      setConfig({
        apiKey: voiceConfig.apiKey || '',
        voiceId: voiceConfig.voiceId || 'eVXYtPVYB9wDoz9NVTIy',
        stability: voiceConfig.stability || 0.5,
        similarity: voiceConfig.similarity || 0.75,
        style: voiceConfig.style || 0,
        speakerBoost: voiceConfig.speakerBoost !== undefined ? voiceConfig.speakerBoost : true
      });
      setVoices(defaultVoices);
    } catch (error) {
      console.error('Erro ao carregar configuração:', error);
      setVoices(defaultVoices);
    }
  };

  // Auto-save with debounce
  const debouncedSave = useDebouncedCallback(
    async (newConfig: ElevenLabsConfig) => {
      setSaving(true);
      try {
        await voiceAPI.saveConfig({
          apiKey: newConfig.apiKey,
          voiceId: newConfig.voiceId,
          stability: newConfig.stability,
          similarity: newConfig.similarity,
          style: newConfig.style,
          speakerBoost: newConfig.speakerBoost
        });
        
        toast.success(t('voice.save'));
      } catch (error) {
        toast.error(t('common:toast.config_error'));
      } finally {
        setSaving(false);
      }
    },
    1000 // 1 segundo de debounce
  );

  // Update config and trigger auto-save
  const updateConfig = useCallback((updates: Partial<ElevenLabsConfig>) => {
    const newConfig = { ...config, ...updates };
    setConfig(newConfig);
    debouncedSave(newConfig);
  }, [config, debouncedSave]);

  const testVoice = async () => {
    setTestingVoice(true);
    try {
      const filename = await voiceAPI.testVoice(
        'Teste de voz do runinhas. Sem tilts, só timing!',
        config.voiceId,
        {
          stability: config.stability,
          similarity_boost: config.similarity,
          style: config.style,
          use_speaker_boost: config.speakerBoost
        }
      );
      
      console.log('🎵 Test voice filename:', filename);
      
      // Fetch audio as blob and create object URL
      const { audioAPI } = await import('@/services/api-wails');
      const blob = await audioAPI.getAudioBlob(filename);
      const blobUrl = URL.createObjectURL(blob);
      
      const audio = new Audio(blobUrl);
      
      audio.onloadeddata = () => {
        console.log('✅ Test audio loaded successfully');
      };
      
      audio.onended = () => {
        console.log('✅ Test audio playback ended');
        URL.revokeObjectURL(blobUrl); // Clean up
        setTestingVoice(false);
      };
      
      audio.onerror = (e) => {
        console.error('❌ Failed to play test audio:', e);
        console.error('Audio error details:', audio.error);
        URL.revokeObjectURL(blobUrl); // Clean up
        toast.error('Erro ao reproduzir áudio');
        setTestingVoice(false);
      };
      
      await audio.play();
      console.log('🎵 Test audio.play() called successfully');
      
      toast.success('Teste de voz iniciado!');
    } catch (error) {
      console.error('❌ Test voice error:', error);
      toast.error('Erro ao testar voz. Verifique sua API Key.');
      setTestingVoice(false);
    }
  };

  const fetchVoices = async () => {
    if (!config.apiKey) {
      return;
    }
    
    setLoading(true);
    try {
      const voices = await voiceAPI.getVoices();
      if (voices.length > 0) {
        setVoices(voices);
        toast.success('Vozes carregadas com sucesso!');
      } else {
        setVoices(defaultVoices);
        toast.warning('Nenhuma voz encontrada. Usando vozes padrão.');
      }
    } catch (error) {
      toast.error('Erro ao buscar vozes');
      setVoices(defaultVoices);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Pro Feature Warning */}
      {!isPro && (
        <div className="p-4 bg-gradient-to-r from-purple-50 to-pink-50 border-2 border-purple-200 rounded-xl flex items-start gap-3">
          <Lock className="w-5 h-5 text-purple-600 mt-0.5 flex-shrink-0" />
          <div className="flex-1">
            <h3 className="font-semibold text-purple-900 mb-1">
              {t('common:version.upgrade_message')}
            </h3>
            <p className="text-sm text-purple-700">
              Na versão gratuita, todas as vozes usam áudios padrão pré-gravados.
            </p>
          </div>
        </div>
      )}
      
      {/* Header */}
      <div className="flex items-center space-x-3 mb-6">
        <div className="p-2 bg-gradient-to-br from-purple-500 to-pink-500 rounded-lg">
          <Settings className="w-6 h-6 text-white" />
        </div>
        <div>
          <h2 className="text-2xl font-bold text-gray-800">Configurações de Voz</h2>
          <p className="text-sm text-gray-600">Configure o ElevenLabs para síntese de voz</p>
        </div>
      </div>

      {/* API Key Section */}
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Key className="w-5 h-5 text-gray-600" />
            <h3 className="text-lg font-semibold text-gray-800">{t('voice.api_key')}</h3>
          </div>
          <button
            onClick={() => setShowApiKey(!showApiKey)}
            className="text-sm text-blue-600 hover:text-blue-700"
          >
            {showApiKey ? 'Ocultar' : 'Mostrar'}
          </button>
        </div>
        
        <div className="space-y-3">
          <input
            type={showApiKey ? 'text' : 'password'}
            value={config.apiKey}
            onChange={(e) => updateConfig({ apiKey: e.target.value })}
            placeholder={t('voice.api_key_placeholder')}
            className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
          <div className="flex items-center space-x-2 text-sm text-gray-600">
            <AlertCircle className="w-4 h-4" />
            <span>
              Obtenha sua API Key em{' '}
              <a 
                href="https://elevenlabs.io/api" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline"
              >
                elevenlabs.io/api
              </a>
            </span>
          </div>
        </div>
      </div>

      {/* Voice Selection */}
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Volume2 className="w-5 h-5 text-gray-600" />
            <h3 className="text-lg font-semibold text-gray-800">{t('voice.voice_id')}</h3>
          </div>
          <button
            onClick={fetchVoices}
            disabled={loading}
            className="flex items-center space-x-2 px-3 py-1 text-sm bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors disabled:opacity-50"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            <span>{t('common:buttons.test_audio')}</span>
          </button>
        </div>
        
        <div className="grid grid-cols-1 gap-2">
          {voices.map((voice) => (
            <label
              key={voice.voice_id}
              className={`flex items-center p-3 rounded-lg border cursor-pointer transition-all ${
                config.voiceId === voice.voice_id
                  ? 'border-purple-500 bg-purple-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <input
                type="radio"
                name="voice"
                value={voice.voice_id}
                checked={config.voiceId === voice.voice_id}
                onChange={(e) => updateConfig({ voiceId: e.target.value })}
                className="mr-3 text-purple-600 focus:ring-purple-500"
              />
              <span className="flex-1 text-gray-700">{voice.name}</span>
              {config.voiceId === voice.voice_id && (
                <Check className="w-5 h-5 text-purple-600" />
              )}
            </label>
          ))}
        </div>
      </div>

      {/* Voice Settings */}
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
        <div className="flex items-center space-x-2 mb-4">
          <Sliders className="w-5 h-5 text-gray-600" />
          <h3 className="text-lg font-semibold text-gray-800">{t('voice.title')}</h3>
        </div>
        
        <div className="space-y-4">
          {/* Stability */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-700">
                {t('voice.stability')}
              </label>
              <span className="text-sm text-gray-600">
                {Math.round(config.stability * 100)}%
              </span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={config.stability * 100}
              onChange={(e) => updateConfig({ stability: Number(e.target.value) / 100 })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
            />
            <p className="text-xs text-gray-500 mt-1">
              Maior estabilidade = voz mais consistente
            </p>
          </div>

          {/* Similarity */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-700">
                Similaridade
              </label>
              <span className="text-sm text-gray-600">
                {Math.round(config.similarity * 100)}%
              </span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={config.similarity * 100}
              onChange={(e) => updateConfig({ similarity: Number(e.target.value) / 100 })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
            />
            <p className="text-xs text-gray-500 mt-1">
              Maior similaridade = mais fiel à voz original
            </p>
          </div>

          {/* Style */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-700">
                Estilo/Emoção
              </label>
              <span className="text-sm text-gray-600">
                {Math.round(config.style * 100)}%
              </span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={config.style * 100}
              onChange={(e) => updateConfig({ style: Number(e.target.value) / 100 })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
            />
            <p className="text-xs text-gray-500 mt-1">
              Aumenta a expressividade e emoção na voz
            </p>
          </div>

          {/* Speaker Boost */}
          <div className="flex items-center justify-between">
            <div>
              <label className="text-sm font-medium text-gray-700">
                Speaker Boost
              </label>
              <p className="text-xs text-gray-500">
                Melhora a clareza e qualidade da voz
              </p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={config.speakerBoost}
                onChange={(e) => updateConfig({ speakerBoost: e.target.checked })}
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
            </label>
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex items-center justify-between">
        <button
          onClick={testVoice}
          disabled={testingVoice || loading}
          className="flex items-center justify-center space-x-2 px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:from-blue-600 hover:to-cyan-600 transition-all disabled:opacity-50"
        >
          <Volume2 className={`w-5 h-5 ${testingVoice ? 'animate-pulse' : ''}`} />
          <span>{testingVoice ? 'Testando...' : 'Testar Voz'}</span>
        </button>
        
        {saving && (
          <div className="flex items-center space-x-2 text-sm text-gray-600">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-purple-600"></div>
            <span>Salvando...</span>
          </div>
        )}
      </div>
    </div>
  );
};

export default ElevenLabsSettings;
