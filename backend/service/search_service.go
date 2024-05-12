package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

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

type (
	FuzzMarkerSearch      = dto.FuzzMarkerSearch
	MarkerSearchResponse  = dto.MarkerSearchResponse
	ElasticsearchResponse = dto.ElasticsearchResponse
	MarkerIndexData       = dto.MarkerIndexData
)

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

	resp, err := s.sendRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonByte))
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

// not bulk action
func (s *ZincSearchService) InsertMarkerIndex(indexBody MarkerIndexData) error {
	// Marshal the value to JSON
	jsonByte, err := json.Marshal(indexBody)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/api/markers/_doc", s.ZincConfig.ZincAPI)
	resp, err := s.sendRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonByte))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (s *ZincSearchService) DeleteMarkerIndex(markerID string) error {
	if markerID == "" {
		return fmt.Errorf("missing marker id")
	}
	markerIndexID, err := s.getMarkerIndexID(markerID)
	if err != nil {
		return fmt.Errorf("getting marker index ID: %w", err)
	}
	if markerIndexID == "" {
		return nil // No marker found, no need to delete
	}

	return s.deleteDocument(markerIndexID)
}

func (s *ZincSearchService) getMarkerIndexID(markerID string) (string, error) {
	query := fmt.Sprintf(`
	{
		"search_type": "match",
		"query": {
			"term": "%s",
			"field": "markerId"
		},
		"sort_fields": ["-@timestamp"],
		"from": 0,
		"max_results": 1,
		"_source": []
	}`, markerID)

	reqURL := fmt.Sprintf("%s/api/markers/_search", s.ZincConfig.ZincAPI)
	resp, err := s.sendRequest(http.MethodPost, reqURL, strings.NewReader(query))
	if err != nil {
		return "", fmt.Errorf("searching for marker index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search response status code: %d", resp.StatusCode)
	}

	var zincResp ElasticsearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if zincResp.Hits.Total.Value == 0 {
		return "", nil // No result found
	}

	return zincResp.Hits.Hits[0].ID, nil
}

func (s *ZincSearchService) deleteDocument(docID string) error {
	reqURL := fmt.Sprintf("%s/api/markers/_doc/%s", s.ZincConfig.ZincAPI, docID)
	resp, err := s.sendRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("deleting document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete response status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *ZincSearchService) sendRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}
	req.SetBasicAuth(s.ZincConfig.ZincUser, s.ZincConfig.ZincPassword)
	req.Header.Set("Content-Type", "application/json")
	return s.HTTPClient.Do(req)
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
