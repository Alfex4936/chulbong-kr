package dto

// ElasticsearchResponse represents the entire JSON response structure
type ElasticsearchResponse struct {
	Took     int      `json:"took"`
	TimedOut bool     `json:"timed_out"`
	Shards   Shards   `json:"_shards"`
	Hits     HitsData `json:"hits"`
}

// Shards struct represents the "_shards" JSON object
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// HitsData struct represents the "hits" JSON object
type HitsData struct {
	Total    TotalValue  `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []HitDetail `json:"hits"`
}

// TotalValue struct captures the nested "total" JSON object within "hits"
type TotalValue struct {
	Value int `json:"value"`
}

// HitDetail struct represents each element in the "hits" array
type HitDetail struct {
	Index     string     `json:"_index"`
	Type      string     `json:"_type"`
	ID        string     `json:"_id"`
	Score     float64    `json:"_score"`
	Timestamp string     `json:"@timestamp"`
	Source    ZincMarker `json:"_source"`
}

// ZincMarker struct represents the "_source" nested JSON object within each "hit"
type ZincMarker struct {
	Timestamp string `json:"@timestamp,omitempty"`
	Address   string `json:"address"`
	MarkerID  int    `json:"markerId"`
}

// FuzzSearch represents the structure of the search
type FuzzMarkerSearch struct {
	SearchType   string   `json:"search_type"`
	Query        Query    `json:"query"`
	From         int      `json:"from"`
	MaxResults   int      `json:"max_results"`
	SourceFields []string `json:"_source"`
}

// Query represents the "query" part of the search request
type Query struct {
	Term string `json:"term"`
}

type MarkerSearchResponse struct {
	Took    int          `json:"took"`
	Markers []ZincMarker `json:"markers"`
}

type MarkerIndexData struct {
	MarkerID    int    `json:"markerId"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Address     string `json:"address"` // such as Korean: 경기도 부천시 소사구 경인로29번길 32, 우성아파트
	FullAddress string `json:"fullAddress"`
}
