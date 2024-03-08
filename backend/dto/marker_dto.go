package dto

import "chulbong-kr/models"

type MarkerRequest struct {
	MarkerID    int     `json:"markerId,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	PhotoURL    string  `json:"photoUrl,omitempty"`
}

type MarkerResponse struct {
	MarkerID    int     `json:"markerId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	// Username    string   `json:"username"`
	// UserID      int      `json:"userId"`
	// PhotoURLs   []string `json:"photoUrls"`
}

type QueryParams struct {
	Latitude  float64 `query:"latitude"`
	Longitude float64 `query:"longitude"`
	Distance  int     `query:"distance"`
	N         int     `query:"n"`
	Page      int     `query:"page"`
}

type MarkerWithDistance struct {
	MarkerSimple
	Description string  `json:"description" db:"Description"`
	Distance    float64 `json:"distance" db:"distance"` // Distance in meters
}

type MarkerWithDislike struct {
	models.Marker
	Username     string `db:"Username"`
	DislikeCount int    `db:"DislikeCount"`
}

type MarkerSimple struct {
	MarkerID  int     `json:"markerId" db:"MarkerID"`
	Latitude  float64 `json:"latitude" db:"Latitude"`
	Longitude float64 `json:"longitude" db:"Longitude"`
}
