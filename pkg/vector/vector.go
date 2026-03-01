package vector

import (
	"math"
)

// Embedding — векторное представление текста
type Embedding []float32

// VectorStore — интерфейс для векторной БД
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

// CosineSimilarity — косинусное сходство между векторами
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

// GenerateEmbedding — заглушка: в продакшене использовать локальную модель
// TODO: Интеграция с llama.cpp / ollama для локальной генерации эмбеддингов
// Пример: https://github.com/ollama/ollama/blob/main/docs/api.md#generate-embeddings
func GenerateEmbedding(text string) (Embedding, error) {
	// Временно: возвращаем нулевой вектор фиксированной размерности
	// В продакшене: HTTP-запрос к ollama: POST /api/embeddings {model: "bge-m3", prompt: text}
	dim := 384 // bge-small размерность
	return make(Embedding, dim), nil
}

// SemanticSearch — семантический поиск по тексту (заглушка)
func SemanticSearch(store VectorStore, query string, limit int) ([]Result, error) {
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}
	return store.Search(queryEmbedding, limit), nil
}

// ContextualAnalysis — контекстуальный анализ (заглушка)
// TODO: Интеграция с LLM для извлечения тем, эмоций, намерений
func ContextualAnalysis(messages []string) map[string]interface{} {
	// В продакшене: промпт к локальной LLM:
	// "Extract: topics (list), emotions (list), intent (string), urgency (0-1) from: {messages}"
	return map[string]interface{}{
		"topics":    []string{},
		"emotions":  []string{},
		"intent":    "",
		"urgency":   0.0,
		"note":      "LLM integration pending",
	}
}

// MockVectorStore — реализация для тестов (in-memory)
type MockVectorStore struct {
	vectors map[string]Embedding
	meta    map[string]map[string]interface{}
}

func NewMockVectorStore() *MockVectorStore {
	return &MockVectorStore{
		vectors: make(map[string]Embedding),
		meta:    make(map[string]map[string]interface{}),
	}
}

func (m *MockVectorStore) Upsert(id string, embedding Embedding, metadata map[string]interface{}) error {
	m.vectors[id] = embedding
	m.meta[id] = metadata
	return nil
}

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
	// Сортировка по убыванию сходства (упрощённо)
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

func (m *MockVectorStore) Delete(id string) error {
	delete(m.vectors, id)
	delete(m.meta, id)
	return nil
}
