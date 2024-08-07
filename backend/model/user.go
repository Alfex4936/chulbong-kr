package model

import (
	"database/sql"
	"time"
)

// User corresponds to the Users table in the database
type User struct {
	CreatedAt    time.Time      `json:"-" db:"CreatedAt"`
	UpdatedAt    time.Time      `json:"-" db:"UpdatedAt"`
	UserID       int            `json:"userId" db:"UserID"` // TODO: UUID?
	Username     string         `json:"username" db:"Username"`
	Email        string         `json:"email" db:"Email"`
	Role         string         `json:"-" db:"Role"`
	PasswordHash sql.NullString `json:"-" db:"PasswordHash"` // Can be empty for OAuth2 users
	Provider     sql.NullString `json:"-" db:"Provider"`     // e.g., "google", "kakao"
	ProviderID   sql.NullString `json:"-" db:"ProviderID"`   // Unique ID from the OAuth provider
}
