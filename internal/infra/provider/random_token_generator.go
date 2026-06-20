package provider

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomTokenGenerator struct {
	size int
}

func NewRandomTokenGenerator(size int) RandomTokenGenerator {
	if size <= 0 {
		size = 32
	}
	return RandomTokenGenerator{size: size}
}

func (g RandomTokenGenerator) Generate() (string, error) {
	buf := make([]byte, g.size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
