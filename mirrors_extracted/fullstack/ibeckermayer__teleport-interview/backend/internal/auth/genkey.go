package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomBytes and generateRandomString are copied from
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// generateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// SessionID is a 32 byte, base64 encoded, cryptographically secure random string
type SessionID string

// newSessionID creates a new SessionID. It will return an error if the system's secure random
// number generator fails to function correctly, in which case the caller should not continue.
func newSessionID() (SessionID, error) {
	s, err := generateRandomString(32)
	return SessionID(s), err
}

// Key is a 32 byte, base64 encoded, cryptographically secure random string
type Key string

// NewKey creates a new API Key. It will return an error if the system's secure random
// number generator fails to function correctly, in which case the caller should not continue.
func NewKey() (Key, error) {
	s, err := generateRandomString(32)
	return Key(s), err
}
