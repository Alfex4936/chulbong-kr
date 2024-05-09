package service

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/Alfex4936/chulbong-kr/dto"

	"github.com/goccy/go-json"
)

type ZincSearchService struct {
	ZincConfig *config.ZincSearchConfig
	HTTPClient *http.Client
}

func NewZincSearchService(zconfig *config.ZincSearchConfig, httpClient *http.Client) *ZincSearchService {
	return &ZincSearchService{
		ZincConfig: zconfig,
		HTTPClient: httpClient,
	}
}

type FuzzMarkerSearch = dto.FuzzMarkerSearch
type MarkerSearchResponse = dto.MarkerSearchResponse
type ElasticsearchResponse = dto.ElasticsearchResponse

// SearchMarkerAddress calls ZincSearch (ElasticSearch-like) client.
func (s *ZincSearchService) SearchMarkerAddress(term string) (dto.MarkerSearchResponse, error) {
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

	reqURL := fmt.Sprintf("%s/api/markers/_search", s.ZincConfig.ZincAPI)
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonByte))
	if err != nil {
		return apiResponse, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(s.ZincConfig.ZincUser, s.ZincConfig.ZincPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
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
