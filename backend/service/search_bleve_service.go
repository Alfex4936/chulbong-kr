package service

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/blevesearch/bleve/v2"
	bleve_search "github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
)

// Map of Hangul initial consonant Unicode values to their corresponding Korean consonants.
var (
	initialConsonantMap = map[rune]rune{
		0x1100: 'ㄱ', 0x1101: 'ㄲ', 0x1102: 'ㄴ', 0x1103: 'ㄷ', 0x1104: 'ㄸ',
		0x1105: 'ㄹ', 0x1106: 'ㅁ', 0x1107: 'ㅂ', 0x1108: 'ㅃ', 0x1109: 'ㅅ',
		0x110A: 'ㅆ', 0x110B: 'ㅇ', 0x110C: 'ㅈ', 0x110D: 'ㅉ', 0x110E: 'ㅊ',
		0x110F: 'ㅋ', 0x1110: 'ㅌ', 0x1111: 'ㅍ', 0x1112: 'ㅎ',
	}

	validInitialConsonants = map[rune]bool{
		'ㄱ': true, 'ㄲ': true, 'ㄴ': true, 'ㄷ': true, 'ㄸ': true,
		'ㄹ': true, 'ㅁ': true, 'ㅂ': true, 'ㅃ': true, 'ㅅ': true,
		'ㅆ': true, 'ㅇ': true, 'ㅈ': true, 'ㅉ': true, 'ㅊ': true,
		'ㅋ': true, 'ㅌ': true, 'ㅍ': true, 'ㅎ': true,
	}

	doubleConsonants = map[rune][]rune{
		'ㄳ': {'ㄱ', 'ㅅ'}, 'ㄵ': {'ㄴ', 'ㅈ'}, 'ㄶ': {'ㄴ', 'ㅎ'}, 'ㄺ': {'ㄹ', 'ㄱ'},
		'ㄻ': {'ㄹ', 'ㅁ'}, 'ㄼ': {'ㄹ', 'ㅂ'}, 'ㄽ': {'ㄹ', 'ㅅ'}, 'ㄾ': {'ㄹ', 'ㅌ'},
		'ㄿ': {'ㄹ', 'ㅍ'}, 'ㅀ': {'ㄹ', 'ㅎ'}, 'ㅄ': {'ㅂ', 'ㅅ'},
	}
)

type BleveSearchService struct {
	Index bleve.Index
	// Path string
}

func NewBleveSearchService(index bleve.Index) *BleveSearchService {
	return &BleveSearchService{Index: index}
}

// SearchMarkerAddress calls bleve (Lucene-like) search
// func (s *BleveSearchService) SearchMarkerAddress(term string) (dto.MarkerSearchResponse, error) {
// 	var response dto.MarkerSearchResponse
// 	searchQuery := bleve.NewFuzzyQuery(term)
// 	searchRequest := bleve.NewSearchRequest(searchQuery)
// 	searchRequest.From = 0
// 	searchRequest.Size = 10
// 	searchRequest.Fields = []string{"address"} // or *

// 	searchResults, err := s.Index.Search(searchRequest)
// 	if err != nil {
// 		return response, fmt.Errorf("error searching index")
// 	}

// 	response = MarkerSearchResponse{
// 		Took:    int(searchResults.Took.Milliseconds()),
// 		Markers: make([]dto.ZincMarker, 0, len(searchResults.Hits)),
// 	}

// 	// Extract relevant fields from search results
// 	for _, hit := range searchResults.Hits {
// 		var marker dto.ZincMarker
// 		intID, _ := strconv.Atoi(hit.ID)
// 		marker.MarkerID = intID
// 		marker.Address = hit.Fields["address"].(string)
// 		response.Markers = append(response.Markers, marker)
// 	}

// 	return response, nil
// }

// SearchMarkerAddress calls bleve (Lucene-like) search
func (s *BleveSearchService) SearchMarkerAddress(t string) (dto.MarkerSearchResponse, error) {
	// index.Index("test", Marker{Address: "석원", MarkerID: 123, FullAddress: "경기도 석원동 123-456"})
	// index.Delete("test")
	response := dto.MarkerSearchResponse{Markers: make([]dto.ZincMarker, 0)}

	// Split the search term by spaces
	terms := strings.Fields(t)
	results := make(chan *bleve_search.DocumentMatch, 100)
	done := make(chan struct{})
	tookTimes := make(chan time.Duration, len(terms))

	workerCount := 5
	var wg sync.WaitGroup
	wg.Add(workerCount)
	termChan := make(chan string, len(terms))

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for term := range termChan {
				performSearch(s.Index, term, results, tookTimes)
			}
		}()
	}

	for _, term := range terms {
		termChan <- term
	}
	close(termChan)

	go func() {
		wg.Wait()
		close(results)
		close(tookTimes)
		done <- struct{}{}
	}()

	var allResults []*bleve_search.DocumentMatch
	var totalTook time.Duration
	for result := range results {
		if result != nil {
			allResults = append(allResults, result)
		}
	}
	for took := range tookTimes {
		totalTook += took
	}
	<-done

	// Sort and ensure diverse results
	// diverseResults := ensureDiversity(allResults, 10)

	sortResultsByScore(allResults)

	// if len(searchResult.Hits) == 0 {
	// 	return response, nil
	// }

	if len(allResults) == 0 {
		// If no results, try fuzzy search with controlled fuzziness
		for _, term := range terms {
			fuzzyQuery := bleve.NewFuzzyQuery(term)
			fuzzyQuery.Fuzziness = 1
			searchRequest := bleve.NewSearchRequest(fuzzyQuery)
			searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
			searchRequest.Size = 10
			searchResult, err := s.Index.Search(searchRequest)
			if err != nil {
				return response, fmt.Errorf("error fuzzy searching marker")
			}

			for _, hit := range searchResult.Hits {
				allResults = append(allResults, hit)
			}
		}
	} // end of if len(allResults) == 0

	response.Took = int(totalTook.Milliseconds())
	response.Markers = make([]dto.ZincMarker, 0, len(allResults))

	// Extract relevant fields from search results
	for _, hit := range allResults {
		var marker dto.ZincMarker
		intID, _ := strconv.Atoi(hit.ID)
		marker.MarkerID = intID
		marker.Address = hit.Fields["fullAddress"].(string)
		response.Markers = append(response.Markers, marker)
	}

	return response, nil
}

func (s *BleveSearchService) AutoComplete(term string) ([]string, error) {
	var suggestions []string

	prefixQuery := bleve.NewPrefixQuery(term)
	searchRequest := bleve.NewSearchRequest(prefixQuery)
	searchRequest.Fields = []string{"fullAddress", "address", "province", "city"}
	searchRequest.Size = 10 // Limit the number of suggestions

	searchResult, err := s.Index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error performing autocomplete search: %v", err)
	}

	for _, hit := range searchResult.Hits {
		suggestions = append(suggestions, hit.Fields["fullAddress"].(string))
	}

	return suggestions, nil
}

func (s *BleveSearchService) InsertMarkerIndex(indexBody MarkerIndexData) error {
	province, city, rest := splitAddress(indexBody.Address)
	indexBody.Province = province
	indexBody.City = city
	indexBody.FullAddress = indexBody.Address
	indexBody.Address = rest
	indexBody.InitialConsonants = extractInitialConsonants(indexBody.FullAddress)
	err := s.Index.Index(strconv.Itoa(indexBody.MarkerID), indexBody)
	if err != nil {
		return fmt.Errorf("error indexing marker")
	}
	return nil
}

func (s *BleveSearchService) DeleteMarkerIndex(markerId string) error {
	return s.Index.Delete(markerId)
}

func performSearch(index bleve.Index, term string, results chan<- *bleve_search.DocumentMatch, tookTimes chan<- time.Duration) {
	var queries []query.Query

	// Exact match queries with higher boosts
	matchQueryFullAddress := query.NewMatchQuery(term)
	matchQueryFullAddress.SetField("fullAddress")
	matchQueryFullAddress.Analyzer = "koCJKEdgeNgram"
	matchQueryFullAddress.SetBoost(25.0)
	queries = append(queries, matchQueryFullAddress)

	// Phrase match query for exact phrases
	matchPhraseQueryFullAddress := query.NewMatchPhraseQuery(term)
	matchPhraseQueryFullAddress.SetField("fullAddress")
	matchPhraseQueryFullAddress.Analyzer = "koCJKEdgeNgram"
	matchPhraseQueryFullAddress.SetBoost(20.0)
	queries = append(queries, matchPhraseQueryFullAddress)

	// Wildcard queries with lower boosts
	wildcardQueryFullAddress := query.NewWildcardQuery("*" + term + "*")
	wildcardQueryFullAddress.SetField("fullAddress")
	wildcardQueryFullAddress.SetBoost(10.0)
	queries = append(queries, wildcardQueryFullAddress)

	// Additional fields and queries
	brokenConsonants := segmentConsonants(term)
	matchQueryInitialConsonants := query.NewMatchQuery(brokenConsonants)
	matchQueryInitialConsonants.SetField("initialConsonants")
	matchQueryInitialConsonants.Analyzer = "koCJKEdgeNgram"
	matchQueryInitialConsonants.SetBoost(15.0)
	queries = append(queries, matchQueryInitialConsonants)

	wildcardQueryInitialConsonants := query.NewWildcardQuery("*" + brokenConsonants + "*")
	wildcardQueryInitialConsonants.SetField("initialConsonants")
	wildcardQueryInitialConsonants.SetBoost(5.0)
	queries = append(queries, wildcardQueryInitialConsonants)

	standardizedProvince := standardizeProvince(term)
	if standardizedProvince != term {
		matchQueryProvince := query.NewMatchQuery(standardizedProvince)
		matchQueryProvince.SetField("province")
		matchQueryProvince.Analyzer = "koCJKEdgeNgram"
		matchQueryProvince.SetBoost(1.5)
		queries = append(queries, matchQueryProvince)
	} else {
		prefixQueryCity := query.NewPrefixQuery(term)
		prefixQueryCity.SetField("city")
		prefixQueryCity.SetBoost(10.0)
		queries = append(queries, prefixQueryCity)

		matchQueryCity := query.NewMatchQuery(term)
		matchQueryCity.SetField("city")
		matchQueryCity.Analyzer = "koCJKEdgeNgram"
		matchQueryCity.SetBoost(10.0)
		queries = append(queries, matchQueryCity)

		matchQueryDistrict := query.NewMatchQuery(term)
		matchQueryDistrict.SetField("district")
		matchQueryDistrict.Analyzer = "koCJKEdgeNgram"
		matchQueryDistrict.SetBoost(5.0)
		queries = append(queries, matchQueryDistrict)

		prefixQueryAddr := query.NewPrefixQuery(term)
		prefixQueryAddr.SetField("address")
		prefixQueryAddr.SetBoost(5.0)
		queries = append(queries, prefixQueryAddr)

		wildcardQueryAddr := query.NewWildcardQuery("*" + term + "*")
		wildcardQueryAddr.SetField("address")
		wildcardQueryAddr.SetBoost(2.0)
		queries = append(queries, wildcardQueryAddr)
	}

	disjunctionQuery := bleve.NewDisjunctionQuery(queries...)
	searchRequest := bleve.NewSearchRequest(disjunctionQuery)
	searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
	searchRequest.Size = 10
	searchRequest.SortBy([]string{"_score", "markerId"})

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("Error performing search: %v", err)
		return
	}

	tookTimes <- searchResult.Took
	for _, hit := range searchResult.Hits {
		results <- hit
	}
}

// extractInitialConsonants extracts the initial consonants from a Korean string.
//
// ex) "부산 해운대구 좌동 1395" -> "ㅂㅅㅎㅇㄷㄱㅈㄷ"
func extractInitialConsonants(s string) string {
	var initials []rune
	for _, r := range s {
		if unicode.Is(unicode.Hangul, r) {
			initial := (r - 0xAC00) / 28 / 21
			if mapped, exists := initialConsonantMap[0x1100+initial]; exists {
				initials = append(initials, mapped)
			}
		}
	}
	return string(initials)
}

// Split the user input into valid Korean initial consonants, breaking double consonants where necessary
//
// ex) "앍돍ㅄㄳ산" -> "앍돍ㅂㅅㄱㅅ산"
func segmentConsonants(input string) string {
	var result []rune

	for _, r := range input {
		// Check if the character is a valid initial consonant
		if validInitialConsonants[r] {
			result = append(result, r)
		} else if components, found := doubleConsonants[r]; found {
			// If it's a double consonant, break it into its components
			result = append(result, components...)
		} else {
			// If it's not a valid consonant or double consonant, add it as is
			result = append(result, r)
		}
	}

	return string(result)
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

func ensureDiversity(results []*bleve_search.DocumentMatch, limit int) []*bleve_search.DocumentMatch {
	addressSet := make(map[string]struct{})
	var diverseResults []*bleve_search.DocumentMatch

	for _, result := range results {
		if len(diverseResults) >= limit {
			break
		}
		address := result.Fields["fullAddress"].(string)
		if _, exists := addressSet[address]; !exists {
			addressSet[address] = struct{}{}
			diverseResults = append(diverseResults, result)
		}
	}

	return diverseResults
}

func sortResultsByScore(results []*bleve_search.DocumentMatch) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}
