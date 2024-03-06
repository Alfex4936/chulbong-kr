package models

import (
	"time"
)

// Marker corresponds to the Markers table in the database
type Marker struct {
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	UserID      *int      `json:"userId" db:"UserID"`
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

// // Custom JSON marshaling to handle UserID
// func (m Marker) MarshalJSON() ([]byte, error) {
// 	type Alias Marker            // Prevents recursion
// 	var userID interface{} = nil // Default to null in JSON
// 	if m.UserID.Valid {
// 		userID = m.UserID.Int64
// 	}

// 	return json.Marshal(&struct {
// 		UserID interface{} `json:"userId"`
// 		*Alias
// 	}{
// 		UserID: userID,
// 		Alias:  (*Alias)(&m),
// 	})
// }

// // MarshalJSON custom JSON marshaling for MarkerWithPhotos
// func (mwp MarkerWithPhotos) MarshalJSON() ([]byte, error) {
// 	// Marshal Marker first to handle UserID properly
// 	type markerAlias Marker // Create an alias of Marker to use default Marshaling
// 	data, err := json.Marshal(&struct {
// 		UserID interface{} `json:"userId"`
// 		markerAlias
// 	}{
// 		UserID:      nil, // Start with UserID as nil
// 		markerAlias: markerAlias(mwp.Marker),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Override UserID if it's valid
// 	if mwp.UserID.Valid {
// 		data, err = json.Marshal(&struct {
// 			UserID int64 `json:"userId"`
// 			markerAlias
// 		}{
// 			UserID:      mwp.UserID.Int64,
// 			markerAlias: markerAlias(mwp.Marker),
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// Now, unmarshal it to a map to combine with other fields of MarkerWithPhotos
// 	var baseMap map[string]interface{}
// 	err = json.Unmarshal(data, &baseMap)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add additional fields from MarkerWithPhotos
// 	baseMap["username"] = mwp.Username
// 	baseMap["photos"] = mwp.Photos
// 	baseMap["dislikeCount"] = mwp.DislikeCount

// 	// Marshal combined map to JSON
// 	return json.Marshal(baseMap)
// }
