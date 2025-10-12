import { useState, useEffect } from "react";
import { useLoading } from "@/hooks/useLoading";
import { useDebouncedCallback } from "@/hooks/useDebounce";
import { useToast } from "@/hooks/useToast";
import { cache } from "@/services/cache";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { EventCard } from "@/components/shared/EventCard";
import { ConfigTab } from "@/components/ConfigTab";
import ElevenLabsSettings from "@/components/ElevenLabsSettings";
import { WelcomeModal } from "@/components/WelcomeModal";
import { OnboardingTutorial } from "@/components/OnboardingTutorial";
import { ToastContainer } from "@/components/ToastContainer";
import {
  Coins,
  Zap,
  Droplet,
  Brain,
  Clock,
  Play,
  Square,
  Download,
  Palette,
  Package,
  Sun,
  Shield,
  Settings,
  Volume2,
} from "lucide-react";
import {
  InstallGSI,
  IsGSIInstalled,
  IsDotaInstalled,
  IsServerRunning,
  StartEmbeddedServer,
  StopEmbeddedServer
} from '../wailsjs/go/main/App';
// Use api-wails.ts em vez de api.ts para evitar problemas de CORS
import { timingAPI, messageAPI } from "@/services/api-wails";
import logoBlue from "@/assets/logo-runinha-blue.svg";
import logoPink from "@/assets/logo-runinha-pink.svg";

function App() {
  const { withLoading } = useLoading();
  const toast = useToast();
  
  const [gsiInstalled, setGsiInstalled] = useState(false);
  const [serverRunning, setServerRunning] = useState(false);
  const [currentTheme, setCurrentTheme] = useState<"blue" | "pink">("blue");
  const [customMessages, setCustomMessages] = useState<Record<string, string>>({});
  const [showWelcomeModal, setShowWelcomeModal] = useState(false);
  const [showOnboarding, setShowOnboarding] = useState(false);

  // Definições de tema (mantém o existente)
  const themes = {
    blue: {
      name: "Azul",
      background: "from-blue-600 via-indigo-600 to-purple-600",
      navbar: "from-blue-50/95 via-indigo-50/95 to-purple-50/95",
      navbarBorder: "border-blue-200/50",
      logo: logoBlue,
      accent: "blue",
      accentNum: 500,
      tabActive: "bg-blue-500/20",
      iconColor: "text-blue-600",
      iconBg: "bg-blue-50",
      iconMain: "text-blue-600",
      sliderThumb: "bg-blue-500",
      sliderTrack: "bg-blue-200",
      titleText: "text-blue-900",
      subtitleText: "text-blue-700",
      buttonTheme: "text-blue-700 hover:bg-blue-100/50",
      gradient: "data-[state=checked]:from-blue-500 data-[state=checked]:to-purple-500",
    },
    pink: {
      name: "Rosa",
      background: "from-pink-500 via-rose-500 to-fuchsia-600",
      navbar: "from-pink-50/95 via-rose-50/95 to-fuchsia-50/95",
      navbarBorder: "border-pink-200/50",
      logo: logoPink,
      accent: "rose",
      accentNum: 500,
      tabActive: "bg-rose-500/20",
      iconColor: "text-rose-600",
      iconBg: "bg-rose-50",
      iconMain: "text-rose-600",
      sliderThumb: "bg-rose-500",
      sliderTrack: "bg-rose-200",
      titleText: "text-rose-900",
      subtitleText: "text-rose-700",
      buttonTheme: "text-rose-700 hover:bg-rose-100/50",
      gradient: "data-[state=checked]:from-pink-500 data-[state=checked]:to-fuchsia-500",
    },
  };

  const theme = themes[currentTheme];

  // Configurações de timing
  const timings = [
    {
      key: "stack_timing",
      name: "Stack Timing",
      description: "Avisar para stackar camps (XX:53)",
      icon: <Package className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 5,
      max: 15,
      step: 1,
    },
    {
      key: "day_night_cycle",
      name: "Ciclo Dia/Noite",
      description: "Mudanças de ciclo (a cada 5min)",
      icon: <Sun className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 30,
      step: 5,
    },
    {
      key: "catapult_timing",
      name: "Catapulta",
      description: "Spawn de catapultas (a cada 5min)",
      icon: <Shield className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 30,
      step: 5,
    },
  ];

  // Configurações de runas
  const runes = [
    {
      key: "bounty_rune",
      name: "Runa de Recompensa",
      description: "0:00, depois a cada 3min",
      icon: <Coins className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 90,
      step: 5,
    },
    {
      key: "power_rune",
      name: "Runa de Poder",
      description: "A cada 2min (rios)",
      icon: <Zap className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 90,
      step: 5,
    },
    {
      key: "water_rune",
      name: "Runa de Água",
      description: "2:00 e 4:00",
      icon: <Droplet className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 50,
      step: 5,
    },
    {
      key: "wisdom_rune",
      name: "Runa de Sabedoria",
      description: "7:00, depois a cada 7min",
      icon: <Brain className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 90,
      step: 5,
    },
  ];

  // Estados para runas
  const [runeStates, setRuneStates] = useState<Record<string, { enabled: boolean; value: number }>>({});
  
  // Estados para timings
  const [timingStates, setTimingStates] = useState<Record<string, { enabled: boolean; value: number }>>({});

  // Carregar configurações iniciais
  useEffect(() => {
    const initApp = async () => {
      try {
        // 1. Verificar se Dota 2 está instalado
        const dotaCheck = await IsDotaInstalled();
        if (!dotaCheck.installed) {
          toast.error(dotaCheck.message || "Dota 2 não encontrado. Por favor, instale o Dota 2 via Steam.");
          // Podemos mostrar um modal específico aqui se quiser
          return;
        }
        
        // 2. Verificar se o arquivo GSI .cfg já existe
        const gsiInstalled = await IsGSIInstalled();
        setGsiInstalled(gsiInstalled);
        
        // 3. Se GSI não está instalado, mostrar modal de boas-vindas
        if (!gsiInstalled) {
          setShowWelcomeModal(true);
        }
        
        // 4. Carregar outras configurações
        checkServerStatus();
        loadRuneConfigs();
        loadTimingConfigs();
      } catch (error) {
        console.error("Erro ao inicializar app:", error);
      }
    };
    
    initApp();
  }, []);

  const checkServerStatus = async () => {
    let running = false;
    try {
      running = await IsServerRunning();
    } catch (error) {
      console.error('Error checking server status:', error);
    } finally {
      setServerRunning(running);
    }
  };

  const loadRuneConfigs = async () => {
    const states: Record<string, { enabled: boolean; value: number }> = {};
    const messages: Record<string, string> = {};
    
    for (const rune of runes) {
      const enabled = await timingAPI.getEnabled(rune.key);
      const value = await timingAPI.getValue(rune.key, "warning_seconds");
      states[rune.key] = { enabled, value: value || rune.min };
      
      // Load custom message if available
      const customMsg = await messageAPI.get(rune.key);
      if (customMsg) {
        messages[rune.key] = customMsg;
      }
    }
    
    setRuneStates(states);
    setCustomMessages(messages);
  };

  const loadTimingConfigs = async () => {
    const states: Record<string, { enabled: boolean; value: number }> = {};
    const messages: Record<string, string> = {};
    
    for (const timing of timings) {
      const enabled = await timingAPI.getEnabled(timing.key);
      const value = await timingAPI.getValue(timing.key, "warning_seconds");
      states[timing.key] = { enabled, value: value || timing.min };
      
      // Load custom message if available
      const customMsg = await messageAPI.get(timing.key);
      if (customMsg) {
        messages[timing.key] = customMsg;
      }
    }
    
    setTimingStates(states);
    setCustomMessages(prev => ({ ...prev, ...messages }));
  };

  const handleRuneEnabledChange = async (key: string, enabled: boolean) => {
    try {
      // Atualiza o estado local imediatamente para feedback visual
      setRuneStates(prev => ({
        ...prev,
        [key]: { ...prev[key], enabled }
      }));
      
      // Depois persiste no backend
      await timingAPI.setEnabled(key, enabled);
    } catch (error) {
      console.error('Failed to update rune:', error);
      // Reverte o estado local em caso de erro
      setRuneStates(prev => ({
        ...prev,
        [key]: { ...prev[key], enabled: !enabled }
      }));
    }
  };

  // Usa debounce para evitar múltiplas chamadas ao backend
  const debouncedSetValue = useDebouncedCallback(
    async (key: string, value: number) => {
      try {
        await timingAPI.setValue(key, "warning_seconds", value);
        // Salva no cache para offline mode
        cache.set(`rune_${key}_value`, value, 60 * 60 * 1000); // 1 hora
        toast.success(`Runa ${key} atualizada: ${value}s`);
      } catch (error) {
        toast.error('Erro ao salvar valor da runa');
      }
    },
    500 // 500ms de delay (aumentado para garantir)
  );

  const handleRuneValueChange = async (key: string, value: number) => {
    try {
      // Atualiza o estado local imediatamente
      setRuneStates(prev => ({
        ...prev,
        [key]: { ...prev[key], value }
      }));
      
      // Persiste no backend com debounce
      debouncedSetValue(key, value);
    } catch (error) {
      console.error('Failed to update rune value:', error);
      toast.error('Erro ao salvar configuração');
    }
  };

  const handleTimingEnabledChange = async (key: string, enabled: boolean) => {
    try {
      // Atualiza o estado local imediatamente
      setTimingStates(prev => ({
        ...prev,
        [key]: { ...prev[key], enabled }
      }));
      
      // Depois persiste no backend
      await timingAPI.setEnabled(key, enabled);
    } catch (error) {
      console.error('Failed to update timing:', error);
      // Reverte o estado local em caso de erro
      setTimingStates(prev => ({
        ...prev,
        [key]: { ...prev[key], enabled: !enabled }
      }));
    }
  };

  // Usa o mesmo debounce para timings
  const debouncedSetTimingValue = useDebouncedCallback(
    async (key: string, value: number) => {
      try {
        await timingAPI.setValue(key, "warning_seconds", value);
        cache.set(`timing_${key}_value`, value, 60 * 60 * 1000);
        toast.success(`Timing ${key} atualizado: ${value}s`);
      } catch (error) {
        toast.error('Erro ao salvar timing');
      }
    },
    500
  );

  const handleTimingValueChange = async (key: string, value: number) => {
    try {
      // Atualiza o estado local imediatamente
      setTimingStates(prev => ({
        ...prev,
        [key]: { ...prev[key], value }
      }));
      
      // Persiste com debounce
      debouncedSetTimingValue(key, value);
    } catch (error) {
      console.error('Failed to update timing value:', error);
      toast.error('Erro ao salvar timing');
    }
  };

  const handleMessageChange = async (key: string, message: string) => {
    try {
      await messageAPI.set(key, message);
      setCustomMessages(prev => ({ ...prev, [key]: message }));
      toast.success('Mensagem personalizada salva');
    } catch (error) {
      toast.error('Erro ao salvar mensagem');
    }
  };

  const handleInstallGSI = async () => {
    const result = await withLoading('install-gsi', async () => {
      return await InstallGSI();
    });
    
    if (result.success) {
      setGsiInstalled(true);
      toast.success('GSI instalado com sucesso!');
    } else {
      toast.error(result.message || 'Erro na instalação');
    }
  };

  const handleStartServer = async () => {
    try {
      await withLoading('start-server', async () => {
        await StartEmbeddedServer();
      });
      setServerRunning(true);
      toast.success('Servidor GSI rodando na porta 3001');
    } catch (error) {
      toast.error(`Erro ao iniciar servidor: ${error}`);
    }
  };

  const handleStopServer = async () => {
    try {
      await withLoading('stop-server', async () => {
        await StopEmbeddedServer();
      });
      setServerRunning(false);
      toast.info('Servidor GSI encerrado');
    } catch (error) {
      toast.error(`Erro ao parar servidor: ${error}`);
    }
  };


  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme.background} transition-all duration-500`}>
      {/* Navbar */}
      <nav className={`sticky top-0 z-50 bg-gradient-to-r ${theme.navbar} backdrop-blur-xl border-b ${theme.navbarBorder} shadow-lg transition-all duration-500`}>
        <div className="max-w-7xl mx-auto px-8 py-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="relative">
                <img src={theme.logo} alt="Runinhas" className="w-12 h-12 drop-shadow-lg" />
                {serverRunning && (
                  <span className="absolute -top-1 -right-1 flex h-3 w-3">
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                    <span className="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
                  </span>
                )}
              </div>
              <div>
                <h1 className={`text-2xl font-black tracking-tight ${theme.titleText} transition-colors duration-500`}>
                  runinhas
                </h1>
                <p className="text-xs text-gray-500 -mt-1">Dota 2 Assistant</p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              {/* Theme toggle */}
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setCurrentTheme(currentTheme === "blue" ? "pink" : "blue")}
                className={`${theme.buttonTheme} hover:scale-105 transition-transform`}
              >
                <Palette className="w-4 h-4 mr-2" />
                {theme.name}
              </Button>

              {/* Server controls */}
              {!gsiInstalled ? (
                <Button
                  onClick={handleInstallGSI}
                  size="sm"
                  className="bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-600 hover:to-teal-700 text-white shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105"
                >
                  <Download className="w-4 h-4 mr-2" />
                  Instalar GSI
                </Button>
              ) : (
                <Button
                  onClick={serverRunning ? handleStopServer : handleStartServer}
                  size="sm"
                  className={`${
                    serverRunning
                      ? 'bg-gradient-to-r from-red-500 to-rose-600 hover:from-red-600 hover:to-rose-700'
                      : 'bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-600 hover:to-teal-700'
                  } text-white shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105 min-w-[120px]`}
                >
                  {serverRunning ? (
                    <>
                      <Square className="w-4 h-4 mr-2 animate-pulse" />
                      <span className="animate-pulse">Parar</span>
                    </>
                  ) : (
                    <>
                      <Play className="w-4 h-4 mr-2" />
                      Iniciar
                    </>
                  )}
                </Button>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Main content */}
      <main className="max-w-5xl mx-auto p-6">
        <Tabs defaultValue="runes" className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="runes">
              <Coins className="w-4 h-4 mr-2" />
              Runas
            </TabsTrigger>
            <TabsTrigger value="timing">
              <Clock className="w-4 h-4 mr-2" />
              Timing
            </TabsTrigger>
            <TabsTrigger value="voice">
              <Volume2 className="w-4 h-4 mr-2" />
              Voz
            </TabsTrigger>
            <TabsTrigger value="config">
              <Settings className="w-4 h-4 mr-2" />
              Config
            </TabsTrigger>
          </TabsList>

          <TabsContent value="runes" className="mt-4">
            <div className="grid gap-4 md:grid-cols-2">
              {runes.map((rune) => (
                <EventCard
                  key={rune.key}
                  event={rune}
                  enabled={runeStates[rune.key]?.enabled || false}
                  value={runeStates[rune.key]?.value || rune.min}
                  customMessage={customMessages[rune.key]}
                  theme={theme}
                  onToggle={(enabled: boolean) => handleRuneEnabledChange(rune.key, enabled)}
                  onValueChange={(value: number) => handleRuneValueChange(rune.key, value)}
                  onMessageChange={(message: string) => handleMessageChange(rune.key, message)}
                />
              ))}
            </div>
          </TabsContent>

          <TabsContent value="timing" className="mt-4">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {timings.map((timing) => (
                <EventCard
                  key={timing.key}
                  event={timing}
                  enabled={timingStates[timing.key]?.enabled || false}
                  value={timingStates[timing.key]?.value || timing.min}
                  customMessage={customMessages[timing.key]}
                  theme={theme}
                  onToggle={(enabled: boolean) => handleTimingEnabledChange(timing.key, enabled)}
                  onValueChange={(value: number) => handleTimingValueChange(timing.key, value)}
                  onMessageChange={(message: string) => handleMessageChange(timing.key, message)}
                />
              ))}
            </div>
          </TabsContent>

          <TabsContent value="voice" className="mt-4">
            <ElevenLabsSettings />
          </TabsContent>
          
          <TabsContent value="config" className="mt-4">
            <ConfigTab theme={theme} />
          </TabsContent>
        </Tabs>
      </main>

      {/* Toast Container para notificações */}
      <ToastContainer />
      
      {/* Welcome Modal - Primeira vez ou GSI não instalado */}
      <WelcomeModal 
        open={showWelcomeModal}
        onClose={() => setShowWelcomeModal(false)}
        onInstallComplete={() => {
          setGsiInstalled(true);
          setShowWelcomeModal(false);
          // Após instalar GSI, mostrar tutorial
          setShowOnboarding(true);
        }}
      />
      
      {/* Onboarding Tutorial - Após instalar GSI na primeira vez */}
      <OnboardingTutorial 
        open={showOnboarding}
        onComplete={() => {
          setShowOnboarding(false);
          // Marcar que não é mais primeira vez
          toast.success('Tutorial concluído! Aproveite o runinhas!');
        }}
      />
    </div>
  );
}

export default App;
