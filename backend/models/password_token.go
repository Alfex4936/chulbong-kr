package models

import "time"

type PasswordToken struct {
	TokenID   int       `db:"TokenID"`
	Token     string    `db:"Token"`
	Email     string    `db:"Email"`
	Verified  bool      `db:"Verified"`
	ExpiresAt time.Time `db:"ExpiresAt"`
	CreatedAt time.Time `db:"CreatedAt"`
}
