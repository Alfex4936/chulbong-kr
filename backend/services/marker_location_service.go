package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"sync"
	"time"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/utils"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

var KAKAO_STATIC = os.Getenv("KAKAO_STATIC_MAP")

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
	RedisStore.Do(context.Background(), RedisStore.B().Zmscore().Key("marker_clicks").Member(markerIDs...).Build())

	scores, err := RedisStore.Do(context.Background(), RedisStore.B().Zmscore().Key("marker_clicks").Member(markerIDs...).Build()).AsZScores()
	if err != nil {
		return nil, err
	}

	rankedMarkers := make([]dto.MarkerWithDistance, 0, len(scores))
	for i, zscore := range scores {
		if zscore.Score > float64(MIN_CLICK_RANK) {
			nearbyMarkers[i].Distance = zscore.Score // Repurpose Distance to store score
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

func SaveOfflineMap(lat, lng float64) (string, error) {
	if !utils.IsInSouthKoreaPrecisely(lat, lng) {
		return "", fmt.Errorf("only allowed in South Korea")
	}

	// 0. Get Address of lat/lng
	address, _ := FetchAddressFromAPI(lat, lng)
	if address == "-2" {
		return "", fmt.Errorf("address not found")
	}
	if address == "-1" {
		address = "대한민국 철봉 지도"
	}

	// 1. Convert them into WCONGNAMUL
	mapWcon := utils.ConvertWGS84ToWCONGNAMUL(lat, lng)

	// 2. Get the static map image (base_map_blah.png)
	// temporarily download from fmt.Sprintf("%s&MX=%f%MY=%f", KAKAO_STATIC, map_wcon.X, map_wcon.Y)
	tempDir, err := os.MkdirTemp("", "chulbongkr-*") // Use "" for the system default temp directory
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory")
	}
	// defer os.RemoveAll(tempDir)

	baseImageFile := fmt.Sprintf("base_map-%s.png", uuid.New().String())
	baseImageFilePath := path.Join(tempDir, baseImageFile)
	utils.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f", KAKAO_STATIC, mapWcon.X, mapWcon.Y), baseImageFilePath)

	// 3. Load all close markers nearby map lat/lng
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	// 1280*720 500m
	nearbyMarkers, total, err := FindClosestNMarkersWithinDistance(lat, lng, 500, 15, 0) // meter, pageSize, offset
	if err != nil {
		return "", fmt.Errorf("failed to find nearby markers")
	}
	if total == 0 {
		return "", nil // Return nil to signify no markers in the area, reducing slice allocation
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(nearbyMarkers))

	// temporarily download each
	tempImagePath := path.Join(tempDir, "marker_images")

	os.Mkdir(tempImagePath, os.FileMode(0755))

	for i, marker := range nearbyMarkers {
		wg.Add(1)
		go func(i int, marker dto.MarkerWithDistance) {
			defer wg.Done()
			markerWcon := utils.ConvertWGS84ToWCONGNAMUL(marker.Latitude, marker.Longitude)
			markerFile := path.Join(tempImagePath, fmt.Sprintf("marker-%d.png", i))
			err := utils.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f&CX=%f&CY=%f", KAKAO_STATIC, mapWcon.X, mapWcon.Y, markerWcon.X, markerWcon.Y), markerFile)
			if err != nil {
				errors <- fmt.Errorf("failed to download marker %d: %w", i, err)
				return
			}
		}(i, marker)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return "", err
		}
	}

	// there will be len(nearbyMarkers) + 1 (base map) png files. but doesn't matter if some missing

	// 4. Overlay them
	resultImagePath, err := utils.OverlayImages(baseImageFilePath, tempImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to overlay images")
	}

	// 5. Make PDF
	downloadPath, err := utils.GenerateMapPDF(resultImagePath, tempDir, address)
	if err != nil {
		return "", fmt.Errorf("failed to make pdf file: " + err.Error())
	}

	// Schedule cleanup for 5mins later
	// go func() {
	// 	time.Sleep(5 * time.Minute)
	// 	os.RemoveAll(tempDir)
	// }()

	return downloadPath, nil
}

// SaveOfflineMap2 draws markers with go rather than download images
func SaveOfflineMap2(lat, lng float64) (string, error) {
	if !utils.IsInSouthKoreaPrecisely(lat, lng) {
		return "", fmt.Errorf("only allowed in South Korea")
	}

	// 0. Get Address of lat/lng
	address, _ := FetchAddressFromAPI(lat, lng)
	if address == "-2" {
		return "", fmt.Errorf("address not found")
	}
	if address == "-1" {
		address = "대한민국 철봉 지도"
	}

	// 1. Convert them into WCONGNAMUL
	mapWcon := utils.ConvertWGS84ToWCONGNAMUL(lat, lng)

	// 2. Get the static map image (base_map_blah.png)
	// temporarily download from fmt.Sprintf("%s&MX=%f%MY=%f", KAKAO_STATIC, map_wcon.X, map_wcon.Y)
	tempDir, err := os.MkdirTemp("", "chulbongkr-*") // Use "" for the system default temp directory
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory")
	}
	// defer os.RemoveAll(tempDir)

	baseImageFile := fmt.Sprintf("base_map-%s.png", uuid.New().String())
	baseImageFilePath := path.Join(tempDir, baseImageFile)
	utils.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f", KAKAO_STATIC, mapWcon.X, mapWcon.Y), baseImageFilePath)

	// 3. Load all close markers nearby map lat/lng
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	// 1280*720 500m
	nearbyMarkers, total, err := FindClosestNMarkersWithinDistance(lat, lng, 700, 30, 0) // meter, pageSize, offset
	if err != nil {
		return "", fmt.Errorf("failed to find nearby markers")
	}
	if total == 0 {
		return "", nil // Return nil to signify no markers in the area, reducing slice allocation
	}

	markers := make([]utils.WCONGNAMULCoord, len(nearbyMarkers))
	for i, marker := range nearbyMarkers {
		markers[i] = utils.ConvertWGS84ToWCONGNAMUL(marker.Latitude, marker.Longitude)
	}

	// 4. Place them
	resultImagePath, err := utils.PlaceMarkersOnImage(baseImageFilePath, markers, mapWcon.X, mapWcon.Y)
	if err != nil {
		return "", fmt.Errorf("failed to overlay images")
	}

	os.Remove(baseImageFilePath) // Remove base image file

	// 5. Make PDF
	downloadPath, err := utils.GenerateMapPDF(resultImagePath, tempDir, address)
	if err != nil {
		return "", fmt.Errorf("failed to make pdf file: " + err.Error())
	}

	return downloadPath, nil
}
