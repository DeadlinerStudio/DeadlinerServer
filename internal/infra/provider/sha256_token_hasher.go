package provider

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256TokenHasher struct{}

func NewSHA256TokenHasher() SHA256TokenHasher {
	return SHA256TokenHasher{}
}

func (SHA256TokenHasher) Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
