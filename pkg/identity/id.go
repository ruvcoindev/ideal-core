package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID создаёт детерминированный ID из даты рождения + соли
func GenerateID(birthDate time.Time, salt string) string {
	data := []byte(birthDate.Format("20060102") + salt)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:16]) // 16 байт = 32 hex-символа
}

// GenerateKeyPair генерирует пару ключей Ed25519
func GenerateKeyPair() (PublicKey, PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return PublicKey(pub), PrivateKey(priv), nil
}

// PublicKey — обёртка над ed25519.PublicKey
type PublicKey ed25519.PublicKey

// PrivateKey — обёртка над ed25519.PrivateKey
type PrivateKey ed25519.PrivateKey

// Sign подписывает сообщение приватным ключом
func (pk PrivateKey) Sign(message []byte) []byte {
	return ed25519.Sign(ed25519.PrivateKey(pk), message)
}

// Verify проверяет подпись публичным ключом
func (pk PublicKey) Verify(message, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pk), message, signature)
}

// Sign — пакетная функция для совместимости с тестами
func Sign(priv PrivateKey, message []byte) []byte {
	return ed25519.Sign(ed25519.PrivateKey(priv), message)
}

// Verify — пакетная функция для совместимости с тестами
func Verify(pub PublicKey, message, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pub), message, signature)
}

// ToHex возвращает hex-представление публичного ключа
func (pk PublicKey) ToHex() string {
	return hex.EncodeToString([]byte(pk))
}

// DeriveYggdrasilIP преобразует публичный ключ в IPv6 (формат Yggdrasil)
func DeriveYggdrasilIP(pubKey PublicKey) string {
	// Упрощённая эвристика: первые 16 байт SHA-512 хеша ключа
	hash := sha256.Sum256([]byte(pubKey))
	ip := hash[:16]
	return fmt.Sprintf("200:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7],
		ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15])
}
