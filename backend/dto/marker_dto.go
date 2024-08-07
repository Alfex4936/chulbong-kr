package dto

import (
	"time"

	"github.com/Alfex4936/chulbong-kr/model"
)

type MarkerRequest struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MarkerID    int     `json:"markerId,omitempty"`
	Description string  `json:"description"`
	PhotoURL    string  `json:"photoUrl,omitempty"`
}

type MarkerResponse struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MarkerID    int     `json:"markerId"`
	Description string  `json:"description"`
	// Username    string   `json:"username"`
	// UserID      int      `json:"userId"`
	// PhotoURLs   []string `json:"photoUrls"`
}

type QueryParams struct {
	Latitude  float64 `query:"latitude"`
	Longitude float64 `query:"longitude"`
	Distance  int     `query:"distance"`
	PageSize  int     `query:"n"`
	Page      int     `query:"page"`
}

type MarkerWithDistance struct {
	MarkerSimple
	Distance    float64 `json:"distance" db:"distance"` // Distance in meters
	Description string  `json:"description" db:"Description"`
	Address     *string `json:"address,omitempty" db:"Address"`
}

type MarkerWithDislike struct {
	model.Marker
	DislikeCount int    `db:"DislikeCount"`
	Username     string `db:"Username"`
}

type MarkerSimple struct {
	Latitude  float64 `json:"latitude" db:"Latitude"`
	Longitude float64 `json:"longitude" db:"Longitude"`
	MarkerID  int     `json:"markerId" db:"MarkerID"`
}

type MarkerSimpleWithDescrption struct {
	Latitude    float64   `json:"latitude" db:"Latitude"`
	Longitude   float64   `json:"longitude" db:"Longitude"`
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	CreatedAt   time.Time `json:"-" db:"CreatedAt"`
	Description string    `json:"description" db:"Description"`
	Address     string    `json:"address,omitempty" db:"Address"`
}

type MarkerSimpleWithAddr struct {
	Latitude  float64 `json:"latitude" db:"Latitude"`
	Longitude float64 `json:"longitude" db:"Longitude"`
	MarkerID  int     `json:"markerId" db:"MarkerID"`
	Address   string  `json:"address,omitempty" db:"Address"`
}

type MarkerRSS struct {
	UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
	MarkerID  int       `json:"markerId" db:"MarkerID"`
	Address   string    `json:"address,omitempty" db:"Address"`
}

type FacilityQuantity struct {
	FacilityID int `json:"facilityId"`
	Quantity   int `json:"quantity"`
}

type FacilityRequest struct {
	MarkerID   int                `json:"markerId"`
	Facilities []FacilityQuantity `json:"facilities"`
}

type MarkerRank struct {
	Clicks   int
	MarkerID string
}

type MarkerGroup struct {
	CentralMarker MarkerSimple         // 중심 마커
	NearbyMarkers []MarkerWithDistance // 중심 마커 주변의 마커들
}

type WaterAPIResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Water     bool    `json:"water"`
}

type MarkerNewPicture struct {
	PhotoURL string `json:"photoURL" db:"PhotoURL"`
	MarkerID int    `json:"markerId" db:"MarkerID"`
}

type MarkersWithUsernames struct {
	model.Marker
	Username      string `db:"Username"`
	DislikeCount  int    `db:"DislikeCount"`
	FavoriteCount int    `db:"FavoriteCount"`
}
