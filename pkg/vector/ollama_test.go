package vector

import (
	"os"
	"testing"
)

func TestOllamaEmbedding_Integration(t *testing.T) {
	// Пропускаем тест, если Ollama не доступен (CI/CD)
	if os.Getenv("OLLAMA_TEST") != "1" {
		t.Skip("Set OLLAMA_TEST=1 to run Ollama integration test")
	}

	client := NewOllamaEmbeddingClient("http://localhost:11434", "bge-m3")
	
	// Проверка доступности
	if !client.IsAvailable() {
		t.Skip("Ollama server not available at http://localhost:11434")
	}
	
	// Тестовый запрос
	text := "финансовые проблемы, ресурс, долг"
	embedding, err := client.GenerateEmbedding(text)
	if err != nil {
		t.Fatalf("GenerateEmbedding failed: %v", err)
	}
	
	// Проверка размерности (bge-m3 = 1024)
	if len(embedding) != 1024 {
		t.Errorf("Expected embedding length 1024, got %d", len(embedding))
	}
	
	// Проверка, что вектор не нулевой
	var sum float32
	for _, v := range embedding {
		sum += v * v
	}
	if sum < 0.001 {
		t.Error("Embedding appears to be zero vector")
	}
	
	t.Logf("✅ Generated embedding for %q (length: %d)", text, len(embedding))
}
