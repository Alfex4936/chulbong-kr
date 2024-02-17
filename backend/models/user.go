package models

import (
	"database/sql"
	"time"
)

// User corresponds to the Users table in the database
type User struct {
	UserID       int            `json:"userId" db:"UserID"`
	Username     string         `json:"username" db:"Username"`
	Email        string         `json:"email" db:"Email"`
	PasswordHash sql.NullString `json:"-" db:"PasswordHash"`                  // Can be empty for OAuth2 users
	Provider     sql.NullString `json:"provider,omitempty" db:"Provider"`     // e.g., "google", "kakao"
	ProviderID   sql.NullString `json:"providerId,omitempty" db:"ProviderID"` // Unique ID from the OAuth provider
	CreatedAt    time.Time      `json:"-" db:"CreatedAt"`
	UpdatedAt    time.Time      `json:"-" db:"UpdatedAt"`
}
