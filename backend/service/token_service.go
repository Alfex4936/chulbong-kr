package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/Alfex4936/chulbong-kr/util"
	"go.uber.org/fx"

	"github.com/jmoiron/sqlx"
)

const (
	countTokenQuery = "SELECT COUNT(*) FROM OpaqueTokens WHERE UserID = ? AND ExpiresAt > NOW()"

	// mysql8.0 does not support LIMIT in subqueries for DELETE statements
	deleteTokenQueryOld = `
WITH cte AS (
	SELECT TokenID FROM OpaqueTokens WHERE UserID = ? ORDER BY CreatedAt ASC LIMIT ?
)
DELETE FROM OpaqueTokens WHERE TokenID IN (SELECT TokenID FROM cte)`

	deleteTokenCTEQuery = `
WITH cte AS (
	SELECT TokenID
	FROM OpaqueTokens
	WHERE UserID = ? AND ExpiresAt > NOW()
	ORDER BY CreatedAt ASC
	LIMIT 1
)
DELETE FROM OpaqueTokens
WHERE TokenID IN (SELECT TokenID FROM cte);
`

	insertOpaqueTokenQuery = "INSERT INTO OpaqueTokens (UserID, OpaqueToken, ExpiresAt) VALUES (?, ?, ?)"

	saveOrUpdateOpaqueTokenQuery = `
INSERT INTO OpaqueTokens (UserID, OpaqueToken, ExpiresAt) 
VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE 
OpaqueToken=VALUES(OpaqueToken), ExpiresAt=VALUES(ExpiresAt), UpdatedAt=NOW()`

	deleteOpaqueTokenQuery = "DELETE FROM OpaqueTokens WHERE UserID = ? AND OpaqueToken = ?"

	deleteExpiredOpaqueTokenQuery   = "DELETE FROM OpaqueTokens WHERE ExpiresAt < NOW()"
	deleteExpiredPasswordTokenQuery = "DELETE FROM PasswordTokens WHERE ExpiresAt < NOW()"

	insertOrUpdatePasswordTokenQuery = `
INSERT INTO PasswordTokens (Token, Email, Verified, ExpiresAt, CreatedAt)
VALUES (?, ?, FALSE, ?, NOW())
ON DUPLICATE KEY UPDATE Token=VALUES(Token), ExpiresAt=VALUES(ExpiresAt), Verified=FALSE`

	getValidatedPasswordTokenQuery = "SELECT ExpiresAt FROM PasswordTokens WHERE Token = ? AND Email = ? AND ExpiresAt > NOW() LIMIT 1"

	updatePasswordTokenQuery = "UPDATE PasswordTokens SET Verified = TRUE WHERE Token = ? AND Email = ?"

	getVerifiedPasswordTokenQuery = "SELECT Verified FROM PasswordTokens WHERE Email = ? AND ExpiresAt > NOW() AND Verified = TRUE LIMIT 1"
)

type TokenService struct {
	DB        *sqlx.DB
	TokenUtil *util.TokenUtil

	insertOpaqueTokenStmt *sqlx.Stmt
	deleteOpaqueTokenStmt *sqlx.Stmt
	countOpaqueTokenStmt  *sqlx.Stmt
}

func NewTokenService(db *sqlx.DB, tokenUtil *util.TokenUtil) *TokenService {
	prepareInsertOpaqueTokenStmt, _ := db.Preparex(insertOpaqueTokenQuery)
	prepareDeleteOpaqueTokenStmt, _ := db.Preparex(deleteTokenCTEQuery)
	prepareCountOpaqueTokenStmt, _ := db.Preparex(countTokenQuery)

	return &TokenService{
		DB:        db,
		TokenUtil: tokenUtil,

		insertOpaqueTokenStmt: prepareInsertOpaqueTokenStmt,
		deleteOpaqueTokenStmt: prepareDeleteOpaqueTokenStmt,
		countOpaqueTokenStmt:  prepareCountOpaqueTokenStmt,
	}
}

func RegisterTokenServiceLifecycle(lifecycle fx.Lifecycle, service *TokenService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			_ = service.countOpaqueTokenStmt.Close()
			_ = service.deleteOpaqueTokenStmt.Close()
			_ = service.insertOpaqueTokenStmt.Close()
			return nil
		},
	})
}

// GenerateAndSaveToken generates a new token for a user and saves it in the database.
func (s *TokenService) GenerateAndSaveToken(userID int) (string, error) {
	token, err := s.TokenUtil.GenerateOpaqueToken(s.TokenUtil.Config.TokenLength) // a secure, random token.
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(s.TokenUtil.Config.TokenExpirationTime) // Use the global duration.

	// Start a transaction so insert+delete remain consistent
	tx, err := s.DB.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 1) Insert the new token
	// _, err = tx.Exec(insertOpaqueTokenQuery, userID, token, expiresAt)
	_, err = tx.Stmtx(s.insertOpaqueTokenStmt).Exec(userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	// 2) Check token counts
	var tokenCount int
	// if err := tx.Get(&tokenCount, countTokenQuery, userID); err != nil {
	if err := tx.Stmtx(s.countOpaqueTokenStmt).Get(&tokenCount, userID); err != nil {
		return "", err
	}

	// 3) If the token count is at or above the limit, delete the oldest tokens
	if tokenCount > s.TokenUtil.Config.TokenConcurrentMax {
		tx.Stmtx(s.deleteOpaqueTokenStmt).Exec(userID)
	}

	// Commit everything
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token, nil
}

// func (s *TokenService) EnforceTokenLimit(userID int) error {
// 	var tokenCount int

// 	err := s.DB.QueryRow(countTokenQuery, userID).Scan(&tokenCount)
// 	if err != nil {
// 		return err
// 	}

// 	// If the token count is at or above the limit, delete the oldest tokens
// 	if tokenCount >= s.TokenUtil.Config.TokenConcurrentMax {
// 		tokensToDelete := tokenCount - s.TokenUtil.Config.TokenConcurrentMax + 1

// 		_, err := s.DB.Exec(deleteTokenQuery, userID, tokensToDelete)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (s *TokenService) SaveOpaqueToken(userID int, token string, expiresAt time.Time) error {
	_, err := s.DB.Exec(insertOpaqueTokenQuery, userID, token, expiresAt)
	return err
}

// SaveOrUpdateOpaqueToken stores a new opaque token in the database
func (s *TokenService) SaveOrUpdateOpaqueToken(userID int, token string, expiresAt time.Time) error {
	// Attempt to update the token if it exists for the user
	_, err := s.DB.Exec(saveOrUpdateOpaqueTokenQuery, userID, token, expiresAt)
	return err
}

// DeleteOpaqueToken removes an opaque token from the database for a user
// func (s *TokenService) DeleteOpaqueToken(userID int) error {
// 	query := "DELETE FROM OpaqueTokens WHERE UserID = ?"
// 	_, err := s.DB.Exec(query, userID)
// 	return err
// }

func (s *TokenService) DeleteOpaqueToken(userID int, token string) error {
	_, err := s.DB.Exec(deleteOpaqueTokenQuery, userID, token)
	return err
}

func (s *TokenService) DeleteExpiredTokens() error {
	_, err := s.DB.Exec(deleteExpiredOpaqueTokenQuery)
	return err
}

func (s *TokenService) DeleteExpiredPasswordTokens() error {
	_, err := s.DB.Exec(deleteExpiredPasswordTokenQuery)
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
	_, err = s.DB.Exec(
		insertOrUpdatePasswordTokenQuery,
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
	err = tx.QueryRow(getValidatedPasswordTokenQuery, token, email).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Token not found or expired
		}
		return false, err // Database or other error
	}

	// Update the Verified status
	_, err = tx.Exec(updatePasswordTokenQuery, token, email)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil // Token is valid, not expired, and now marked as verified
}

func (s *TokenService) IsTokenVerified(email string) (bool, error) {
	var verified bool
	err := s.DB.Get(&verified, getVerifiedPasswordTokenQuery, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No verified token found
		}
		return false, err // An error occurred
	}
	return verified, nil // A verified token exists
}
