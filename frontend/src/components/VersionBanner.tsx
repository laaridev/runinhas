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
      className={`
        fixed bottom-0 left-0 right-0 z-40
        px-4 py-2 text-xs font-medium
        backdrop-blur-md border-t
        ${
          isFree
            ? 'bg-gradient-to-r from-emerald-500/90 via-green-500/90 to-teal-500/90 text-white border-emerald-400/30 shadow-lg shadow-emerald-500/20'
            : 'bg-gradient-to-r from-purple-500/90 via-pink-500/90 to-rose-500/90 text-white border-purple-400/30 shadow-lg shadow-purple-500/20'
        }
      `}
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        {/* Left - Version Badge */}
        <div className="flex items-center gap-3">
          {isFree ? (
            <div className="flex items-center gap-2 px-2.5 py-1 rounded-lg bg-white/20 backdrop-blur-sm border border-white/30">
              <Lock className="w-3 h-3" />
              <span className="font-bold tracking-wide">FREE</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 px-2.5 py-1 rounded-lg bg-white/20 backdrop-blur-sm border border-white/30">
              <Crown className="w-3 h-3 animate-pulse" />
              <span className="font-bold tracking-wide">PRO</span>
              <Sparkles className="w-3 h-3" />
            </div>
          )}
          
          {/* Version Number */}
          <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-white/15 backdrop-blur-sm border border-white/20">
            <Code2 className="w-3 h-3 opacity-80" />
            <span className="font-mono font-semibold tracking-wider opacity-90">v1.0.0</span>
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
