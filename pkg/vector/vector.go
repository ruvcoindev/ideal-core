package vector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// Embedding — векторное представление текста (массив float32)
type Embedding []float32

// ============================================================================
// OLLAMA EMBEDDING CLIENT
// ============================================================================

// OllamaEmbeddingClient — клиент для генерации эмбеддингов через локальный Ollama
type OllamaEmbeddingClient struct {
	host   string
	model  string
	client *http.Client
}

// NewOllamaEmbeddingClient создаёт клиента для работы с Ollama API
func NewOllamaEmbeddingClient(host, model string) *OllamaEmbeddingClient {
	return &OllamaEmbeddingClient{
		host:  host,
		model: model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// GenerateEmbedding генерирует эмбеддинг текста через Ollama API
// Возвращает вектор из 1024 float32 для модели bge-m3
func (c *OllamaEmbeddingClient) GenerateEmbedding(text string) (Embedding, error) {
	req := map[string]string{
		"model":  c.model,
		"prompt": text,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/embeddings", c.host)
	resp, err := c.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Embedding, nil
}

// IsAvailable проверяет, доступен ли Ollama сервер
func (c *OllamaEmbeddingClient) IsAvailable() bool {
	resp, err := c.client.Get(c.host + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ============================================================================
// VECTOR STORE INTERFACE
// ============================================================================

// VectorStore — интерфейс для векторного хранилища
type VectorStore interface {
	Upsert(id string, embedding Embedding, metadata map[string]interface{}) error
	Search(query Embedding, limit int) []Result
	Delete(id string) error
}

// Result — результат векторного поиска
type Result struct {
	ID         string
	Similarity float32
	Metadata   map[string]interface{}
}

// CosineSimilarity вычисляет косинусное сходство между двумя векторами
func CosineSimilarity(a, b Embedding) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// ============================================================================
// MOCK VECTOR STORE (для тестов и разработки)
// ============================================================================

// MockVectorStore — in-memory реализация VectorStore для тестов
type MockVectorStore struct {
	vectors map[string]Embedding
	meta    map[string]map[string]interface{}
}

// NewMockVectorStore создаёт новое in-memory хранилище
func NewMockVectorStore() *MockVectorStore {
	return &MockVectorStore{
		vectors: make(map[string]Embedding),
		meta:    make(map[string]map[string]interface{}),
	}
}

// Upsert добавляет или обновляет вектор в хранилище
func (m *MockVectorStore) Upsert(id string, embedding Embedding, metadata map[string]interface{}) error {
	m.vectors[id] = embedding
	m.meta[id] = metadata
	return nil
}

// Search ищет ближайшие векторы к запросу по косинусному сходству
func (m *MockVectorStore) Search(query Embedding, limit int) []Result {
	var results []Result
	for id, vec := range m.vectors {
		sim := CosineSimilarity(query, vec)
		results = append(results, Result{
			ID:         id,
			Similarity: sim,
			Metadata:   m.meta[id],
		})
	}
	// Сортировка по убыванию сходства
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Similarity > results[i].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}

// Delete удаляет вектор по ID
func (m *MockVectorStore) Delete(id string) error {
	delete(m.vectors, id)
	delete(m.meta, id)
	return nil
}

// ============================================================================
// HELPER FUNCTIONS (для тестов)
// ============================================================================

// GenerateEmbeddingStub — заглушка для тестов (не использует Ollama)
func GenerateEmbeddingStub(text string) (Embedding, error) {
	// Возвращает нулевой вектор фиксированной размерности
	return make(Embedding, 384), nil
}
