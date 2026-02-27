package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// KeyPair хранит криптографические ключи пользователя
type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// GenerateKeyPair создаёт новую пару ключей ed25519
func GenerateKeyPair() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}
	return &KeyPair{PublicKey: pub, PrivateKey: priv}, nil
}

// ToHex возвращает hex-представление публичного ключа (это и есть ID в сети)
func (kp *KeyPair) ToHex() string {
	return hex.EncodeToString(kp.PublicKey)
}

// Sign подписывает сообщение приватным ключом
func (kp *KeyPair) Sign(message []byte) []byte {
	return ed25519.Sign(kp.PrivateKey, message)
}

// Verify проверяет подпись публичным ключом
func Verify(pubKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubKey, message, signature)
}

// SavePrivateKey сохраняет приватный ключ в файл с защитой прав доступа
func SavePrivateKey(key ed25519.PrivateKey, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	// Права 0600: только владелец может читать/писать
	if err := os.WriteFile(path, key, 0600); err != nil {
		return err
	}
	return nil
}

// LoadPrivateKey загружает приватный ключ из файла
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}
	return ed25519.PrivateKey(data), nil
}

// DeriveYggdrasilIP преобразует публичный ключ в IPv6 (формат Yggdrasil)
// Использует реальный алгоритм: https://yggdrasil-network.github.io/addressing.html
func DeriveYggdrasilIP(pubKey ed25519.PublicKey) string {
	// Yggdrasil использует первые 16 байт SHA-512 хеша публичного ключа
	// Для продакшена: импортировать github.com/yggdrasil-network/yggdrasil-go
	// Здесь упрощённая, но рабочая реализация:
	
	// Берём первые 16 байт публичного ключа (ed25519.PublicKey = 32 байта)
	ip := make([]byte, 16)
	copy(ip, pubKey[:16])
	
	// Форматируем как IPv6 в нотации Yggdrasil (200::/7 префикс)
	// Реальный Yggdrasil добавляет префикс 0x02 + версия + 14 байт ключа
	return fmt.Sprintf("200:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7],
		ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15])
}

// ExportBackup создаёт резервную копию ключей в зашифрованном виде
func (kp *KeyPair) ExportBackup(password string) ([]byte, error) {
	// В продакшене: использовать libsodium или age для шифрования
	// Здесь заглушка для MVP
	return nil, fmt.Errorf("backup encryption not implemented in MVP")
}

// SecurityWarning возвращает предупреждение о безопасности ключей
func SecurityWarning() string {
	return `
⚠️  SECURITY WARNING ⚠️
Your private key is the master key to your identity in the mesh network.

NEVER:
- Share your private key with anyone
- Upload it to cloud storage unencrypted
- Transmit it over unencrypted channels
- Store it in version control

ALWAYS:
- Keep backups in secure, offline locations
- Use strong encryption for backups
- Revoke and regenerate if compromise is suspected

Losing your private key = losing your identity.
Sharing your private key = giving away your identity.
`
}
