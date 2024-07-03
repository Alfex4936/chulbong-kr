package service

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/Alfex4936/chulbong-kr/util"

	"github.com/jmoiron/sqlx"
)

type TokenService struct {
	DB        *sqlx.DB
	TokenUtil *util.TokenUtil
	TokenMax  int
}

func NewTokenService(db *sqlx.DB, tokenUtil *util.TokenUtil) *TokenService {
	return &TokenService{
		DB:        db,
		TokenUtil: tokenUtil,
		TokenMax:  2,
	}
}

// GenerateAndSaveToken generates a new token for a user and saves it in the database.
func (s *TokenService) GenerateAndSaveToken(userID int) (string, error) {
	token, err := s.TokenUtil.GenerateOpaqueToken(16) // a secure, random token.
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(s.TokenUtil.Config.TokenExpirationTime) // Use the global duration.

	// Ensure the token limit is enforced
	err = s.EnforceTokenLimit(userID)
	if err != nil {
		return "", err
	}

	err = s.SaveOpaqueToken(userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) EnforceTokenLimit(userID int) error {
	var tokenCount int

	countQuery := `SELECT COUNT(*) FROM OpaqueTokens WHERE UserID = ?`
	err := s.DB.QueryRow(countQuery, userID).Scan(&tokenCount)
	if err != nil {
		return err
	}

	// If the token count is at or above the limit, delete the oldest tokens
	if tokenCount >= s.TokenMax {
		tokensToDelete := tokenCount - s.TokenMax + 1

		// mysql8.0 does not support LIMIT in subqueries for DELETE statements
		deleteQuery := `
        WITH cte AS (
            SELECT TokenID FROM OpaqueTokens WHERE UserID = ? ORDER BY CreatedAt ASC LIMIT ?
        )
        DELETE FROM OpaqueTokens WHERE TokenID IN (SELECT TokenID FROM cte)
        `
		_, err := s.DB.Exec(deleteQuery, userID, tokensToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TokenService) SaveOpaqueToken(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO OpaqueTokens (UserID, OpaqueToken, ExpiresAt) VALUES (?, ?, ?)`
	_, err := s.DB.Exec(query, userID, token, expiresAt)
	return err
}

// SaveOrUpdateOpaqueToken stores a new opaque token in the database
func (s *TokenService) SaveOrUpdateOpaqueToken(userID int, token string, expiresAt time.Time) error {
	// Attempt to update the token if it exists for the user
	query := `
    INSERT INTO OpaqueTokens (UserID, OpaqueToken, ExpiresAt) 
    VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE 
    OpaqueToken=VALUES(OpaqueToken), ExpiresAt=VALUES(ExpiresAt), UpdatedAt=NOW()`
	_, err := s.DB.Exec(query, userID, token, expiresAt)
	return err
}

// DeleteOpaqueToken removes an opaque token from the database for a user
// func (s *TokenService) DeleteOpaqueToken(userID int) error {
// 	query := "DELETE FROM OpaqueTokens WHERE UserID = ?"
// 	_, err := s.DB.Exec(query, userID)
// 	return err
// }

func (s *TokenService) DeleteOpaqueToken(userID int, token string) error {
	query := "DELETE FROM OpaqueTokens WHERE UserID = ? AND OpaqueToken = ?"
	_, err := s.DB.Exec(query, userID, token)
	return err
}

func (s *TokenService) DeleteExpiredTokens() error {
	query := `DELETE FROM OpaqueTokens WHERE ExpiresAt < NOW()`
	_, err := s.DB.Exec(query)
	return err
}

func (s *TokenService) DeleteExpiredPasswordTokens() error {
	query := `DELETE FROM PasswordTokens WHERE ExpiresAt < NOW()`
	_, err := s.DB.Exec(query)
	return err
}

// GenerateToken generates a secure random token that is 6 digits long
func GenerateSixDigitToken() (string, error) {
	// Define the maximum value (999999) for a 6-digit number
	max := big.NewInt(999999)

	// Generate a random number between 0 and max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format the number as a 6-digit string with leading zeros if necessary
	token := fmt.Sprintf("%06d", n.Int64())

	return token, nil
}

func (s *TokenService) GenerateAndSaveSignUpToken(email string) (string, error) {
	token, err := GenerateSixDigitToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	// Attempt to insert or update the token for the user
	_, err = s.DB.Exec(`
        INSERT INTO PasswordTokens (Token, Email, Verified, ExpiresAt, CreatedAt)
        VALUES (?, ?, FALSE, ?, NOW())
        ON DUPLICATE KEY UPDATE Token=VALUES(Token), ExpiresAt=VALUES(ExpiresAt), Verified=FALSE`,
		token, email, expiresAt)
	if err != nil {
		return "", fmt.Errorf("error saving or updating token: %w", err)
	}

	return token, nil
}

func (s *TokenService) ValidateToken(token string, email string) (bool, error) {
	// Start transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var expiresAt time.Time
	err = tx.QueryRow("SELECT ExpiresAt FROM PasswordTokens WHERE Token = ? AND Email = ? AND ExpiresAt > NOW() LIMIT 1", token, email).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Token not found or expired
		}
		return false, err // Database or other error
	}

	// Update the Verified status
	_, err = tx.Exec("UPDATE PasswordTokens SET Verified = TRUE WHERE Token = ? AND Email = ?", token, email)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil // Token is valid, not expired, and now marked as verified
}

func (s *TokenService) IsTokenVerified(email string) (bool, error) {
	var verified bool
	err := s.DB.Get(&verified, "SELECT Verified FROM PasswordTokens WHERE Email = ? AND ExpiresAt > NOW() AND Verified = TRUE LIMIT 1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No verified token found
		}
		return false, err // An error occurred
	}
	return verified, nil // A verified token exists
}
