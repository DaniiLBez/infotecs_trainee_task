package hasher

import (
	"crypto/sha256"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{salt: salt}
}

func (h *SHA256Hasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
