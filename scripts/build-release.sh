#!/bin/bash

# Script para build local e preparação de release
# Uso: ./scripts/build-release.sh v1.0.0

set -e

VERSION=${1:-v1.0.0}
BUILD_DIR="build/bin"
RELEASE_DIR="release"

echo "🚀 Building Runinhas $VERSION"
echo ""

# Limpar builds anteriores
echo "🧹 Cleaning previous builds..."
rm -rf build/
rm -rf $RELEASE_DIR
mkdir -p $RELEASE_DIR

# Build para Linux
echo "🐧 Building for Linux..."
wails build -clean -platform linux/amd64

# Renomear com versão
cp $BUILD_DIR/runinhas $RELEASE_DIR/runinhas-linux-amd64
echo "✅ Linux build: $RELEASE_DIR/runinhas-linux-amd64"

# Build para Windows
echo "🪟 Building for Windows..."
wails build -clean -platform windows/amd64

# Renomear com versão
cp $BUILD_DIR/runinhas.exe $RELEASE_DIR/runinhas-windows-amd64.exe
echo "✅ Windows build: $RELEASE_DIR/runinhas-windows-amd64.exe"

# Gerar checksums
echo ""
echo "🔐 Generating checksums..."
cd $RELEASE_DIR

sha256sum runinhas-linux-amd64 > runinhas-linux-amd64.sha256
sha256sum runinhas-windows-amd64.exe > runinhas-windows-amd64.exe.sha256

echo "✅ Checksums generated"
echo ""

# Listar arquivos
echo "📦 Release files ready:"
ls -lh

cd ..

echo ""
echo "🎉 Build complete!"
echo ""
echo "📋 Next steps:"
echo "1. Go to: https://github.com/laaridev/runinhas/releases/new"
echo "2. Choose tag: $VERSION"
echo "3. Upload files from: $RELEASE_DIR/"
echo "4. Publish release"
echo ""
echo "Files to upload:"
echo "  - runinhas-linux-amd64"
echo "  - runinhas-linux-amd64.sha256"
echo "  - runinhas-windows-amd64.exe"
echo "  - runinhas-windows-amd64.exe.sha256"
