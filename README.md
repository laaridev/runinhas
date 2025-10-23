<div align="center">
  <img src="https://i.imgur.com/CFu2M7H.png" alt="Runinhas" width="280"/>
  
  # Runinhas
  
  ### *"sem tilts, só timing"* ⚡
  
  Alertas de voz em tempo real para timings críticos do Dota 2
  
</div>

---

## 🏗️ Arquitetura

```
Dota 2 (GSI) → Event Bus → Consumers → Voice Handler → ElevenLabs TTS → Alertas
```

**Stack Tecnológica:**
- **Backend:** Go 1.24 + Wails v2 (Event-driven architecture)
- **Frontend:** React 18 + TypeScript + Tailwind CSS + shadcn/ui
- **TTS:** ElevenLabs API
- **Performance:** -40% CPU • <1% event drops • Zero memory leaks

---

## 💡 O Problema Resolvido

Você está focado em last-hits, olhando o mapa, calculando o próximo movimento... **e esquece que a runa de poder spawna em 15 segundos.**

Quando lembra, o mid inimigo já pegou.

**O Runinhas soluciona isso:** alertas de voz automáticos nos timings exatos que você configurar. Sem precisar ficar olhando o relógio, sem precisar calcular mentalmente.

---

## ⚡ Funcionalidades

### 🎤 Alertas de Voz Customizáveis

Avisos em português com voz natural (ElevenLabs):
- ⏰ **Runas de Bounty** (0:00, 3:00, 5:00...)
- 💎 **Runas de Poder** (6:00, 8:00, 10:00...)
- 💧 **Runas de Água** (2:00, 4:00, 6:00...)
- 📚 **Runas de Sabedoria** (7:00, 14:00, 21:00...)
- 📦 **Stacks de Neutrals** (timing perfeito)
- 🏰 **Catapult Waves** (5:00, 10:00, 15:00...)
- 🌓 **Ciclos Day/Night**

### ⚙️ Configuração Total

- Ajuste **quando ser avisado** (5s, 10s, 15s antes...)
- Escolha entre **10+ vozes diferentes**
- Personalize **mensagens dos alertas**
- Configure **tom, velocidade e intensidade** da voz
- Ative/desative alertas individualmente

### 🎨 Interface Moderna

- **2 Temas:** Azul ou Rosa
- **Glassmorphism UI:** Design elegante
- **Auto-save:** Configurações salvas automaticamente
- **Simples de usar:** Interface intuitiva

### 🔒 Privacidade

- Zero telemetria
- 100% local (exceto TTS)
- Código aberto
- Gratuito

---

## 🚀 Como Usar

### Instalação

1. Baixe o executável para seu sistema (Windows/Linux)
2. Execute o arquivo
3. Siga o assistente de configuração
4. Configure sua API key do ElevenLabs (opcional para voz)
5. Pronto!

### Durante o Jogo

Mantenha o Runinhas aberto em segundo plano. Os alertas tocam automaticamente nos timings configurados.

### ElevenLabs (Opcional)

Para ter alertas de voz:
1. Crie conta grátis no [ElevenLabs](https://elevenlabs.io)
2. Copie sua API key
3. Cole no Runinhas em "Configurações de Voz"

**Nota:** Plano gratuito do ElevenLabs tem créditos mensais suficientes para uso casual.

---

## ❓ FAQ

**O Runinhas pode causar ban?**  
Não. Usa apenas GSI (Game State Integration), funcionalidade oficial da Valve.

**Funciona em que modos?**  
Todos: Ranked, Unranked, Turbo, Arcade, Bot matches.

**Precisa de internet?**  
Apenas se usar alertas de voz (ElevenLabs). O resto é 100% local.

**Quanto custa?**  
Gratuito e open source. ElevenLabs tem plano free com créditos mensais.

**Funciona no Linux?**  
Sim, Windows e Linux totalmente suportados.

---

## 📜 Licença

MIT License - Código aberto e gratuito.

---

<div align="center">

Desenvolvido para a comunidade Dota 2 ❤️

**[Reportar Bug](https://github.com/laaridev/runinhas/issues)** • **[Sugerir Feature](https://github.com/laaridev/runinhas/issues)**

</div>
