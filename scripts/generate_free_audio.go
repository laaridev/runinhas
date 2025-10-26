package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	API_KEY  = "sk_90e96a28a9713070de866b56ad41c849242ecba7a6ffbd3d"
	VOICE_ID = "eVXYtPVYB9wDoz9NVTIy"
)

// Mensagens genéricas para versão FREE (sem placeholder {seconds})
var messages = map[string]string{
	"bounty_rune_warning.mp3":        "Runa de Recompensa em alguns segundos",
	"power_rune_warning.mp3":         "Runa de Poder em alguns segundos",
	"wisdom_rune_warning.mp3":        "Runa de Sabedoria em alguns segundos",
	"water_rune_warning.mp3":         "Runa de Água em alguns segundos",
	"stack_timing_warning.mp3":       "Hora de stackar em alguns segundos",
	"catapult_timing_warning.mp3":    "Catapulta chegando em alguns segundos",
	"day_night_cycle_warning.mp3":    "Mudança de ciclo em alguns segundos",
}

type ElevenLabsRequest struct {
	Text          string                 `json:"text"`
	ModelID       string                 `json:"model_id"`
	VoiceSettings map[string]interface{} `json:"voice_settings"`
}

func generateAudio(text, filename string) error {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", VOICE_ID)

	payload := ElevenLabsRequest{
		Text:    text,
		ModelID: "eleven_multilingual_v2",
		VoiceSettings: map[string]interface{}{
			"stability":         0.5,
			"similarity_boost":  0.75,
			"style":             0.0,
			"use_speaker_boost": true,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", API_KEY)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ElevenLabs API error %d: %s", resp.StatusCode, string(body))
	}

	// Ler o áudio
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Salvar arquivo
	if err := os.WriteFile(filename, audioData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("🎵 Gerando áudios para versão FREE com ElevenLabs...")
	fmt.Println()

	// Determinar diretório de saída
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Erro ao obter diretório atual: %v\n", err)
		os.Exit(1)
	}

	// Se estamos em scripts/, subir um nível
	if filepath.Base(scriptDir) == "scripts" {
		scriptDir = filepath.Dir(scriptDir)
	}

	outputDir := filepath.Join(scriptDir, "backend", "assets", "audio")

	// Criar diretório se não existir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Erro ao criar diretório: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📁 Diretório: %s\n", outputDir)
	fmt.Printf("🎤 Voice ID: %s\n", VOICE_ID)
	fmt.Println()

	successCount := 0
	totalFiles := len(messages)

	for filename, text := range messages {
		outputPath := filepath.Join(outputDir, filename)

		fmt.Printf("⏳ Gerando: %s...", filename)

		if err := generateAudio(text, outputPath); err != nil {
			fmt.Printf(" ❌ ERRO: %v\n", err)
			continue
		}

		// Obter tamanho do arquivo
		fileInfo, err := os.Stat(outputPath)
		if err != nil {
			fmt.Printf(" ❌ ERRO ao ler arquivo: %v\n", err)
			continue
		}

		sizeKB := float64(fileInfo.Size()) / 1024
		fmt.Printf(" ✅ (%.1f KB)\n", sizeKB)
		successCount++

		// Pequeno delay para não sobrecarregar a API
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf("🎉 Concluído! %d/%d arquivos gerados\n", successCount, totalFiles)
	fmt.Println()
	fmt.Println("📦 Próximo passo:")
	fmt.Println("   wails build")
	fmt.Println("   Os áudios serão automaticamente embutidos no binário via go:embed")
}
