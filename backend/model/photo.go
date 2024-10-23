package model

import (
	"time"
)

// Photo corresponds to the Photos table in the database
type Photo struct {
	UploadedAt   time.Time `json:"uploadedAt" db:"UploadedAt"`
	PhotoID      int       `json:"photoId" db:"PhotoID"`
	MarkerID     int       `json:"markerId" db:"MarkerID"`
	PhotoURL     string    `json:"photoUrl" db:"PhotoURL"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty" db:"ThumbnailURL"`
}
