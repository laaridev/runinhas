import { useState } from 'react';
import { Mic, Headphones, Volume2, AlertCircle, CheckCircle2 } from 'lucide-react';
import { useVirtualMic } from '@/hooks/useVirtualMic';
import { useToast } from '@/hooks/useToast';
import { Switch } from './ui/switch';

interface VirtualMicSettingsProps {
  theme: any;
}

export function VirtualMicSettings({ theme }: VirtualMicSettingsProps) {
  const { enabled, device, detected, loading, toggleVirtualMic, detectVirtualMic } = useVirtualMic();
  const toast = useToast();
  const [isToggling, setIsToggling] = useState(false);

  const handleToggle = async (checked: boolean) => {
    if (!detected && checked) {
      toast.warning('Microfone virtual não detectado no sistema');
      return;
    }

    setIsToggling(true);
    const success = await toggleVirtualMic(checked);
    setIsToggling(false);

    if (success) {
      if (checked) {
        toast.success('Saída para microfone virtual ativada! Seu time ouvirá os avisos.');
      } else {
        toast.info('Saída para microfone virtual desativada');
      }
    } else {
      toast.error('Erro ao alterar configuração');
    }
  };

  const handleDetect = async () => {
    const found = await detectVirtualMic();
    if (found) {
      toast.success('Microfone virtual detectado!');
    } else {
      toast.warning('Nenhum microfone virtual encontrado. Configure um no sistema primeiro.');
    }
  };

  if (loading) {
    return (
      <div className="p-4 bg-gray-50 dark:bg-gray-900/50 rounded-lg border border-gray-200 dark:border-gray-800">
        <p className="text-sm text-gray-500">Carregando...</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center gap-2">
        <Mic className={`w-5 h-5 ${theme.iconColor}`} />
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          Saída para Microfone Virtual
        </h3>
      </div>

      {/* Description */}
      <p className="text-sm text-gray-600 dark:text-gray-400">
        Reproduza os avisos simultaneamente no seu microfone virtual. Assim seu time poderá ouvir os avisos via Discord, TeamSpeak, etc.
      </p>

      {/* Status Card */}
      <div className={`p-4 rounded-lg border ${
        detected 
          ? 'bg-emerald-50 dark:bg-emerald-900/20 border-emerald-200 dark:border-emerald-800'
          : 'bg-amber-50 dark:bg-amber-900/20 border-amber-200 dark:border-amber-800'
      }`}>
        <div className="flex items-start gap-3">
          {detected ? (
            <CheckCircle2 className="w-5 h-5 text-emerald-600 dark:text-emerald-400 mt-0.5" />
          ) : (
            <AlertCircle className="w-5 h-5 text-amber-600 dark:text-amber-400 mt-0.5" />
          )}
          <div className="flex-1">
            <p className={`text-sm font-medium ${
              detected 
                ? 'text-emerald-800 dark:text-emerald-200'
                : 'text-amber-800 dark:text-amber-200'
            }`}>
              {detected ? 'Microfone virtual detectado' : 'Microfone virtual não encontrado'}
            </p>
            {detected && (
              <p className="text-xs text-emerald-600 dark:text-emerald-400 mt-1">
                Dispositivo: <code className="bg-emerald-100 dark:bg-emerald-900/50 px-1 py-0.5 rounded">{device}</code>
              </p>
            )}
            {!detected && (
              <p className="text-xs text-amber-600 dark:text-amber-400 mt-1">
                Configure um microfone virtual no sistema (ex: virt_mic no PulseAudio)
              </p>
            )}
          </div>
          <button
            onClick={handleDetect}
            className="text-xs px-2 py-1 rounded bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            Detectar
          </button>
        </div>
      </div>

      {/* Toggle Control */}
      <div className={`p-4 rounded-lg border transition-all ${
        enabled
          ? 'bg-gradient-to-r from-emerald-50 to-green-50 dark:from-emerald-900/20 dark:to-green-900/20 border-emerald-200 dark:border-emerald-800'
          : 'bg-gray-50 dark:bg-gray-900/50 border-gray-200 dark:border-gray-800'
      }`}>
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-1">
              <Volume2 className={`w-4 h-4 ${enabled ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-500'}`} />
              <p className="text-sm font-medium text-gray-900 dark:text-white">
                Ativar Saída Dupla
              </p>
            </div>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              Áudios tocarão nos fones {enabled && '+ microfone virtual'}
            </p>
          </div>
          <Switch
            checked={enabled}
            onCheckedChange={handleToggle}
            disabled={isToggling || !detected}
            className={theme.gradient}
          />
        </div>
      </div>

      {/* Audio Flow Visualization */}
      {enabled && (
        <div className="p-4 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <p className="text-xs font-medium text-blue-900 dark:text-blue-200 mb-3">
            Fluxo de Áudio Ativo:
          </p>
          <div className="flex items-center gap-2 text-xs">
            <div className="flex items-center gap-1.5 px-2 py-1 bg-white dark:bg-gray-800 rounded border border-blue-200 dark:border-blue-700">
              <Volume2 className="w-3 h-3 text-blue-600 dark:text-blue-400" />
              <span className="text-gray-700 dark:text-gray-300">Evento</span>
            </div>
            <span className="text-blue-500">→</span>
            <div className="flex items-center gap-1.5 px-2 py-1 bg-white dark:bg-gray-800 rounded border border-emerald-200 dark:border-emerald-700">
              <Headphones className="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
              <span className="text-gray-700 dark:text-gray-300">Fones</span>
            </div>
            <span className="text-blue-500">+</span>
            <div className="flex items-center gap-1.5 px-2 py-1 bg-white dark:bg-gray-800 rounded border border-purple-200 dark:border-purple-700">
              <Mic className="w-3 h-3 text-purple-600 dark:text-purple-400" />
              <span className="text-gray-700 dark:text-gray-300">Mic Virtual</span>
            </div>
          </div>
        </div>
      )}

      {/* Help Text */}
      <div className="p-3 bg-gray-50 dark:bg-gray-900/50 rounded-lg border border-gray-200 dark:border-gray-800">
        <p className="text-xs text-gray-600 dark:text-gray-400">
          <strong>Dica:</strong> No Linux, use PulseAudio para criar um microfone virtual (virt_mic). 
          No Windows, use VB-Audio Cable ou software similar.
        </p>
      </div>
    </div>
  );
}
