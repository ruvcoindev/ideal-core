package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	UserID     string
	CreatedAt  time.Time
	LastActive time.Time
}

type AuthManager struct {
	sessions       map[string]*Session
	mu             sync.RWMutex
	sessionTimeout time.Duration
}

func NewAuthManager(timeout time.Duration) *AuthManager {
	return &AuthManager{
		sessions:       make(map[string]*Session),
		sessionTimeout: timeout,
	}
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (am *AuthManager) Login(userID string) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}
	am.mu.Lock()
	defer am.mu.Unlock()
	now := time.Now()
	am.sessions[token] = &Session{
		UserID:     userID,
		CreatedAt:  now,
		LastActive: now,
	}
	return token, nil
}

func (am *AuthManager) Validate(token string) (*Session, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	session, exists := am.sessions[token]
	if !exists {
		return nil, false
	}
	if time.Since(session.LastActive) > am.sessionTimeout {
		return nil, false
	}
	return session, true
}

func (am *AuthManager) Logout(token string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.sessions, token)
}
