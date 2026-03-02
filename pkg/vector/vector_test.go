package vector

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	// identical vectors → similarity = 1
	a := Embedding{1, 2, 3}
	b := Embedding{1, 2, 3}
	sim := CosineSimilarity(a, b)
	if math.Abs(float64(sim)-1.0) > 0.001 {
		t.Errorf("CosineSimilarity identical: got %.3f, want 1.0", sim)
	}
	
	// orthogonal vectors → similarity = 0
	c := Embedding{1, 0, 0}
	d := Embedding{0, 1, 0}
	sim2 := CosineSimilarity(c, d)
	if sim2 != 0 {
		t.Errorf("CosineSimilarity orthogonal: got %.3f, want 0", sim2)
	}
	
	// opposite vectors → similarity = -1
	e := Embedding{1, 0, 0}
	f := Embedding{-1, 0, 0}
	sim3 := CosineSimilarity(e, f)
	if sim3 != -1 {
		t.Errorf("CosineSimilarity opposite: got %.3f, want -1", sim3)
	}
}

func TestMockVectorStore(t *testing.T) {
	store := NewMockVectorStore()
	
	// Upsert
	err := store.Upsert("id1", Embedding{1, 2, 3}, map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("Upsert failed: %v", err)
	}
	
	// Search
	results := store.Search(Embedding{1, 2, 3}, 10)
	if len(results) != 1 {
		t.Errorf("Search: got %d results, want 1", len(results))
	}
	if results[0].ID != "id1" {
		t.Errorf("Search: got ID %q, want %q", results[0].ID, "id1")
	}
	
	// Delete
	err = store.Delete("id1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	
	results = store.Search(Embedding{1, 2, 3}, 10)
	if len(results) != 0 {
		t.Errorf("Search after delete: got %d results, want 0", len(results))
	}
}

func TestGenerateEmbeddingStub(t *testing.T) {
	emb, err := GenerateEmbeddingStub("test text")
	if err != nil {
		t.Errorf("GenerateEmbeddingStub failed: %v", err)
	}
	if len(emb) != 384 {
		t.Errorf("GenerateEmbeddingStub: got length %d, want 384", len(emb))
	}
}
