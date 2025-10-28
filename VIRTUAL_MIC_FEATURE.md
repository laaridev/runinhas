# 🎤 Virtual Microphone Output - Feature Plan

## 📋 Objetivo
Permitir que os áudios do Runinhas sejam reproduzidos simultaneamente em:
1. **Saída de som padrão** (speakers/fones) - para o usuário ouvir
2. **Microfone virtual** (virt_mic) - para o time ouvir via Discord/Teamspeak/etc

## 🔍 Sistema de Áudio Detectado

### Dispositivos de Saída (Sinks):
```
✅ alsa_output.pci-0000_00_1f.3.analog-stereo (Saída padrão - RUNNING)
✅ virt_sink (Sink virtual - IDLE)
```

### Dispositivos de Entrada (Sources):
```
✅ alsa_output.pci-0000_00_1f.3.analog-stereo.monitor (Monitor da saída)
✅ alsa_input.pci-0000_00_1f.3.analog-stereo (Microfone físico)
✅ virt_sink.monitor (Monitor do sink virtual)
✅ virt_mic (MICROFONE VIRTUAL - RUNNING) ← TARGET!
```

### Configuração Atual:
- **Default Sink:** alsa_output.pci-0000_00_1f.3.analog-stereo
- **Default Source:** virt_mic

## 🎯 Implementação Necessária

### 1. Backend (Go)

#### a) Biblioteca de Áudio
Usar **Beep** ou **Oto** para controle de áudio em Go:
```go
// Opção 1: github.com/faiface/beep
// Pros: Fácil, múltiplos outputs
// Cons: Mais pesado

// Opção 2: github.com/hajimehoshi/oto
// Pros: Leve, baixo nível
// Cons: Mais trabalho manual

// Opção 3: Chamar PulseAudio diretamente
// Pros: Controle total
// Cons: Linux-only
```

#### b) Novo Package: `audio/player.go`
```go
type AudioPlayer struct {
    defaultOutput  OutputDevice
    virtualMic     OutputDevice
    dualOutputMode bool // true = tocar em ambos
}

func (ap *AudioPlayer) Play(audioData []byte) error {
    if ap.dualOutputMode {
        // Tocar simultaneamente em:
        // 1. Default sink (speakers)
        // 2. virt_mic (microfone virtual)
    } else {
        // Tocar apenas no default sink
    }
}
```

#### c) Config Update
```go
// backend/config/config.go
type Config struct {
    // ... existing fields
    
    // Virtual Mic Output
    VirtualMicEnabled bool   `json:"virtual_mic_enabled"`
    VirtualMicDevice  string `json:"virtual_mic_device"` // "virt_mic"
}
```

#### d) Endpoint API
```go
// GET /api/audio/devices - Lista dispositivos disponíveis
// POST /api/audio/settings - Atualiza config de virtual mic
```

### 2. Frontend (React + TypeScript)

#### a) Novo Tab/Seção: "Áudio"
Local: `frontend/src/components/AudioTab.tsx`

Elementos:
```tsx
- [x] Toggle: "Ativar saída para microfone virtual"
- [ ] Dropdown: Selecionar dispositivo (auto-detect "virt_mic")
- [ ] Test Button: Testar áudio no microfone virtual
- [ ] Volume Slider: Volume do microfone virtual (independente)
- [ ] Status: "Microfone virtual detectado ✅" ou "Não encontrado ❌"
```

#### b) Settings Hook
```typescript
// hooks/useAudioSettings.ts
export function useAudioSettings() {
  const [virtualMicEnabled, setVirtualMicEnabled] = useState(false);
  const [virtualMicDevice, setVirtualMicDevice] = useState("virt_mic");
  
  const toggleVirtualMic = async (enabled: boolean) => {
    await SetAudioSettings({ virtualMicEnabled: enabled });
  };
  
  return { virtualMicEnabled, toggleVirtualMic };
}
```

### 3. Fluxo de Execução

```
Evento do Jogo
    ↓
Voice Handler (backend)
    ↓
Audio Player (novo)
    ├─→ Default Sink (sempre)
    └─→ Virtual Mic (se enabled)
```

### 4. Comandos PulseAudio

Para tocar áudio no microfone virtual:
```bash
# Método 1: paplay direto no virt_mic
paplay --device=virt_mic audio.mp3

# Método 2: parec + pacat (loopback)
parec --device=virt_sink.monitor | pacat --device=virt_mic

# Método 3: module-loopback (permanente)
pactl load-module module-loopback source=virt_sink.monitor sink=virt_mic
```

## 📦 Dependências Necessárias

### Go Modules:
```bash
# Opção 1: Beep (recomendado)
go get -u github.com/faiface/beep
go get -u github.com/faiface/beep/mp3
go get -u github.com/faiface/beep/speaker

# Opção 2: Oto (mais leve)
go get -u github.com/hajimehoshi/oto/v2

# Opção 3: PulseAudio bindings
go get -u github.com/mafredri/pulseaudio
```

### Sistema (já tem):
- ✅ PulseAudio instalado
- ✅ Microfone virtual criado (virt_mic)

## 🔧 Implementação em Etapas

### Fase 1: Detecção ✅
- [x] Identificar dispositivos disponíveis
- [x] Detectar virt_mic
- [x] Criar branch feature/virtual-mic-output

### Fase 2: Backend (Audio Player)
- [ ] Criar package `backend/audio/`
- [ ] Implementar AudioPlayer com dual output
- [ ] Integrar com VoiceHandler existente
- [ ] Adicionar config de virtual mic
- [ ] Criar endpoint API

### Fase 3: Frontend (Settings UI)
- [ ] Criar AudioTab.tsx
- [ ] Toggle para ativar virtual mic
- [ ] Teste de áudio
- [ ] Salvar preferências

### Fase 4: Testes
- [ ] Testar áudio em speakers
- [ ] Testar áudio no virt_mic
- [ ] Testar dual output
- [ ] Verificar sincronização
- [ ] Testar volume independente

### Fase 5: Cross-platform
- [ ] Linux (PulseAudio) ✅
- [ ] Windows (VB-Audio Cable detection)
- [ ] Fallback se virtual mic não existir

## ⚠️ Considerações

### Performance:
- Áudio tocado 2x simultaneamente
- Possível latência no virtual mic
- CPU adicional mínimo

### Compatibilidade:
- Linux: PulseAudio (✅ detectado)
- Windows: Requer VB-Audio Cable ou similar
- macOS: Requer BlackHole ou Loopback

### UX:
- Auto-detectar virtual mic na primeira vez
- Se não encontrar, esconder opção
- Mensagem clara: "Para que seu time ouça os avisos"

## 🎨 UI Mock (AudioTab)

```
┌─────────────────────────────────────────────┐
│  ⚙️ Configurações de Áudio                  │
├─────────────────────────────────────────────┤
│                                              │
│  🔊 Saída Padrão                            │
│  ✅ Ativo (Speakers/Fones)                  │
│  ────────────────●──── 100%                 │
│                                              │
│  🎤 Microfone Virtual                       │
│  [Toggle ON/OFF] ← NOVA FEATURE             │
│                                              │
│  Dispositivo: virt_mic (Virtual_Microphone) │
│  Status: ✅ Detectado e funcionando         │
│  ────────────────●──── 80%                  │
│                                              │
│  ℹ️  Com essa opção ativa, seu time poderá │
│     ouvir os avisos via Discord/Teamspeak   │
│                                              │
│  [Testar Áudio no Microfone Virtual]        │
│                                              │
└─────────────────────────────────────────────┘
```

## 📝 Próximos Passos

1. **Escolher biblioteca de áudio** (Recomendo: Beep)
2. **Implementar AudioPlayer básico**
3. **Testar dual output no Linux**
4. **Criar UI no frontend**
5. **Integrar com sistema existente**

---

**Branch:** `feature/virtual-mic-output`  
**Status:** 🟡 Em planejamento  
**Dispositivo Target:** `virt_mic` (Virtual_Microphone)  
**Sistema:** PulseAudio (Linux) ✅
