import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Button } from './ui/button';
import { Shield, Info, CheckCircle, AlertTriangle, FileCode, ExternalLink, Lock, Users, Loader2, RefreshCw, Zap, Settings } from 'lucide-react';
import { InstallGSI } from '../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import logoBlue from '@/assets/logo-runinha-blue.svg';
import { motion } from 'framer-motion';

interface WelcomeModalProps {
  open: boolean;
  onClose: () => void;
  onInstallComplete: () => void;
}

export function WelcomeModal({ open, onClose, onInstallComplete }: WelcomeModalProps) {
  const { t } = useTranslation('welcome');
  const [step, setStep] = useState<'welcome' | 'installing' | 'success' | 'error'>('welcome');
  const [errorMessage, setErrorMessage] = useState('');

  const handleInstall = async () => {
    setStep('installing');

    try {
      const result = await InstallGSI();
      if (result.success) {
        setStep('success');
        setTimeout(() => {
          onInstallComplete();
          onClose();
        }, 2500);
      } else {
        setErrorMessage(result.message || t('status.installing'));
        setStep('error');
      }
    } catch (error) {
      setErrorMessage(t('status.installing'));
      setStep('error');
    }
  };

  const openLink = (url: string) => {
    BrowserOpenURL(url);
  };

  return (
    <Dialog open={open} onOpenChange={() => { }}>
      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto bg-white dark:bg-gray-900 border border-purple-200/30 dark:border-purple-800/30 shadow-2xl shadow-purple-500/10"
        onPointerDownOutside={(e) => e.preventDefault()}>
        {step === 'welcome' && (
          <>
            <DialogHeader className="space-y-6 text-center pb-4">
              <div className="flex flex-col items-center space-y-5">
                <motion.div
                  className="relative"
                  initial={{ scale: 0.9, opacity: 0 }}
                  animate={{
                    scale: 1,
                    opacity: 1
                  }}
                  transition={{
                    duration: 0.6,
                    ease: "easeOut"
                  }}
                >
                  <motion.div
                    className="absolute inset-0 bg-gradient-to-br from-purple-400 via-blue-500 to-indigo-500 rounded-full blur-3xl opacity-30"
                    animate={{
                      scale: [1, 1.15, 1],
                      opacity: [0.25, 0.4, 0.25]
                    }}
                    transition={{
                      duration: 4,
                      repeat: Infinity,
                      ease: "easeInOut",
                      repeatType: "reverse"
                    }}
                  />
                  <motion.img
                    src={logoBlue}
                    alt="Runinhas Logo"
                    className="h-32 w-32 drop-shadow-2xl relative z-10"
                    animate={{
                      y: [0, -8, 0]
                    }}
                    transition={{
                      duration: 3,
                      repeat: Infinity,
                      ease: "easeInOut",
                      repeatType: "reverse"
                    }}
                  />
                </motion.div>

                <motion.div
                  className="space-y-2"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.2, duration: 0.5 }}
                >
                  <DialogTitle className="text-3xl font-bold bg-gradient-to-r from-purple-600 via-blue-600 to-purple-600 bg-clip-text text-transparent bg-[length:200%_auto] animate-gradient">
                    {t('title')}
                  </DialogTitle>
                </motion.div>
              </div>

            </DialogHeader>

            <motion.div
              className="space-y-4 py-4"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.4, duration: 0.6 }}
            >
              {/* Boas-vindas */}
              <div className="flex gap-3 p-4 bg-white/50 dark:bg-white/5 backdrop-blur-sm rounded-xl border border-purple-200/40 dark:border-purple-800/30 shadow-sm hover:shadow-lg hover:border-purple-300/60 dark:hover:border-purple-700/50 transition-all duration-300">
                <CheckCircle className="h-6 w-6 text-purple-600 dark:text-purple-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-3">
                  <div>
                    <h4 className="font-semibold text-sm mb-1.5">{t('how_it_works.title')}</h4>
                    <p className="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">
                      {t('description')}
                    </p>
                  </div>

                  <div className="space-y-2 pl-1">
                    <div className="flex items-start gap-2">
                      <Zap className="h-4 w-4 text-purple-500 dark:text-purple-400 mt-0.5 flex-shrink-0" />
                      <p className="text-xs text-gray-600 dark:text-gray-400">
                        Interpreta eventos do jogo e te avisa com voz natural sobre momentos críticos
                      </p>
                    </div>
                    <div className="flex items-start gap-2">
                      <Settings className="h-4 w-4 text-blue-500 dark:text-blue-400 mt-0.5 flex-shrink-0" />
                      <p className="text-xs text-gray-600 dark:text-gray-400">
                        Avisos sobre spawn de runas, tempo ideal para stacks e transição dia/noite
                      </p>
                    </div>
                    <div className="flex items-start gap-2">
                      <Shield className="h-4 w-4 text-green-500 dark:text-green-400 mt-0.5 flex-shrink-0" />
                      <p className="text-xs text-gray-600 dark:text-gray-400">
                        Processo local e seguro usando apenas recursos oficialmente suportados pela Valve
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Segurança VAC */}
              <div className="flex gap-3 p-4 bg-white/50 dark:bg-white/5 backdrop-blur-sm rounded-xl border border-green-200/40 dark:border-green-800/30 shadow-sm hover:shadow-lg hover:border-green-300/60 dark:hover:border-green-700/50 transition-all duration-300">
                <Shield className="h-6 w-6 text-green-600 dark:text-green-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-2">
                  <h4 className="font-semibold text-sm flex items-center gap-2">
                    100% Seguro - Sem risco de VAC Ban
                    <span className="text-xs bg-green-100 dark:bg-green-900/50 text-green-700 dark:text-green-300 px-2 py-0.5 rounded-full">
                      Oficial
                    </span>
                  </h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {t('is_safe.description')}
                  </p>
                  <ul className="text-xs text-gray-600 dark:text-gray-400 space-y-1 ml-4">
                    <li>• Usado em todos os torneios oficiais (TI, Majors, DPC)</li>
                    <li>• Integrado em plataformas como Twitch e YouTube</li>
                    <li>• Milhares de streamers usam diariamente</li>
                  </ul>
                </div>
              </div>

              {/* O que é GSI */}
              <div className="flex gap-3 p-4 bg-white/50 dark:bg-white/5 backdrop-blur-sm rounded-xl border border-blue-200/40 dark:border-blue-800/30 shadow-sm hover:shadow-lg hover:border-blue-300/60 dark:hover:border-blue-700/50 transition-all duration-300">
                <Info className="h-6 w-6 text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-2">
                  <h4 className="font-semibold text-sm">{t('what_is_gsi.title')}</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {t('what_is_gsi.description')}
                  </p>
                  <div className="flex flex-wrap gap-2 mt-2">
                    <button
                      onClick={() => openLink('https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration')}
                      className="inline-flex items-center gap-1 text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      <ExternalLink className="h-3 w-3" />
                      Documentação Oficial Valve
                    </button>
                    <button
                      onClick={() => openLink('https://github.com/antonpup/Dota2GSI')}
                      className="inline-flex items-center gap-1 text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      <ExternalLink className="h-3 w-3" />
                      Exemplos no GitHub
                    </button>
                  </div>
                </div>
              </div>

              {/* Privacidade */}
              <div className="flex gap-3 p-4 bg-white/50 dark:bg-white/5 backdrop-blur-sm rounded-xl border border-purple-200/40 dark:border-purple-800/30 shadow-sm hover:shadow-lg hover:border-purple-300/60 dark:hover:border-purple-700/50 transition-all duration-300">
                <Lock className="h-6 w-6 text-purple-600 dark:text-purple-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-2">
                  <h4 className="font-semibold text-sm">Privacidade e Fair Play</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    O GSI <span className="font-semibold">envia apenas dados do SEU herói</span>.
                    Não há acesso a informações de inimigos ou aliados que você não poderia ver normalmente no jogo.
                  </p>
                  <div className="bg-purple-100/50 dark:bg-purple-900/20 rounded p-2 mt-2">
                    <p className="text-xs text-purple-700 dark:text-purple-300">
                      <Users className="h-3 w-3 inline mr-1" />
                      <span className="font-semibold">Dados enviados:</span> Tempo de jogo, status das runas,
                      cooldowns do seu herói, seus itens e habilidades.
                    </p>
                  </div>
                </div>
              </div>

              {/* Como funciona */}
              <div className="flex gap-3 p-4 bg-white/50 dark:bg-white/5 backdrop-blur-sm rounded-xl border border-gray-200/40 dark:border-gray-700/30 shadow-sm hover:shadow-lg hover:border-gray-300/60 dark:hover:border-gray-600/50 transition-all duration-300">
                <FileCode className="h-6 w-6 text-gray-600 dark:text-gray-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-2">
                  <h4 className="font-semibold text-sm">Como funciona a instalação?</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Para o Runinhas funcionar, ele irá criar um arquivo de configuração (.cfg) do GSI dentro da pasta do dota 2 na pasta do Dota 2 que
                    permite que o jogo envie dados para o aplicativo em tempo real.
                  </p>
                  <div className="bg-gray-100 dark:bg-gray-900/50 rounded p-2 mt-2 font-mono">
                    <p className="text-xs text-gray-600 dark:text-gray-400 break-all">
                      Steam/steamapps/common/dota 2 beta/game/dota/cfg/gamestate_integration/
                    </p>
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 italic">
                    * Você pode remover este arquivo a qualquer momento para desativar o GSI
                  </p>
                </div>
              </div>

              {/* Aviso de permissão */}
              <div className="flex gap-3 p-4 bg-amber-50/70 dark:bg-amber-950/20 backdrop-blur-sm rounded-xl border-2 border-amber-300/60 dark:border-amber-700/40 shadow-sm hover:shadow-lg hover:border-amber-400/70 dark:hover:border-amber-600/50 transition-all duration-300">
                <AlertTriangle className="h-6 w-6 text-amber-600 dark:text-amber-400 mt-0.5 flex-shrink-0" />
                <div className="space-y-1">
                  <h4 className="font-semibold text-sm">Permissão necessária</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    O Runinhas precisa criar um arquivo na pasta do Dota 2.
                    Isso é necessário apenas uma vez e você pode revisar o arquivo criado a qualquer momento.
                  </p>
                </div>
              </div>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.8, duration: 0.5 }}
            >
              <DialogFooter className="flex justify-center pt-4 pb-2">
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Button
                    onClick={handleInstall}
                    size="lg"
                    className="bg-gradient-to-r from-purple-600 via-blue-600 to-purple-600 hover:from-purple-700 hover:via-blue-700 hover:to-purple-700 text-white px-12 py-6 text-base font-semibold shadow-xl hover:shadow-2xl transition-all duration-300 bg-[length:200%_auto] animate-gradient"
                  >
                    <CheckCircle className="h-5 w-5 mr-2" />
                    {t('buttons.install')}
                  </Button>
                </motion.div>
              </DialogFooter>
            </motion.div>
          </>
        )}

        {step === 'installing' && (
          <div className="py-12 text-center space-y-4">
            <div className="inline-flex p-4 bg-gradient-to-br from-purple-100 to-blue-100 dark:from-purple-950/50 dark:to-blue-950/50 rounded-full">
              <Loader2 className="h-10 w-10 text-blue-600 dark:text-blue-400 animate-spin" />
            </div>
            <h3 className="text-xl font-semibold">{t('status.installing')}</h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Configurando a integração com o Dota 2
            </p>
            <div className="text-xs text-gray-500 dark:text-gray-400 animate-pulse">
              Criando arquivo de configuração na pasta do jogo...
            </div>
          </div>
        )}

        {step === 'success' && (
          <div className="py-12 text-center space-y-4">
            <div className="inline-flex p-4 bg-green-100 dark:bg-green-950/50 rounded-full">
              <CheckCircle className="h-8 w-8 text-green-600 dark:text-green-400" />
            </div>
            <h3 className="text-lg font-semibold">{t('status.installed')}</h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Tudo pronto! Agora inicie o Dota 2 e o Runinhas funcionará automaticamente
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">
              Fechando em 2 segundos...
            </p>
          </div>
        )}

        {step === 'error' && (
          <>
            <DialogHeader className="space-y-2">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-red-100 dark:bg-red-950/50 rounded-lg">
                  <AlertTriangle className="h-6 w-6 text-red-600 dark:text-red-400" />
                </div>
                <DialogTitle className="text-xl">{t('status.not_installed')}</DialogTitle>
              </div>
            </DialogHeader>
            <div className="py-6 space-y-4">
              <div className="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                <p className="text-sm text-red-800 dark:text-red-200 font-medium">
                  {errorMessage}
                </p>
              </div>
              <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
                <p className="font-medium">Possíveis soluções:</p>
                <ul className="list-disc list-inside space-y-1 ml-2">
                  <li>Verifique se o Dota 2 está instalado via Steam</li>
                  <li>Execute o aplicativo como administrador</li>
                  <li>Certifique-se que o Steam está aberto</li>
                </ul>
              </div>
            </div>
            <DialogFooter className="flex gap-2 justify-center">
              <Button
                variant="outline"
                onClick={() => setStep('welcome')}
                className="px-6"
              >
                {t('buttons.close')}
              </Button>
              <Button
                onClick={handleInstall}
                className="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white px-6"
              >
                <RefreshCw className="h-4 w-4 mr-2" />
                {t('buttons.install')}
              </Button>
            </DialogFooter>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
