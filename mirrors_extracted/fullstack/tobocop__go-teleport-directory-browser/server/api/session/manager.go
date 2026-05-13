package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

const CookieName = "id"
const ExpiresIn = 2 * time.Hour

type Manager interface {
	NewSession() (string, error)
}


type inMemoryManager struct {
	sessions map[string]time.Time
	mutex *sync.Mutex
}

func NewInMemoryManager() Manager {
	return &inMemoryManager{
		sessions: make(map[string]time.Time),
		mutex: &sync.Mutex{},
	}
}

func (i inMemoryManager) NewSession() (string, error)  {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	sessionId := base64.URLEncoding.EncodeToString(b)
	i.mutex.Lock()
	i.sessions[sessionId] = time.Now().Add(ExpiresIn)
	i.mutex.Unlock()
	return sessionId, nil
}
