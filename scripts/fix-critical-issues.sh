#!/bin/bash

# Script para corrigir problemas críticos identificados na análise
# Execute com: ./scripts/fix-critical-issues.sh

echo "🔧 Iniciando correção de problemas críticos..."

# 1. Backup dos arquivos que serão modificados
echo "📦 Criando backup..."
mkdir -p backups/$(date +%Y%m%d_%H%M%S)
cp -r backend/events backups/$(date +%Y%m%d_%H%M%S)/
cp -r backend/consumers backups/$(date +%Y%m%d_%H%M%S)/

# 2. Aumentar buffer do Event Bus
echo "🔄 Ajustando buffer do Event Bus..."
sed -i 's/make(chan TickEvent, 50)/make(chan TickEvent, 100)/g' backend/events/bus.go

# 3. Criar diretório para testes
echo "🧪 Criando estrutura de testes..."
mkdir -p backend/tests/{events,consumers,handlers}
mkdir -p frontend/src/__tests__/{components,hooks,services}

# 4. Instalar ferramentas de teste Go
echo "📚 Instalando ferramentas de teste..."
go get -u github.com/stretchr/testify
go get -u github.com/golang/mock/mockgen

# 5. Criar arquivo de métricas básico
echo "📊 Criando estrutura de métricas..."
cat > backend/metrics/metrics.go << 'EOF'
package metrics

import (
    "sync/atomic"
    "time"
)

type Metrics struct {
    EventsProcessed uint64
    EventsDropped   uint64
    CacheHits       uint64
    CacheMisses     uint64
    StartTime       time.Time
}

var Instance = &Metrics{
    StartTime: time.Now(),
}

func (m *Metrics) IncrementProcessed() {
    atomic.AddUint64(&m.EventsProcessed, 1)
}

func (m *Metrics) IncrementDropped() {
    atomic.AddUint64(&m.EventsDropped, 1)
}

func (m *Metrics) GetStats() map[string]interface{} {
    return map[string]interface{}{
        "events_processed": atomic.LoadUint64(&m.EventsProcessed),
        "events_dropped":   atomic.LoadUint64(&m.EventsDropped),
        "cache_hits":       atomic.LoadUint64(&m.CacheHits),
        "cache_misses":     atomic.LoadUint64(&m.CacheMisses),
        "uptime_seconds":   time.Since(m.StartTime).Seconds(),
    }
}
EOF

# 6. Criar script de limpeza de cache
echo "🗑️ Criando script de limpeza de cache..."
cat > scripts/clean-cache.sh << 'EOF'
#!/bin/bash
# Limpa arquivos de cache mais antigos que 7 dias

CACHE_DIRS=(
    "$HOME/.cache/runinhas/voice"
    "$HOME/.config/runinhas/cache"
    "$LOCALAPPDATA/Runinhas/Cache/voice"
)

for dir in "${CACHE_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "Limpando cache em: $dir"
        find "$dir" -type f -mtime +7 -delete
        echo "Cache limpo!"
    fi
done
EOF
chmod +x scripts/clean-cache.sh

# 7. Adicionar git hooks para qualidade
echo "🪝 Configurando git hooks..."
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Run go fmt and vet before commit

echo "Running go fmt..."
gofmt -w backend/

echo "Running go vet..."
go vet ./backend/...

echo "Running frontend lint..."
cd frontend && npm run lint
EOF
chmod +x .git/hooks/pre-commit

# 8. Criar exemplo de teste
echo "✅ Criando exemplo de teste..."
cat > backend/tests/events/bus_test.go << 'EOF'
package events_test

import (
    "testing"
    "time"
    "dota-gsi/backend/events"
)

func TestEventBus_Subscribe(t *testing.T) {
    bus := events.NewEventBus()
    defer bus.Close()
    
    ch := bus.Subscribe()
    if ch == nil {
        t.Fatal("Subscribe returned nil channel")
    }
}

func TestEventBus_Publish(t *testing.T) {
    bus := events.NewEventBus()
    defer bus.Close()
    
    ch := bus.Subscribe()
    
    event := events.TickEvent{
        RawJSON: []byte(`{"test": "data"}`),
        Time:    time.Now(),
    }
    
    go bus.Publish(event)
    
    select {
    case received := <-ch:
        if string(received.RawJSON) != string(event.RawJSON) {
            t.Fatal("Received event does not match published")
        }
    case <-time.After(1 * time.Second):
        t.Fatal("Timeout waiting for event")
    }
}
EOF

echo "✨ Correções aplicadas!"
echo ""
echo "📋 Próximos passos:"
echo "1. Revisar as mudanças no git diff"
echo "2. Executar os testes: go test ./backend/tests/..."
echo "3. Implementar as otimizações restantes do relatório"
echo "4. Configurar CI/CD com os testes"
echo ""
echo "📊 Veja o relatório completo em: ANALYSIS_REPORT.md"
