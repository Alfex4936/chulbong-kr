package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	gounicode "unicode"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/char/html"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	"github.com/blevesearch/bleve/v2/analysis/token/edgengram"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	"github.com/blevesearch/bleve/v2/analysis/token/unicodenorm"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	_ "github.com/blevesearch/bleve/v2/index/upsidedown/store/boltdb"
	"github.com/blevesearch/bleve/v2/search/highlight/format/html"
	"github.com/blevesearch/bleve/v2/search/query"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// MarkerDB represents the structure of the marker data from DB.
type MarkerDB struct {
	MarkerID int    `json:"markerId"`
	Address  string `json:"address"`
}

type Marker struct {
	MarkerID          int    `json:"markerId"`
	Province          string `json:"province"`
	City              string `json:"city"`
	Address           string `json:"address"` // such as Korean: 경기도 부천시 소사구 경인로29번길 32, 우성아파트
	FullAddress       string `json:"fullAddress"`
	InitialConsonants string `json:"initialConsonants"` // 초성
}

// Load environment variables from .env file
func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func saveJson() error {
	// Database connection parameters
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}
	defer db.Close()

	// Query to select all rows from the Markers table
	selectSQL := `SELECT MarkerID, Address FROM Markers`

	// Execute the query
	rows, err := db.Query(selectSQL)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Prepare data for JSON
	var markers []MarkerDB
	for rows.Next() {
		var marker MarkerDB
		err := rows.Scan(&marker.MarkerID, &marker.Address)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
		markers = append(markers, marker)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %v", err)
	}

	// Write to JSON file
	file, err := os.Create("markers.json")
	if err != nil {
		return fmt.Errorf("error creating JSON file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(markers)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	fmt.Println("Data has been written to markers.json")

	return nil
}

func main() {
	// Load environment variables
	err := loadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if len(os.Args) < 2 {
		saveJson()
		log.Println("Fetching completed.")
	}

	// Sample JSON data
	markers, _ := getMarkersFromJson("markers.json")

	// Create a new Bleve index
	// Define index path
	indexPath := "markers.bleve"

	// Define custom mapping
	indexMapping := bleve.NewIndexMapping()
	err = indexMapping.AddCustomTokenFilter("edge_ngram_min_1_max_4",
		map[string]interface{}{
			"type": edgengram.Name,
			"min":  1.0,
			"max":  4.0,
		})
	if err != nil {
		log.Fatalf("Error creating custom analyzer: %v", err)
	}

	err = indexMapping.AddCustomTokenFilter("unicodeNormalizer",
		map[string]interface{}{
			"type": unicodenorm.Name,
			"form": unicodenorm.NFKC,
		})
	if err != nil {
		log.Fatalf("Error creating custom token filter: %v", err)
	}

	// normalize -> CJK bigrams -> edge ngrams -> lowercase
	err = indexMapping.AddCustomAnalyzer("koCJKEdgeNgram",
		map[string]interface{}{
			"type":         custom.Name,
			"tokenizer":    unicode.Name,
			"char_filters": []string{html.Name},
			"token_filters": []string{
				"unicodeNormalizer",
				"cjk_bigram",
				"edge_ngram_min_1_max_4",
				lowercase.Name,
			},
		})
	if err != nil {

		log.Fatalf("Error creating custom analyzer: %v", err)
	}

	markerMapping := bleve.NewDocumentMapping()

	markerIDMapping := bleve.NewNumericFieldMapping()
	markerMapping.AddFieldMappingsAt("markerId", markerIDMapping)

	provinceMapping := bleve.NewTextFieldMapping()
	provinceMapping.Analyzer = "koCJKEdgeNgram"
	provinceMapping.Store = true
	provinceMapping.IncludeTermVectors = true
	markerMapping.AddFieldMappingsAt("province", provinceMapping)

	cityMapping := bleve.NewTextFieldMapping()
	cityMapping.Analyzer = "koCJKEdgeNgram"
	cityMapping.Store = true
	cityMapping.IncludeTermVectors = true
	markerMapping.AddFieldMappingsAt("city", cityMapping)

	addressMapping := bleve.NewTextFieldMapping()
	addressMapping.Analyzer = "koCJKEdgeNgram"
	addressMapping.Store = true
	addressMapping.IncludeTermVectors = true
	markerMapping.AddFieldMappingsAt("address", addressMapping)

	fullAddressMapping := bleve.NewTextFieldMapping()
	fullAddressMapping.Store = true
	fullAddressMapping.IncludeTermVectors = true
	markerMapping.AddFieldMappingsAt("fullAddress", fullAddressMapping)

	initialConsonantsMapping := bleve.NewTextFieldMapping()
	initialConsonantsMapping.Analyzer = "koCJKEdgeNgram"
	initialConsonantsMapping.Store = true
	initialConsonantsMapping.IncludeTermVectors = true
	markerMapping.AddFieldMappingsAt("initialConsonants", initialConsonantsMapping)

	indexMapping.AddDocumentMapping("marker", markerMapping)

	// Create a new Bleve index with custom settings
	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		if len(os.Args) < 2 {
			search("영통역")
		} else {
			searchTerm := strings.Join(os.Args[1:], " ")
			search(searchTerm)
		}

		log.Fatalf("Error creating index: %v", err)
	}

	// Index the markers
	for i, marker := range markers {
		province, city, rest := splitAddress(marker.Address)
		marker.Province = province
		marker.City = city
		marker.FullAddress = marker.Address
		marker.Address = rest
		marker.InitialConsonants = extractInitialConsonants(marker.FullAddress)
		if i == 10 {
			log.Printf("🌍 address: %s (%s)", marker.FullAddress, marker.InitialConsonants)
		}
		err = index.Index(strconv.Itoa(marker.MarkerID), marker)
		if err != nil {
			log.Fatalf("Error indexing document: %v", err)
		}
	}

	log.Println("Indexing completed")

	index.Close()

	// Perform a search using DisjunctionQuery for more comprehensive matching
	search("해운대")

}

func search(t string) {
	index, _ := bleve.Open("markers.bleve")
	defer index.Close()

	// index.Index("test", Marker{Address: "석원", MarkerID: 123, FullAddress: "경기도 석원동 123-456"})
	// index.Delete("test")

	// Capture the start time
	start := time.Now()

	// Split the search term by spaces
	terms := strings.Fields(t)
	var queries []query.Query

	// Add a phrase query for the full search term
	// phraseQuery := query.NewMatchPhraseQuery(t)
	// phraseQuery.SetBoost(3.0)
	// phraseQuery.Analyzer = "koCJKEdgeNgram"
	// queries = append(queries, phraseQuery)

	for _, term := range terms {
		// Add a MatchQuery for the full search term
		matchQuery := query.NewMatchQuery(term)
		matchQuery.SetField("initialConsonants")
		matchQuery.Analyzer = "koCJKEdgeNgram"
		matchQuery.SetBoost(10.0)
		queries = append(queries, matchQuery)

		// Use WildcardQuery for more flexible matches
		wildcardQuery := query.NewWildcardQuery("*" + term + "*")
		wildcardQuery.SetField("initialConsonants")
		wildcardQuery.SetBoost(5.0)
		queries = append(queries, wildcardQuery)

		standardizedProvince := standardizeProvince(term)
		if standardizedProvince != term {
			// If the term is a province, use a lower boost
			matchQuery := query.NewMatchQuery(standardizedProvince)
			matchQuery.SetField("province")
			matchQuery.Analyzer = "koCJKEdgeNgram"
			matchQuery.SetBoost(1.5)
			queries = append(queries, matchQuery)
		} else {
			// Use PrefixQuery for cities and regions
			prefixQueryCity := query.NewPrefixQuery(term)
			prefixQueryCity.SetField("city")
			prefixQueryCity.SetBoost(10.0)
			queries = append(queries, prefixQueryCity)

			// Use MatchPhraseQuery for detailed matches in full address
			matchPhraseQueryFull := query.NewMatchPhraseQuery(term)
			matchPhraseQueryFull.SetField("fullAddress")
			matchPhraseQueryFull.Analyzer = "koCJKEdgeNgram"
			matchPhraseQueryFull.SetBoost(5.0)
			queries = append(queries, matchPhraseQueryFull)

			// Use WildcardQuery for more flexible matches
			wildcardQueryFull := query.NewWildcardQuery("*" + term + "*")
			wildcardQueryFull.SetField("fullAddress")
			wildcardQueryFull.SetBoost(2.0)
			queries = append(queries, wildcardQueryFull)

			// Additional PrefixQuery and WildcardQuery for other fields
			prefixQueryAddr := query.NewPrefixQuery(term)
			prefixQueryAddr.SetField("address")
			prefixQueryAddr.SetBoost(5.0)
			queries = append(queries, prefixQueryAddr)

			wildcardQueryAddr := query.NewWildcardQuery("*" + term + "*")
			wildcardQueryAddr.SetField("address")
			wildcardQueryAddr.SetBoost(2.0)
			queries = append(queries, wildcardQueryAddr)

			// Use MatchQuery for city and district to catch all matches
			matchQueryCity := query.NewMatchQuery(term)
			matchQueryCity.SetField("city")
			matchQueryCity.Analyzer = "koCJKEdgeNgram"
			matchQueryCity.SetBoost(5.0)
			queries = append(queries, matchQueryCity)

			matchQueryDistrict := query.NewMatchQuery(term)
			matchQueryDistrict.SetField("district")
			matchQueryDistrict.Analyzer = "koCJKEdgeNgram"
			matchQueryDistrict.SetBoost(5.0)
			queries = append(queries, matchQueryDistrict)
		}
	}

	disjunctionQuery := bleve.NewDisjunctionQuery(queries...)
	// conjunctionQuery := bleve.NewConjunctionQuery(disjunctionQuery)

	searchRequest := bleve.NewSearchRequest(disjunctionQuery)
	searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
	searchRequest.Size = 10 // Limit the number of results
	searchResult, err := index.Search(searchRequest)
	searchRequest.SortBy([]string{"_score", "markerId"})
	if err != nil {
		log.Fatalf("Error performing search: %v", err)
	}

	if len(searchResult.Hits) > 0 {
		for _, hit := range searchResult.Hits {
			log.Printf("Search Result: %+v", hit)
			log.Printf("Document ID: %v", hit.ID)
			log.Printf("Document Score: %v", hit.Score)
			if fullAddress, ok := hit.Fields["fullAddress"]; ok {
				log.Printf("Full Address: %v", fullAddress)
			}
			// if province, ok := hit.Fields["province"]; ok {
			// 	log.Printf("Province: %v", province)
			// }
			// if city, ok := hit.Fields["city"]; ok {
			// 	log.Printf("City: %v", city)
			// }
			// if address, ok := hit.Fields["address"]; ok {
			// 	log.Printf("Address: %v", address)
			// }
		}
	} else {
		log.Println("No search results found")
	}

	if len(searchResult.Hits) > 0 {
		for _, hit := range searchResult.Hits {
			log.Printf("Search Result: %+v", hit)
			log.Printf("Document ID: %v", hit.ID)
			log.Printf("Document Score: %v", hit.Score)
			if fullAddress, ok := hit.Fields["fullAddress"]; ok {
				log.Printf("Full Address: %v", fullAddress)
			}
		}
	} else {
		log.Println("No search results found")

		// If no results, try fuzzy search with controlled fuzziness
		for _, term := range terms {
			fuzzyQuery := bleve.NewFuzzyQuery(term)
			fuzzyQuery.Fuzziness = 1
			searchRequest = bleve.NewSearchRequest(fuzzyQuery)
			searchRequest.Fields = []string{"fullAddress", "address", "province", "city"}
			searchRequest.Size = 10
			searchResult, err = index.Search(searchRequest)
			if err != nil {
				log.Fatalf("Error performing fuzzy search: %v", err)
			}
			if len(searchResult.Hits) > 0 {
				for _, hit := range searchResult.Hits {
					log.Printf("Search Result: %+v", hit)
					log.Printf("Document ID: %v", hit.ID)
					log.Printf("Document Score: %v", hit.Score)
					if fullAddress, ok := hit.Fields["fullAddress"]; ok {
						log.Printf("Full Address: %v", fullAddress)
					}
				}
			}
		}
	}

	// Capture the end time
	end := time.Now()
	// Calculate the duration
	duration := end.Sub(start)

	// Log the duration
	log.Printf("📆Function runtime: %v", duration)
}

func standardizeProvince(province string) string {
	switch province {
	case "경기", "경기도":
		return "경기도"
	case "서울", "서울특별시":
		return "서울특별시"
	case "부산", "부산광역시":
		return "부산광역시"
	case "대구", "대구광역시":
		return "대구광역시"
	case "인천", "인천광역시":
		return "인천광역시"
	case "제주", "제주특별자치도", "제주도":
		return "제주특별자치도"
	case "대전", "대전광역시":
		return "대전광역시"
	case "울산", "울산광역시":
		return "울산광역시"
	case "광주", "광주광역시":
		return "광주광역시"
	case "세종", "세종특별자치시":
		return "세종특별자치시"
	case "강원", "강원도", "강원특별자치도":
		return "강원특별자치도"
	case "경남", "경상남도":
		return "경상남도"
	case "경북", "경상북도":
		return "경상북도"
	case "전북", "전북특별자치도":
		return "전북특별자치도"
	case "충남", "충청남도":
		return "충청남도"
	case "충북", "충청북도":
		return "충청북도"
	case "전남", "전라남도":
		return "전라남도"
	default:
		return province
	}
}

func splitAddress(address string) (string, string, string) {
	parts := strings.Fields(address)
	if len(parts) < 2 {
		return "", "", address
	}
	province := standardizeProvince(parts[0])
	city := parts[1]
	rest := strings.Join(parts[2:], " ")
	return province, city, rest
}

func getMarkersFromJson(filepath string) ([]Marker, error) {
	var markerData []Marker

	file, err := os.Open(filepath)
	if err != nil {
		return markerData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&markerData)
	if err != nil {
		return markerData, err
	}

	return markerData, nil
}

// Map of Hangul initial consonant Unicode values to their corresponding Korean consonants.
var initialConsonantMap = map[rune]rune{
	0x1100: 'ㄱ', 0x1101: 'ㄲ', 0x1102: 'ㄴ', 0x1103: 'ㄷ', 0x1104: 'ㄸ',
	0x1105: 'ㄹ', 0x1106: 'ㅁ', 0x1107: 'ㅂ', 0x1108: 'ㅃ', 0x1109: 'ㅅ',
	0x110A: 'ㅆ', 0x110B: 'ㅇ', 0x110C: 'ㅈ', 0x110D: 'ㅉ', 0x110E: 'ㅊ',
	0x110F: 'ㅋ', 0x1110: 'ㅌ', 0x1111: 'ㅍ', 0x1112: 'ㅎ',
}

// ExtractInitialConsonants extracts the initial consonants from a Korean string.
// ex) "부산 해운대구 좌동 1395" -> "ㅂㅅㅎㅇㄷㄱㅈㄷ"
func extractInitialConsonants(s string) string {
	var initials []rune
	for _, r := range s {
		if gounicode.Is(gounicode.Hangul, r) {
			initial := (r - 0xAC00) / 28 / 21
			if mapped, exists := initialConsonantMap[0x1100+initial]; exists {
				initials = append(initials, mapped)
			}
		}
	}
	return string(initials)
}

/*
my json looks like

[
  {
    "markerId": 1,
    "address": "경북 포항시 북구 창포동 655"
  },
  {
    "markerId": 2,
    "address": "서울 노원구 하계동 250"
  },
]
*/
