# 🎮 Runinhas - Dota 2 GSI Assistant

> *"sem tilts, só timing"* ⚡

Sistema de Game State Integration para Dota 2 com interface moderna, sistema de temas dinâmicos e avisos em tempo real.

## ✨ Features

- 🎨 **2 Temas Dinâmicos:** Azul e Rosa com transições suaves
- ⏰ **Avisos de Runas:** Bounty, Power, Water e Wisdom
- 📦 **Timing Perfeito:** Stack, Day/Night e Catapult
- 🔊 **Text-to-Speech:** Avisos por voz em tempo real
- 💾 **Auto-save:** Todas configurações salvas automaticamente
- 🎯 **Interface Moderna:** Design elegante com glassmorphism
- ⚡ **Performance:** Build único com backend + frontend embutido

## 📁 Estrutura do Projeto

```
dota-gsi/
├── backend/                    # Backend em Go
│   ├── config/                # Gerenciamento de configuração
│   ├── consumers/             # Event consumers (rune, timing, map, hero)
│   ├── handlers/              # Event handlers (voice/TTS)
│   ├── server/                # Servidor HTTP + endpoints
│   ├── installer/             # Instalador GSI
│   ├── events/                # Event bus
│   └── utils/                 # Utilitários
├── frontend/                   # Frontend React + TypeScript
│   ├── src/
│   │   ├── components/        # UI components (shadcn/ui)
│   │   ├── services/          # API services
│   │   ├── types/             # TypeScript types
│   │   └── App.tsx            # Componente principal
│   └── package.json
├── build/                      # Recursos de build
│   ├── appicon.png            # Ícone do app
│   └── windows/               # Recursos Windows (.ico, manifest)
├── scripts/                    # Scripts de build
│   └── build.sh               # Build Linux + Windows
├── app.go                      # Wails App (bindings Go ↔ Frontend)
├── main.go                     # Entry point
└── wails.json                  # Configuração do Wails
```

## 🚀 Instalação e Build

### Pré-requisitos
- **Go** 1.21+ ([instalar](https://go.dev/dl/))
- **Node.js** 18+ e npm ([instalar](https://nodejs.org/))
- **Wails** v2 ([instalar](https://wails.io/docs/gettingstarted/installation))

### 1. Clonar o Repositório
```bash
git clone https://github.com/seu-usuario/dota-gsi.git
cd dota-gsi
```

### 2. Instalar Dependências
```bash
# Frontend
cd frontend
npm install
cd ..

# Backend (automático com go build)
```

### 3. Build para Produção
```bash
# Build completo (recomendado)
./scripts/build.sh

# Ou manualmente
cd frontend && npm run build && cd ..
~/go/bin/wails build -tags desktop,production
```

### 4. Executar
```bash
# Linux
./build/bin/runinhas-linux

# Windows
./build/bin/runinhas.exe
```

### Desenvolvimento
```bash
# Modo dev (Wails + hot reload)
wails dev

# Ou apenas frontend
cd frontend && npm run dev
```

## 📂 Configurações do Sistema

O app cria automaticamente as configurações nos diretórios do sistema:

- **Linux**: `~/.config/runinhas/config.json`
- **Windows**: `%APPDATA%\Runinhas\config.json`
- **Cache de áudio**: `~/.cache/runinhas/voice/` (Linux) ou `%LOCALAPPDATA%\Runinhas\Cache\voice\` (Windows)

## 🗂️ Artefatos de Build

- **build/bin/**: Binários executáveis gerados
- **build/windows/**: Recursos Windows (ícone, manifest, metadados)
- **frontend/dist/**: Bundle do frontend (usado pelo Wails)

## 🎨 Stack Tecnológica

### Backend
- Go + Wails v2
- Logrus (logging)
- ElevenLabs (TTS)

### Frontend
- React 18 + TypeScript
- Vite
- Tailwind CSS
- shadcn/ui

## 📦 Build Final

Executável único com backend + frontend embutido.
