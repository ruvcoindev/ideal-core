package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

// HardwareProfile описывает доступные ресурсы
type HardwareProfile struct {
	CPUCores    int
	RAMGB       float64
	HasCUDA     bool
	GPUMemGB    float64
}

// DetectHardware определяет профиль железа
func DetectHardware() HardwareProfile {
	profile := HardwareProfile{
		CPUCores: runtime.NumCPU(),
		RAMGB:    detectRAMGB(),
		HasCUDA:  detectCUDA(),
	}
	if profile.HasCUDA {
		profile.GPUMemGB = detectGPUMemGB()
	}
	return profile
}

func detectRAMGB() float64 {
	// Упрощённо: читаем /proc/meminfo
	// В продакшене: использовать gopsutil
	return 16.0 // заглушка, заменить на реальный парсинг
}

func detectCUDA() bool {
	// Проверка наличия nvidia-smi
	_, err := os.Stat("/usr/bin/nvidia-smi")
	return err == nil
}

func detectGPUMemGB() float64 {
	// Заглушка
	return 0
}

// RecommendedModel возвращает оптимальную модель под железо
func RecommendedModel(hw HardwareProfile) string {
	if hw.HasCUDA && hw.GPUMemGB >= 8 {
		return "qwen2.5:7b" // GPU-режим
	}
	if hw.RAMGB >= 16 && hw.CPUCores >= 8 {
		return "qwen2.5:3b" // CPU, среднее железо
	}
	return "qwen2.5:1.5b" // CPU, консервативный режим
}

// OllamaConfig — конфигурация клиента
type OllamaConfig struct {
	Host        string
	Model       string
	EmbedModel  string // отдельная модель для эмбеддингов
	Timeout     time.Duration
	CPUThreads  int // количество потоков для CPU-инференса
}

// DefaultConfigForHardware создаёт конфиг под текущее железо
func DefaultConfigForHardware() OllamaConfig {
	hw := DetectHardware()
	return OllamaConfig{
		Host:        "http://localhost:11434",
		Model:       RecommendedModel(hw),
		EmbedModel:  "bge-m3",
		Timeout:     120 * time.Second, // дольше для CPU
		CPUThreads:  hw.CPUCores - 2,   // оставить 2 ядра системе
	}
}

// Client — клиент для Ollama
type Client struct {
	config OllamaConfig
	http   *http.Client
}

func NewClient(cfg OllamaConfig) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

// GenerateEmbedding генерирует эмбеддинг
func (c *Client) GenerateEmbedding(text string) ([]float32, error) {
	req := map[string]string{
		"model":  c.config.EmbedModel,
		"prompt": text,
	}
	body, _ := json.Marshal(req)
	resp, err := c.http.Post(c.config.Host+"/api/embeddings", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}
	
	var result struct{ Embedding []float32 }
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Embedding, nil
}

// GenerateText генерирует текст с CPU-оптимизациями
func (c *Client) GenerateText(prompt string) (string, error) {
	req := map[string]interface{}{
		"model":  c.config.Model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"num_thread": c.config.CPUThreads, // используем доступные ядра
			"num_predict": 512, // ограничиваем длину ответа для скорости
		},
	}
	body, _ := json.Marshal(req)
	resp, err := c.http.Post(c.config.Host+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}
	
	var result struct{ Response string }
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Response, nil
}

// IsAvailable проверяет доступность Ollama
func (c *Client) IsAvailable() bool {
	resp, err := c.http.Get(c.config.Host + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
