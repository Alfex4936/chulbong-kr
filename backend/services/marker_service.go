package services

import (
	"errors"
	"fmt"
	"math"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
)

// CreateMarker creates a new marker in the database after checking for nearby markers
func CreateMarker(markerDto *dto.MarkerRequest, userId int) (*models.Marker, error) {
	// First, check if there is a nearby marker
	nearby, err := IsMarkerNearby(markerDto.Latitude, markerDto.Longitude)
	if err != nil {
		return nil, err // Return any error encountered
	}
	if nearby {
		return nil, errors.New("a marker is already nearby")
	}

	// Insert the new marker
	const insertQuery = `INSERT INTO Markers (UserID, Latitude, Longitude, Description, CreatedAt, UpdatedAt) 
                         VALUES (?, ?, ?, ?, NOW(), NOW())`
	res, err := database.DB.Exec(insertQuery, userId, markerDto.Latitude, markerDto.Longitude, markerDto.Description)
	if err != nil {
		return nil, fmt.Errorf("inserting marker: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert ID: %w", err)
	}

	// Create a marker model instance to hold the full marker information
	marker := &models.Marker{
		MarkerID:    int(id),
		UserID:      userId,
		Latitude:    markerDto.Latitude,
		Longitude:   markerDto.Longitude,
		Description: markerDto.Description,
		// CreatedAt and UpdatedAt will be populated in the next step
	}

	// Fetch the newly created marker to populate all fields, including CreatedAt and UpdatedAt
	const selectQuery = `SELECT CreatedAt, UpdatedAt FROM Markers WHERE MarkerID = ?`
	err = database.DB.QueryRow(selectQuery, marker.MarkerID).Scan(&marker.CreatedAt, &marker.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("fetching created marker: %w", err)
	}

	return marker, nil
}

// GetMarker retrieves a single marker and its associated photo by the marker's ID
func GetMarker(markerID int) (*models.MarkerWithPhoto, error) {
	const query = `
	SELECT m.*, p.PhotoID, p.PhotoURL, p.UploadedAt 
	FROM Markers m
	LEFT JOIN Photos p ON m.MarkerID = p.MarkerID
	WHERE m.MarkerID = ?`

	var markerWithPhoto models.MarkerWithPhoto
	err := database.DB.Get(&markerWithPhoto, query, markerID)
	if err != nil {
		return nil, fmt.Errorf("fetching marker with photo: %w", err)
	}

	return &markerWithPhoto, nil
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
