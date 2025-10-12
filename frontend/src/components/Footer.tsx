import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Wifi, WifiOff } from "lucide-react";
import { EventsOn } from "../../wailsjs/runtime/runtime";

interface FooterProps {
  theme: {
    iconColor: string;
    iconBg: string;
    navbar: string;
    navbarBorder: string;
    iconMain: string;
  };
}

export function Footer({ theme }: FooterProps) {
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    // Escutar status de conexão
    const unsubConnection = EventsOn("game:connection", (data: any) => {
      setConnected(data.connected);
    });

    // Cleanup
    return () => {
      unsubConnection();
    };
  }, []);

  return (
    <footer
      className={`bg-gradient-to-r ${theme.navbar} backdrop-blur-xl px-8 py-4 border-t-2 ${theme.navbarBorder} sticky bottom-0 z-50 transition-all duration-500`}
    >
      <div className="flex items-center justify-center max-w-7xl mx-auto">
        {/* Status de Conexão */}
        <div className="flex items-center gap-2.5">
          {connected ? (
            <>
              <Wifi
                className={`w-5 h-5 ${theme.iconMain} transition-colors duration-500`}
              />
              <Badge className="bg-emerald-500/30 text-emerald-700 border-2 border-emerald-400 font-semibold text-sm px-4 py-1.5">
                Conectado
              </Badge>
            </>
          ) : (
            <>
              <WifiOff
                className={`w-5 h-5 text-gray-400 transition-colors duration-500`}
              />
              <Badge className="bg-gray-500/20 text-gray-600 border-2 border-gray-400/50 font-semibold text-sm px-4 py-1.5">
                Desconectado
              </Badge>
            </>
          )}
        </div>
      </div>
    </footer>
  );
}
