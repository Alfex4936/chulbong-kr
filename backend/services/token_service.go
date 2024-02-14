package services

import (
	"chulbong-kr/database"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateOpaqueToken creates a random token
func GenerateOpaqueToken() (string, error) {
	bytes := make([]byte, 16) // Generates a 128-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SaveOrUpdateOpaqueToken stores a new opaque token in the database
func SaveOrUpdateOpaqueToken(email, token string, expiresAt time.Time) error {
	// Attempt to update the token if it exists for the user
	query := `
    INSERT INTO OpaqueTokens (Email, OpaqueToken, ExpiresAt) 
    VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE 
    OpaqueToken=VALUES(OpaqueToken), ExpiresAt=VALUES(ExpiresAt), UpdatedAt=NOW()`
	_, err := database.DB.Exec(query, email, token, expiresAt)
	return err
}

func DeleteExpiredTokens() error {
	query := `DELETE FROM OpaqueTokens WHERE ExpiresAt < NOW()`
	_, err := database.DB.Exec(query)
	return err
}
