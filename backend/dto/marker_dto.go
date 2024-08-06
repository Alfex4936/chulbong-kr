package dto

import (
	"time"

	"github.com/Alfex4936/chulbong-kr/model"
)

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
	PageSize  int     `query:"n"`
	Page      int     `query:"page"`
}

type MarkerWithDistance struct {
	MarkerSimple
	Description string  `json:"description" db:"Description"`
	Distance    float64 `json:"distance" db:"distance"` // Distance in meters
	Address     *string `json:"address,omitempty" db:"Address"`
}

type MarkerWithDislike struct {
	model.Marker
	Username     string `db:"Username"`
	DislikeCount int    `db:"DislikeCount"`
}

type MarkerSimple struct {
	MarkerID  int     `json:"markerId" db:"MarkerID"`
	Latitude  float64 `json:"latitude" db:"Latitude"`
	Longitude float64 `json:"longitude" db:"Longitude"`
}

type MarkerSimpleWithDescrption struct {
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	Latitude    float64   `json:"latitude" db:"Latitude"`
	Longitude   float64   `json:"longitude" db:"Longitude"`
	Description string    `json:"description" db:"Description"`
	CreatedAt   time.Time `json:"-" db:"CreatedAt"`
	Address     string    `json:"address,omitempty" db:"Address"`
}

type MarkerSimpleWithAddr struct {
	MarkerID  int     `json:"markerId" db:"MarkerID"`
	Latitude  float64 `json:"latitude" db:"Latitude"`
	Longitude float64 `json:"longitude" db:"Longitude"`
	Address   string  `json:"address,omitempty" db:"Address"`
}

type MarkerRSS struct {
	MarkerID  int       `json:"markerId" db:"MarkerID"`
	UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
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
	MarkerID string
	Clicks   int
}

type MarkerGroup struct {
	CentralMarker MarkerSimple         // 중심 마커
	NearbyMarkers []MarkerWithDistance // 중심 마커 주변의 마커들
}

type WaterAPIResponse struct {
	Water     bool    `json:"water"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MarkerNewPicture struct {
	MarkerID int    `json:"markerId" db:"MarkerID"`
	PhotoURL string `json:"photoURL" db:"PhotoURL"`
}
