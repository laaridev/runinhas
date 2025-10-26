#!/usr/bin/env bash
set -euo pipefail

# Build script for runinhas - Linux and Windows binaries
# Creates single-file executables with embedded frontend and backend

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
cd "$ROOT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Add Go bin to PATH for wails
export PATH="$HOME/go/bin:$PATH"

echo -e "${BLUE}🚀 Building runinhas for Linux and Windows${NC}"
echo "================================================"

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo -e "${RED}❌ Wails not found. Please install it first:${NC}"
    echo "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Clean old builds
echo -e "${YELLOW}🧹 Cleaning old builds...${NC}"
rm -rf build/
rm -rf frontend/dist/

# Install frontend dependencies
echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
cd frontend
npm install
cd ..

# Build frontend
echo -e "${BLUE}🎨 Building frontend...${NC}"
cd frontend
npx vite build
cd ..

# Build Linux binary
echo -e "${GREEN}🐧 Building Linux binary...${NC}"
wails build \
    -platform linux/amd64 \
    -tags desktop,production \
    -s \
    -o runinhas-linux

# Build Windows binary
echo -e "${GREEN}🪟 Building Windows binary...${NC}"
wails build \
    -platform windows/amd64 \
    -tags desktop,production \
    -s \
    -o runinhas.exe

# Show build results
echo ""
echo -e "${GREEN}✅ Build complete!${NC}"
echo "================================================"
echo -e "${BLUE}📁 Build artifacts:${NC}"
ls -lh build/bin/ | grep -E "runinhas"

# Calculate sizes
LINUX_SIZE=$(du -h build/bin/runinhas-linux | cut -f1)
WINDOWS_SIZE=$(du -h build/bin/runinhas.exe | cut -f1)

echo ""
echo -e "${BLUE}📊 Binary sizes:${NC}"
echo -e "  Linux:   ${GREEN}$LINUX_SIZE${NC} (build/bin/runinhas-linux)"
echo -e "  Windows: ${GREEN}$WINDOWS_SIZE${NC} (build/bin/runinhas.exe)"

echo ""
echo -e "${YELLOW}📝 Features:${NC}"
echo "  • Single executable file"
echo "  • No external dependencies"
echo "  • Embedded frontend and backend"
echo "  • Auto-creates config in system directories:"
echo "    - Linux: ~/.config/runinhas/"
echo "    - Windows: %APPDATA%\\Runinhas\\"
echo "  • Fixed window size: 850x620"
echo "  • No console window on Windows"

echo ""
echo -e "${GREEN}🎉 Ready for distribution!${NC}"
