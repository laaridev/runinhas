import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader } from '../ui/card';
import { Switch } from '../ui/switch';
import { Slider } from '../ui/slider';
import { Button } from '../ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '../ui/dialog';
import { Textarea } from '../ui/textarea';
import { Label } from '../ui/label';
import { Edit3, Volume2, Mic, Activity, Loader2 } from 'lucide-react';
import { audioAPI, messageAPI } from '@/services/api-wails';

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
  const [isPlaying, setIsPlaying] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [tempMessage, setTempMessage] = useState(customMessage || `${event.name} em {seconds} segundos`);
  const [audioExists, setAudioExists] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [pendingValue, setPendingValue] = useState<number | null>(null);

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

  // Auto-generate audio when value changes (if audio already exists)
  React.useEffect(() => {
    if (pendingValue !== null && enabled && audioExists) {
      const timer = setTimeout(async () => {
        setIsGenerating(true);
        try {
          await audioAPI.generate(event.key);
          setAudioExists(true);
        } catch (error) {
          console.error('Failed to generate audio:', error);
        } finally {
          setIsGenerating(false);
          setPendingValue(null);
        }
      }, 500); // Debounce de 500ms para não gerar a cada movimento do slider
      
      return () => clearTimeout(timer);
    }
  }, [pendingValue, event.key, enabled, audioExists]);

  const handleToggle = () => onToggle(!enabled);

  const handlePreview = async () => {
    if (isPlaying) return;
    
    setIsPlaying(true);
    try {
      await audioAPI.preview(event.key);
    } catch (error) {
      console.error('Preview failed:', error);
    }
    setTimeout(() => setIsPlaying(false), 2000);
  };

  const handleGenerateAudio = async () => {
    setIsGenerating(true);
    try {
      await audioAPI.generate(event.key);
      setAudioExists(true);
    } catch (error) {
      console.error('Failed to generate audio:', error);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleValueChange = (newValue: number[]) => {
    const val = newValue[0];
    onValueChange(val);
    setPendingValue(val);
  };

  const handleSaveMessage = async () => {
    if (onMessageChange && tempMessage !== customMessage) {
      onMessageChange(tempMessage);
      
      // Save to backend
      try {
        await messageAPI.set(event.key, tempMessage);
        
        // Regenerate audio with new message
        setIsGenerating(true);
        await audioAPI.generate(event.key);
        setAudioExists(true);
      } catch (error) {
        console.error('Failed to save message:', error);
      } finally {
        setIsGenerating(false);
      }
    }
    setShowEditDialog(false);
  };

  const handleResetMessage = async () => {
    const defaultMessage = `${event.name} em {seconds} segundos`;
    setTempMessage(defaultMessage);
    if (onMessageChange) {
      onMessageChange(defaultMessage);
    }
    
    // Save to backend
    try {
      await messageAPI.set(event.key, defaultMessage);
      
      // Regenerate audio with default message
      setIsGenerating(true);
      await audioAPI.generate(event.key);
      setAudioExists(true);
    } catch (error) {
      console.error('Failed to reset message:', error);
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <>
      <Card className={`${enabled ? 'border-2 border-purple-400/50' : 'border-gray-200'} bg-white/95 backdrop-blur-xl transition-all duration-300 hover:shadow-2xl hover:shadow-purple-500/20`}>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className={`p-2.5 rounded-xl ${theme.iconBg} transition-all duration-500`}>
                {React.cloneElement(event.icon as React.ReactElement, {
                  className: `w-6 h-6 ${theme.iconColor} transition-colors duration-500`
                })}
              </div>
              <div>
                <h3 className="text-base font-semibold text-gray-900">{event.name}</h3>
                <CardDescription className="text-xs text-gray-500">
                  {event.description}
                </CardDescription>
              </div>
            </div>
            <Switch
              checked={enabled}
              onCheckedChange={handleToggle}
              className={`data-[state=checked]:bg-gradient-to-r ${theme.gradient}`}
            />
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-medium text-gray-600">
                Avisar com {value} segundos de antecedência
              </Label>
              <span className={`text-sm font-bold ${theme.iconColor} transition-colors duration-500`}>
                {value}s
              </span>
            </div>
            <Slider
              value={[value]}
              onValueChange={handleValueChange}
              min={event.min}
              max={event.max}
              step={event.step}
              disabled={!enabled}
              className={`cursor-pointer [&>span:first-child]:${theme.sliderTrack} [&>span>span]:${theme.sliderThumb} [&>span:last-child]:border-2 [&>span:last-child]:${theme.sliderThumb.replace('bg-', 'border-')}`}
            />
          </div>

          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handlePreview}
              disabled={!enabled || isPlaying || !audioExists}
              className="flex-1 bg-white border-gray-300 hover:bg-gray-50 hover:border-purple-500/50 transition-all duration-300"
            >
              {isPlaying ? (
                <>
                  <Activity className="w-4 h-4 mr-2 animate-pulse" />
                  Reproduzindo...
                </>
              ) : (
                <>
                  <Volume2 className="w-4 h-4 mr-2" />
                  Testar Voz
                </>
              )}
            </Button>

            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowEditDialog(true)}
              disabled={!enabled}
              className="flex-1 bg-white border-gray-300 hover:bg-gray-50 hover:border-blue-500/50 transition-all duration-300"
            >
              <Edit3 className="w-4 h-4 mr-2" />
              Editar
            </Button>

            {!audioExists && enabled && (
              <Button
                variant="outline"
                size="sm"
                onClick={handleGenerateAudio}
                disabled={isGenerating}
                className="bg-white border-gray-300 hover:bg-gray-50 hover:border-green-500/50 transition-all duration-300"
              >
                {isGenerating ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <Mic className="w-4 h-4" />
                )}
              </Button>
            )}
          </div>

          {isGenerating && (
            <div className="flex items-center gap-2 text-sm text-blue-400 animate-pulse">
              <Loader2 className="w-4 h-4 animate-spin" />
              Gerando áudio...
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Message Dialog */}
      <Dialog open={showEditDialog} onOpenChange={setShowEditDialog}>
        <DialogContent className="bg-white border-2 border-gray-200">
          <DialogHeader>
            <DialogTitle className="text-xl font-bold text-gray-900 flex items-center gap-2">
              <Edit3 className="w-5 h-5 text-purple-500" />
              Editar Mensagem - {event.name}
            </DialogTitle>
            <DialogDescription className="text-gray-600">
              Personalize a mensagem de voz. Use <code className="bg-gray-100 px-1.5 py-0.5 rounded text-purple-600">{'{seconds}'}</code> para o tempo.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="message" className="text-gray-700">Mensagem</Label>
              <Textarea
                id="message"
                value={tempMessage}
                onChange={(e) => setTempMessage(e.target.value)}
                placeholder={`${event.name} em {seconds} segundos`}
                className="bg-white border-gray-300 text-gray-900 min-h-[100px] focus:border-purple-500"
              />
              <p className="text-xs text-gray-500">
                Exemplo: "{event.name} em {'{seconds}'} segundos, prepara!"
              </p>
            </div>
          </div>
          <DialogFooter className="flex justify-between">
            <Button
              variant="ghost"
              onClick={handleResetMessage}
              className="text-gray-600 hover:text-gray-900"
            >
              Restaurar Padrão
            </Button>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setShowEditDialog(false)}>
                Cancelar
              </Button>
              <Button onClick={handleSaveMessage} className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600">
                Salvar
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
