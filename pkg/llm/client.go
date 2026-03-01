package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaConfig — конфигурация клиента
type OllamaConfig struct {
	Host        string        // "http://localhost:11434"
	Model       string        // "bge-m3" для эмбеддингов, "qwen3" для генерации
	Timeout     time.Duration // таймаут запросов
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() OllamaConfig {
	return OllamaConfig{
		Host:    "http://localhost:11434",
		Model:   "bge-m3",
		Timeout: 30 * time.Second,
	}
}

// Client — клиент для взаимодействия с Ollama API
type Client struct {
	config OllamaConfig
	http   *http.Client
}

// NewClient создаёт нового клиента
func NewClient(cfg OllamaConfig) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// EmbeddingRequest — запрос на генерацию эмбеддинга
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse — ответ с эмбеддингом
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerateEmbedding генерирует векторное представление текста через Ollama
func (c *Client) GenerateEmbedding(text string) ([]float32, error) {
	req := EmbeddingRequest{
		Model:  c.config.Model,
		Prompt: text,
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/embeddings", c.config.Host)
	resp, err := c.http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}
	
	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	
	return result.Embedding, nil
}

// GenerateRequest — запрос на генерацию текста
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse — ответ с генерацией
type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// GenerateText генерирует текст через Ollama (для анализа, рекомендаций)
func (c *Client) GenerateText(prompt string) (string, error) {
	req := GenerateRequest{
		Model:  "qwen3", // или другая модель для генерации
		Prompt: prompt,
		Stream: false,
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/generate", c.config.Host)
	resp, err := c.http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}
	
	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	
	return result.Response, nil
}

// AnalyzeContext анализирует переписку и извлекает инсайты
func (c *Client) AnalyzeContext(messages []string) (map[string]interface{}, error) {
	prompt := fmt.Sprintf(`Проанализируй переписку и извлеки:
1. Темы (список строк)
2. Эмоции (список строк)
3. Намерение автора (строка)
4. Срочность (число 0-1)

Переписка:
%s

Ответь в формате JSON.`, messages)
	
	response, err := c.GenerateText(prompt)
	if err != nil {
		return nil, err
	}
	
	// Парсим JSON-ответ (упрощённо)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return map[string]interface{}{
			"raw_response": response,
			"parse_error":  err.Error(),
		}, nil
	}
	
	return result, nil
}

// IsAvailable проверяет доступность Ollama
func (c *Client) IsAvailable() bool {
	resp, err := c.http.Get(c.config.Host + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
