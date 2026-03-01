package vector

import (
	"fmt"
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
	if len(a) != len(b) {
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
	return dot / (sqrt(normA) * sqrt(normB))
}

func sqrt(x float32) float32 {
	return float32(fmt.Sprintf("%.4f", math.Sqrt(float64(x))))
}

// GenerateEmbedding — заглушка: в продакшене использовать локальную модель (Qwen3, bge-m3)
func GenerateEmbedding(text string) (Embedding, error) {
	// TODO: Интеграция с llama.cpp / ollama для локальной генерации эмбеддингов
	// Временно: возвращаем случайный вектор (для тестов)
	return make(Embedding, 384), nil // 384 — размерность bge-small
}

// SemanticSearch — семантический поиск по тексту
func SemanticSearch(store VectorStore, query string, limit int) ([]Result, error) {
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}
	return store.Search(queryEmbedding, limit), nil
}

// ContextualAnalysis — контекстуальный анализ (заглушка)
func ContextualAnalysis(messages []string) map[string]interface{} {
	// TODO: Интеграция с LLM для извлечения тем, эмоций, намерений
	return map[string]interface{}{
		"topics":    []string{},
		"emotions":  []string{},
		"intent":    "",
		"urgency":   0.0,
	}
}
