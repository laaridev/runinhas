import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Card, CardContent, CardDescription, CardHeader } from '../ui/card';
import { Switch } from '../ui/switch';
import { Slider } from '../ui/slider';
import { Button } from '../ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '../ui/dialog';
import { Textarea } from '../ui/textarea';
import { Label } from '../ui/label';
import { Edit3, Volume2, Mic, Activity, Loader2 } from 'lucide-react';
import { audioAPI, messageAPI, timingAPI } from '@/services/api-wails';
import { useToast } from '@/hooks/useToast';
import { useDebouncedCallback } from '@/hooks/useDebounce';
import { withRetry } from '@/utils/retry';
import { TIMING, RETRY_CONFIG, ANIMATION } from '@/constants/defaults';
import { TOAST_MESSAGES } from '@/constants/messages';

interface EventCardProps {
  event: {
    key: string;
    name: string;
    description: string;
    icon: React.ReactNode;
    min: number;
    max: number;
    step: number;
  };
  enabled: boolean;
  value: number;
  customMessage?: string;
  onToggle: (enabled: boolean) => void;
  onValueChange: (value: number) => void;
  onMessageChange?: (message: string) => void;
  theme: any;
}

export function EventCard({ 
  event, 
  enabled, 
  value, 
  onToggle, 
  onValueChange, 
  customMessage, 
  onMessageChange, 
  theme 
}: EventCardProps) {
  const toast = useToast();
  const [isPlaying, setIsPlaying] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [tempMessage, setTempMessage] = useState(customMessage || `${event.name} em {seconds} segundos`);
  const [audioExists, setAudioExists] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);

  // Update tempMessage when customMessage changes
  React.useEffect(() => {
    if (customMessage) {
      setTempMessage(customMessage);
    }
  }, [customMessage]);

  // Check if audio exists when component mounts or message changes
  React.useEffect(() => {
    const checkAudio = async () => {
      const exists = await audioAPI.check(event.key);
      setAudioExists(exists);
    };
    checkAudio();
  }, [event.key, customMessage]);

  const handleToggle = () => onToggle(!enabled);

  const handlePreview = async () => {
    if (isPlaying) return;
    
    setIsPlaying(true);
    try {
      const filename = `${event.key}_warning.mp3`;
      
      // Check if audio exists with retry
      const exists = await withRetry(
        () => audioAPI.check(event.key),
        { maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO, delayMs: RETRY_CONFIG.DELAY_MS }
      );
      
      if (!exists) {
        // If doesn't exist, generate once with retry
        toast.info(TOAST_MESSAGES.GENERATING_FIRST_TIME);
        await withRetry(
          () => audioAPI.generate(event.key),
          {
            maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO,
            delayMs: RETRY_CONFIG.BACKOFF_DELAY,
            onRetry: (attempt) => {
              toast.warning(TOAST_MESSAGES.RETRY_GENERATING(attempt, RETRY_CONFIG.MAX_RETRIES_AUDIO));
            }
          }
        );
      }
      
      // Play the existing audio with cache bypass
      const filenameWithTimestamp = `${filename}?t=${Date.now()}`;
      const blob = await audioAPI.getAudioBlob(filenameWithTimestamp);
      const blobUrl = URL.createObjectURL(blob);
      const audio = new Audio(blobUrl);
      
      audio.onended = () => {
        URL.revokeObjectURL(blobUrl);
        setIsPlaying(false);
      };
      
      audio.onerror = () => {
        URL.revokeObjectURL(blobUrl);
        setIsPlaying(false);
        toast.error(TOAST_MESSAGES.AUDIO_PLAY_ERROR);
      };
      
      await audio.play();
    } catch (error: any) {
      console.error('Preview error:', error);
      toast.error(`${TOAST_MESSAGES.AUDIO_PREVIEW_ERROR}: ${error.message || 'Tente novamente'}`);
      setIsPlaying(false);
    }
  };

  const handleGenerateAudio = async () => {
    setIsGenerating(true);
    try {
      await withRetry(
        () => audioAPI.generate(event.key),
        {
          maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO,
          delayMs: RETRY_CONFIG.BACKOFF_DELAY,
          onRetry: (attempt) => {
            toast.warning(TOAST_MESSAGES.RETRY_GENERATING(attempt, RETRY_CONFIG.MAX_RETRIES_AUDIO));
          }
        }
      );
      setAudioExists(true);
      toast.success(TOAST_MESSAGES.AUDIO_GENERATED);
    } catch (error: any) {
      console.error('Failed to generate audio:', error);
      toast.error(`${TOAST_MESSAGES.AUDIO_ERROR}: ${error.message || ''}`);
    } finally {
      setIsGenerating(false);
    }
  };

  // Debounced audio generation - only generate after user stops moving slider
  const debouncedGenerateAudio = useDebouncedCallback(async (val: number) => {
    try {
      setIsGenerating(true);
      
      // Save value to backend with retry
      await withRetry(
        () => timingAPI.setValue(event.key, "warning_seconds", val),
        {
          maxRetries: RETRY_CONFIG.MAX_RETRIES_SAVE,
          delayMs: RETRY_CONFIG.DELAY_MS,
          onRetry: (attempt) => {
            console.log(`Retry ${attempt} - saving value for ${event.key}`);
          }
        }
      );
      
      // Small delay to ensure backend has processed the update
      await new Promise(resolve => setTimeout(resolve, TIMING.API_RETRY_DELAY));
      
      // Generate audio with retry
      const response = await withRetry(
        () => audioAPI.generate(event.key),
        {
          maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO,
          delayMs: RETRY_CONFIG.BACKOFF_DELAY,
          onRetry: (attempt) => {
            toast.warning(TOAST_MESSAGES.RETRY_GENERATING(attempt, RETRY_CONFIG.MAX_RETRIES_AUDIO));
          }
        }
      );
      
      const data = JSON.parse(response);
      
      if (data.filename) {
        // Add timestamp to bypass any cache
        const filenameWithTimestamp = `${data.filename}?t=${Date.now()}`;
        
        // Play the generated audio immediately
        const blob = await audioAPI.getAudioBlob(filenameWithTimestamp);
        const blobUrl = URL.createObjectURL(blob);
        const audio = new Audio(blobUrl);
        
        audio.onended = () => {
          URL.revokeObjectURL(blobUrl);
        };
        
        await audio.play();
      }
      
      setAudioExists(true);
      toast.success(TOAST_MESSAGES.AUDIO_GENERATED);
    } catch (error) {
      console.error('Failed to generate audio:', error);
      toast.error(TOAST_MESSAGES.AUDIO_ERROR);
    } finally {
      setIsGenerating(false);
    }
  }, TIMING.DEBOUNCE_DELAY);

  const handleValueChange = (newValue: number[]) => {
    const val = newValue[0];
    onValueChange(val);
    
    // Debounced generation
    debouncedGenerateAudio(val);
  };

  const handleSaveMessage = async () => {
    // Validação: verificar se {seconds} está presente na mensagem
    if (!tempMessage.includes('{seconds}')) {
      toast.warning(TOAST_MESSAGES.MESSAGE_MISSING_PLACEHOLDER);
      return;
    }

    // Validação: mensagem não pode estar vazia
    if (!tempMessage.trim()) {
      toast.error(TOAST_MESSAGES.MESSAGE_EMPTY);
      return;
    }

    if (onMessageChange && tempMessage !== customMessage) {
      onMessageChange(tempMessage);
      
      // Save to backend
      try {
        setIsGenerating(true);
        
        // Save message with retry
        await withRetry(
          () => messageAPI.set(event.key, tempMessage),
          {
            maxRetries: RETRY_CONFIG.MAX_RETRIES_SAVE,
            delayMs: RETRY_CONFIG.DELAY_MS,
            onRetry: (attempt) => {
              toast.info(TOAST_MESSAGES.RETRY_SAVING(attempt, RETRY_CONFIG.MAX_RETRIES_SAVE));
            }
          }
        );
        
        // Small delay to ensure backend has processed the update
        await new Promise(resolve => setTimeout(resolve, TIMING.API_RETRY_DELAY));
        
        // Regenerate audio with new message (with retry)
        const response = await withRetry(
          () => audioAPI.generate(event.key),
          {
            maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO,
            delayMs: RETRY_CONFIG.BACKOFF_DELAY,
            onRetry: (attempt) => {
              toast.warning(TOAST_MESSAGES.RETRY_GENERATING(attempt, RETRY_CONFIG.MAX_RETRIES_AUDIO));
            }
          }
        );
        
        const data = JSON.parse(response);
        
        // Play the new audio
        if (data.filename) {
          // Add timestamp to bypass any cache
          const filenameWithTimestamp = `${data.filename}?t=${Date.now()}`;
          const blob = await audioAPI.getAudioBlob(filenameWithTimestamp);
          const blobUrl = URL.createObjectURL(blob);
          const audio = new Audio(blobUrl);
          audio.onended = () => URL.revokeObjectURL(blobUrl);
          await audio.play();
        }
        
        setAudioExists(true);
        
        // Success toast
        toast.success(TOAST_MESSAGES.MESSAGE_SAVED);
        
        // Close dialog after successful generation
        setShowEditDialog(false);
      } catch (error: any) {
        console.error('Failed to save message:', error);
        toast.error(`${TOAST_MESSAGES.SAVE_ERROR}: ${error.message || 'Tente novamente'}`);
      } finally {
        setIsGenerating(false);
      }
    } else {
      setShowEditDialog(false);
    }
  };

  const handleResetMessage = async () => {
    const defaultMessage = `${event.name} em {seconds} segundos`;
    setTempMessage(defaultMessage);
    if (onMessageChange) {
      onMessageChange(defaultMessage);
    }
    
    // Save to backend
    try {
      setIsGenerating(true);
      
      // Save with retry
      await withRetry(
        () => messageAPI.set(event.key, defaultMessage),
        {
          maxRetries: RETRY_CONFIG.MAX_RETRIES_SAVE,
          delayMs: RETRY_CONFIG.DELAY_MS,
          onRetry: (attempt) => {
            toast.info(TOAST_MESSAGES.RESTORING_DEFAULT(attempt, RETRY_CONFIG.MAX_RETRIES_SAVE));
          }
        }
      );
      
      // Small delay to ensure backend has processed the update
      await new Promise(resolve => setTimeout(resolve, TIMING.API_RETRY_DELAY));
      
      // Regenerate audio with default message
      const response = await withRetry(
        () => audioAPI.generate(event.key),
        {
          maxRetries: RETRY_CONFIG.MAX_RETRIES_AUDIO,
          delayMs: RETRY_CONFIG.BACKOFF_DELAY,
          onRetry: (attempt) => {
            toast.warning(TOAST_MESSAGES.GENERATING_DEFAULT(attempt, RETRY_CONFIG.MAX_RETRIES_AUDIO));
          }
        }
      );
      
      const data = JSON.parse(response);
      
      // Play the new audio
      if (data.filename) {
        const blob = await audioAPI.getAudioBlob(data.filename);
        const blobUrl = URL.createObjectURL(blob);
        const audio = new Audio(blobUrl);
        audio.onended = () => URL.revokeObjectURL(blobUrl);
        await audio.play();
      }
      
      setAudioExists(true);
      toast.success(TOAST_MESSAGES.MESSAGE_RESTORED);
    } catch (error: any) {
      console.error('Failed to reset message:', error);
      toast.error(`${TOAST_MESSAGES.RESTORE_ERROR}: ${error.message || 'Tente novamente'}`);
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <>
      <Card className={`
        group relative overflow-hidden
        ${enabled 
          ? 'border-2 shadow-lg shadow-purple-500/20 bg-gradient-to-br from-white via-purple-50/30 to-blue-50/30' 
          : 'border border-gray-200/60 bg-white/80'
        }
        backdrop-blur-xl 
        transition-all duration-500 ease-out
        hover:shadow-2xl hover:shadow-purple-500/30 hover:-translate-y-1
        rounded-2xl
      `}>
        {/* Gradient overlay for enabled state */}
        {enabled && (
          <div className="absolute inset-0 bg-gradient-to-br from-purple-500/5 via-transparent to-blue-500/5 pointer-events-none" />
        )}
        
        <CardHeader className="pb-3 relative">
          <div className="flex items-start justify-between gap-3">
            <div className="flex items-start gap-3 flex-1 min-w-0">
              {/* Icon with glow effect */}
              <div className={`
                relative p-2.5 rounded-xl 
                ${enabled ? theme.iconBg + ' ring-2 ring-purple-200/50' : theme.iconBg}
                transition-all duration-500
                group-hover:scale-110 group-hover:rotate-3
              `}>
                {enabled && (
                  <div className={`absolute inset-0 rounded-xl ${theme.iconBg} blur-lg opacity-50`} />
                )}
                {React.cloneElement(event.icon as React.ReactElement, {
                  className: `w-5 h-5 ${theme.iconColor} transition-all duration-500 relative z-10`
                })}
              </div>
              
              {/* Title and description */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <h3 className="text-base font-bold text-gray-900 tracking-tight leading-tight">{event.name}</h3>
                  {enabled && (
                    <span className="px-1.5 py-0.5 text-[10px] font-semibold text-purple-700 bg-purple-100 rounded-full animate-pulse leading-none">
                      Ativo
                    </span>
                  )}
                </div>
                <CardDescription className="text-xs text-gray-600 leading-snug">
                  {event.description}
                </CardDescription>
              </div>
            </div>
            
            {/* Toggle switch */}
            <Switch
              checked={enabled}
              onCheckedChange={handleToggle}
              disabled={isGenerating}
              className={`
                data-[state=checked]:bg-gradient-to-r ${theme.gradient}
                transition-all duration-300
                shadow-md data-[state=checked]:shadow-purple-300/50
                flex-shrink-0
                disabled:opacity-50 disabled:cursor-not-allowed
              `}
            />
          </div>
        </CardHeader>

        <CardContent className="space-y-4 relative">
          {/* Slider section */}
          <div className="space-y-2.5">
            <div className="flex items-center justify-between">
              <Label className={`text-xs font-semibold transition-all duration-300 ${isGenerating ? 'text-gray-400' : 'text-gray-600'}`}>
                Antecedência do aviso
              </Label>
              <span className={`
                text-sm font-bold
                transition-all duration-500
                px-2.5 py-1 rounded-lg
                min-w-[45px] text-center
                ${isGenerating 
                  ? 'bg-gray-100 text-gray-500 ring-1 ring-gray-200' 
                  : enabled 
                    ? `${theme.iconColor} ${theme.iconBg} ring-1 ring-purple-200`
                    : 'bg-gray-50 text-gray-500'
                }
              `}>
                {value}s
              </span>
            </div>
            <Slider
              value={[value]}
              onValueChange={handleValueChange}
              min={event.min}
              max={event.max}
              step={event.step}
              disabled={!enabled || isGenerating}
              className={`
                cursor-pointer 
                [&>span:first-child]:${theme.sliderTrack} 
                [&>span>span]:${theme.sliderThumb} 
                [&>span:last-child]:border-2 
                [&>span:last-child]:${theme.sliderThumb.replace('bg-', 'border-')}
                [&>span:last-child]:shadow-lg
                [&>span:last-child]:shadow-purple-300/50
                transition-all duration-300
                disabled:opacity-50 disabled:cursor-not-allowed
              `}
            />
          </div>

          {/* Action buttons or Generation indicator with smooth transition */}
          <AnimatePresence mode="wait">
            {isGenerating ? (
              <motion.div
                key="generating"
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: ANIMATION.BUTTON_TRANSITION / 1000 }}
                className="flex items-center justify-center gap-2 p-2.5 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 border border-purple-200/30"
              >
                <Loader2 className="w-4 h-4 animate-spin text-purple-600" />
                <span className="text-xs font-semibold text-purple-700">Gerando áudio com IA...</span>
              </motion.div>
            ) : (
              <motion.div
                key="buttons"
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: ANIMATION.BUTTON_TRANSITION / 1000 }}
                className="flex gap-2"
              >
              <Button
                variant="outline"
                size="sm"
                onClick={handlePreview}
                disabled={!enabled || isPlaying || !audioExists}
                className={`
                  flex-1 h-9
                  bg-gradient-to-r from-white to-gray-50
                  border-2 ${enabled ? 'border-purple-200/60' : 'border-gray-200'}
                  hover:border-purple-400 hover:bg-purple-50/50
                  hover:shadow-md hover:shadow-purple-200/50
                  transition-all duration-300
                  font-medium text-sm
                  ${isPlaying ? 'ring-2 ring-purple-300 ring-offset-1' : ''}
                `}
              >
                {isPlaying ? (
                  <>
                    <Activity className="w-3.5 h-3.5 mr-1.5 animate-pulse text-purple-600" />
                    <span className="text-purple-700">Reproduzindo...</span>
                  </>
                ) : (
                  <>
                    <Volume2 className="w-3.5 h-3.5 mr-1.5" />
                    Testar Voz
                  </>
                )}
              </Button>

              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowEditDialog(true)}
                disabled={!enabled}
                className={`
                  flex-1 h-9
                  bg-gradient-to-r from-white to-gray-50
                  border-2 ${enabled ? 'border-blue-200/60' : 'border-gray-200'}
                  hover:border-blue-400 hover:bg-blue-50/50
                  hover:shadow-md hover:shadow-blue-200/50
                  transition-all duration-300
                  font-medium text-sm
                `}
              >
                <Edit3 className="w-3.5 h-3.5 mr-1.5" />
                Editar
              </Button>

              {!audioExists && enabled && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleGenerateAudio}
                  disabled={isGenerating}
                  className={`
                    h-9 px-3
                    bg-gradient-to-r from-white to-gray-50
                    border-2 border-green-200/60
                    hover:border-green-400 hover:bg-green-50/50
                    hover:shadow-md hover:shadow-green-200/50
                    transition-all duration-300
                  `}
                >
                  {isGenerating ? (
                    <Loader2 className="w-3.5 h-3.5 animate-spin text-green-600" />
                  ) : (
                    <Mic className="w-3.5 h-3.5 text-green-600" />
                  )}
                </Button>
              )}
              </motion.div>
            )}
          </AnimatePresence>
        </CardContent>
      </Card>

      {/* Edit Message Dialog - Theme Adaptive */}
      <Dialog open={showEditDialog} onOpenChange={setShowEditDialog}>
        <DialogContent className="bg-white border-2 border-gray-200 rounded-2xl shadow-2xl max-w-2xl">
          <DialogHeader className="space-y-4">
            <DialogTitle className="text-2xl font-bold flex items-center gap-3">
              <div className={`p-2.5 rounded-xl ${theme.iconBg} border-2 ${theme.navbarBorder} shadow-md`}>
                <Edit3 className={`w-5 h-5 ${theme.iconColor}`} />
              </div>
              <span className={`${theme.titleText}`}>
                Editar Mensagem
              </span>
            </DialogTitle>
            <DialogDescription className="text-gray-600 text-sm leading-relaxed">
              Personalize a mensagem de voz para <span className={`font-bold ${theme.iconColor}`}>{event.name}</span>.
              Use <code className={`${theme.iconBg} px-2 py-0.5 rounded-md ${theme.iconColor} font-mono text-xs font-bold`}>{'{seconds}'}</code> para incluir o tempo.
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4 py-4">
            <div className="space-y-3">
              <Label htmlFor="message" className="text-sm font-bold text-gray-700 flex items-center gap-2">
                <span className={`w-1.5 h-1.5 rounded-full ${theme.iconColor.replace('text-', 'bg-')}`}></span>
                Mensagem Personalizada
              </Label>
              <Textarea
                id="message"
                value={tempMessage}
                onChange={(e) => setTempMessage(e.target.value)}
                placeholder={`${event.name} em {seconds} segundos`}
                disabled={isGenerating}
                className={`
                  text-gray-900 
                  min-h-[100px] rounded-xl
                  transition-all duration-300
                  placeholder:text-gray-400
                  text-sm leading-relaxed
                  resize-none
                  border-2 focus:ring-2 focus:ring-offset-0
                  disabled:opacity-60 disabled:cursor-not-allowed disabled:bg-gray-50
                `}
                style={{
                  borderColor: isGenerating ? '#E5E7EB' : (theme.iconColor.includes('blue') ? '#DBEAFE' : '#FCE7F3'),
                  outlineColor: theme.iconColor.includes('blue') ? '#3B82F6' : '#EC4899',
                  backgroundColor: isGenerating ? '#F9FAFB' : '#FFFFFF'
                }}
              />
              <div className={`flex items-start gap-2.5 p-3 rounded-xl ${theme.iconBg} border ${theme.navbarBorder}`}>
                <div className={`w-1.5 h-1.5 rounded-full ${theme.iconColor.replace('text-', 'bg-')} mt-1.5 flex-shrink-0`}></div>
                <p className={`text-xs ${theme.iconColor} leading-relaxed`}>
                  <span className="font-bold">Dica:</span> "{event.name} em {'{seconds}'} segundos, prepara!" ou "Faltam {'{seconds}'} para {event.name}!"
                </p>
              </div>
            </div>
          </div>
          
          <DialogFooter className="flex flex-col sm:flex-row gap-3 sm:justify-between pt-4 border-t border-gray-200">
            <Button
              variant="ghost"
              onClick={handleResetMessage}
              disabled={isGenerating}
              className="text-gray-600 hover:text-gray-900 hover:bg-gray-100 transition-all duration-300 font-semibold text-sm disabled:opacity-50"
            >
              🔄 Restaurar Padrão
            </Button>
            <div className="flex gap-2.5">
              <Button 
                variant="outline" 
                onClick={() => setShowEditDialog(false)}
                disabled={isGenerating}
                className="border-2 border-gray-300 hover:border-gray-400 hover:bg-gray-100 transition-all duration-300 font-semibold text-sm disabled:opacity-50"
              >
                Cancelar
              </Button>
              <Button 
                onClick={handleSaveMessage}
                disabled={isGenerating}
                style={{
                  background: theme.iconColor.includes('blue') 
                    ? 'linear-gradient(to right, #3B82F6, #8B5CF6)' 
                    : 'linear-gradient(to right, #EC4899, #8B5CF6)'
                }}
                className={`
                  shadow-md hover:shadow-lg
                  transition-all duration-300
                  font-bold text-sm
                  hover:scale-105 hover:brightness-110
                  text-white
                  disabled:opacity-90 disabled:cursor-not-allowed
                  min-w-[140px]
                  border-0
                `}
              >
                {isGenerating ? (
                  <div className="flex items-center gap-2">
                    <Loader2 className="w-4 h-4 animate-spin" />
                    <span>Gerando...</span>
                  </div>
                ) : (
                  <div className="flex items-center gap-2">
                    <span>✨</span>
                    <span>Salvar e Gerar</span>
                  </div>
                )}
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
