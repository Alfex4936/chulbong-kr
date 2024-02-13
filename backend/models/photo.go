package models

import "time"

// Photo corresponds to the Photos table in the database
type Photo struct {
	PhotoID    int       `json:"photoId" db:"PhotoID"`
	MarkerID   int       `json:"markerId" db:"MarkerID"`
	PhotoURL   string    `json:"photoUrl" db:"PhotoURL"`
	UploadedAt time.Time `json:"uploadedAt" db:"UploadedAt"`
}
