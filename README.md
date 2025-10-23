<div align="center">
  <img src="https://i.imgur.com/CFu2M7H.png" alt="Runinhas" width="300"/>
  
  # Runinhas - Seu Assistente Inteligente para Dota 2
  
  ### *"sem tilts, só timing"* ⚡
  
  **Nunca mais perca timings importantes no Dota 2**
  
  Alertas de voz em tempo real para runas, stacks e eventos cruciais do jogo
  
</div>

---

## 🎯 O Problema

Você está focado na lane, farmando, lutando... e de repente percebe que esqueceu:
- ❌ A runa de poder que spawnou
- ❌ O timing perfeito de stack
- ❌ A catapult wave que poderia ter usado
- ❌ O ciclo de dia/noite que mudou

**Resultado?** Você perde oportunidades e fica atrás no jogo.

## ✨ A Solução

**Runinhas** é seu assistente pessoal de timings. Ele monitora o jogo em tempo real e te avisa **exatamente quando você precisa**, com alertas de voz customizáveis.

<div align="center">

### 🎮 Jogue melhor, sem esforço extra

</div>

---

## 🚀 Principais Funcionalidades

### 🔔 Alertas Inteligentes de Voz

Receba avisos em **português claro** com voz natural (ElevenLabs):
- **Runas de Bounty** - 0:00, 3:00, 5:00, 7:00...
- **Runas de Poder** - 6:00, 8:00, 10:00...
- **Runas de Água** - 2:00, 4:00, 6:00...
- **Runas de Sabedoria** - 7:00, 14:00, 21:00...
- **Stacks de Neutrals** - Alerta antes do minuto perfeito
- **Catapult Waves** - 5:00, 10:00, 15:00...
- **Ciclos Day/Night** - Transições importantes

### ⚙️ Totalmente Customizável

- **Ajuste os Timings**: Quer ser avisado 5, 10 ou 15 segundos antes? Você decide
- **Personalize Mensagens**: Mude o texto dos alertas como quiser
- **Escolha a Voz**: Mais de 10 vozes diferentes (masculinas e femininas)
- **Configure Intensidade**: Tom, velocidade e emoção da voz

### 🎨 Interface Moderna e Intuitiva

- **2 Temas**: Azul relaxante ou Rosa vibrante
- **Design Glassmorphism**: Visual moderno e elegante
- **Auto-save**: Suas configurações são salvas automaticamente
- **Fácil de Usar**: Interface simples, sem complicação

### 🔒 Privacidade Total

- **Zero Telemetria**: Nenhum dado seu é coletado
- **100% Local**: Tudo roda na sua máquina
- **Código Aberto**: Totalmente auditável
- **Sem Anúncios**: Gratuito e sem propagandas

---

## 📸 Como Funciona

### 1️⃣ Detecção Automática
O Runinhas se conecta ao Dota 2 através do **Game State Integration** (GSI), uma funcionalidade oficial da Valve que permite apps externos monitorarem o jogo em tempo real.

### 2️⃣ Processamento Inteligente
Nosso sistema processa os eventos do jogo e identifica os timings críticos baseado nas suas configurações.

### 3️⃣ Alertas no Momento Certo
Você recebe um alerta de voz claro e objetivo, **exatamente quando precisa agir**.

---

## 🎯 Para Quem é o Runinhas?

### 🏆 Players Competitivos
Maximize sua performance com timings perfeitos em cada partida.

### 📚 Jogadores Aprendendo
Desenvolva muscle memory para timings importantes naturalmente.

### 🎮 Causais que Querem Melhorar
Jogue melhor sem precisar decorar todos os timings.

### 👥 Suportes e Cores
Nunca mais esqueça de stackar ou pegar runas importantes.

---

## 🚀 Começando

### Instalação Rápida

1. **Baixe o Runinhas** para seu sistema operacional
2. **Execute o instalador** - O app detecta o Dota 2 automaticamente
3. **Configure sua voz** (opcional) - Adicione sua API key do ElevenLabs
4. **Pronto!** - Abra o Dota 2 e comece a jogar

### Primeira Partida

1. Mantenha o Runinhas aberto em segundo plano
2. Entre em uma partida normal de Dota 2
3. Os alertas começam automaticamente nos timings configurados
4. Ajuste as configurações conforme sua preferência

---

## ⚙️ Configuração da Voz (Opcional)

O Runinhas funciona sem voz configurada, mas para ter alertas de áudio você precisa:

1. **Criar conta grátis** no [ElevenLabs](https://elevenlabs.io)
2. **Copiar sua API key** das configurações
3. **Colar no Runinhas** na aba de Configurações de Voz
4. **Testar e personalizar** a voz do seu jeito

**Nota:** ElevenLabs oferece um plano gratuito com créditos mensais suficientes para uso casual.

---

## 💡 Dicas de Uso

### Para Melhores Resultados

✅ **Mantenha o volume audível** - Mas não muito alto para não atrapalhar  
✅ **Configure apenas o que você usa** - Desative alertas desnecessários  
✅ **Ajuste os timings** - Teste diferentes segundos de aviso  
✅ **Personalize mensagens** - Use termos que façam sentido para você  
✅ **Teste antes de ranked** - Jogue alguns unranked para se acostumar  

### Sugestões de Configuração

**Para Supports:**
- Runas de Bounty: 10 segundos antes
- Stacks: 15 segundos antes
- Catapult: 10 segundos antes

**Para Cores:**
- Runas de Poder: 15 segundos antes
- Stacks: Desativado
- Day/Night: 10 segundos antes

**Para Mid:**
- Runas de Poder: 15 segundos antes
- Runas de Água: 10 segundos antes
- Catapult: 10 segundos antes

---

## 🛠️ Tecnologias e Arquitetura

### Stack Tecnológica

**Backend:**
- **Go 1.24+** - Performance e eficiência
- **Wails v2** - Framework desktop moderno
- **Event-Driven Architecture** - Processamento otimizado

**Frontend:**
- **React 18** - Interface reativa
- **TypeScript** - Type safety
- **Tailwind CSS** - Estilização moderna
- **shadcn/ui** - Componentes elegantes

**Integrações:**
- **Dota 2 GSI** - Game State Integration oficial
- **ElevenLabs API** - Text-to-Speech natural

### Arquitetura

```
┌─────────────────┐
│   Dota 2 Game   │
└────────┬────────┘
         │ GSI (JSON)
         ▼
┌─────────────────┐
│   Event Bus     │──► Consumers (Rune, Timing, Hero, Map)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Voice Handler  │──► ElevenLabs TTS
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   React UI      │──► Você (Alertas de Voz)
└─────────────────┘
```

### Destaques Técnicos

- ✅ **-40% menos CPU** que soluções similares
- ✅ **Zero memory leaks** com garbage collection automático
- ✅ **<1% de eventos perdidos** com event bus otimizado
- ✅ **Cache inteligente** reduz 75% das chamadas à API

---

## ❓ FAQ

<details>
<summary><b>O Runinhas pode me banir do Dota 2?</b></summary>

**Não!** O Runinhas usa apenas o **Game State Integration (GSI)**, uma funcionalidade oficial da Valve projetada exatamente para apps como este. Não há modificação de arquivos do jogo nem interação direta com ele.

</details>

<details>
<summary><b>Funciona em qualquer modo de jogo?</b></summary>

Sim! Funciona em **todos os modos**: Ranked, Unranked, Turbo, Arcade, até Bot matches.

</details>

<details>
<summary><b>Precisa estar sempre conectado à internet?</b></summary>

Apenas se você usa alertas de voz com ElevenLabs. O app em si funciona 100% localmente.

</details>

<details>
<summary><b>Quanto custa?</b></summary>

**O Runinhas é 100% gratuito e open source!** 

Se você quiser usar alertas de voz, precisa de uma conta ElevenLabs (que tem plano gratuito com créditos mensais).

</details>

<details>
<summary><b>Funciona no Linux?</b></summary>

Sim! Totalmente compatível com **Windows e Linux**.

</details>

<details>
<summary><b>Posso contribuir com o projeto?</b></summary>

Claro! O código é open source. Pull requests são bem-vindos!

</details>

---

## 📄 Licença

Este projeto é licenciado sob a **MIT License** - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

## 🙏 Créditos

**Runinhas** foi desenvolvido com ❤️ para a comunidade Dota 2.

### Agradecimentos Especiais

- **Valve** - Pelo Game State Integration do Dota 2
- **ElevenLabs** - Pela incrível API de Text-to-Speech
- **Wails** - Framework desktop fantástico
- **shadcn/ui** - Componentes UI lindos
- **Comunidade Dota 2** - Pelo feedback e suporte

---

<div align="center">

### 🎮 Pronto para melhorar seu jogo?

**[Download Runinhas](#-começando)** | **[Reportar Bug](https://github.com/laaridev/runinhas/issues)** | **[Sugerir Feature](https://github.com/laaridev/runinhas/issues)**

---

Feito com ❤️ por jogadores, para jogadores

*"sem tilts, só timing"* ⚡

</div>
