package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	mrand "math/rand"
)

// GenerateOpaqueToken creates a random token
func GenerateOpaqueToken() (string, error) {
	bytes := make([]byte, 16) // Generates a 128-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s := make([]rune, n)
	for i := range s {
		s[i] = rune(letters[mrand.Intn(len(letters))])
	}
	return string(s)
}

func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
