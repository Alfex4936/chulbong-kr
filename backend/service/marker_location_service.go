package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	earthRadius = 6371000

	existNearbyMarkerQuery = `
SELECT EXISTS (
    SELECT 1 
    FROM Markers
    WHERE ST_Within(Location, ST_Buffer(ST_GeomFromText(?, 4326), ?))
) AS Nearby;
`

	// Using the optimized query with bounding box
	findClosestMarkersQuery = `
SELECT MarkerID, 
       ST_X(Location) AS Latitude, 
       ST_Y(Location) AS Longitude, 
       Description, 
       ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) AS distance, 
       Address
FROM Markers
WHERE MBRContains(
    ST_SRID(
        ST_GeomFromText(
            CONCAT('POLYGON((',
                   ?, ' ', ?, ',', 
                   ?, ' ', ?, ',', 
                   ?, ' ', ?, ',', 
                   ?, ' ', ?, ',', 
                   ?, ' ', ?, 
                   '))')), 
        4326), 
    Location)
  AND ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?
ORDER BY distance ASC`
)

type MarkerLocationService struct {
	DB              *sqlx.DB
	Config          *config.AppConfig
	KakaoConfig     *config.KakaoConfig
	Redis           *RedisService
	MapUtil         *util.MapUtil
	FacilityService *MarkerFacilityService
}

func NewMarkerLocationService(
	db *sqlx.DB,
	config *config.AppConfig,
	kakaoConfig *config.KakaoConfig,
	redis *RedisService,
	mapUtil *util.MapUtil,
	facilityService *MarkerFacilityService,
) *MarkerLocationService {
	return &MarkerLocationService{
		DB:              db,
		Config:          config,
		KakaoConfig:     kakaoConfig,
		Redis:           redis,
		MapUtil:         mapUtil,
		FacilityService: facilityService,
	}
}

// meters_per_degree = 40075000 / 360 / 1000
// IsMarkerNearby checks if there's a marker within n meters of the given latitude and longitude
func (s *MarkerLocationService) IsMarkerNearby(lat, long float64, bufferDistanceMeters int) (bool, error) {
	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	// Execute the query
	var nearby bool
	err := s.DB.Get(&nearby, existNearbyMarkerQuery, point, bufferDistanceMeters)
	if err != nil {
		return false, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	return nearby, nil
}

// FindClosestNMarkersWithinDistance
func (s *MarkerLocationService) FindClosestNMarkersWithinDistance(lat, long float64, distance, pageSize, offset int) ([]dto.MarkerWithDistance, int, error) {
	// Calculate bounding box
	radLat := lat * math.Pi / 180
	radDist := float64(distance) / earthRadius
	minLat := lat - radDist*180/math.Pi
	maxLat := lat + radDist*180/math.Pi
	minLon := long - radDist*180/(math.Pi*math.Cos(radLat))
	maxLon := long + radDist*180/(math.Pi*math.Cos(radLat))

	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	var allMarkers []dto.MarkerWithDistance
	err := s.DB.Select(&allMarkers, findClosestMarkersQuery, point, minLon, minLat, maxLon, minLat, maxLon, maxLat, minLon, maxLat, minLon, minLat, point, distance)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	// Implementing pagination in application logic
	markers := paginateMarkers(allMarkers, pageSize, offset)
	total := len(allMarkers)

	return markers, total, nil
}

func (s *MarkerLocationService) FindRankedMarkersInCurrentArea(lat, long float64, distance, limit int) ([]dto.MarkerWithDistance, error) {
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	nearbyMarkers, total, err := s.FindClosestNMarkersWithinDistance(lat, long, distance, limit, 0)
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

	ctx := context.Background()
	floatMin := float64(MinClickRank)

	result, _ := s.Redis.Core.Client.Do(ctx, s.Redis.Core.Client.B().Zmscore().Key("marker_clicks").Member(markerIDs...).Build()).AsFloatSlice()
	rankedMarkers := make([]dto.MarkerWithDistance, 0, len(result))
	for i, score := range result {
		if score > floatMin { // Include markers with score > minScore
			nearbyMarkers[i].Distance = score
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

func (s *MarkerLocationService) SaveOfflineMap(lat, lng float64) (string, error) {
	if !s.MapUtil.IsInSouthKoreaPrecisely(lat, lng) {
		return "", fmt.Errorf("only allowed in South Korea")
	}

	// 0. Get Address of lat/lng
	address, _ := s.FacilityService.FetchAddressFromAPI(lat, lng)
	if address == "-2" {
		return "", fmt.Errorf("address not found")
	}
	if address == "-1" {
		address = "대한민국 철봉 지도"
	}

	// 1. Convert them into WCONGNAMUL
	mapWcon := util.ConvertWGS84ToWCONGNAMUL(lat, lng)

	// 2. Get the static map image (base_map_blah.png)
	// temporarily download from fmt.Sprintf("%s&MX=%f%MY=%f", KAKAO_STATIC, map_wcon.X, map_wcon.Y)
	tempDir, err := os.MkdirTemp("", "chulbongkr-*") // Use "" for the system default temp directory
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory")
	}
	// defer os.RemoveAll(tempDir)

	baseImageFile := fmt.Sprintf("base_map-%s.png", uuid.New().String())
	baseImageFilePath := path.Join(tempDir, baseImageFile)
	util.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f", s.KakaoConfig.KakaoStaticMap, mapWcon.X, mapWcon.Y), baseImageFilePath)

	// 3. Load all close markers nearby map lat/lng
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	// 1280*720 500m
	nearbyMarkers, total, err := s.FindClosestNMarkersWithinDistance(lat, lng, 500, 15, 0) // meter, pageSize, offset
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
			markerWcon := util.ConvertWGS84ToWCONGNAMUL(marker.Latitude, marker.Longitude)
			markerFile := path.Join(tempImagePath, fmt.Sprintf("marker-%d.png", i))
			err := util.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f&CX=%f&CY=%f", s.KakaoConfig.KakaoStaticMap, mapWcon.X, mapWcon.Y, markerWcon.X, markerWcon.Y), markerFile)
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
	resultImagePath, err := util.OverlayImages(baseImageFilePath, tempImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to overlay images")
	}

	// 5. Make PDF
	downloadPath, err := util.GenerateMapPDF(resultImagePath, tempDir, address, nearbyMarkers[0].MarkerID)
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
func (s *MarkerLocationService) SaveOfflineMap2(lat, lng float64) (string, error) {
	if !s.MapUtil.IsInSouthKoreaPrecisely(lat, lng) {
		return "", fmt.Errorf("only allowed in South Korea")
	}

	// 0. Get Address of lat/lng
	address, _ := s.FacilityService.FetchAddressFromAPI(lat, lng)
	if address == "-2" {
		return "", fmt.Errorf("address not found")
	}
	if address == "-1" {
		address = "대한민국 철봉 지도"
	}

	// 1. Convert them into WCONGNAMUL
	mapWcon := util.ConvertWGS84ToWCONGNAMUL(lat, lng)

	// 2. Get the static map image (base_map_blah.png)
	// temporarily download from fmt.Sprintf("%s&MX=%f%MY=%f", KAKAO_STATIC, map_wcon.X, map_wcon.Y)
	tempDir, err := os.MkdirTemp("", "chulbongkr-*") // Use "" for the system default temp directory
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory")
	}
	// defer os.RemoveAll(tempDir)

	baseImageFile := fmt.Sprintf("base_map-%s.png", uuid.New().String())
	baseImageFilePath := path.Join(tempDir, baseImageFile)
	util.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f", s.KakaoConfig.KakaoStaticMap, mapWcon.X, mapWcon.Y), baseImageFilePath)

	// 3. Load all close markers nearby map lat/lng
	// Predefine capacity for slices based on known limits to avoid multiple allocations
	// 1280*720 500m
	nearbyMarkers, total, err := s.FindClosestNMarkersWithinDistance(lat, lng, 700, 30, 0) // meter, pageSize, offset
	log.Printf("🎙️ Found %+v markers", nearbyMarkers)
	if err != nil {
		return "", fmt.Errorf("failed to find nearby markers")
	}
	if total == 0 {
		return "", nil // Return nil to signify no markers in the area, reducing slice allocation
	}

	markers := make([]util.WCONGNAMULCoord, len(nearbyMarkers))
	for i, marker := range nearbyMarkers {
		markers[i] = util.ConvertWGS84ToWCONGNAMUL(marker.Latitude, marker.Longitude)
	}

	// 4. Place them
	resultImagePath, err := util.PlaceMarkersOnImage(baseImageFilePath, markers, mapWcon.X, mapWcon.Y)
	if err != nil {
		return "", fmt.Errorf("failed to overlay images")
	}

	os.Remove(baseImageFilePath) // Remove base image file

	// 5. Make PDF
	downloadPath, err := util.GenerateMapPDF(resultImagePath, tempDir, address, nearbyMarkers[0].MarkerID)
	if err != nil {
		return "", fmt.Errorf("failed to make pdf file: " + err.Error())
	}

	return downloadPath, nil
}

// Simple pagination helper function
func paginateMarkers(markers []dto.MarkerWithDistance, pageSize, offset int) []dto.MarkerWithDistance {
	if offset >= len(markers) {
		// Calculate the starting index of the last possible page
		lastPageOffset := len(markers) - pageSize
		if lastPageOffset < 0 {
			lastPageOffset = 0
		}
		return markers[lastPageOffset:]
	}
	end := offset + pageSize
	if end > len(markers) {
		end = len(markers)
	}
	return markers[offset:end]
}

func (s *MarkerLocationService) TestDynamic(latitude, longitude, zoomScale float64, width, height int64) {
	nearbyMarkers, total, err := s.FindClosestNMarkersWithinDistance(latitude, longitude, 700, 30, 0) // meter, pageSize, offset
	if err != nil {
		return
	}
	if total == 0 {
		return
	}

	markers := make([]util.WCONGNAMULCoord, len(nearbyMarkers))
	for i, marker := range nearbyMarkers {
		markers[i] = util.ConvertWGS84ToWCONGNAMUL(marker.Latitude, marker.Longitude)
	}
	mapWcon := util.ConvertWGS84ToWCONGNAMUL(latitude, longitude)
	baseImageFile := fmt.Sprintf("base_map-%s.png", uuid.New().String())
	baseImageFilePath := path.Join("./tests", baseImageFile)
	log.Printf("📆 %s", baseImageFilePath)

	static := fmt.Sprintf("https://spi.maps.daum.net/map2/map/imageservice?IW=%d&IH=%d&SCALE=%f&service=open", width, height, zoomScale)
	util.DownloadFile(fmt.Sprintf("%s&MX=%f&MY=%f", static, mapWcon.X, mapWcon.Y), baseImageFilePath)

	resultImagePath, _ := util.PlaceMarkersOnImageDynamic(baseImageFilePath, markers, mapWcon.X, mapWcon.Y, zoomScale)
	fmt.Println(resultImagePath)
}
