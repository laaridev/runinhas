#!/bin/bash
# Script para compilar e atualizar release v1.0.0 automaticamente
# Uso: ./scripts/update-release.sh

set -e

echo "🔨 Compilando binários..."
echo ""

# Build Linux
echo "📦 Building Linux binary..."
wails build -clean -platform linux/amd64

if [ -f "build/bin/runinhas" ]; then
    LINUX_SIZE=$(du -h build/bin/runinhas | cut -f1)
    echo "✅ Linux build OK ($LINUX_SIZE)"
else
    echo "❌ Linux build failed"
    exit 1
fi

# Build Windows
echo ""
echo "📦 Building Windows binary..."

if command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then
    # MinGW disponível - build Windows
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
        wails build -clean -platform windows/amd64
    
    if [ -f "build/bin/runinhas.exe" ]; then
        WIN_SIZE=$(du -h build/bin/runinhas.exe | cut -f1)
        echo "✅ Windows build OK ($WIN_SIZE)"
    else
        echo "❌ Windows build failed"
    fi
else
    echo "❌ MinGW não instalado, pulando build Windows"
    echo ""
    echo "   Para compilar Windows no Linux, instale:"
    echo "   sudo pacman -S mingw-w64-gcc"
    echo ""
    echo "   Ou compile diretamente no Windows:"
    echo "   wails build"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Verificar GitHub CLI
if ! command -v gh >/dev/null 2>&1; then
    echo "⚠️  GitHub CLI (gh) não encontrado."
    echo ""
    echo "Instale com:"
    echo "  sudo pacman -S github-cli"
    echo "  gh auth login"
    echo ""
    echo "📁 Binários compilados em: build/bin/"
    echo "   Faça upload manual em:"
    echo "   https://github.com/laaridev/runinhas/releases/tag/v1.0.0"
    exit 0
fi

# Atualizar release
echo "🚀 Atualizando release v1.0.0..."
echo ""

# Deletar assets antigos (se existirem)
echo "🗑️  Removendo binários antigos..."
gh release delete-asset v1.0.0 runinhas --yes 2>/dev/null || true
gh release delete-asset v1.0.0 runinhas.exe --yes 2>/dev/null || true

# Upload novos binários
echo "📤 Fazendo upload dos binários..."

if [ -f "build/bin/runinhas" ]; then
    gh release upload v1.0.0 build/bin/runinhas --clobber
    echo "   ✅ runinhas (Linux)"
fi

if [ -f "build/bin/runinhas.exe" ]; then
    gh release upload v1.0.0 build/bin/runinhas.exe --clobber
    echo "   ✅ runinhas.exe (Windows)"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🎉 Release v1.0.0 atualizada com sucesso!"
echo ""
echo "🔗 Ver release:"
echo "   https://github.com/laaridev/runinhas/releases/tag/v1.0.0"
echo ""
