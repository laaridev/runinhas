# 🔧 Correção: API Key não atualizava em tempo real

## 🐛 Problema Identificado

Quando você atualizava a API Key do ElevenLabs no frontend, ela era salva corretamente no arquivo `config.json`, mas o **VoiceHandler continuava usando a API key antiga em memória**, resultando em erros de quota mesmo com créditos disponíveis.

### Por que acontecia?

1. **Frontend salvava**: API key era salva no JSON ✅
2. **Backend lia arquivo**: Config era persistida corretamente ✅  
3. **VoiceHandler não atualizava**: O handler mantinha a API key da inicialização ❌

## ✅ Solução Aplicada

### Arquivos Modificados

#### 1. `/backend/server/config_endpoints.go`

**Adicionado**: Atualização do VoiceHandler em memória após salvar config

```go
// Update VoiceHandler settings in memory if voice config was updated
if voiceData, ok := updates["voice"].(map[string]interface{}); ok {
    if vh, ok := s.voiceHandler.(*handlers.VoiceHandler); ok {
        // Extrai valores da configuração
        apiKey := ""
        voiceID := "eVXYtPVYB9wDoz9NVTIy"
        stability := 0.5
        // ... outros parâmetros
        
        // Atualiza handler em memória
        vh.UpdateSettings(apiKey, voiceID, stability, similarity, style, speakerBoost)
        s.logger.Info("✅ VoiceHandler settings updated in memory")
    }
}
```

#### 2. `/backend/server/elevenlabs_handler.go`

**Melhorado**: Endpoint de teste agora sempre carrega a API key mais recente

```go
// Load current config to get API key
cfg, err := config.Load()
if err != nil {
    http.Error(w, "Failed to load config", http.StatusInternalServerError)
    return
}

apiKey := cfg.ElevenLabsAPIKey
if apiKey == "" && cfg.Game != nil && cfg.Game.Voice != nil {
    if key, ok := cfg.Game.Voice["apiKey"].(string); ok {
        apiKey = key
    }
}

// Atualiza VoiceHandler com API key fresca
vh.UpdateSettings(apiKey, voiceID, ...)
```

## 🎯 Como Funciona Agora

### Fluxo Completo

1. **Usuário cola nova API Key** no frontend (ElevenLabsSettings.tsx)
2. **Frontend salva via debounce** (1 segundo após digitar)
3. **Backend recebe POST** `/api/config` com nova config
4. **Backend salva no JSON** ✅
5. **Backend atualiza VoiceHandler** em memória ✅ **[NOVO!]**
6. **Próxima requisição usa** a nova API key imediatamente ✅

### Endpoints Corrigidos

- ✅ `POST /api/config` - Salva e atualiza handler
- ✅ `POST /api/elevenlabs/test` - Sempre usa API key mais recente
- ✅ `POST /api/elevenlabs/config` - Atualiza handler após salvar

## 🧪 Como Testar

1. **Abra o app** e vá em Configurações de Voz
2. **Cole uma nova API Key** (com créditos)
3. **Aguarde 1 segundo** (debounce automático)
4. **Clique em "Testar Voz"**
5. ✅ Deve funcionar imediatamente com a nova key

### Verificar nos Logs

Procure por:
```
✅ VoiceHandler settings updated in memory
Configuration saved successfully
```

## 🔄 Reinicialização Não Necessária

**Antes**: Precisava reiniciar o app para usar nova API key  
**Agora**: Atualização instantânea, sem reinicialização

## 📋 Checklist de Funcionamento

- ✅ API key salva no config.json
- ✅ VoiceHandler atualizado em memória
- ✅ Teste de voz funciona imediatamente
- ✅ Eventos do jogo usam nova API key
- ✅ Sem necessidade de reiniciar

## 🐛 Troubleshooting

### Se ainda der erro de quota:

1. **Verifique o config.json**:
```bash
cat ~/.config/runinhas/config.json | grep apiKey
```

2. **Confira os logs** durante o salvamento:
```bash
# Deve aparecer:
✅ VoiceHandler settings updated in memory
```

3. **Teste a API key manualmente**:
```bash
curl -X POST https://api.elevenlabs.io/v1/voices \
  -H "xi-api-key: SUA_API_KEY"
```

### Se a API key não atualizar:

1. **Certifique-se** que esperou 1 segundo (debounce)
2. **Verifique** se apareceu "Configurações salvas automaticamente"
3. **Reinicie** o servidor se necessário (última opção)

## 📝 Notas Técnicas

- Usa **singleton pattern** na configuração
- **UpdateSettings()** já existia no VoiceHandler
- Frontend usa **debounce de 1 segundo** para evitar salvamentos excessivos
- Configuração é **thread-safe** (sync.Once + sync.RWMutex)

---

*Correção aplicada em: 2025-10-23*  
*Problema reportado pelo usuário: API key não atualizava, erro de quota persistia*
