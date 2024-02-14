package services

import (
	"errors"
	"math"

	"chulbong-kr/database"
	"chulbong-kr/models"
)

// CreateMarker creates a new marker in the database after checking for nearby markers
func CreateMarker(marker *models.Marker) error {
	// First, check if there is a nearby marker
	nearby, err := IsMarkerNearby(marker.Latitude, marker.Longitude)
	if err != nil {
		return err // Return any error encountered
	}
	if nearby {
		return errors.New("a marker is already nearby")
	}

	// If no nearby marker, proceed to insert the new marker
	const query = `INSERT INTO Markers (UserID, Latitude, Longitude, Description) 
                   VALUES (?, ?, ?, ?)`
	res, err := database.DB.Exec(query, marker.UserID, marker.Latitude, marker.Longitude, marker.Description)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	marker.MarkerID = int(id)
	return nil
}

// GetMarker retrieves a single marker by its ID
func GetMarker(markerID int) (*models.Marker, error) {
	const query = `SELECT * FROM Markers WHERE MarkerID = ?`
	var marker models.Marker
	err := database.DB.Get(&marker, query, markerID)
	if err != nil {
		return nil, err
	}
	return &marker, nil
}

// UpdateMarker updates an existing marker's latitude, longitude, and description
func UpdateMarker(marker *models.Marker) error {
	const query = `UPDATE Markers SET Latitude = ?, Longitude = ?, Description = ?, UpdatedAt = NOW() 
                   WHERE MarkerID = ?`
	_, err := database.DB.Exec(query, marker.Latitude, marker.Longitude, marker.Description, marker.MarkerID)
	return err
}

// IsMarkerNearby checks if there's a marker within 5 meters of the given latitude and longitude
func IsMarkerNearby(lat, long float64) (bool, error) {
	const query = `SELECT Latitude, Longitude FROM Markers`
	var markers []models.Marker
	err := database.DB.Select(&markers, query)
	if err != nil {
		return false, err
	}
	for _, m := range markers {
		if distance(lat, long, m.Latitude, m.Longitude) <= 5 {
			return true, nil
		}
	}
	return false, nil
}

// distance calculates the distance between two geographic coordinates in meters
func distance(lat1, long1, lat2, long2 float64) float64 {
	var deltaLat = (lat2 - lat1) * (math.Pi / 180)
	var deltaLong = (long2 - long1) * (math.Pi / 180)
	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(deltaLong/2)*math.Sin(deltaLong/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return 6371000 * c // Earth radius in meters
}
