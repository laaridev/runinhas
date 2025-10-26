import { Sparkles, Crown, Lock, Code2 } from 'lucide-react';
import { useAppMode } from '@/hooks/useAppMode';
import { UpgradeProButton } from './UpgradeProButton';

export function VersionBanner() {
  const { appMode, loading } = useAppMode();

  if (loading) {
    return null;
  }

  const isFree = appMode.mode === 'free';

  return (
    <footer
      className="
        fixed bottom-0 left-0 right-0 z-40
        bg-white/95 backdrop-blur-lg
        border-t border-gray-200/60 shadow-sm
        px-4 py-2.5 text-xs font-medium
      "
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        {/* Left - Version Badge */}
        <div className="flex items-center gap-3">
          {isFree ? (
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gradient-to-r from-emerald-500 to-green-500 text-white shadow-md shadow-emerald-500/30 border border-emerald-400">
              <Lock className="w-3.5 h-3.5" />
              <span className="font-bold tracking-wide">FREE</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-md shadow-purple-500/30 border border-purple-400">
              <Crown className="w-3.5 h-3.5 animate-pulse" />
              <span className="font-bold tracking-wide">PRO</span>
              <Sparkles className="w-3.5 h-3.5" />
            </div>
          )}
          
          {/* Version Number */}
          <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-gray-50 border border-gray-200 shadow-sm">
            <Code2 className="w-3.5 h-3.5 text-gray-600" />
            <span className="font-mono font-semibold tracking-wider text-gray-700">v1.0.0</span>
          </div>
        </div>

        {/* Right - Upgrade Button */}
        {isFree && (
          <UpgradeProButton />
        )}
      </div>
    </footer>
  );
}
