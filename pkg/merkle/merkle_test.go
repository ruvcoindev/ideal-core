package merkle

import "testing"

func TestBuildTree(t *testing.T) {
	hashes := []string{"a", "b", "c"}
	root := BuildTree(hashes)
	
	if root == nil {
		t.Error("Expected non-nil root")
	}
	if len(root.Hash) != 64 { // SHA256 в hex = 64 символа
		t.Errorf("Expected 64-char hash, got %d", len(root.Hash))
	}
}

func TestHashRelation(t *testing.T) {
	hash := HashRelation("user1", "user2", "MOTHER", 1708473600)
	if len(hash) != 64 {
		t.Errorf("Expected 64-char hash, got %d", len(hash))
	}
	// Детерминизм: одинаковые входы = одинаковый хеш
	hash2 := HashRelation("user1", "user2", "MOTHER", 1708473600)
	if hash != hash2 {
		t.Error("Hash should be deterministic")
	}
}
