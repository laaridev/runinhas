import { useState } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Button } from './ui/button';
import { 
  ArrowRight, 
  ArrowLeft,
  Volume2, 
  Sliders, 
  Edit3, 
  Mic,
  Timer,
  Sparkles,
  CheckCircle2,
  Zap,
  Settings
} from 'lucide-react';

interface OnboardingTutorialProps {
  open: boolean;
  onComplete: () => void;
}

interface Step {
  title: string;
  subtitle: string;
  icon: React.ReactNode;
  content: React.ReactNode;
}

export function OnboardingTutorial({ open, onComplete }: OnboardingTutorialProps) {
  const [currentStep, setCurrentStep] = useState(0);

  const steps: Step[] = [
    {
      title: "Como o Runinhas funciona?",
      subtitle: "Entenda a mágica por trás",
      icon: <Sparkles className="h-6 w-6 text-purple-500" />,
      content: (
        <div className="space-y-4">
          <div className="bg-gradient-to-r from-purple-50 to-blue-50 dark:from-purple-950/20 dark:to-blue-950/20 p-4 rounded-lg">
            <p className="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">
              O <span className="font-semibold text-purple-600 dark:text-purple-400">Runinhas</span> fica de olho no seu jogo em tempo real!
            </p>
          </div>
          
          <div className="space-y-3">
            <div className="flex gap-3 items-start">
              <div className="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg mt-0.5">
                <Timer className="h-4 w-4 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <h4 className="font-medium text-sm mb-1">Monitora o tempo do jogo</h4>
                <p className="text-xs text-gray-600 dark:text-gray-400">
                  Sabe exatamente quando cada runa vai spawnar, quando vem catapulta, quando muda o ciclo dia/noite
                </p>
              </div>
            </div>

            <div className="flex gap-3 items-start">
              <div className="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg mt-0.5">
                <Zap className="h-4 w-4 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <h4 className="font-medium text-sm mb-1">Avisa você antes</h4>
                <p className="text-xs text-gray-600 dark:text-gray-400">
                  Com X segundos de antecedência (você escolhe!), ele te avisa por áudio que algo importante vai acontecer
                </p>
              </div>
            </div>

            <div className="flex gap-3 items-start">
              <div className="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg mt-0.5">
                <Volume2 className="h-4 w-4 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <h4 className="font-medium text-sm mb-1">Fala com você</h4>
                <p className="text-xs text-gray-600 dark:text-gray-400">
                  Usa inteligência artificial pra criar áudios naturais que não atrapalham sua gameplay
                </p>
              </div>
            </div>
          </div>

        </div>
      )
    },
    {
      title: "Configurando os avisos",
      subtitle: "Personalize do seu jeito",
      icon: <Sliders className="h-6 w-6 text-blue-500" />,
      content: (
        <div className="space-y-4">
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/20 dark:to-indigo-950/20 p-4 rounded-lg">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              Cada timing pode ser configurado individualmente! 
            </p>
          </div>

          <div className="space-y-3">
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-10 h-5 bg-gradient-to-r from-emerald-400 to-emerald-500 rounded-full flex items-center justify-end px-0.5">
                  <div className="w-4 h-4 bg-white rounded-full shadow-sm" />
                </div>
                <span className="text-sm font-medium">Ativar/Desativar</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Use o toggle para ligar ou desligar cada aviso. Simples assim!
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <div className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                  <div className="w-2/3 h-full bg-gradient-to-r from-purple-500 to-blue-500 rounded-full" />
                </div>
                <span className="text-sm font-medium text-blue-600 dark:text-blue-400">45s</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Arraste o slider para escolher com quantos segundos de antecedência quer ser avisado
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <CheckCircle2 className="h-4 w-4 text-green-500" />
                <span className="text-sm font-medium">Auto-save</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Todas as mudanças são salvas automaticamente. Sem botão de salvar!
              </p>
            </div>
          </div>

          <div className="bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
            <p className="text-xs text-blue-800 dark:text-blue-200">
              <span className="font-medium">Dica:</span> Comece com 30-45 segundos e ajuste conforme sua preferência!
            </p>
          </div>
        </div>
      )
    },
    {
      title: "Gerando e ouvindo áudios",
      subtitle: "Crie avisos personalizados",
      icon: <Mic className="h-6 w-6 text-green-500" />,
      content: (
        <div className="space-y-4">
          <div className="bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-950/20 dark:to-emerald-950/20 p-4 rounded-lg">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              Cada aviso tem seu próprio áudio personalizado!
            </p>
          </div>

          <div className="space-y-3">
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <Button size="sm" variant="ghost" className="h-7 px-2">
                  <Mic className="h-3 w-3 mr-1" />
                  Gerar Áudio
                </Button>
                <span className="text-xs text-gray-500">Primeira vez</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Na primeira vez, clique para gerar o áudio com IA. Demora uns segundinhos!
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <Button size="sm" variant="ghost" className="h-7 px-2">
                  <Volume2 className="h-3 w-3 mr-1" />
                  Preview
                </Button>
                <span className="text-xs text-gray-500">Testar áudio</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Depois de gerado, use o Preview para ouvir como ficou. Não gostou? Edite o texto!
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <div className="p-1.5 bg-blue-100 dark:bg-blue-900/30 rounded">
                  <Settings className="h-3 w-3 text-blue-600 dark:text-blue-400" />
                </div>
                <span className="text-sm font-medium">Auto-geração</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Quando você muda os segundos, o áudio é regerado automaticamente com o novo tempo!
              </p>
            </div>
          </div>
        </div>
      )
    },
    {
      title: "Personalizando mensagens",
      subtitle: "Fale do seu jeito",
      icon: <Edit3 className="h-6 w-6 text-purple-500" />,
      content: (
        <div className="space-y-4">
          <div className="bg-gradient-to-r from-purple-50 to-pink-50 dark:from-purple-950/20 dark:to-pink-950/20 p-4 rounded-lg">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              Cansou do texto padrão? Mude para o que quiser!
            </p>
          </div>

          <div className="space-y-3">
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <Button size="sm" variant="ghost" className="h-7 px-2">
                  <Edit3 className="h-3 w-3 mr-1" />
                  Editar Fala
                </Button>
                <span className="text-xs text-gray-500">Personalizar</span>
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Clique em "Editar Fala" para abrir o editor de mensagem
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="font-mono text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded mb-2">
                "Runa de Sabedoria em <span className="text-blue-500">{'{seconds}'}</span> segundos"
              </div>
              <p className="text-xs text-gray-600 dark:text-gray-400">
                Use <code className="bg-gray-200 dark:bg-gray-700 px-1 rounded">{'{seconds}'}</code> onde quiser que apareça o tempo
              </p>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <div className="space-y-1 mb-2">
                <p className="text-xs font-medium">Exemplos criativos:</p>
                <p className="text-xs text-gray-600 dark:text-gray-400 italic">"Wisdom em {'{seconds}'}, prepare o time."</p>
                <p className="text-xs text-gray-600 dark:text-gray-400 italic">"Faltam {'{seconds}'} segundos pra runa do poder"</p>
                <p className="text-xs text-gray-600 dark:text-gray-400 italic">"Vai anoitecer em {'{seconds}'} segundos"</p>
              </div>
            </div>
          </div>

          <div className="bg-purple-50 dark:bg-purple-950/20 border border-purple-200 dark:border-purple-800 rounded-lg p-3">
            <p className="text-xs text-purple-800 dark:text-purple-200">
              <span className="font-medium">Importante:</span> Depois de editar, o áudio é gerado automaticamente com o novo texto!
            </p>
          </div>
        </div>
      )
    },
    {
      title: "Tudo pronto!",
      subtitle: "Hora de jogar!",
      icon: <CheckCircle2 className="h-6 w-6 text-green-500" />,
      content: (
        <div className="space-y-4">
          <div className="bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 dark:from-blue-950/20 dark:via-purple-950/20 dark:to-pink-950/20 rounded-lg p-4 space-y-3">
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-green-500 mt-0.5" />
              <p className="text-sm text-gray-700 dark:text-gray-300">
                GSI instalado e funcionando
              </p>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-green-500 mt-0.5" />
              <p className="text-sm text-gray-700 dark:text-gray-300">
                Você sabe configurar os timings
              </p>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-green-500 mt-0.5" />
              <p className="text-sm text-gray-700 dark:text-gray-300">
                Sabe gerar e personalizar áudios
              </p>
            </div>
          </div>

          <div className="space-y-3">
            <div className="bg-emerald-50 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-800 rounded-lg p-3">
              <p className="text-sm text-emerald-800 dark:text-emerald-200 font-medium mb-1">
                Próximos passos:
              </p>
              <ul className="text-xs text-emerald-700 dark:text-emerald-300 space-y-1 ml-4">
                <li>1. Abra o Dota 2</li>
                <li>2. Entre em uma partida</li>
                <li>3. O Runinhas vai detectar automaticamente!</li>
              </ul>
            </div>

            <div className="bg-amber-50 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3">
              <p className="text-xs text-amber-800 dark:text-amber-200">
                <span className="font-medium">Lembrete:</span> Mantenha o Runinhas aberto durante o jogo. 
              </p>
            </div>
          </div>
        </div>
      )
    }
  ];

  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      onComplete();
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const currentStepData = steps[currentStep];

  return (
    <Dialog open={open} onOpenChange={() => {}}>
      <DialogContent 
        className="sm:max-w-[600px] max-h-[85vh] overflow-hidden flex flex-col" 
        onPointerDownOutside={(e) => e.preventDefault()}
      >
        {/* Progress indicator */}
        <div className="flex items-center justify-center gap-2 pt-2">
          {steps.map((_, index) => (
            <div
              key={index}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                index === currentStep 
                  ? 'w-8 bg-gradient-to-r from-purple-500 to-blue-500' 
                  : index < currentStep
                  ? 'w-4 bg-green-500'
                  : 'w-4 bg-gray-200 dark:bg-gray-700'
              }`}
            />
          ))}
        </div>

        <DialogHeader className="space-y-3">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-gradient-to-br from-purple-100 to-blue-100 dark:from-purple-900/30 dark:to-blue-900/30 rounded-xl">
              {currentStepData.icon}
            </div>
            <div>
              <DialogTitle className="text-xl">
                {currentStepData.title}
              </DialogTitle>
              <DialogDescription className="text-sm">
                {currentStepData.subtitle}
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto py-4 px-1">
          {currentStepData.content}
        </div>

        <DialogFooter className="flex items-center justify-between gap-2">
          <div className="text-xs text-gray-500 dark:text-gray-400">
            Passo {currentStep + 1} de {steps.length}
          </div>
          
          <div className="flex gap-2">
            {currentStep > 0 && (
              <Button
                variant="outline"
                onClick={handlePrevious}
                size="sm"
                className="px-3"
              >
                <ArrowLeft className="h-4 w-4 mr-1" />
                Anterior
              </Button>
            )}
            
            <Button
              onClick={handleNext}
              size="sm"
              className={`px-4 ${
                currentStep === steps.length - 1
                  ? 'bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700'
                  : 'bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700'
              } text-white`}
            >
              {currentStep === steps.length - 1 ? (
                <>
                  Começar
                  <Sparkles className="h-4 w-4 ml-1" />
                </>
              ) : (
                <>
                  Próximo
                  <ArrowRight className="h-4 w-4 ml-1" />
                </>
              )}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
