package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/nacl/sign"
)

const (
	NonceSize = 24
	KeySize   = 32
	Overhead  = 16
)

// EncryptKey шифрует данные с использованием ключа и nonce
func EncryptKey(data, key, nonce []byte) ([]byte, error) {
	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("nonce must be %d bytes", NonceSize)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("key must be %d bytes", KeySize)
	}
	encrypted := secretbox.Seal(nil, data, (*[NonceSize]byte)(nonce), (*[KeySize]byte)(key))
	return encrypted, nil
}

// DecryptKey расшифровывает данные
func DecryptKey(encrypted, key, nonce []byte) ([]byte, error) {
	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("nonce must be %d bytes", NonceSize)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("key must be %d bytes", KeySize)
	}
	decrypted, ok := secretbox.Open(nil, encrypted, (*[NonceSize]byte)(nonce), (*[KeySize]byte)(key))
	if !ok {
		return nil, errors.New("decryption failed: invalid key or nonce, or data tampered")
	}
	return decrypted, nil
}

// GenerateNonce генерирует криптографически безопасный nonce
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// DeriveKeyFromPassword деривирует ключ из пароля (SHA-512 для MVP)
func DeriveKeyFromPassword(password, salt []byte) ([KeySize]byte, error) {
	hash := sha512.Sum512(append(password, salt...))
	var key [KeySize]byte
	copy(key[:], hash[:KeySize])
	return key, nil
}

// ExportEncryptedBackup экспортирует приватный ключ в зашифрованном виде
func (kp *KeyPair) ExportEncryptedBackup(password string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, err
	}
	key, err := DeriveKeyFromPassword([]byte(password), salt)
	if err != nil {
		return nil, err
	}
	encrypted, err := EncryptKey(kp.PrivateKey, key[:], nonce)
	if err != nil {
		return nil, err
	}
	result := append(salt, nonce...)
	result = append(result, encrypted...)
	return result, nil
}

// ImportEncryptedBackup импортирует приватный ключ из зашифрованного бэкапа
func ImportEncryptedBackup(backup []byte, password string) (*KeyPair, error) {
	if len(backup) < 16+NonceSize+Overhead {
		return nil, errors.New("backup too short")
	}
	salt := backup[:16]
	nonce := backup[16 : 16+NonceSize]
	encrypted := backup[16+NonceSize:]
	key, err := DeriveKeyFromPassword([]byte(password), salt)
	if err != nil {
		return nil, err
	}
	privateKey, err := DecryptKey(encrypted, key[:], nonce)
	if err != nil {
		return nil, err
	}
	pubKey := privateKey[32:]
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  pubKey,
	}, nil
}

// SaveEncryptedBackup сохраняет зашифрованный бэкап в файл
func SaveEncryptedBackup(backup []byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, backup, 0600)
}

// LoadEncryptedBackup загружает зашифрованный бэкап из файла
func LoadEncryptedBackup(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// SignMessage подписывает сообщение с использованием Ed25519
func (kp *KeyPair) SignMessage(message []byte) (signed []byte) {
	return sign.Sign(nil, message, (*[64]byte)(kp.PrivateKey))
}

// VerifyMessage проверяет подпись сообщения
func VerifyMessage(signed []byte, publicKey []byte) (message []byte, ok bool) {
	var pub [32]byte
	copy(pub[:], publicKey)
	return sign.Open(nil, signed, &pub)
}

// EncryptForRecipient шифрует сообщение для получателя (упрощённо: симметричное)
func EncryptForRecipient(message, sharedKey []byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, err
	}
	encrypted, err := EncryptKey(message, sharedKey, nonce)
	if err != nil {
		return nil, err
	}
	return append(nonce, encrypted...), nil
}

// DecryptFromSender расшифровывает сообщение от отправителя
func DecryptFromSender(payload, sharedKey []byte) ([]byte, error) {
	if len(payload) < NonceSize {
		return nil, errors.New("payload too short")
	}
	nonce := payload[:NonceSize]
	encrypted := payload[NonceSize:]
	return DecryptKey(encrypted, sharedKey, nonce)
}

// DeriveSharedKey деривирует общий ключ для симметричного шифрования (упрощённо)
func DeriveSharedKey(pubKey1, pubKey2 []byte) [KeySize]byte {
	combined := append(pubKey1, pubKey2...)
	hash := sha512.Sum512(combined)
	var key [KeySize]byte
	copy(key[:], hash[:KeySize])
	return key
}
