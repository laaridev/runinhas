import { useState } from 'react';
import { Sparkles, Loader2 } from 'lucide-react';
import { Button } from './ui/button';
import { SetMode } from '../../wailsjs/go/main/App';
import { useToast } from '@/hooks/useToast';

export function UpgradeProButton() {
  const [isUpgrading, setIsUpgrading] = useState(false);
  const toast = useToast();

  const handleUpgrade = async () => {
    setIsUpgrading(true);
    
    try {
      // For now, directly activate PRO mode without checkout
      // TODO: Integrate with Lemon Squeezy checkout
      // const checkout = window.open(import.meta.env.VITE_LS_CHECKOUT_URL, "_blank", "width=600,height=800");
      
      // Simulate license key (in production, this comes from webhook)
      const licenseKey = `RUNINHAS-${Date.now().toString(36).toUpperCase()}`;
      
      // Call backend to activate PRO mode
      // @ts-ignore - SetMode bindings will be regenerated on build
      await SetMode("pro", licenseKey);
      
      toast.success(
        "Modo PRO ativado!",
        "🎉 Agora você tem acesso a todos os recursos premium"
      );
    } catch (error: any) {
      console.error('Upgrade error:', error);
      toast.error(
        "Erro ao ativar PRO",
        error.message || "Tente novamente"
      );
    } finally {
      setIsUpgrading(false);
    }
  };

  return (
    <Button
      onClick={handleUpgrade}
      disabled={isUpgrading}
      className="
        bg-gradient-to-r from-pink-500 via-purple-500 to-blue-500
        hover:from-pink-600 hover:via-purple-600 hover:to-blue-600
        text-white font-bold
        shadow-lg hover:shadow-xl
        transition-all duration-300
        hover:scale-105
        disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100
      "
    >
      {isUpgrading ? (
        <>
          <Loader2 className="w-4 h-4 mr-2 animate-spin" />
          Ativando...
        </>
      ) : (
        <>
          <Sparkles className="w-4 h-4 mr-2" />
          Upgrade PRO ⭐
        </>
      )}
    </Button>
  );
}
