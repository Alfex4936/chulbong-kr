package services

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/dto/kakao"
	"chulbong-kr/models"
	"chulbong-kr/util"
)

const (
	KAKAO_GEOCODE      = "https://dapi.kakao.com/v2/local/geo/address.json"
	KAKAO_COORD2ADDR   = "https://dapi.kakao.com/v2/local/geo/coord2address.json"
	KAKAO_COORD2REGION = "https://dapi.kakao.com/v2/local/geo/coord2regioncode.json"
	KAKAO_WEATHER      = "https://map.kakao.com/api/dapi/point/weather?inputCoordSystem=WCONGNAMUL&outputCoordSystem=WCONGNAMUL&version=2&service=map.daum.net"
)

var (
	KakaoAK = os.Getenv("KAKAO_AK")

	HTTPClient = &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging requests indefinitely
	}

	IsWaterURL = os.Getenv("IS_WATER_API")
	IsWaterKEY = os.Getenv("IS_WATER_API_KEY")
	CkURL      = os.Getenv("CK_URL")
)

// GetFacilitiesByMarkerID retrieves facilities for a given marker ID.
func GetFacilitiesByMarkerID(markerID int) ([]models.Facility, error) {
	facilities := make([]models.Facility, 0)
	query := `SELECT FacilityID, MarkerID, Quantity FROM MarkerFacilities WHERE MarkerID = ?`
	err := database.DB.Select(&facilities, query, markerID)
	if err != nil {
		return nil, err
	}
	return facilities, nil
}

func SetMarkerFacilities(markerID int, facilities []dto.FacilityQuantity) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing facilities for the marker
	if _, err := tx.Exec("DELETE FROM MarkerFacilities WHERE MarkerID = ?", markerID); err != nil {
		return err
	}

	// Insert new facilities with quantities for the marker
	for _, fq := range facilities {
		if _, err := tx.Exec("INSERT INTO MarkerFacilities (FacilityID, MarkerID, Quantity) VALUES (?, ?, ?)", fq.FacilityID, markerID, fq.Quantity); err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// UpdateMarkersAddresses fetches all markers, updates their addresses using an external API, and returns the updated list.
func UpdateMarkersAddresses() ([]dto.MarkerSimpleWithAddr, error) {
	const markerQuery = `SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude FROM Markers;`

	var markers []dto.MarkerSimpleWithAddr
	err := database.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	for i := range markers {
		address, err := FetchAddressFromAPI(markers[i].Latitude, markers[i].Longitude)
		if err != nil || address == "" {
			// If there's an error fetching the address or the address is not found, skip updating this marker.
			continue
		}

		markers[i].Address = address
		// Update the marker's address in the database.
		if err := UpdateMarkerAddress(markers[i].MarkerID, address); err != nil {
			// Log or handle the error based on your application's error handling policy
			fmt.Printf("Failed to update address for marker %d: %v\n", markers[i].MarkerID, err)
		}
	}

	return markers, nil
}

// UpdateMarkerAddress updates the address of a marker by its ID.
func UpdateMarkerAddress(markerID int, address string) error {
	query := `UPDATE Markers SET Address = ? WHERE MarkerID = ?`
	_, err := database.DB.Exec(query, address, markerID)
	return err
}

// FetchAddressFromAPI queries the external API to get the address for a given latitude and longitude.
func FetchAddressFromAPI(latitude, longitude float64) (string, error) {
	reqURL := fmt.Sprintf("%s?x=%f&y=%f", KAKAO_COORD2ADDR, longitude, latitude)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "-1", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+KakaoAK)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "-1", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "-1", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp kakao.KakaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "-1", fmt.Errorf("unmarshalling response: %w", err)
	}

	if len(apiResp.Documents) == 0 {
		// log.Print("No address found for the given coordinates")
		return "", nil // Returning nil error to indicate absence of data rather than a failure
	}

	doc := apiResp.Documents[0]
	if doc.Address != nil {
		return doc.Address.AddressName, nil
	}
	if doc.RoadAddress != nil {
		return doc.RoadAddress.AddressName, nil
	}

	return "", nil // Address data is empty but no error occurred
}

// FetchXYFromAPI queries the external API to get the latitude and longitude for a given address.
func FetchXYFromAPI(address string) (float64, float64, error) {
	reqURL := fmt.Sprintf("%s?query=%s", KAKAO_GEOCODE, address)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+KakaoAK)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0.0, 0.0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp kakao.KakaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0.0, 0.0, fmt.Errorf("unmarshalling response: %w", err)
	}

	if len(apiResp.Documents) == 0 {
		// log.Print("No address found for the given coordinates")
		return 0.0, 0.0, nil // Returning nil error to indicate absence of data rather than a failure
	}

	doc := apiResp.Documents[0]
	if doc.X != nil && doc.Y != nil {

		latitude, err := strconv.ParseFloat(*doc.X, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid latitude")
		}

		longitude, err := strconv.ParseFloat(*doc.Y, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid longitude")
		}

		return latitude, longitude, nil
	}

	return 0.0, 0.0, nil // Address data is empty but no error occurred
}

// FetchRegionFromAPI queries the external API to get the address for a given latitude and longitude.
func FetchRegionFromAPI(latitude, longitude float64) (string, error) {
	reqURL := fmt.Sprintf("%s?x=%f&y=%f", KAKAO_COORD2REGION, longitude, latitude)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "-1", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+KakaoAK)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "-1", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "-1", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp kakao.KakaoRegionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "-1", fmt.Errorf("unmarshalling response: %w", err)
	}

	if len(apiResp.Documents) == 0 {
		// log.Print("No address found for the given coordinates")
		return "", nil // Returning nil error to indicate absence of data rather than a failure
	}

	doc := apiResp.Documents[0]
	if doc.AddressName == "북한" {
		return "-2", nil
	} else if doc.AddressName == "일본" {
		return "-2", nil
	}

	return doc.AddressName, nil // Address data is empty but no error occurred
}

// FetchWeatherFromAddress
func FetchWeatherFromAddress(latitude, longitude float64) (*kakao.WeatherRequest, error) {
	wcongnamul := util.ConvertWGS84ToWCONGNAMUL(latitude, longitude)

	reqURL := fmt.Sprintf("%s&x=%f&y=%f", KAKAO_WEATHER, wcongnamul.X, wcongnamul.Y)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Referer", reqURL)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp kakao.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	if apiResp.Codes.ResultCode != "OK" {
		// log.Print("No weather found for this address")
		return nil, fmt.Errorf("no weather found for this address")
	}

	icon := fmt.Sprintf("https://t1.daumcdn.net/localimg/localimages/07/2018/pc/weather/ico_weather%s.png", apiResp.WeatherInfos.Current.IconId)
	weatherRequest := kakao.WeatherRequest{
		Temperature: apiResp.WeatherInfos.Current.Temperature,
		Desc:        apiResp.WeatherInfos.Current.Desc,
		IconImage:   icon,
		Humidity:    apiResp.WeatherInfos.Current.Humidity,
		Rainfall:    apiResp.WeatherInfos.Current.Rainfall,
		Snowfall:    apiResp.WeatherInfos.Current.Snowfall,
	}
	return &weatherRequest, nil
}

// FetchRegionWaterInfo checks if latitude/longitude is in the water possibly.
func FetchRegionWaterInfo(latitude, longitude float64) (bool, error) {
	reqURL := fmt.Sprintf("%s?latitude=%f&longitude=%f&rapidapi-key=%s", IsWaterURL, latitude, longitude, IsWaterKEY)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp dto.WaterAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return false, fmt.Errorf("unmarshalling response: %w", err)
	}

	return apiResp.Water, nil
}

type DataItem struct {
	Date          string  `json:"date"`
	ChulbongCount int     `json:"chulbong_count"`
	Longitude     float64 `json:"longitude"`
	IsAble        int     `json:"is_able"`
	ID            string  `json:"id"`
	PyeongCount   int     `json:"pyeong_count"`
	Latitude      float64 `json:"latitude"`
}

func FetchLatestMarkers(thresholdDateString string) ([]DataItem, error) {
	req, err := http.NewRequest(http.MethodGet, CkURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data []DataItem
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	var filteredData []DataItem
	thresholdDate, err := time.Parse("2006-1-2", thresholdDateString)
	if err != nil {
		return nil, fmt.Errorf("parsing threshold date: %w", err)
	}

	for _, item := range data {
		itemDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date from data: %w", err)
		}
		if itemDate.Equal(thresholdDate) || itemDate.After(thresholdDate) {
			filteredData = append(filteredData, item)
		}
	}
	return filteredData, nil
}
