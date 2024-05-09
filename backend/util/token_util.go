package util

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand/v2"
	"time"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/gofiber/fiber/v2"
)

type TokenUtil struct {
	Config *config.AppConfig
}

func NewTokenUtil(config *config.AppConfig) *TokenUtil {
	return &TokenUtil{Config: config}
}

// GenerateOpaqueToken creates a random token
func (t *TokenUtil) GenerateOpaqueToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (t *TokenUtil) GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s := make([]rune, n)
	for i := range s {
		s[i] = rune(letters[mrand.IntN(len(letters))])
	}
	return string(s)
}

func (t *TokenUtil) GenerateLoginCookie(value string) fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     t.Config.LoginTokenCookie,
		Value:    value,                              // The token generated for the user
		Expires:  time.Now().Add(24 * 7 * time.Hour), // Set the cookie to expire in 7 days
		HTTPOnly: true,                               // Ensure the cookie is not accessible through client-side scripts
		Secure:   true,                               // Ensure the cookie is sent over HTTPS
		SameSite: "Lax",                              // Lax, None, or Strict. Lax is a reasonable default
		Path:     "/",                                // Scope of the cookie
	}

	if t.Config.IsProduction == "production" {
		cookie.Domain = ".k-pullup.com" // Allow cookie to be shared across all subdomains
	}
	return cookie
}

func (t *TokenUtil) ClearLoginCookie() fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     t.Config.LoginTokenCookie,
		Value:    "",                         // The token generated for the user
		Expires:  time.Now().Add(-time.Hour), // Set the cookie to be expired
		HTTPOnly: true,                       // Ensure the cookie is not accessible through client-side scripts
		Secure:   true,                       // Ensure the cookie is sent over HTTPS
		SameSite: "Lax",                      // Lax, None, or Strict. Lax is a reasonable default
		Path:     "/",                        // Scope of the cookie
	}

	if t.Config.IsProduction == "production" {
		cookie.Domain = ".k-pullup.com" // Allow cookie to be shared across all subdomains
	}
	return cookie
}
