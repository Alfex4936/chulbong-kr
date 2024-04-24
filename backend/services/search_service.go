package services

import (
	"bytes"
	"chulbong-kr/dto"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goccy/go-json"
)

var (
	zincApi      = os.Getenv("ZINCSEARCH_URL")
	zincUser     = os.Getenv("ZINCSEARCH_USER")
	zincPassword = os.Getenv("ZINCSEARCH_PASSWORD")
)

type FuzzMarkerSearch = dto.FuzzMarkerSearch
type MarkerSearchResponse = dto.MarkerSearchResponse
type ElasticsearchResponse = dto.ElasticsearchResponse

// SearchMarkerAddress calls ZincSearch (ElasticSearch-like) client.
func SearchMarkerAddress(term string) (dto.MarkerSearchResponse, error) {
	var apiResponse MarkerSearchResponse

	body := FuzzMarkerSearch{
		SearchType:   "fuzzy",
		Query:        dto.Query{Term: term},
		From:         0,
		MaxResults:   10,
		SourceFields: []string{},
	}

	// Marshal the value to JSON
	jsonByte, err := json.Marshal(body)
	if err != nil {
		return apiResponse, err
	}
	log.Printf("📆 %s", string(jsonByte))

	reqURL := fmt.Sprintf("%s/api/markers/_search", zincApi)
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonByte))
	if err != nil {
		return apiResponse, fmt.Errorf("creating request: %w", err)
	}

	log.Printf(" 📆came")

	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return apiResponse, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apiResponse, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var zincResp ElasticsearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResp); err != nil {
		return apiResponse, fmt.Errorf("unmarshalling response: %w", err)
	}

	customResp := MarkerSearchResponse{
		Took:    zincResp.Took,
		Markers: convertToMarkers(zincResp.Hits.Hits),
	}

	return customResp, nil
}

func convertToMarkers(response []dto.HitDetail) []dto.ZincMarker {
	markers := make([]dto.ZincMarker, 0, len(response))

	for _, hit := range response {
		marker := dto.ZincMarker{
			MarkerID: hit.Source.MarkerID,
			Address:  hit.Source.Address,
		}
		markers = append(markers, marker)
	}

	return markers
}
