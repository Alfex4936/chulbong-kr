package repository

import (
	"chulbong-kr/database"
	"chulbong-kr/models"
	"database/sql"
	"fmt"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetUserById retrieves a user by their email address
func (u *UserRepository) GetUserById(userID int) (*models.User, error) {
	var user models.User

	// Define the query to select the user
	query := `SELECT UserID, Username, Email, PasswordHash, Provider, ProviderID, CreatedAt, UpdatedAt FROM Users WHERE UserID = ?`

	// Execute the query
	err := database.DB.Get(&user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with userID %d", userID)
		}
		return nil, fmt.Errorf("error fetching user by userID: %w", err)
	}

	return &user, nil
}
