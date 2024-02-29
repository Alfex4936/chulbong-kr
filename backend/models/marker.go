package models

import "time"

// Marker corresponds to the Markers table in the database
type Marker struct {
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	UserID      int       `json:"userId" db:"UserID"`
	Latitude    float64   `json:"latitude" db:"Latitude"`
	Longitude   float64   `json:"longitude" db:"Longitude"`
	Description string    `json:"description" db:"Description"`
	CreatedAt   time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"UpdatedAt"`
}

// MarkerWithPhoto includes information about the marker and its associated photo.
type MarkerWithPhoto struct {
	Marker
	Photo Photo `json:"photo,omitempty"` // Embedded Photo struct
}

type MarkerWithPhotos struct {
	Marker
	Username     string  `json:"username,omitempty"`
	Photos       []Photo `json:"photos,omitempty"`
	DislikeCount int     `json:"dislikeCount,omitempty"`
}
