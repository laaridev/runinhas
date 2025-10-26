import { Sparkles, Lock } from 'lucide-react';
import { useAppMode } from '@/hooks/useAppMode';
import { UpgradeProButton } from './UpgradeProButton';

export function VersionBanner() {
  const { appMode, loading } = useAppMode();

  if (loading) {
    return null; // Don't show banner while loading
  }

  const isFree = appMode.mode === 'free';

  return (
    <div
      className={`
        px-4 py-2.5 text-sm font-medium
        ${
          isFree
            ? 'bg-gradient-to-r from-gray-100 to-gray-200 text-gray-700 border-b border-gray-300'
            : 'bg-gradient-to-r from-purple-500 to-pink-500 text-white border-b border-purple-600'
        }
      `}
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        <div className="flex items-center gap-2">
          {isFree ? (
            <>
              <Lock className="w-4 h-4" />
              <span>🔒 Versão Free</span>
            </>
          ) : (
            <>
              <Sparkles className="w-4 h-4 animate-pulse" />
              <span>⭐ Modo PRO ativo</span>
            </>
          )}
        </div>
        
        {isFree && (
          <UpgradeProButton />
        )}
      </div>
    </div>
  );
}
