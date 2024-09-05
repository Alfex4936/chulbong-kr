package service

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/dkssud"
	"github.com/blevesearch/bleve/v2"
	bleve_search "github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	gocache "github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
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

	documentMatchPool = &sync.Pool{
		New: func() any {
			slice := make([]*bleve_search.DocumentMatch, 0, 100)
			return &slice
		},
	}

	termsPool = &sync.Pool{
		New: func() any {
			slice := make([]string, 0, 30)
			return &slice
		},
	}
)

type BleveSearchService struct {
	Index  bleve.Index
	Shards []bleve.Index

	LocalCacheStorage *ristretto_store.RistrettoStore
	Logger            *zap.Logger
	DB                *sqlx.DB
	GetAllMarkersStmt *sqlx.Stmt

	searchCache *gocache.Cache[dto.MarkerSearchResponse]
}

func NewBleveSearchService(
	index bleve.Index, shards []bleve.Index,
	localCacheStorage *ristretto_store.RistrettoStore, logger *zap.Logger,
	db *sqlx.DB) *BleveSearchService {
	searchCache := gocache.New[dto.MarkerSearchResponse](localCacheStorage)

	getMarkerStmt, _ := db.Preparex("SELECT MarkerID, Address FROM Markers")

	return &BleveSearchService{Index: index, Shards: shards,
		searchCache: searchCache, Logger: logger, DB: db,
		GetAllMarkersStmt: getMarkerStmt,
	}
}

func RegisteBleveLifecycle(lifecycle fx.Lifecycle, service *BleveSearchService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			service.GetAllMarkersStmt.Close()
			return nil
		},
	})
}

// SearchMarkerAddress calls bleve (Lucene-like) search
func (s *BleveSearchService) SearchMarkerAddress(t string) (dto.MarkerSearchResponse, error) {
	// t is already trimmed
	cacheKey := fmt.Sprintf("search:%s", t)
	cachedResponse, err := s.searchCache.Get(context.Background(), cacheKey)
	if err == nil {
		return cachedResponse, nil
	}

	// 쿼티 한글? -> 한글로 변환 (ex. "rudrleh" -> "경기도")
	if dkssud.IsQwertyHangul(t) {
		t = dkssud.QwertyToHangul(t)
	}

	response := dto.MarkerSearchResponse{Markers: make([]dto.ZincMarker, 0)}

	// Get a pointer to a slice from the pool
	termsPtr := termsPool.Get().(*[]string)
	terms := (*termsPtr)[:0] // Reset the slice

	// Split the search term by spaces and append to the pooled slice
	terms = append(terms, strings.Fields(t)...)
	if len(terms) == 0 {
		termsPool.Put(termsPtr)
		return response, nil
	}

	terms[0] = standardizeInitials(terms[0])

	results := make(chan *bleve_search.DocumentMatch, len(terms)*10)
	tookTimes := make(chan time.Duration, len(terms))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	termChan := make(chan string, len(terms))

	// -- Use runtime.GOMAXPROCS(0) to determine the number of workers
	workerCount := runtime.NumCPU() // Limit goroutines to number of CPU cores
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for term := range termChan {
				select {
				case <-ctx.Done():
					return
				default:
					performSearch(s.Index, term, results, tookTimes)
				}
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
	}()

	// Get a pointer to a slice from the pool
	allResultsPtr := documentMatchPool.Get().(*[]*bleve_search.DocumentMatch)
	allResults := (*allResultsPtr)[:0] // Reset the slice
	var totalTook time.Duration

	for {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
			} else if result != nil {
				allResults = append(allResults, result)
			}
		case took, ok := <-tookTimes:
			if !ok {
				tookTimes = nil
			} else {
				totalTook += took
			}
		}

		if results == nil && tookTimes == nil {
			break
		}
	}

	// Sort and ensure diverse results
	// diverseResults := ensureDiversity(allResults, 10)

	if len(allResults) == 0 {
		// If no results, try fuzzy search with controlled fuzziness
		fuzzyResults, fuzzyTook := performFuzzySearch(s.Index, terms)
		allResults = fuzzyResults
		totalTook += fuzzyTook
	}

	sortResultsByScore(allResults)
	response.Took = int(totalTook.Milliseconds())
	response.Markers = extractMarkers(allResults)

	// Cache the response
	s.searchCache.Set(context.Background(), cacheKey, response)

	// Reset the slice and return the pointer to the pool
	*allResultsPtr = allResults[:0]
	documentMatchPool.Put(allResultsPtr)

	// Return terms slice to the pool
	*termsPtr = terms[:0]
	termsPool.Put(termsPtr)

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
	// Compute which shard to use based on the marker ID (or any other key)
	shardIndex := indexBody.MarkerID % len(s.Shards)
	selectedShard := s.Shards[shardIndex]

	province, city, rest := splitAddress(indexBody.Address)
	indexBody.Province = province
	indexBody.City = city
	indexBody.FullAddress = indexBody.Address
	indexBody.Address = rest
	indexBody.InitialConsonants = ExtractInitialConsonants(indexBody.FullAddress)

	err := selectedShard.Index(strconv.Itoa(indexBody.MarkerID), indexBody)
	if err != nil {
		return fmt.Errorf("error indexing marker: %v", err)
	}

	// Invalidate search cache
	s.InvalidateCache()

	defer s.Logger.Info("New Marker indexed", zap.Int("markerID", indexBody.MarkerID), zap.String("address", indexBody.FullAddress))

	return nil
}

func (s *BleveSearchService) DeleteMarkerIndex(markerId int) error {
	// Compute which shard to use based on the marker ID (same as in InsertMarkerIndex)
	shardIndex := markerId % len(s.Shards)
	selectedShard := s.Shards[shardIndex]

	// Perform the deletion operation on the selected shard
	err := selectedShard.Delete(strconv.Itoa(markerId))
	if err != nil {
		return fmt.Errorf("error deleting marker: %v", err)
	}

	// Invalidate search cache
	s.InvalidateCache()

	return nil
}

func (s *BleveSearchService) InvalidateCache() {
	s.searchCache.Clear(context.Background())
}

func (s *BleveSearchService) CheckIndexes() error {
	// Step 1: Fetch all valid marker IDs from the database along with their addresses
	markers, err := s.GetAllMarkers()
	if err != nil {
		return fmt.Errorf("error fetching valid marker IDs from database: %v", err)
	}

	// Convert marker data to a map for quick lookup by MarkerID
	markerDataMap := make(map[int]dto.MarkerIndexData, len(markers))
	for _, marker := range markers {
		markerDataMap[marker.MarkerID] = dto.MarkerIndexData{
			MarkerID: marker.MarkerID,
			Address:  marker.Address,
		}
	}

	// Step 2: Execute a MatchAllQuery to retrieve all documents from the index
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 10000

	searchResult, err := s.Index.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("error executing search: %v", err)
	}

	// Step 3: Iterate through all hits (documents) in the result
	for _, hit := range searchResult.Hits {
		markerID, err := strconv.Atoi(hit.ID)
		if err != nil {
			s.Logger.Error("Failed to convert document ID to markerID", zap.String("docID", hit.ID), zap.Error(err))
			continue
		}

		// Check if this marker exists in the database
		if _, existsInDB := markerDataMap[markerID]; !existsInDB {
			// Step 4: If the MarkerID doesn't exist in the database, delete it from the index
			err = s.DeleteMarkerIndex(markerID)
			if err != nil {
				s.Logger.Error("Failed to delete orphaned index", zap.Int("markerID", markerID), zap.Error(err))
			} else {
				s.Logger.Info("Orphaned index deleted", zap.Int("markerID", markerID))
			}
		}
	}

	// Step 5: Check if each valid marker from the database exists in the index, and index it if missing
	for markerID, markerData := range markerDataMap {
		exists, err := s.MarkerExists(markerID)
		if err != nil {
			s.Logger.Error("Error checking marker existence in index", zap.Int("markerID", markerID), zap.Error(err))
			continue
		}

		// If the marker doesn't exist in the index but is in the database, index it
		if !exists {
			err = s.InsertMarkerIndex(markerData)
			if err != nil {
				s.Logger.Error("Failed to index marker", zap.Int("markerID", markerID), zap.Error(err))
			} else {
				s.Logger.Info("Marker indexed", zap.Int("markerID", markerID))
			}
		}
	}

	return nil
}

// GetAllMarkers now returns a simplified list of markers
func (s *BleveSearchService) GetAllMarkers() ([]dto.MarkerOnlyWithAddr, error) {
	var markers []dto.MarkerOnlyWithAddr

	err := s.GetAllMarkersStmt.Select(&markers)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	// go s.MarkerLocationService.Redis.AddGeoMarkers(markers)

	return markers, nil
}

func (s *BleveSearchService) MarkerExists(markerID int) (bool, error) {
	// Iterate over all shards
	for _, shard := range s.Shards {
		// Fetch the document from the shard
		document, err := shard.Document(strconv.Itoa(markerID))
		if err != nil {
			return false, fmt.Errorf("error fetching document for markerID %d: %v", markerID, err)
		}

		// If the document is found, return true
		if document != nil {
			return true, nil
		}
	}

	// If no document was found in any shard, return false
	return false, nil
}

func performSearch(index bleve.Index, term string, results chan<- *bleve_search.DocumentMatch, tookTimes chan<- time.Duration) {
	var queries []query.Query
	const analyzer = "koCJKEdgeNgram"

	// Use a buffer pool for reduced allocations
	bufferPool := sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 256) // Increased initial capacity
			return &b
		},
	}

	bufferPtr := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(bufferPtr)
	buffer := (*bufferPtr)[:0]

	// Helper function to reset and reuse the buffer
	resetBuffer := func() {
		buffer = buffer[:0]
	}

	// Helper function to append string to buffer and get string back
	appendString := func(s string) string {
		resetBuffer()
		buffer = append(buffer, s...)
		return string(buffer)
	}

	// Exact match queries with higher boosts
	matchQueryFullAddress := query.NewMatchQuery(appendString(term))
	matchQueryFullAddress.SetField("fullAddress")
	matchQueryFullAddress.Analyzer = analyzer
	matchQueryFullAddress.SetBoost(25.0)
	queries = append(queries, matchQueryFullAddress)

	// Phrase match query for exact phrases
	matchPhraseQueryFullAddress := query.NewMatchPhraseQuery(appendString(term))
	matchPhraseQueryFullAddress.SetField("fullAddress")
	matchPhraseQueryFullAddress.Analyzer = analyzer
	matchPhraseQueryFullAddress.SetBoost(20.0)
	queries = append(queries, matchPhraseQueryFullAddress)

	// Wildcard queries with lower boosts
	resetBuffer()
	buffer = append(buffer, '*')
	buffer = append(buffer, term...)
	buffer = append(buffer, '*')
	wildcardQueryFullAddress := query.NewWildcardQuery(string(buffer))
	wildcardQueryFullAddress.SetField("fullAddress")
	wildcardQueryFullAddress.SetBoost(15.0)
	queries = append(queries, wildcardQueryFullAddress)

	prefixQueryFullAddress := query.NewPrefixQuery(appendString(term))
	prefixQueryFullAddress.SetField("fullAddress")
	prefixQueryFullAddress.SetBoost(35.0)
	queries = append(queries, prefixQueryFullAddress)

	// Additional fields and queries
	brokenConsonants := appendString(SegmentConsonants(term))
	matchQueryInitialConsonants := query.NewMatchQuery(brokenConsonants)
	matchQueryInitialConsonants.SetField("initialConsonants")
	matchQueryInitialConsonants.Analyzer = analyzer
	matchQueryInitialConsonants.SetBoost(15.0)
	queries = append(queries, matchQueryInitialConsonants)

	resetBuffer()
	buffer = append(buffer, '*')
	buffer = append(buffer, brokenConsonants...)
	buffer = append(buffer, '*')
	wildcardQueryInitialConsonants := query.NewWildcardQuery(string(buffer))
	wildcardQueryInitialConsonants.SetField("initialConsonants")
	wildcardQueryInitialConsonants.SetBoost(7.0)
	queries = append(queries, wildcardQueryInitialConsonants)

	prefixQueryInitialConsonants := query.NewPrefixQuery(appendString(brokenConsonants))
	prefixQueryInitialConsonants.SetField("initialConsonants")
	prefixQueryInitialConsonants.SetBoost(25.0)
	queries = append(queries, prefixQueryInitialConsonants)

	standardizedProvince := appendString(standardizeProvince(term))
	if standardizedProvince != term {
		matchQueryProvince := query.NewPrefixQuery(standardizedProvince)
		matchQueryProvince.SetField("province")
		matchQueryProvince.SetBoost(3.0)
		queries = append(queries, matchQueryProvince)
	} else {
		prefixQueryCity := query.NewPrefixQuery(appendString(term))
		prefixQueryCity.SetField("city")
		prefixQueryCity.SetBoost(10.0)
		queries = append(queries, prefixQueryCity)

		matchQueryCity := query.NewMatchQuery(appendString(term))
		matchQueryCity.SetField("city")
		matchQueryCity.Analyzer = analyzer
		matchQueryCity.SetBoost(10.0)
		queries = append(queries, matchQueryCity)

		matchQueryDistrict := query.NewMatchQuery(appendString(term))
		matchQueryDistrict.SetField("district")
		matchQueryDistrict.Analyzer = analyzer
		matchQueryDistrict.SetBoost(5.0)
		queries = append(queries, matchQueryDistrict)

		prefixQueryAddr := query.NewPrefixQuery(appendString(term))
		prefixQueryAddr.SetField("address")
		prefixQueryAddr.SetBoost(10.0)
		queries = append(queries, prefixQueryAddr)

		resetBuffer()
		buffer = append(buffer, '*')
		buffer = append(buffer, term...)
		buffer = append(buffer, '*')
		wildcardQueryAddr := query.NewWildcardQuery(string(buffer))
		wildcardQueryAddr.SetField("address")
		wildcardQueryAddr.SetBoost(5.0)
		queries = append(queries, wildcardQueryAddr)
	}

	disjunctionQuery := bleve.NewDisjunctionQuery(queries...)
	searchRequest := bleve.NewSearchRequestOptions(disjunctionQuery, 15, 0, false)
	searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
	searchRequest.Size = 10
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.SortBy([]string{"_score", "markerId"})

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return
	}

	tookTimes <- searchResult.Took
	for _, hit := range searchResult.Hits {
		results <- hit
	}
}

// Perform search with facets
func performSearchFacet(index bleve.Index, term string, results chan<- *bleve_search.DocumentMatch, tookTimes chan<- time.Duration) {
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
	brokenConsonants := SegmentConsonants(term)
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
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
	searchRequest.Size = 10
	searchRequest.SortBy([]string{"_score", "markerId"})

	// Add facet request for province
	facetRequest := bleve.NewFacetRequest("province", 10)
	searchRequest.AddFacet("province_facets", facetRequest)

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("Error performing search: %v", err)
		return
	}

	tookTimes <- searchResult.Took
	for _, hit := range searchResult.Hits {
		results <- hit
	}

	// Print facet results
	if facet, found := searchResult.Facets["province_facets"]; found {
		fmt.Printf("Facet Results for 'province':\n")
		fmt.Printf("Total: %d\n", facet.Total)
		fmt.Printf("Missing: %d\n", facet.Missing)
		for _, term := range facet.Terms.Terms() {
			fmt.Printf("Term: %s, Count: %d\n", term.Term, term.Count)
		}
	}
}

func performFuzzySearch(index bleve.Index, terms []string) ([]*bleve_search.DocumentMatch, time.Duration) {
	var allResults []*bleve_search.DocumentMatch
	var totalTook time.Duration

	for _, term := range terms {
		fuzzyQuery := bleve.NewFuzzyQuery(term)
		fuzzyQuery.Fuzziness = 1
		searchRequest := bleve.NewSearchRequest(fuzzyQuery)
		searchRequest.Fields = []string{"fullAddress", "address", "province", "city", "initialConsonants"}
		searchRequest.Size = 10
		searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
		searchResult, err := index.Search(searchRequest)
		if err != nil {
			log.Printf("Error fuzzy searching marker: %v", err)
			continue
		}
		allResults = append(allResults, searchResult.Hits...)
		totalTook += searchResult.Took
	}

	return allResults, totalTook
}

func extractMarkers(allResults []*bleve_search.DocumentMatch) []dto.ZincMarker {
	markers := make([]dto.ZincMarker, 0, len(allResults))
	for _, hit := range allResults {
		intID, _ := strconv.Atoi(hit.ID)
		markers = append(markers, dto.ZincMarker{
			MarkerID: intID,
			Address:  hit.Fields["fullAddress"].(string),
		})
	}
	return markers
}

// extractInitialConsonants extracts the initial consonants from a Korean string.
//
// ex) "부산 해운대구 좌동 1395" -> "ㅂㅅㅎㅇㄷㄱㅈㄷ"
func ExtractInitialConsonants(s string) string {
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
func SegmentConsonants(input string) string {
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
	case "경기", "경기도", "ㄱㄱㄷ":
		return "경기도"
	case "서울", "서울특별시", "ㅅㅇ", "ㅅㅇㅌㅂㅅ":
		return "서울특별시"
	case "부산", "부산광역시", "ㅄ":
		return "부산광역시"
	case "대구", "대구광역시", "ㄷㄱ":
		return "대구광역시"
	case "인천", "인천광역시", "ㅇㅊ":
		return "인천광역시"
	case "제주", "제주특별자치도", "제주도", "ㅈㅈㄷ":
		return "제주특별자치도"
	case "대전", "대전광역시":
		return "대전광역시"
	case "울산", "울산광역시":
		return "울산광역시"
	case "광주", "광주광역시":
		return "광주광역시"
	case "세종", "세종특별자치시":
		return "세종특별자치시"
	case "강원", "강원도", "강원특별자치도", "ㄱㅇㄷ":
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

func standardizeInitials(initials string) string {
	switch initials {
	case "ㄱㄱ":
		return "ㄱㄱㄷ"
	case "ㅅㅇ", "ㅅㅇㅅ":
		return "ㅅㅇㅌㅂㅅ"
	case "ㅄ", "ㅄㅅ", "ㅂㅅ":
		return "ㅂㅅㄱㅇㅅ"
	case "ㄷㄱ":
		return "ㄷㄱㄱㅇㅅ"
	case "ㅇㅊㅅ":
		return "ㅇㅊㄱㅇㅅ"
	case "ㅈㅈㄷ", "ㅈㅈ":
		return "ㅈㅈㅌㅂㅈㅊㄷ"
	case "ㄷㅈ":
		return "ㄷㅈㄱㅇㅅ"
	case "ㅇㅅ":
		return "ㅇㅅㄱㅇㅅ"
	case "ㄱㅈ":
		return "ㄱㅈㄱㅇㅅ"
	case "ㅅㅈㅅ":
		return "ㅅㅈㅌㅂㅈㅊㅅ"
	case "ㄱㅇㄷ":
		return "ㄱㅇㅌㅂㅈㅊㄷ"
	case "ㄳㄴㄷ", "ㄱㅅㄴㄷ":
		return "ㄱㅅㄴㄷ"
	case "ㄳㅂㄷ", "ㄱㅅㅂㄷ":
		return "ㄱㅅㅂㄷ"
	case "ㅈㅂ":
		return "ㅈㅂㅌㅂㅈㅊㄷ"
	case "ㅊㄴ":
		return "ㅊㅊㄴㄷ"
	case "ㅊㅂ":
		return "ㅊㅊㅂㄷ"
	case "ㅈㄴ":
		return "ㅈㄹㄴㄷ"
	default:
		return initials
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
