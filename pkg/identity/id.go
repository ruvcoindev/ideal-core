// Package identity — идентификация пользователей
package identity

import (
	"crypto/rand"
	"encoding/hex"
)

// ID — уникальный идентификатор пользователя
//
// Почему string, а не int:
// - UUID более безопасен (непредсказуем)
// - Не раскрывает порядок создания
// - Легко генерировать распределённо
type ID string

// NewID генерирует новый уникальный идентификатор
//
// Формат: 32-символьный hex (128 бит энтропии)
// Пример: "a3f5c8d2e1b4f7a9c6d8e2f1a5b3c7d9"
func NewID() ID {
	bytes := make([]byte, 16) // 128 бит
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback для тестов
		return ID("00000000000000000000000000000000")
	}
	return ID(hex.EncodeToString(bytes))
}

// String возвращает строковое представление ID
func (id ID) String() string {
	return string(id)
}

// IsValid проверяет валидность ID
func (id ID) IsValid() bool {
	return len(id) == 32 && id != "00000000000000000000000000000000"
}

// Empty возвращает пустой ID
func Empty() ID {
	return ID("")
}
