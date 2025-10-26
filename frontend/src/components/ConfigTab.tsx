import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Switch } from './ui/switch';
import { Label } from './ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { ScrollArea } from './ui/scroll-area';
import { 
  Settings, 
  History, 
  Download, 
  Upload, 
  Trash2, 
  Volume2,
  Bell,
  MessageSquare,
  Globe,
  Clock,
  CheckCircle,
  AlertCircle,
  Info,
  Lock
} from 'lucide-react';
import { useToast } from '../hooks/useToast';

interface EventHistoryItem {
  id: string;
  type: 'rune' | 'timing' | 'hero' | 'map';
  title: string;
  message: string;
  timestamp: Date;
  status: 'success' | 'warning' | 'info';
}

interface ConfigTabProps {
  theme: any;
}

export function ConfigTab({ theme }: ConfigTabProps) {
  const toast = useToast();
  const [voiceEnabled, setVoiceEnabled] = useState(true);
  
  // Mock histórico de eventos
  const [eventHistory] = useState<EventHistoryItem[]>([
    {
      id: '1',
      type: 'rune',
      title: 'Runa de Recompensa',
      message: 'Runa em 30 segundos',
      timestamp: new Date(Date.now() - 60000),
      status: 'success'
    },
    {
      id: '2',
      type: 'timing',
      title: 'Stack Timing',
      message: 'Hora de stackar camps',
      timestamp: new Date(Date.now() - 120000),
      status: 'warning'
    },
    {
      id: '3',
      type: 'hero',
      title: 'Ultimate Pronto',
      message: 'Sua ultimate está pronta',
      timestamp: new Date(Date.now() - 180000),
      status: 'info'
    }
  ]);

  const handleExportConfig = () => {
    // Implementar export real
    const config = {
      handlers: { voiceEnabled },
      // ... outras configs
    };
    const blob = new Blob([JSON.stringify(config, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'runinhas-config.json';
    a.click();
    
    toast.success("Arquivo salvo como runinhas-config.json", "Configurações exportadas");
  };

  const handleImportConfig = () => {
    // Implementar import real
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = (e: any) => {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          try {
            JSON.parse(e.target?.result as string);
            // TODO: Aplicar configurações importadas
            toast.success("Suas configurações foram restauradas com sucesso", "Configurações importadas");
          } catch (err) {
            toast.error("Arquivo inválido ou corrompido", "Erro ao importar");
          }
        };
        reader.readAsText(file);
      }
    };
    input.click();
  };

  const handleClearHistory = () => {
    toast.success("Todo o histórico de eventos foi removido", "Histórico limpo");
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'warning':
        return <AlertCircle className="w-4 h-4 text-yellow-500" />;
      case 'info':
        return <Info className="w-4 h-4 text-blue-500" />;
      default:
        return null;
    }
  };

  const formatTime = (date: Date) => {
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    
    if (minutes < 1) return 'Agora mesmo';
    if (minutes < 60) return `${minutes}min atrás`;
    
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h atrás`;
    
    return date.toLocaleDateString('pt-BR');
  };

  return (
    <div className="space-y-4">
      <Tabs defaultValue="handlers" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="handlers">
            <Settings className="w-4 h-4 mr-2" />
            Handlers
          </TabsTrigger>
          <TabsTrigger value="history">
            <History className="w-4 h-4 mr-2" />
            Histórico
          </TabsTrigger>
          <TabsTrigger value="advanced">
            <Globe className="w-4 h-4 mr-2" />
            Avançado
          </TabsTrigger>
        </TabsList>

        <TabsContent value="handlers" className="space-y-4 mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Handlers de Notificação</CardTitle>
              <CardDescription>
                Configure como você deseja receber os avisos
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Voice Handler */}
              <div className="flex items-center justify-between p-3 rounded-lg border">
                <div className="flex items-center space-x-3">
                  <div className={`p-2 ${theme.iconBg} rounded-lg`}>
                    <Volume2 className={`w-5 h-5 ${theme.iconMain}`} />
                  </div>
                  <div>
                    <Label className="text-sm font-medium">Avisos por Voz</Label>
                    <p className="text-xs text-gray-500">Text-to-Speech com ElevenLabs</p>
                  </div>
                </div>
                <Switch checked={voiceEnabled} onCheckedChange={setVoiceEnabled} />
              </div>

              {/* Discord Handler - PRO ONLY */}
              <div className="flex items-center justify-between p-3 rounded-lg border border-gray-200 bg-gray-50 opacity-60">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gray-200 rounded-lg">
                    <MessageSquare className="w-5 h-5 text-gray-400" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <Label className="text-sm font-medium text-gray-600">Discord Webhook</Label>
                      <Lock className="w-3 h-3 text-purple-500" />
                    </div>
                    <p className="text-xs text-gray-400">Disponível apenas na versão PRO</p>
                  </div>
                </div>
                <Switch disabled checked={false} />
              </div>

              {/* System Notification - PRO ONLY */}
              <div className="flex items-center justify-between p-3 rounded-lg border border-gray-200 bg-gray-50 opacity-60">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gray-200 rounded-lg">
                    <Bell className="w-5 h-5 text-gray-400" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <Label className="text-sm font-medium text-gray-600">Notificações do Sistema</Label>
                      <Lock className="w-3 h-3 text-purple-500" />
                    </div>
                    <p className="text-xs text-gray-400">Disponível apenas na versão PRO</p>
                  </div>
                </div>
                <Switch disabled checked={false} />
              </div>

              {/* Overlay - PRO ONLY */}
              <div className="flex items-center justify-between p-3 rounded-lg border border-gray-200 bg-gray-50 opacity-60">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gray-200 rounded-lg">
                    <Globe className="w-5 h-5 text-gray-400" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <Label className="text-sm font-medium text-gray-600">Overlay (OBS/Browser)</Label>
                      <Lock className="w-3 h-3 text-purple-500" />
                    </div>
                    <p className="text-xs text-gray-400">Disponível apenas na versão PRO</p>
                  </div>
                </div>
                <Switch disabled checked={false} />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="history" className="space-y-4 mt-4">
          <Card>
            <CardHeader>
              <div className="flex justify-between items-center">
                <div>
                  <CardTitle>Histórico de Eventos</CardTitle>
                  <CardDescription>
                    Últimos avisos e notificações enviados
                  </CardDescription>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleClearHistory}
                  className={theme.buttonTheme}
                >
                  <Trash2 className="w-4 h-4 mr-1" />
                  Limpar
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <ScrollArea className="h-[400px] pr-4">
                <div className="space-y-3">
                  {eventHistory.map((event) => (
                    <div
                      key={event.id}
                      className="flex items-start space-x-3 p-3 rounded-lg border hover:bg-gray-50 transition-colors"
                    >
                      {getStatusIcon(event.status)}
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                          <p className="text-sm font-medium">{event.title}</p>
                          <Badge variant="outline" className="text-xs">
                            {event.type}
                          </Badge>
                        </div>
                        <p className="text-xs text-gray-500">{event.message}</p>
                        <div className="flex items-center text-xs text-gray-400">
                          <Clock className="w-3 h-3 mr-1" />
                          {formatTime(event.timestamp)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="advanced" className="space-y-4 mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Configurações Avançadas</CardTitle>
              <CardDescription>
                Backup, restauração e configurações avançadas
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  className={`flex-1 ${theme.buttonTheme}`}
                  onClick={handleExportConfig}
                >
                  <Download className="w-4 h-4 mr-2" />
                  Exportar Config
                </Button>
                <Button
                  variant="outline"
                  className={`flex-1 ${theme.buttonTheme}`}
                  onClick={handleImportConfig}
                >
                  <Upload className="w-4 h-4 mr-2" />
                  Importar Config
                </Button>
              </div>

              <div className="space-y-2 p-3 bg-gray-50 rounded-lg">
                <Label className="text-sm font-medium">Informações do Sistema</Label>
                <div className="space-y-1 text-xs text-gray-600">
                  <p>• Porta GSI: 3001</p>
                  <p>• Config Path: backend/config.json</p>
                  <p>• Voice Cache: ./voice-cache</p>
                  <p>• Server Binary: build/server/dota-gsi-server</p>
                </div>
              </div>

              <div className="space-y-2 p-3 bg-yellow-50 rounded-lg border border-yellow-200">
                <Label className="text-sm font-medium text-yellow-800">Debug Mode</Label>
                <p className="text-xs text-yellow-700">
                  Ative para salvar ticks GSI em debug/gsi_ticks/
                </p>
                <Switch />
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
