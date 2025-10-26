<div align="center">
  <img src="https://i.imgur.com/CFu2M7H.png" alt="Runinhas" width="320"/>

# 🎮 Runinhas

### Alertas de voz em tempo real para timings críticos do Dota 2

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://reactjs.org/)
[![Wails](https://img.shields.io/badge/Wails-v2-DF0000?style=for-the-badge)](https://wails.io/)

[📥 Download](#-instalação) • [📖 Docs](#-features) • [🐛 Report Bug](https://github.com/laaridev/runinhas/issues)

---

</div>

## 📖 Sobre

**Runinhas** é um app desktop que monitora o Dota 2 em tempo real através do **Game State Integration (GSI)** oficial da Valve e dispara alertas de voz para timings críticos do jogo.

**Nunca mais perca:** runas, stacks, catapult waves ou ciclos de dia/noite.

### 🆓 Versão FREE vs ⭐ PRO

| Feature | FREE (Padrão) | PRO |
|---------|---------------|-----|
| **Alertas de Voz** | ✅ Vozes padrão integradas | ✅ Vozes personalizadas (ElevenLabs) |
| **Mensagens** | ❌ Genéricas fixas | ✅ Customizáveis com `{seconds}` |
| **Vozes Disponíveis** | 1 voz PT-BR (embutida) | 10+ vozes (ElevenLabs) |
| **Configuração de Voz** | ❌ Bloqueado | ✅ Tom, velocidade, estilo |
| **Custo** | 🆓 Grátis (sempre) | Requer API key ElevenLabs* |

_* ElevenLabs tem plano free com 10k caracteres/mês (suficiente para uso casual)_

**Modo FREE:**
- Áudios genéricos embutidos no binário (~200KB)
- Sem necessidade de API key
- Funciona 100% offline
- Ideal para testar o app

**Modo PRO:**
- Ative com uma API key do [ElevenLabs](https://elevenlabs.io)
- Mensagens dinâmicas: "Runa em **30** segundos"
- Escolha voz, tom, velocidade
- Edite mensagens em tempo real

<br/>

---

## ⚡ Features

<table>
<tr>
<td width="50%">

### 🎯 Alertas Suportados

| Evento | Timing |
|--------|--------|
| ⏰ **Bounty Runes** | 0:00, 3:00, 6:00... (cada 3min) |
| 💎 **Power Runes** | 6:00, 8:00, 10:00... (cada 2min) |
| 💧 **Water Runes** | 2:00 e 4:00 apenas |
| 📚 **Wisdom Runes** | 7:00, 14:00, 21:00... (cada 7min) |
| 📦 **Stacks** | Timing perfeito (:53) |
| 🏰 **Catapults** | 5:00, 10:00, 15:00... (cada 5min) |
| 🌓 **Day/Night** | Ciclo de 10min (5min cada) |

</td>
<td width="50%">

### ⚙️ Configuração

- 🎚️ **Antecedência ajustável** (5s ~ 60s)
- 🔔 **Ative/desative** alertas individualmente
- 🎨 **2 temas** (Azul/Rosa)
- 💾 **Auto-save** em todas alterações
- 🌍 **Multi-idioma** (PT-BR/EN)
- 🎤 **Vozes customizadas** (modo PRO)

<br/>

### 🔒 Segurança

✅ **100% Seguro - Sem Risco de Ban**

- Usa apenas **GSI oficial** da Valve
- ❌ Zero modificação de arquivos
- ❌ Zero injeção de código
- ❌ Zero leitura de memória

</td>
</tr>
</table>

<br/>

---

## 🏗️ Arquitetura Event-Driven

<p align="center">
<pre>
┌─────────────────────────────────────────────────────────────┐
│                         Dota 2 GSI                          │
│                    (Game State Integration)                 │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP POST (JSON)
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                     GSI Server (:3001)                       │
│              Recebe ticks do jogo a cada 500ms               │
└──────────────────────┬──────────────────────────────────────┘
                       │ Publica no Event Bus
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      Event Bus (Go)                          │
│            Canal buffered (100 eventos na fila)              │
└───────┬──────────────────────────────┬──────────────────────┘
        │                              │
        ▼                              ▼
┌──────────────────┐          ┌──────────────────┐
│  Rune Consumer   │          │ Timing Consumer  │
│                  │          │                  │
│ • Bounty Runes   │          │ • Stacks         │
│ • Power Runes    │          │ • Catapults      │
│ • Water Runes    │          │ • Day/Night      │
│ • Wisdom Runes   │          │                  │
└────────┬─────────┘          └────────┬─────────┘
         │                             │
         │ Emite evento de áudio       │
         ▼                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     Voice Handler                            │
│                                                              │
│  FREE mode: Usa áudio embedded (go:embed)                    │
│  PRO mode:  Gera com ElevenLabs + cache local                │
└──────────────────────┬──────────────────────────────────────┘
                       │ Wails Event: audio:play
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Frontend (React)                           │
│                Audio Player + Event Queue                    │
└─────────────────────────────────────────────────────────────┘
</pre>
</p>

### Componentes Principais

**Event Bus:**
- Canal Go buffered com capacidade de 100 eventos
- Previne memory leaks com drop tracking
- Métricas em tempo real (events/s, drops)

**Consumers:**
- Processam eventos de forma assíncrona
- Throttle configurável (evita spam)
- Parse único com cache (otimização de CPU)

**Voice Handler:**
- FREE: Serve MP3s do `go:embed` (zero latência)
- PRO: Cache semântico (evita gerar áudio duplicado)
- Garbage collection automático (TTL 7 dias)

<br/>

---

## 🚀 Instalação

### 📥 1. Download

Baixe o executável na seção [**Releases**](https://github.com/laaridev/runinhas/releases):

<div align="center">

| Sistema | Arquivo | Tamanho |
|---------|---------|---------|
| 🐧 **Linux (x64)** | `runinhas` | ~12MB |
| 🪟 **Windows (x64)** | `runinhas.exe` | Em breve |

</div>

### ⚙️ 2. Configuração GSI

O app **cria automaticamente** o arquivo GSI na primeira execução em:

```
📁 steamapps/common/dota 2 beta/game/dota/cfg/gamestate_integration/
└── gamestate_integration_runinhas.cfg
```

<details>
<summary>📄 Ver conteúdo do arquivo GSI (opcional)</summary>

```json
"Runinhas GSI Configuration"
{
  "uri"           "http://localhost:3001/gsi"
  "timeout"       "5.0"
  "buffer"        "0.5"
  "throttle"      "0.5"
  "heartbeat"     "30.0"
  "data"
  {
    "map"         "1"
    "provider"    "1"
  }
}
```

</details>

### ✅ 3. Pronto!

Execute o app e jogue Dota 2. Os alertas tocam automaticamente! 🎵

<br/>

---

## 🛠️ Desenvolvimento

<details>
<summary><b>📚 Stack Tecnológica</b></summary>

<br/>

| Camada | Tecnologia |
|--------|-----------|
| 🖥️ **Desktop Framework** | Wails v2 (Go + WebView) |
| ⚙️ **Backend** | Go 1.24 |
| 🎨 **Frontend** | React 18 + TypeScript + Vite |
| 💅 **UI/Styling** | Tailwind CSS + shadcn/ui |
| 🎤 **TTS** | ElevenLabs API (modo PRO) |
| 🔊 **Audio** | go:embed + HTML5 Audio |

</details>

<details>
<summary><b>🔨 Comandos de Build</b></summary>

<br/>

```bash
# Instalar Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Instalar dependências do frontend
cd frontend && npm install && cd ..

# Dev mode (hot reload)
wails dev

# Build de produção
wails build

# Build multi-plataforma
wails build -platform windows/amd64
wails build -platform linux/amd64
```

</details>

<details>
<summary><b>📁 Estrutura do Projeto</b></summary>

<br/>

```
backend/
├── assets/          # 🎵 Áudios embedded (go:embed)
├── config/          # ⚙️ Config + singleton pattern
├── consumers/       # 📡 Event consumers (Rune, Timing)
├── events/          # 🔄 Event Bus + ParsedTickEvent
├── handlers/        # 🎤 Voice Handler (FREE/PRO)
├── i18n/            # 🌍 Traduções (PT-BR/EN)
├── metrics/         # 📊 Métricas e monitoring
└── server/          # 🌐 HTTP server + endpoints

frontend/src/
├── components/      # 🧩 React components
├── hooks/           # 🎣 useAppMode, useWailsAudioPlayer
├── services/        # 🔌 API Wails bindings
└── i18n/            # 🌍 Traduções do frontend
```

</details>

<br/>

---

## 📜 Licença

MIT License - Código aberto e gratuito.

---

<div align="center">

### ⭐ Desenvolvido para a comunidade Dota 2 ❤️

[![GitHub issues](https://img.shields.io/github/issues/laaridev/runinhas?style=for-the-badge)](https://github.com/laaridev/runinhas/issues)
[![GitHub stars](https://img.shields.io/github/stars/laaridev/runinhas?style=for-the-badge)](https://github.com/laaridev/runinhas/stargazers)

**[🐛 Reportar Bug](https://github.com/laaridev/runinhas/issues)** • **[💡 Sugerir Feature](https://github.com/laaridev/runinhas/issues/new)**

</div>
