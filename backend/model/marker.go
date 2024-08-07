package model

import "time"

// Marker corresponds to the Markers table in the database
type Marker struct {
	Latitude    float64   `json:"latitude" db:"Latitude"`
	Longitude   float64   `json:"longitude" db:"Longitude"`
	CreatedAt   time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"UpdatedAt"`
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	UserID      *int      `json:"userId" db:"UserID"`
	Description string    `json:"description" db:"Description"`
	Address     *string   `json:"address" db:"Address"`
}

// MarkerWithPhoto includes information about the marker and its associated photo.
type MarkerWithPhoto struct {
	Marker
	Photo Photo `json:"photo,omitempty"` // Embedded Photo struct
}

type MarkerWithPhotos struct {
	Marker
	Photos        []Photo `json:"photos,omitempty"`
	Username      string  `json:"username,omitempty"`
	DislikeCount  int     `json:"dislikeCount,omitempty"`
	FavoriteCount int     `json:"favCount,omitempty"`
	IsChulbong    bool    `json:"isChulbong,omitempty"`
	Disliked      bool    `json:"disliked"`
	Favorited     bool    `json:"favorited,omitempty"`
}
