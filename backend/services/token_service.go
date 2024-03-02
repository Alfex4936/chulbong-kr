package services

import (
	"chulbong-kr/database"
	"chulbong-kr/middlewares"
	"crypto/rand"

	"encoding/base64"
	"encoding/hex"
	mrand "math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
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

// GenerateAndSaveToken generates a new token for a user and saves it in the database.
func GenerateAndSaveToken(userID int) (string, error) {
	token, err := GenerateOpaqueToken() // a secure, random token.
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(TOKEN_DURATION) // Use the global duration.
	err = SaveOrUpdateOpaqueToken(userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// SaveOrUpdateOpaqueToken stores a new opaque token in the database
func SaveOrUpdateOpaqueToken(userID int, token string, expiresAt time.Time) error {
	// Attempt to update the token if it exists for the user
	query := `
    INSERT INTO OpaqueTokens (UserID, OpaqueToken, ExpiresAt) 
    VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE 
    OpaqueToken=VALUES(OpaqueToken), ExpiresAt=VALUES(ExpiresAt), UpdatedAt=NOW()`
	_, err := database.DB.Exec(query, userID, token, expiresAt)
	return err
}

func DeleteExpiredTokens() error {
	query := `DELETE FROM OpaqueTokens WHERE ExpiresAt < NOW()`
	_, err := database.DB.Exec(query)
	return err
}

func DeleteExpiredPasswordTokens() error {
	query := `DELETE FROM PasswordTokens WHERE ExpiresAt < NOW()`
	_, err := database.DB.Exec(query)
	return err
}

func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GenerateLoginCookie(value string) fiber.Cookie {
	return fiber.Cookie{
		Name:     middlewares.TOKEN_COOKIE,
		Value:    value,                          // The token generated for the user
		Expires:  time.Now().Add(24 * time.Hour), // Set the cookie to expire in 24 hours
		HTTPOnly: true,                           // Ensure the cookie is not accessible through client-side scripts
		Secure:   true,                           // Ensure the cookie is sent over HTTPS
		SameSite: "Lax",                          // Lax, None, or Strict. Lax is a reasonable default
		Path:     "/",                            // Scope of the cookie
	}
}
