package auth

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashKey returns a stringified (64 bytes) sha256 checksum of a Key.
// This is sufficient for securely storing Key's in the database,
// see reputable StackExchange answer: https://security.stackexchange.com/a/180364
func HashKey(key Key) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
}

// CheckKeyHash checks whether a plaintext Key is represented by a hash
func CheckKeyHash(key Key, hash string) bool {
	return HashKey(key) == hash
}

// HashPassword returns a stringified password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks whether a plaintext password is represented by hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
