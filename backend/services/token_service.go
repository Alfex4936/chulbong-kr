package services

import (
	"chulbong-kr/database"
	"chulbong-kr/utils"

	"time"
)

// GenerateAndSaveToken generates a new token for a user and saves it in the database.
func GenerateAndSaveToken(userID int) (string, error) {
	token, err := utils.GenerateOpaqueToken() // a secure, random token.
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

// DeleteOpaqueToken removes an opaque token from the database for a user
func DeleteOpaqueToken(userID int) error {
	query := "DELETE FROM OpaqueTokens WHERE UserID = ?"
	_, err := database.DB.Exec(query, userID)
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
