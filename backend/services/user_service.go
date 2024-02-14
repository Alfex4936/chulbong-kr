package services

import (
	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// SignUp creates a new user with hashed password
func SignUp(signUpReq *dto.SignUpRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate username from email if not provided
	username := signUpReq.Email
	if signUpReq.Username != nil && *signUpReq.Username != "" {
		username = *signUpReq.Username
	} else {
		emailParts := strings.Split(signUpReq.Email, "@")
		if len(emailParts) > 0 {
			username = emailParts[0]
		}
	}

	// Check if the username is already taken
	var existingUserID int
	checkQuery := `SELECT UserID FROM Users WHERE Username = ? LIMIT 1`
	err = database.DB.QueryRow(checkQuery, username).Scan(&existingUserID)
	if err == nil {
		return nil, errors.New("username already taken")
	} else if err != sql.ErrNoRows {
		// Handle unexpected errors
		return nil, err
	}

	// Proceed with user creation if the username is not taken
	query := `INSERT INTO Users (Username, Email, PasswordHash, CreatedAt, UpdatedAt) VALUES (?, ?, ?, NOW(), NOW())`
	res, err := database.DB.Exec(query, username, signUpReq.Email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// Check if the error message contains 'Email' to determine it's a duplicate email error
			if strings.Contains(mysqlErr.Message, "Email") {
				return nil, fmt.Errorf("%q is already registered", signUpReq.Email)
			}
		}
		return nil, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the newly created user from the database
	var newUser models.User
	query = `SELECT UserID, Username, Email, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`
	err = database.DB.QueryRow(query, userID).Scan(&newUser.UserID, &newUser.Username, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Optionally clear the password hash for security before returning
	// newUser.PasswordHash = ""

	return &newUser, nil
}

// Login checks if a user exists with the given email and password
func Login(email, password string) (*models.User, string, error) {
	user := &models.User{}
	query := `SELECT UserID, Username, Email, PasswordHash FROM Users WHERE Email = ?`
	err := database.DB.Get(user, query, email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Password does not match
		return nil, "", err
	}

	token, err := GenerateOpaqueToken()
	if err != nil {
		return nil, "", err
	}

	expiresAt := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	err = SaveOrUpdateOpaqueToken(user.Email, token, expiresAt)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
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
