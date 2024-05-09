package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/dto/kakao"
	"github.com/Alfex4936/chulbong-kr/model"
	"github.com/Alfex4936/chulbong-kr/util"
)

type MarkerFacilityService struct {
	Config       *config.AppConfig
	KakaoConfig  *config.KakaoConfig
	DB           *sqlx.DB
	HTTPClient   *http.Client
	RedisService *RedisService
}

func NewMarkerFacilityService(
	config *config.AppConfig,
	kakaoConfig *config.KakaoConfig,
	db *sqlx.DB,
	httpClient *http.Client,
	redisService *RedisService) *MarkerFacilityService {
	return &MarkerFacilityService{
		Config:       config,
		KakaoConfig:  kakaoConfig,
		DB:           db,
		HTTPClient:   httpClient,
		RedisService: redisService,
	}
}

// GetFacilitiesByMarkerID retrieves facilities for a given marker ID.
func (s *MarkerFacilityService) GetFacilitiesByMarkerID(markerID int) ([]model.Facility, error) {
	facilities := make([]model.Facility, 0)
	query := `SELECT FacilityID, MarkerID, Quantity FROM MarkerFacilities WHERE MarkerID = ?`
	err := s.DB.Select(&facilities, query, markerID)
	if err != nil {
		return nil, err
	}
	return facilities, nil
}

func (s *MarkerFacilityService) SetMarkerFacilities(markerID int, facilities []dto.FacilityQuantity) error {
	tx, err := s.DB.Beginx()
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
	if err = tx.Commit(); err != nil {
		return err
	}

	s.RedisService.ResetCache(fmt.Sprintf("facilities:%d", markerID))

	return nil
}

// UpdateMarkersAddresses fetches all markers, updates their addresses using an external API, and returns the updated list.
func (s *MarkerFacilityService) UpdateMarkersAddresses() ([]dto.MarkerSimpleWithAddr, error) {
	const markerQuery = `SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude FROM Markers;`

	var markers []dto.MarkerSimpleWithAddr
	err := s.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	for i := range markers {
		address, err := s.FetchAddressFromAPI(markers[i].Latitude, markers[i].Longitude)
		if err != nil || address == "" {
			// If there's an error fetching the address or the address is not found, skip updating this marker.
			continue
		}

		markers[i].Address = address
		// Update the marker's address in the database.
		if err := s.UpdateMarkerAddress(markers[i].MarkerID, address); err != nil {
			// Log or handle the error based on your application's error handling policy
			fmt.Printf("Failed to update address for marker %d: %v\n", markers[i].MarkerID, err)
		}
	}

	return markers, nil
}

// UpdateMarkerAddress updates the address of a marker by its ID.
func (s *MarkerFacilityService) UpdateMarkerAddress(markerID int, address string) error {
	query := `UPDATE Markers SET Address = ? WHERE MarkerID = ?`
	_, err := s.DB.Exec(query, address, markerID)
	return err
}

// FetchAddressFromAPI queries the external API to get the address for a given latitude and longitude.
func (s *MarkerFacilityService) FetchAddressFromAPI(latitude, longitude float64) (string, error) {
	reqURL := fmt.Sprintf("%s?x=%f&y=%f", s.KakaoConfig.KakaoCoord2Addr, longitude, latitude)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "-1", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+s.KakaoConfig.KakaoAK)
	resp, err := s.HTTPClient.Do(req)
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

// FetchAddressFromMap queries the external API to get the address for a given latitude and longitude. (KakaoMap)
func (s *MarkerFacilityService) FetchAddressFromMap(latitude, longitude float64) (string, error) {
	wcongnamul := util.ConvertWGS84ToWCONGNAMUL(latitude, longitude)

	reqURL := fmt.Sprintf("%s&x=%f&y=%f", s.KakaoConfig.KakaoAddressInfo, wcongnamul.X, wcongnamul.Y)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp kakao.KakaoMarkerData
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}

	// Decide which address to return based on the presence of Old and New
	if apiResp.New != nil && *apiResp.New.Name != "" {
		address := *apiResp.New.Name
		if apiResp.New.Building != nil && *apiResp.New.Building != "" {
			address += " " + *apiResp.New.Building
		}
		return address, nil
	} else if apiResp.Old != nil && *apiResp.Old.Name != "" {
		address := *apiResp.Old.Name
		if apiResp.Old.Building != nil && *apiResp.Old.Building != "" {
			address += " " + *apiResp.Old.Building
		}
		return address, nil
	}

	return "", fmt.Errorf("no valid address found")
}

// FetchXYFromAPI queries the external API to get the latitude and longitude for a given address.
func (s *MarkerFacilityService) FetchXYFromAPI(address string) (float64, float64, error) {
	reqURL := fmt.Sprintf("%s?query=%s", s.KakaoConfig.KakaoGeoCode, address)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+s.KakaoConfig.KakaoAK)
	resp, err := s.HTTPClient.Do(req)
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
func (s *MarkerFacilityService) FetchRegionFromAPI(latitude, longitude float64) (string, error) {
	reqURL := fmt.Sprintf("%s?x=%f&y=%f", s.KakaoConfig.KakaoCoord2Region, longitude, latitude)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "-1", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+s.KakaoConfig.KakaoAK)
	resp, err := s.HTTPClient.Do(req)
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
func (s *MarkerFacilityService) FetchWeatherFromAddress(latitude, longitude float64) (*kakao.WeatherRequest, error) {
	wcongnamul := util.ConvertWGS84ToWCONGNAMUL(latitude, longitude)

	reqURL := fmt.Sprintf("%s&x=%f&y=%f", s.KakaoConfig.KakaoWeather, wcongnamul.X, wcongnamul.Y)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Referer", reqURL)
	resp, err := s.HTTPClient.Do(req)
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
func (s *MarkerFacilityService) FetchRegionWaterInfo(latitude, longitude float64) (bool, error) {
	reqURL := fmt.Sprintf("%s?latitude=%f&longitude=%f&rapidapi-key=%s", s.Config.IsWaterURL, latitude, longitude, s.Config.IsWaterKEY)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
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

type CustomString string

func (cs *CustomString) UnmarshalJSON(data []byte) error {
	var asFloat float64
	if err := json.Unmarshal(data, &asFloat); err == nil {
		*cs = CustomString(fmt.Sprintf("%f", asFloat))
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*cs = CustomString(asString)
		return nil
	}

	return errors.New("value must be a float or string")
}

type DataItem struct {
	Date          string       `json:"date"`
	ChulbongCount int          `json:"chulbong_count"`
	Longitude     CustomString `json:"longitude"`
	IsAble        int          `json:"is_able"`
	ID            string       `json:"id"`
	PyeongCount   int          `json:"pyeong_count"`
	Latitude      CustomString `json:"latitude"`
	Address       string       `json:"address,omitempty"`
}

func (s *MarkerFacilityService) FetchLatestMarkers(thresholdDate time.Time) ([]DataItem, error) {
	req, err := http.NewRequest(http.MethodGet, s.Config.CkURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
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
	for _, item := range data {
		itemDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			continue
		}

		if itemDate.Equal(thresholdDate) || itemDate.After(thresholdDate) {
			filteredData = append(filteredData, item)
		}
	}
	return filteredData, nil
}
