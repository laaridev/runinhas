import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { useLoading } from "@/hooks/useLoading";
import { useDebouncedCallback } from "@/hooks/useDebounce";
import { useToast } from "@/hooks/useToast";
import { useWailsAudioPlayer } from "@/hooks/useWailsAudioPlayer";
import { cache } from "@/services/cache";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { EventCard } from "@/components/shared/EventCard";
import { ConfigTab } from "@/components/ConfigTab";
import ElevenLabsSettings from "@/components/ElevenLabsSettings";
import { WelcomeModal } from "@/components/WelcomeModal";
import { OnboardingTutorial } from "@/components/OnboardingTutorial";
import { ToastContainer } from "@/components/ToastContainer";
import { LanguageSelector } from "@/components/LanguageSelector";
import { VersionBanner } from "@/components/VersionBanner";
import { ANIMATION } from "@/constants/defaults";
import { timingAPI, messageAPI } from "@/services/api-wails";
import {
  Coins,
  Zap,
  Droplet,
  Brain,
  Clock,
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
import logoBlue from "@/assets/logo-runinha-blue.svg";
import logoPink from "@/assets/logo-runinha-pink.svg";

function App() {
  const { t } = useTranslation(['common', 'events']);
  const { withLoading } = useLoading();
  const toast = useToast();
  const audioPlayer = useWailsAudioPlayer(); // Initialize Wails audio player
  
  // Log audio player state for debugging
  useEffect(() => {
    if (audioPlayer.isPlaying) {
      console.log('🎵 Playing:', audioPlayer.currentFile);
    }
  }, [audioPlayer]);
  
  const [gsiInstalled, setGsiInstalled] = useState(false);
  const [serverRunning, setServerRunning] = useState(false);
  const [currentTheme, setCurrentTheme] = useState<"blue" | "pink">("blue");
  const [customMessages, setCustomMessages] = useState<Record<string, string>>({});
  const [showWelcomeModal, setShowWelcomeModal] = useState(false);
  const [showOnboarding, setShowOnboarding] = useState(false);

  // Definições de tema (mantém o existente)
  const themes = {
    blue: {
      name: t('common:themes.blue'),
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
      name: t('common:themes.pink'),
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
      name: t('events:stack_timing.name'),
      description: t('events:stack_timing.description'),
      icon: <Package className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 5,
      max: 60,
      step: 1,
    },
    {
      key: "day_night_cycle",
      name: t('events:day_night_cycle.name'),
      description: t('events:day_night_cycle.description'),
      icon: <Sun className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 30,
      step: 5,
    },
    {
      key: "catapult_timing",
      name: t('events:catapult_timing.name'),
      description: t('events:catapult_timing.description'),
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
      name: t('events:bounty_rune.name'),
      description: t('events:bounty_rune.description'),
      icon: <Coins className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 90,
      step: 5,
    },
    {
      key: "power_rune",
      name: t('events:power_rune.name'),
      description: t('events:power_rune.description'),
      icon: <Zap className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 90,
      step: 5,
    },
    {
      key: "water_rune",
      name: t('events:water_rune.name'),
      description: t('events:water_rune.description'),
      icon: <Droplet className={`w-6 h-6 ${theme.iconMain} transition-colors duration-500`} />,
      min: 10,
      max: 50,
      step: 5,
    },
    {
      key: "wisdom_rune",
      name: t('events:wisdom_rune.name'),
      description: t('events:wisdom_rune.description'),
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
    <div className={`min-h-screen bg-gradient-to-br ${theme.background} transition-all duration-500 pb-10`}>
      {/* Clean Light Navbar with Theme Colors */}
      <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur-lg border-b border-gray-200/60 shadow-sm transition-all duration-500">
        <div className="max-w-7xl mx-auto px-8 py-4">
          <div className="flex items-center justify-between">
            
            {/* Left - Status Pill */}
            <div className="flex items-center">
              {serverRunning ? (
                <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-50 border border-emerald-200/60 transition-all duration-300 hover:bg-emerald-100/80">
                  {/* Pulsating Indicator */}
                  <div className="relative flex items-center justify-center">
                    <div className="absolute w-2 h-2 bg-emerald-400 rounded-full animate-ping opacity-75"></div>
                    <div className="relative w-2 h-2 bg-emerald-500 rounded-full"></div>
                  </div>
                  <span className="text-xs font-bold text-emerald-700">Online</span>
                  <div className="w-px h-3 bg-emerald-300"></div>
                  <span className="text-xs font-medium text-emerald-600">Aguardando partida</span>
                </div>
              ) : (
                <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-gray-100 border border-gray-200">
                  <div className="relative w-2 h-2 bg-gray-400 rounded-full"></div>
                  <span className="text-xs font-bold text-gray-500">Offline</span>
                </div>
              )}
            </div>

            {/* Center - Brand Identity */}
            <div className="absolute left-1/2 -translate-x-1/2 flex items-center gap-4 pointer-events-none">
              {/* Logo with premium orbital animation - CENTERED */}
              <div className="relative w-14 h-14 group pointer-events-auto">
                {/* Pulsating glow - centered */}
                <div className={`absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-16 h-16 ${theme.iconBg} rounded-full blur-2xl animate-glow-pulse pointer-events-none`} />
                
                {/* Orbital particles - centered */}
                <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-20 h-20 pointer-events-none">
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full animate-orbit">
                    <div className={`absolute top-0 left-1/2 -translate-x-1/2 w-1.5 h-1.5 rounded-full ${
                      currentTheme === "blue" ? "bg-blue-400" : "bg-pink-400"
                    } shadow-lg`} />
                  </div>
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full animate-orbit" style={{ animationDelay: "-2.67s" }}>
                    <div className={`absolute top-0 left-1/2 -translate-x-1/2 w-1.5 h-1.5 rounded-full ${
                      currentTheme === "blue" ? "bg-purple-400" : "bg-purple-400"
                    } shadow-lg`} />
                  </div>
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full animate-orbit" style={{ animationDelay: "-5.34s" }}>
                    <div className={`absolute top-0 left-1/2 -translate-x-1/2 w-1.5 h-1.5 rounded-full ${
                      currentTheme === "blue" ? "bg-pink-400" : "bg-blue-400"
                    } shadow-lg`} />
                  </div>
                </div>
                
                {/* Logo - centered */}
                <div className="relative w-full h-full flex items-center justify-center transition-all duration-300 group-hover:scale-110">
                  <img 
                    src={theme.logo} 
                    alt="Runinhas" 
                    className="w-full h-full drop-shadow-lg" 
                  />
                </div>
              </div>
              
              {/* Brand name - large and prominent */}
              <h1 className={`text-4xl font-black ${theme.titleText} tracking-tight leading-none transition-colors duration-500 pointer-events-auto`}>
                runinhas
              </h1>
            </div>

            {/* Right - Actions */}
            <div className="flex items-center gap-3">
              
              {/* Language Selector */}
              <LanguageSelector />
              
              {/* Theme Toggle with theme colors */}
              <button
                onClick={() => setCurrentTheme(currentTheme === "blue" ? "pink" : "blue")}
                className={`group relative w-9 h-9 rounded-lg ${theme.iconBg} border ${theme.navbarBorder} flex items-center justify-center transition-all duration-300 hover:shadow-sm hover:scale-105`}
                title={t('common:buttons.switch_theme')}
              >
                <Palette className={`w-3.5 h-3.5 ${theme.iconColor} transition-colors duration-300`} />
              </button>

              {/* Server Control Button */}
              {!gsiInstalled ? (
                <button
                  onClick={handleInstallGSI}
                  className="group relative px-4 h-9 rounded-lg bg-gradient-to-r from-emerald-500 to-teal-500 text-white font-semibold text-sm shadow-sm hover:shadow-md transition-all duration-300 hover:scale-105 flex items-center gap-2"
                >
                  <Download className="w-3.5 h-3.5" />
                  <span>{t('common:buttons.install_gsi')}</span>
                </button>
              ) : (
                <button
                  onClick={serverRunning ? handleStopServer : handleStartServer}
                  className={`group relative px-4 h-9 rounded-lg font-semibold text-sm transition-all duration-300 hover:scale-105 flex items-center gap-2 ${
                    serverRunning
                      ? 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200 shadow-sm'
                      : `bg-gradient-to-r ${theme.gradient} text-white shadow-sm hover:shadow-md`
                  }`}
                >
                  {serverRunning ? (
                    <>
                      <div className="relative w-2 h-2 bg-gray-500 rounded-full"></div>
                      <span>{t('common:buttons.disconnect')}</span>
                    </>
                  ) : (
                    <>
                      <div className="relative flex items-center justify-center">
                        <div className="absolute w-2 h-2 bg-white/40 rounded-full animate-ping"></div>
                        <div className="relative w-2 h-2 bg-white rounded-full"></div>
                      </div>
                      <span>{t('common:buttons.connect')}</span>
                    </>
                  )}
                </button>
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
              {t('common:tabs.runes')}
            </TabsTrigger>
            <TabsTrigger value="timing">
              <Clock className="w-4 h-4 mr-2" />
              {t('common:tabs.timing')}
            </TabsTrigger>
            <TabsTrigger value="voice">
              <Volume2 className="w-4 h-4 mr-2" />
              {t('common:tabs.voice')}
            </TabsTrigger>
            <TabsTrigger value="config">
              <Settings className="w-4 h-4 mr-2" />
              {t('common:tabs.config')}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="runes" className="mt-4">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: ANIMATION.TAB_TRANSITION / 1000, ease: "easeInOut" }}
              className="grid gap-4 md:grid-cols-2"
            >
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
            </motion.div>
          </TabsContent>

          <TabsContent value="timing" className="mt-4">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: ANIMATION.TAB_TRANSITION / 1000, ease: "easeInOut" }}
              className="grid gap-4 md:grid-cols-2 lg:grid-cols-3"
            >
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
            </motion.div>
          </TabsContent>

          <TabsContent value="voice" className="mt-4">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: ANIMATION.TAB_TRANSITION / 1000, ease: "easeInOut" }}
            >
              <ElevenLabsSettings />
            </motion.div>
          </TabsContent>
          
          <TabsContent value="config" className="mt-4">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: ANIMATION.TAB_TRANSITION / 1000, ease: "easeInOut" }}
            >
              <ConfigTab theme={theme} />
            </motion.div>
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

      {/* Version Footer - Fixed at bottom */}
      <VersionBanner />
    </div>
  );
}

export default App;
