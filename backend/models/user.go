package models

import "time"

// User corresponds to the Users table in the database
type User struct {
	UserID       int       `json:"userId" db:"UserID"`
	Username     string    `json:"username" db:"Username"`
	Email        string    `json:"email" db:"Email"`
	PasswordHash string    `json:"-" db:"PasswordHash"` // "-" in json tag to avoid sending it in responses
	CreatedAt    time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt    time.Time `json:"updatedAt" db:"UpdatedAt"`
}
