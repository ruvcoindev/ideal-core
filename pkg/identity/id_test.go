package identity

import (
	"testing"
	"time"
)

func TestGenerateID_Vitaly(t *testing.T) {
	date := time.Date(1974, 10, 15, 0, 0, 0, 0, time.UTC)
	id := GenerateID(date, "ideal-core-v1")
	
	if len(id) != 32 { // 16 байт = 32 hex-символа
		t.Errorf("ID length: got %d, want 32", len(id))
	}
	t.Logf("Vitaly ID: %s", id)
}

func TestKeyPairGeneration(t *testing.T) {
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Key generation failed: %v", err)
	}
	
	// Тест подписи/верификации
	message := []byte("Теория Идеала: канон зафиксирован")
	signature := Sign(priv, message)
	
	if !Verify(pub, message, signature) {
		t.Error("Signature verification failed")
	}
	
	// Негативный тест: подделка сообщения
	if Verify(pub, []byte("подделка"), signature) {
		t.Error("Fake message passed verification")
	}
}
