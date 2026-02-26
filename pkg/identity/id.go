package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// PublicKey представляет публичный ключ пользователя (Ed25519, 32 байта)
type PublicKey [32]byte

// PrivateKey представляет приватный ключ (Ed25519, 64 байта)
type PrivateKey [64]byte

// GenerateKeyPair генерирует пару ключей для нового пользователя
func GenerateKeyPair() (PublicKey, PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return PublicKey{}, PrivateKey{}, err
	}
	var pk PublicKey
	var sk PrivateKey
	copy(pk[:], pub)
	copy(sk[:], priv)
	return pk, sk, nil
}

// GenerateID создаёт детерминированный ID из даты рождения + соли
// Используется для публичной идентификации без раскрытия ключей
func GenerateID(birthDate time.Time, salt string) string {
	data := []byte(birthDate.Format("20060102") + salt)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:16]) // 128-bit ID
}

// Sign подписывает сообщение приватным ключом
func Sign(priv PrivateKey, message []byte) []byte {
	return ed25519.Sign(priv[:], message)
}

// Verify проверяет подпись публичным ключом
func Verify(pub PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pub[:], message, signature)
}

// UserID хранит публичные данные пользователя
type UserID struct {
	ID        string    // Сгенерированный ID (GenerateID)
	PublicKey PublicKey // Для крипто-операций
	BirthDate time.Time // Для расчёта векторов
	CreatedAt time.Time // Время регистрации в системе
}
