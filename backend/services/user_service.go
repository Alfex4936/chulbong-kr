package services

import (
	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var TOKEN_DURATION time.Duration

// SaveUser creates a new user with hashed password
func SaveUser(signUpReq *dto.SignUpRequest) (*models.User, error) {
	var hashedPassword string
	var err error

	// Hash password only for traditional sign-up
	if signUpReq.Password != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(signUpReq.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPassword = string(hashedBytes) // Convert the []byte to a string
	}

	// Generate username from email if not provided
	username := signUpReq.Email
	if signUpReq.Username != nil && *signUpReq.Username != "" {
		username = *signUpReq.Username
	} else {
		emailParts := strings.Split(signUpReq.Email, "@")
		username = emailParts[0]
	}

	// Check if the user is already registered
	var existingUserID int
	checkQuery := `SELECT UserID FROM Users WHERE Email = ? AND (Provider = ? OR Provider IS NULL) LIMIT 1`
	err = database.DB.QueryRow(checkQuery, signUpReq.Email, signUpReq.Provider).Scan(&existingUserID)
	if err == nil {
		return nil, fmt.Errorf("user with email %q is already registered", signUpReq.Email)
	} else if err != sql.ErrNoRows {
		return nil, err // Handle unexpected errors
	}

	// Insert new user into database
	query := `INSERT INTO Users (Username, Email, PasswordHash, Provider, ProviderID, CreatedAt, UpdatedAt) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := database.DB.Exec(query, username, signUpReq.Email, hashedPassword, signUpReq.Provider, signUpReq.ProviderID)
	if err != nil {
		// Handle potential duplicate error
		return nil, fmt.Errorf("error registering user: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the newly created user
	var newUser models.User
	query = `SELECT UserID, Username, Email, Provider, ProviderID, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`
	err = database.DB.QueryRow(query, userID).Scan(&newUser.UserID, &newUser.Username, &newUser.Email, &newUser.Provider, &newUser.ProviderID, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// Login checks if a user exists with the given email and password.
func Login(email, password string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT UserID, Username, Email, PasswordHash FROM Users WHERE Email = ?`
	err := database.DB.Get(user, query, email)
	if err != nil {
		return nil, err // User not found or db error
	}

	// Check if the user was registered through an external provider
	if user.Provider != "website" {
		// The user did not register through the website's traditional sign-up process
		return nil, fmt.Errorf("external provider login not supported here")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Password does not match
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// UpdateUserProfile updates the user's profile information
func UpdateUserProfile(user *models.User, newPassword string) error {
	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	query := `UPDATE Users SET Username = ?, PasswordHash = ?, UpdatedAt = NOW() WHERE UserID = ?`
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
