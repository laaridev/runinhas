#!/usr/bin/env python3
"""
Gerador de áudios padrão para versão FREE do Runinhas
Usa Google TTS (gTTS) - 100% gratuito
"""

import os
from pathlib import Path

try:
    from gtts import gTTS
except ImportError:
    print("❌ gTTS não instalado. Instale com: pip install gtts")
    exit(1)

# Mensagens genéricas para versão FREE
# Sem placeholder {seconds} - são fixas
MESSAGES = {
    "bounty_rune_warning.mp3": "Runa de Recompensa em alguns segundos",
    "power_rune_warning.mp3": "Runa de Poder em alguns segundos",
    "wisdom_rune_warning.mp3": "Runa de Sabedoria em alguns segundos",
    "water_rune_warning.mp3": "Runa de Água em alguns segundos",
    "stack_timing_warning.mp3": "Hora de stackar em alguns segundos",
    "catapult_timing_warning.mp3": "Catapulta chegando em alguns segundos",
    "day_night_cycle_warning.mp3": "Mudança de ciclo em alguns segundos",
}

def generate_audio_files():
    # Caminho relativo ao script
    script_dir = Path(__file__).parent
    output_dir = script_dir.parent / "backend" / "assets" / "audio"
    
    # Criar diretório se não existir
    output_dir.mkdir(parents=True, exist_ok=True)
    
    print("🎵 Gerando áudios para versão FREE...")
    print(f"📁 Diretório: {output_dir}")
    print()
    
    success_count = 0
    
    for filename, text in MESSAGES.items():
        try:
            output_path = output_dir / filename
            
            # Gerar áudio com Google TTS
            # lang='pt-BR' para português brasileiro
            # slow=False para velocidade normal
            tts = gTTS(text=text, lang='pt-br', slow=False)
            tts.save(str(output_path))
            
            # Verificar tamanho do arquivo
            size_kb = output_path.stat().st_size / 1024
            
            print(f"✅ {filename:<30} ({size_kb:.1f} KB)")
            success_count += 1
            
        except Exception as e:
            print(f"❌ Erro ao gerar {filename}: {e}")
    
    print()
    print(f"🎉 Concluído! {success_count}/{len(MESSAGES)} arquivos gerados")
    print()
    print("📦 Próximo passo:")
    print("   wails build")
    print("   Os áudios serão automaticamente embutidos no binário via go:embed")

if __name__ == "__main__":
    generate_audio_files()
