import { Sparkles, Crown, Check } from 'lucide-react';
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
        px-4 py-1.5 text-xs font-medium border-b
        ${
          isFree
            ? 'bg-gradient-to-r from-emerald-50 via-green-50 to-teal-50 dark:from-emerald-900/20 dark:via-green-900/20 dark:to-teal-900/20 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800'
            : 'bg-gradient-to-r from-purple-500 via-pink-500 to-rose-500 text-white border-purple-600'
        }
      `}
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        <div className="flex items-center gap-2">
          {isFree ? (
            <>
              <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-emerald-100 dark:bg-emerald-800/30 border border-emerald-200 dark:border-emerald-700">
                <Check className="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
                <span className="font-semibold text-emerald-700 dark:text-emerald-300">
                  Versão Free
                </span>
              </div>
            </>
          ) : (
            <>
              <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-white/20 backdrop-blur-sm">
                <Crown className="w-3 h-3 animate-pulse" />
                <span className="font-semibold">Modo PRO</span>
                <Sparkles className="w-3 h-3" />
              </div>
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
