import React from 'react';
import { useTranslation } from 'react-i18next';
import { Lock, Key, AlertCircle, Sliders } from 'lucide-react';

const ElevenLabsSettings: React.FC = () => {
  const { t } = useTranslation('settings');

  return (
    <div className="space-y-6">
      {/* Pro Feature Lock */}
      <div className="p-6 bg-gradient-to-br from-purple-50 via-pink-50 to-purple-50 border-2 border-dashed border-purple-300 rounded-2xl flex items-start gap-4">
        <div className="p-3 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl shadow-lg">
          <Lock className="w-7 h-7 text-white" />
        </div>
        <div className="flex-1">
          <h3 className="text-xl font-bold text-purple-900 mb-2">
            Recurso Exclusivo PRO
          </h3>
          <p className="text-sm text-purple-700 leading-relaxed">
            As configurações de voz personalizadas com ElevenLabs estão disponíveis apenas na versão PRO. 
            Na versão FREE, você já possui áudios de alta qualidade pré-gravados para todos os alertas!
          </p>
        </div>
      </div>

      {/* API Key Section - Disabled */}
      <div className="bg-gray-50 rounded-xl p-6 shadow-sm border border-gray-200 opacity-60 pointer-events-none">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Key className="w-5 h-5 text-gray-400" />
            <h3 className="text-lg font-semibold text-gray-600">{t('voice.api_key')}</h3>
          </div>
          <Lock className="w-4 h-4 text-gray-400" />
        </div>
        
        <div className="space-y-3">
          <input
            type="password"
            disabled
            placeholder="sk_••••••••••••••••••••••••••••"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-400 cursor-not-allowed"
          />
          <div className="flex items-center space-x-2 text-sm text-gray-500">
            <AlertCircle className="w-4 h-4" />
            <span>Disponível apenas na versão PRO</span>
          </div>
        </div>
      </div>

      {/* Voice Settings - Disabled */}
      <div className="bg-gray-50 rounded-xl p-6 shadow-sm border border-gray-200 opacity-60 pointer-events-none">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Sliders className="w-5 h-5 text-gray-400" />
            <h3 className="text-lg font-semibold text-gray-600">{t('voice.title')}</h3>
          </div>
          <Lock className="w-4 h-4 text-gray-400" />
        </div>
        
        <div className="space-y-4">
          {/* Stability */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-600">
                {t('voice.stability')}
              </label>
              <span className="text-sm text-gray-500">50%</span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={50}
              disabled
              className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-not-allowed"
            />
            <p className="text-xs text-gray-400 mt-1">
              Maior estabilidade = voz mais consistente
            </p>
          </div>

          {/* Similarity */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-600">
                Similaridade
              </label>
              <span className="text-sm text-gray-500">75%</span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={75}
              disabled
              className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-not-allowed"
            />
            <p className="text-xs text-gray-400 mt-1">
              Maior similaridade = mais fiel à voz original
            </p>
          </div>

          {/* Style */}
          <div>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-gray-600">
                Estilo/Emoção
              </label>
              <span className="text-sm text-gray-500">0%</span>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={0}
              disabled
              className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-not-allowed"
            />
            <p className="text-xs text-gray-400 mt-1">
              Aumenta a expressividade e emoção na voz
            </p>
          </div>

          {/* Speaker Boost */}
          <div className="flex items-center justify-between">
            <div>
              <label className="text-sm font-medium text-gray-600">
                Speaker Boost
              </label>
              <p className="text-xs text-gray-400">
                Melhora a clareza e qualidade da voz
              </p>
            </div>
            <div className="w-11 h-6 bg-gray-300 rounded-full relative">
              <div className="absolute top-[2px] right-[2px] bg-white border-gray-300 border rounded-full h-5 w-5"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ElevenLabsSettings;
