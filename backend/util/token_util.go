package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/gofiber/fiber/v2"
)

var (
	// Precompute letters as a byte slice for faster access
	letters    = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	lettersLen = byte(len(letters))

	bytePool = sync.Pool{
		New: func() interface{} {
			// Initialize with a slice of 64 bytes
			buf := make([]byte, 64)
			return &buf
		},
	}
)

type TokenUtil struct {
	Config *config.AppConfig
}

func NewTokenUtil(config *config.AppConfig) *TokenUtil {
	return &TokenUtil{Config: config}
}

// GenerateOpaqueToken creates a random token using a pooled byte slice.
func (t *TokenUtil) GenerateOpaqueToken(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	// Retrieve a byte slice pointer from the pool.
	bPtr := bytePool.Get().(*[]byte)
	defer bytePool.Put(bPtr)

	// Ensure the slice has enough capacity.
	if cap(*bPtr) < length {
		*bPtr = make([]byte, length)
	} else {
		*bPtr = (*bPtr)[:length]
	}

	// Read random bytes using crypto/rand for cryptographic security.
	if _, err := io.ReadFull(rand.Reader, *bPtr); err != nil {
		return "", err
	}

	// Encode the bytes to a base64 URL-encoded string without padding.
	token := base64.RawURLEncoding.EncodeToString(*bPtr)

	return token, nil
}

// GenerateRandomString creates a random string with optimized performance
func (t *TokenUtil) GenerateRandomString(n int) string {
	if n <= 0 {
		return ""
	}

	// Retrieve a byte slice pointer from the pool.
	bPtr := bytePool.Get().(*[]byte)
	defer bytePool.Put(bPtr)

	// Ensure the slice has enough capacity.
	if cap(*bPtr) < n {
		*bPtr = make([]byte, n)
	} else {
		*bPtr = (*bPtr)[:n]
	}

	// Fill the slice with random letters.
	for i := 0; i < n; i++ {
		(*bPtr)[i] = letters[mrand.Intn(int(lettersLen))]
	}

	return string(*bPtr)
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
