# 🚀 Guia de Configuração do GitHub

## Passos para Configurar o Repositório

### 1. Criar Repositório no GitHub

```bash
# Criar repositório no GitHub primeiro (via web)
# Depois conectar local:
git init
git add .
git commit -m "feat: initial commit with CI/CD setup"
git branch -M main
git remote add origin https://github.com/SEU_USUARIO/dota-gsi.git
git push -u origin main
```

### 2. Configurar Branch Protection (Opcional mas Recomendado)

No GitHub, vá em:
- **Settings** → **Branches** → **Add rule**
- Branch name pattern: `main`
- ✅ Require status checks to pass before merging
  - Selecione: `Lint Backend`, `Lint Frontend`, `Build Test`

### 3. Ativar GitHub Actions

As Actions devem ativar automaticamente no primeiro push. Verifique em:
- **Actions** tab do repositório

Se necessário, vá em **Settings** → **Actions** → **General**:
- ✅ Allow all actions and reusable workflows

### 4. Configurar Secrets (Opcional)

Se quiser usar Codecov ou Semgrep Cloud:

**Settings** → **Secrets and variables** → **Actions** → **New repository secret**

- `CODECOV_TOKEN` - Token do Codecov (opcional)
- `SEMGREP_APP_TOKEN` - Token do Semgrep (opcional, melhora relatórios)

> **Nota:** Ambos funcionam sem tokens em repos públicos, mas com limitações.

### 5. Ativar Security Features

**Settings** → **Security** → **Code security and analysis**

Ative:
- ✅ Dependency graph
- ✅ Dependabot alerts
- ✅ Dependabot security updates
- ✅ Code scanning (CodeQL já está configurado no workflow)

### 6. Criar Primeira Release

#### Opção A: Via Tag (Automático)

```bash
# Crie uma tag semântica
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# O workflow release.yml vai:
# 1. Buildar para Windows e Linux
# 2. Criar checksums SHA256
# 3. Criar release automática no GitHub
```

#### Opção B: Manual (Draft)

1. Vá em **Releases** → **Create a new release**
2. Escolha uma tag: `v1.0.0` (criar nova)
3. Release title: `v1.0.0 - First Release`
4. Descreva as features
5. Salve como **Draft** primeiro
6. Push da tag para rodar o build

### 7. Ajustar URLs no README

Substitua `seu-usuario` pelo seu username do GitHub em:
- Badges no topo
- Links de download
- Links de clone

**Comando rápido:**
```bash
# Substitua SEU_USER pelo seu username
sed -i 's/seu-usuario/SEU_USER/g' README.md
```

## ✅ Checklist Final

Antes de tornar o repo público, verifique:

- [ ] Todos os workflows passando (Actions tab verde)
- [ ] README atualizado com seu username
- [ ] LICENSE presente
- [ ] .gitignore configurado (já está)
- [ ] Secrets sensíveis não commitados (API keys, etc)
- [ ] Primeira release criada
- [ ] Security tab configurada

## 🎯 Estrutura de Tags/Releases

Use **Semantic Versioning**:

- `v1.0.0` - Primeira release estável
- `v1.0.1` - Patch (bugfixes)
- `v1.1.0` - Minor (novas features)
- `v2.0.0` - Major (breaking changes)

Para pré-releases:
- `v1.0.0-beta.1`
- `v1.0.0-rc.1`

## 📊 Verificar Saúde do Projeto

Após algumas horas/dias:

1. **CodeQL**: Settings → Security → Code scanning alerts
2. **Dependabot**: Security → Dependabot alerts
3. **Actions**: Actions tab (todos devem estar verdes)
4. **Badges**: README badges devem mostrar status correto

## 🔧 Troubleshooting

### Actions não rodam
- Verifique Settings → Actions → permitir workflows
- Re-push para forçar

### CodeQL/Semgrep falhando
- Normal na primeira vez
- Aguarde ~5-10 minutos para análise completa
- Verifique logs em Actions tab

### Build release falha
- Verifique dependências: Go 1.24, Node 20
- Logs detalhados em Actions → Release Build
- Pode precisar ajustar paths no workflow

### SmartScreen bloqueia exe
- Normal para novos apps
- Opções: Code signing (caro) ou instruir usuários
- Já documentado no README com aviso

## 🎉 Pronto!

Agora seu repo está profissional com:
- ✅ CI/CD completo
- ✅ Security scanning
- ✅ Releases automatizadas
- ✅ Badges bonitões
- ✅ Zero burocracia de issues/PRs (você que mantém)

Qualquer dúvida, os workflows tem comentários explicativos!
