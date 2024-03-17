package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/dto/kakao"
	"chulbong-kr/models"
)

const KAKAO_COORD2ADDR = "https://dapi.kakao.com/v2/local/geo/coord2address.json"

var KAKAO_AK = os.Getenv("KAKAO_AK")

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
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging the request indefinitely
	}
	reqURL := fmt.Sprintf("%s?x=%f&y=%f", KAKAO_COORD2ADDR, longitude, latitude)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "-1", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "KakaoAK "+KAKAO_AK)
	resp, err := client.Do(req)
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
