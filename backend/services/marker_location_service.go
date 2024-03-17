package services

import (
	"fmt"

	"chulbong-kr/database"
	"chulbong-kr/dto"
)

// meters_per_degree = 40075000 / 360 / 1000
// IsMarkerNearby checks if there's a marker within n meters of the given latitude and longitude
func IsMarkerNearby(lat, long float64, bufferDistanceMeters int) (bool, error) {
	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	query := `
SELECT EXISTS (
    SELECT 1 
    FROM Markers
    WHERE ST_Within(Location, ST_Buffer(ST_GeomFromText(?, 4326), ?))
) AS Nearby;
`

	// Execute the query
	var nearby bool
	err := database.DB.Get(&nearby, query, point, bufferDistanceMeters)
	if err != nil {
		return false, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	return nearby, nil
}

// FindClosestNMarkersWithinDistance
func FindClosestNMarkersWithinDistance(lat, long float64, distance, pageSize, offset int) ([]dto.MarkerWithDistance, int, error) {
	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	// Query to find markers within N meters and calculate total
	query := `
SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude, Description, ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) AS distance, Address
FROM Markers
WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?
ORDER BY distance ASC
LIMIT ? OFFSET ?`

	var markers []dto.MarkerWithDistance
	err := database.DB.Select(&markers, query, point, point, distance, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	// Query to get total count of markers within distance
	countQuery := `
SELECT COUNT(*)
FROM Markers
WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?`

	var total int
	err = database.DB.Get(&total, countQuery, point, distance)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total markers count: %w", err)
	}

	return markers, total, nil
}
