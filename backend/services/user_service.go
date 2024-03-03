package services

import (
	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"chulbong-kr/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var TOKEN_DURATION time.Duration

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	// Define the query to select the user
	query := `SELECT UserID, Username, Email, PasswordHash, Provider, ProviderID, CreatedAt, UpdatedAt FROM Users WHERE Email = ?`

	// Execute the query
	err := database.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
		// if err == sql.ErrNoRows {
		// 	// No user found with the provided email
		// 	return nil, fmt.Errorf("no user found with email %s", email)
		// }
		// // An error occurred during the query execution
		// return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	return &user, nil
}

// SaveUser creates a new user with hashed password
func SaveUser(signUpReq *dto.SignUpRequest) (*models.User, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	hashedPassword, err := hashPassword(signUpReq.Password)
	if err != nil {
		return nil, err
	}

	userID, err := insertUserWithRetry(tx, signUpReq, hashedPassword)
	if err != nil {
		return nil, err
	}

	newUser, err := fetchNewUser(tx, userID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM PasswordTokens WHERE Email = ? AND Verified = TRUE", newUser.Email)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error removing verified token: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login checks if a user exists with the given email and password.
func Login(email, password string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT UserID, Username, Email, PasswordHash, Provider FROM Users WHERE Email = ?`
	err := database.DB.Get(user, query, email)
	if err != nil {
		return nil, err // User not found or db error
	}

	// Check if the user was registered through an external provider
	if user.Provider.Valid && user.Provider.String != "website" {
		// The user did not register through the website's traditional sign-up process
		return nil, fmt.Errorf("external provider login not supported here")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		// Password does not match
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// UpdateUserProfile updates the user's profile information
func UpdateUserProfile(user *models.User, newPassword string) error {
	// Check if a new password is provided and needs hashing
	if newPassword != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		// Update the PasswordHash field to have a valid hashed password string
		user.PasswordHash = sql.NullString{String: string(hashedBytes), Valid: true}
	} else if !user.PasswordHash.Valid {
		// If no new password is provided and the existing password is not valid,
		// ensure PasswordHash is an empty, valid sql.NullString to avoid SQL null constraint issues.
		user.PasswordHash = sql.NullString{String: "", Valid: false}
	}
	// Prepare the SQL query
	query := `UPDATE Users SET Username = ?, PasswordHash = ?, UpdatedAt = NOW() WHERE UserID = ?`
	// Execute the query with the user's information
	_, err := database.DB.Exec(query, user.Username, user.PasswordHash, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserWithRelatedData(ctx context.Context, userID int) error {
	// Begin a transaction
	tx, err := database.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	// Define deletion queries for related tables
	// Note: Order matters due to foreign key constraints
	deletionQueries := []string{
		"DELETE FROM OpaqueTokens WHERE UserID = ?",
		"DELETE FROM Comments WHERE UserID = ?",
		"DELETE FROM MarkerLikes WHERE UserID = ?",
		"DELETE FROM MarkerDislikes WHERE UserID = ?",
		"DELETE FROM Photos WHERE MarkerID IN (SELECT MarkerID FROM Markers WHERE UserID = ?)",
		"UPDATE Markers SET UserID = NULL WHERE UserID = ?", // Set UserID to NULL for Markers instead of deleting
		"DELETE FROM Users WHERE UserID = ?",
	}

	// Execute each deletion query within the transaction
	for _, query := range deletionQueries {
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			tx.Rollback() // Attempt to rollback, but don't override the original error
			return fmt.Errorf("executing deletion query (%s): %w", query, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func ResetPassword(token string, newPassword string) error {
	// Start a transaction
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}

	// Ensure the transaction is rolled back if an error occurs
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var userID int
	// Use the transaction (tx) to perform the query
	err = tx.Get(&userID, "SELECT UserID FROM PasswordResetTokens WHERE Token = ? AND ExpiresAt > NOW()", token)
	if err != nil {
		return err // Token not found or expired
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Use the transaction (tx) to update the user's password
	_, err = tx.Exec("UPDATE Users SET PasswordHash = ? WHERE UserID = ?", string(hashedPassword), userID)
	if err != nil {
		return err
	}

	// Use the transaction (tx) to delete the reset token
	_, err = tx.Exec("DELETE FROM PasswordResetTokens WHERE Token = ?", token)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func GeneratePasswordResetToken(email string) (string, error) {
	user := models.User{}
	err := database.DB.Get(&user, "SELECT UserID FROM Users WHERE Email = ?", email)
	if err != nil {
		return "", err // User not found or db error
	}

	token, err := utils.GenerateOpaqueToken()
	if err != nil {
		return "", err
	}

	_, err = database.DB.Exec(`
    INSERT INTO PasswordResetTokens (UserID, Token, ExpiresAt)
    VALUES (?, ?, ?)
    ON DUPLICATE KEY UPDATE Token = VALUES(Token), ExpiresAt = VALUES(ExpiresAt)`,
		user.UserID, token, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	return token, nil
}

// private
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func generateUsername(signUpReq *dto.SignUpRequest) string {
	if signUpReq.Username != nil && *signUpReq.Username != "" {
		return *signUpReq.Username
	}
	emailParts := strings.Split(signUpReq.Email, "@")
	return emailParts[0]
}

func insertUserWithRetry(tx *sqlx.Tx, signUpReq *dto.SignUpRequest, hashedPassword string) (int64, error) {
	username := generateUsername(signUpReq)
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		res, err := tx.Exec(`INSERT INTO Users (Username, Email, PasswordHash, Provider, ProviderID, Role, CreatedAt, UpdatedAt) VALUES (?, ?, ?, ?, ?, 'user', NOW(), NOW())`,
			username, signUpReq.Email, hashedPassword, signUpReq.Provider, signUpReq.ProviderID)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "for key 'idx_users_username'") {
				username = fmt.Sprintf("%s_%s", username, utils.GenerateRandomString(5))
				continue
			}
			return 0, fmt.Errorf("error registering user: %w", err)
		}
		userID, _ := res.LastInsertId()
		return userID, nil
	}
	return 0, fmt.Errorf("failed to insert user after retries")
}

func fetchNewUser(tx *sqlx.Tx, userID int64) (*models.User, error) {
	var newUser models.User
	query := `SELECT UserID, Username, Email, Provider, ProviderID, Role, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`
	err := tx.QueryRowx(query, userID).StructScan(&newUser)
	if err != nil {
		return nil, fmt.Errorf("error fetching newly created user: %w", err)
	}
	return &newUser, nil
}
