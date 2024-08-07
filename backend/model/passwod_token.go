package model

import "time"

type PasswordToken struct {
	ExpiresAt time.Time `db:"ExpiresAt"`
	CreatedAt time.Time `db:"CreatedAt"`
	TokenID   int       `db:"TokenID"`
	Token     string    `db:"Token"`
	Email     string    `db:"Email"`
	Verified  bool      `db:"Verified"`
}
