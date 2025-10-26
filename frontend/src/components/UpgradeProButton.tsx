import { useState } from 'react';
import { Sparkles, Crown, Zap, Lock, X } from 'lucide-react';
import { Button } from './ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';

export function UpgradeProButton() {
  const [showModal, setShowModal] = useState(false);

  const handleUpgrade = () => {
    setShowModal(true);
  };

  return (
    <>
      <Button
        onClick={handleUpgrade}
        size="sm"
        className="
          bg-gradient-to-r from-pink-500 via-purple-500 to-blue-500
          hover:from-pink-600 hover:via-purple-600 hover:to-blue-600
          text-white font-semibold text-xs
          shadow-md hover:shadow-lg
          transition-all duration-300
          hover:scale-105
          border-0
        "
      >
        <Sparkles className="w-3 h-3 mr-1.5" />
        Upgrade PRO
      </Button>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-2xl">
              <Crown className="w-6 h-6 text-yellow-500" />
              Runinhas PRO
            </DialogTitle>
            <DialogDescription className="text-base pt-2">
              A versão PRO está em desenvolvimento
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Features List */}
            <div className="space-y-3">
              <div className="flex items-start gap-3 p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20">
                <Zap className="w-5 h-5 text-purple-600 dark:text-purple-400 mt-0.5" />
                <div>
                  <p className="font-medium text-sm">Mensagens Personalizadas</p>
                  <p className="text-xs text-muted-foreground">
                    Customize os avisos de eventos
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20">
                <Sparkles className="w-5 h-5 text-blue-600 dark:text-blue-400 mt-0.5" />
                <div>
                  <p className="font-medium text-sm">Voz Sintetizada (ElevenLabs)</p>
                  <p className="text-xs text-muted-foreground">
                    Vozes naturais em português e inglês
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 p-3 rounded-lg bg-pink-50 dark:bg-pink-900/20">
                <Crown className="w-5 h-5 text-pink-600 dark:text-pink-400 mt-0.5" />
                <div>
                  <p className="font-medium text-sm">Recursos Exclusivos</p>
                  <p className="text-xs text-muted-foreground">
                    Novos recursos sendo desenvolvidos
                  </p>
                </div>
              </div>
            </div>

            {/* Coming Soon Message */}
            <div className="bg-gradient-to-r from-yellow-50 to-orange-50 dark:from-yellow-900/20 dark:to-orange-900/20 p-4 rounded-lg border border-yellow-200 dark:border-yellow-800">
              <div className="flex items-start gap-3">
                <Lock className="w-5 h-5 text-yellow-600 dark:text-yellow-400 mt-0.5" />
                <div>
                  <p className="font-semibold text-sm text-yellow-900 dark:text-yellow-100">
                    Em Breve!
                  </p>
                  <p className="text-xs text-yellow-800 dark:text-yellow-200 mt-1">
                    Estamos finalizando a versão PRO. Em breve você poderá fazer upgrade e desbloquear todos os recursos premium.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="flex justify-end">
            <Button
              onClick={() => setShowModal(false)}
              variant="secondary"
              className="gap-2"
            >
              <X className="w-4 h-4" />
              Fechar
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
