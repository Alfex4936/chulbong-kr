package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"chulbong-kr/database"
	"chulbong-kr/dto"

	"github.com/goccy/go-json"
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

func FindRankedMarkersInCurrentArea(lat, long float64, distance, limit int) ([]dto.MarkerWithDistance, error) {
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	nearbyMarkers, total, err := FindClosestNMarkersWithinDistance(lat, long, distance, limit, 0)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, nil // Return nil to signify no markers, reducing slice allocation
	}

	markerIDs := make([]string, len(nearbyMarkers))
	for i, marker := range nearbyMarkers {
		markerIDs[i] = strconv.Itoa(marker.MarkerID)
	}

	// Fetch scores for all markers in one Redis call
	scores, err := RedisStore.ZMScore(context.Background(), "marker_clicks", markerIDs...).Result()
	if err != nil {
		return nil, err
	}

	rankedMarkers := make([]dto.MarkerWithDistance, 0, len(scores))
	for i, score := range scores {
		if score > float64(MIN_CLICK_RANK) {
			nearbyMarkers[i].Distance = score // Repurpose Distance to store score
			rankedMarkers = append(rankedMarkers, nearbyMarkers[i])
		}
	}

	if len(rankedMarkers) == 0 {
		return nil, nil // Return nil to signify no ranked markers
	}

	// The sorting logic remains unchanged, but it's necessary for ranking
	sort.SliceStable(rankedMarkers, func(i, j int) bool {
		return rankedMarkers[i].Distance > rankedMarkers[j].Distance
	})

	// Applying limit after sorting
	if limit > len(rankedMarkers) {
		limit = len(rankedMarkers)
	}

	return rankedMarkers[:limit], nil
}

// GoogleGeoResponse struct to parse the Google Maps Geocoding API response
type GoogleGeoResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

// GeocodeAddress queries the Google Maps Geocoding API to get latitude and longitude for a given address.
func GeocodeAddress(address, apiKey string) (float64, float64, error) {
	var baseURL = "https://maps.googleapis.com/maps/api/geocode/json"

	// Prepare the request URL with query parameters
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", apiKey)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build the complete URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the request
	resp, err := client.Get(requestURL)
	if err != nil {
		return 0, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var geoResp GoogleGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, fmt.Errorf("error decoding response: %w", err)
	}

	if geoResp.Status != "OK" || len(geoResp.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found")
	}

	// Extract latitude and longitude
	lat := geoResp.Results[0].Geometry.Location.Lat
	lng := geoResp.Results[0].Geometry.Location.Lng

	return lat, lng, nil
}
