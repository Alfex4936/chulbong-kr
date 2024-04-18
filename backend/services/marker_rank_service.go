package services

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/axiomhq/hyperloglog"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/redis/rueidis"
	"github.com/zeebo/xxh3"

	"chulbong-kr/database"
	"chulbong-kr/dto"
)

// 클릭 이벤트를 저장할 임시 저장소
var clickEventBuffer = csmap.Create(
	csmap.WithShardCount[int, int](64),
	csmap.WithCustomHasher[int, int](func(key int) uint64 {
		// Convert int to a byte slice
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(key))
		return xxh3.Hash(bs)
	}),
)

var SketchedLocations = csmap.Create(
	csmap.WithShardCount[string, *hyperloglog.Sketch](64),
	csmap.WithCustomHasher[string, *hyperloglog.Sketch](func(key string) uint64 {
		return xxh3.HashString(key)
	}),
)

const RANK_UPDATE_TIME = 3 * time.Minute
const MIN_CLICK_RANK = 5

// 클릭 이벤트를 버퍼에 추가하는 함수
func BufferClickEvent(markerID int) {
	// 현재 클릭 수 조회
	// 마커 ID가 존재하지 않으면 클릭 수를 1로 설정
	clickEventBuffer.SetIfAbsent(markerID, 1)

	actual, ok := clickEventBuffer.Load(markerID)
	if !ok {
		return
	}

	// 마커 ID가 존재하면 클릭 수를 1 증가
	newClicks := actual + 1
	clickEventBuffer.Store(markerID, newClicks)
}

func SaveUniqueVisitor(markerID string, uniqueUser string) {
	if markerID == "" || uniqueUser == "" {
		return
	}

	_, isInt := strconv.Atoi(markerID)
	if isInt != nil {
		return
	}

	SketchedLocations.SetIfAbsent(markerID, hyperloglog.New14())
	sketch, ok := SketchedLocations.Load(markerID)
	if !ok {
		return
	}

	sketch.Insert([]byte(uniqueUser))
}

func GetUniqueVisitorCount(markerID string) int {
	sketch, ok := SketchedLocations.Load(markerID)
	if !ok {
		return 0
	}
	return int(sketch.Estimate())
}

func GetAllUniqueVisitorCounts() map[string]int {
	result := make(map[string]int)

	// Iterate through all items in the concurrent map
	SketchedLocations.Range(func(markerID string, sketch *hyperloglog.Sketch) bool {
		count := int(sketch.Estimate())
		result[markerID] = count
		return true
	})

	return result
}

// 정해진 시간 간격마다 클릭 이벤트 배치 처리를 실행하는 함수
func ProcessClickEventsBatch() {
	// 일정 시간 간격으로 배치 처리 실행
	ticker := time.NewTicker(RANK_UPDATE_TIME)
	defer ticker.Stop() // 함수가 반환될 때 ticker를 정지

	for range ticker.C {
		IncrementMarkerClicks(clickEventBuffer)
		// 처리 후 버퍼 초기화
		clickEventBuffer.Clear()
	}
}

// 마커 방문 시 클릭 수를 파이프라인을 사용하여 증가
func IncrementMarkerClicks(markerClicks *csmap.CsMap[int, int]) {
	ctx := context.Background()

	markerClicks.Range(func(markerID int, clicks int) bool {
		scoreIncrement := float64(clicks)
		// Build and execute the ZINCRBY command for each marker
		zIncrCmd := RedisStore.B().Zincrby().Key("marker_clicks").Increment(scoreIncrement).Member(fmt.Sprintf("%d", markerID)).Build()
		if err := RedisStore.Do(ctx, zIncrCmd).Error(); err != nil {
			log.Printf("Error incrementing clicks for marker %d: %v", markerID, err)
		} else {
			// If successful, delete the marker from the map
			markerClicks.Delete(markerID)
		}
		return true // Continue iterating
	})
}

// 상위 N개 마커 랭킹 조회
func GetTopMarkers(limit int) []dto.MarkerSimpleWithAddr {
	if limit < 3 {
		limit = 5
	}
	// Sorted Set에서 점수(클릭 수)가 높은 순으로 마커 ID 조회
	ctx := context.Background()

	// Convert minClickRank to string and prepare for the ZRangeByScore command
	minScore := strconv.Itoa(MIN_CLICK_RANK + 1) // "+1" to adjust for exclusive minimum

	// Use ZREVRANGEBYSCORE to get marker IDs in descending order based on score
	markerScores, err := RedisStore.Do(ctx, RedisStore.B().Zrevrangebyscore().
		Key("marker_clicks").
		Max("+inf").
		Min(minScore).
		Withscores().
		Limit(0, int64(limit)).
		Build()).AsZScores()

	if err != nil {
		log.Printf("Error retrieving top markers: %v", err)
		return nil
	}

	if len(markerScores) == 0 {
		return []dto.MarkerSimpleWithAddr{} // Early return if no markers are found.
	}

	// Collect all marker IDs from the sorted set result for a batch database query.
	markerIDs := make([]interface{}, len(markerScores))
	for i, markerScore := range markerScores {
		markerIDs[i] = markerScore.Member // Directly use string ID to avoid unnecessary conversions.
		log.Printf("🤣 Marker id: %s and score: %f", markerScore.Member, markerScore.Score)
	}

	// Query to retrieve markers by a set of IDs in a single SQL call.
	var markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude,
        Address
    FROM 
        Markers
    WHERE MarkerID IN (?` + strings.Repeat(",?", len(markerIDs)-1) + `)
    ORDER BY FIELD(MarkerID, ?` + strings.Repeat(",?", len(markerIDs)-1) + `)`

	markerRanks := make([]dto.MarkerSimpleWithAddr, 0, len(markerIDs))
	err = database.DB.Select(&markerRanks, markerQuery, append(markerIDs, markerIDs...)...) // Duplicating markerIDs for both IN and ORDER BY.
	if err != nil {
		log.Printf("Error retrieving markers from DB: %v", err)
		return nil
	}

	return markerRanks
}

func RemoveMarkerClick(markerID int) error {
	ctx := context.Background()

	// Convert markerID to string because Redis sorted set members are strings
	member := fmt.Sprintf("%d", markerID)

	// Remove the marker from the "marker_clicks" sorted set
	err := RedisStore.Do(ctx, RedisStore.B().Zrem().Key("marker_clicks").Member(member).Build()).Error()
	if err != nil {
		log.Printf("Error removing marker click: %v", err)
		return err
	}

	return nil
}

// Admin
func ResetAndRandomizeClickRanking() {
	markers, err := GetAllMarkers()
	if err != nil {
		log.Printf("Error fetching markers: %v", err)
		return
	}

	// Ensure the slice has markers, and if not, there's nothing more to do
	if len(markers) == 0 {
		log.Println("No markers found.")
		return
	}

	// Randomly pick up to 5 marker IDs
	rand.Shuffle(len(markers), func(i, j int) {
		markers[i], markers[j] = markers[j], markers[i]
	})

	numMarkers := rand.IntN(6) + 4 // 4 ~ 9

	selectedMarkers := markers[:numMarkers]

	// atomic
	RedisStore.Dedicated(func(c rueidis.DedicatedClient) error {
		// Start a transaction
		ctx := context.Background()
		c.Do(ctx, c.B().Multi().Build())

		// Delete the existing "marker_clicks" sorted set
		c.Do(ctx, c.B().Del().Key("marker_clicks").Build())

		// Re-populate "marker_clicks" with the selected markers
		zaddCmd := c.B().Zadd().Key("marker_clicks").ScoreMember()
		for _, marker := range selectedMarkers {
			score := float64(10 + rand.IntN(6)) // Random score between 10 and 15
			zaddCmd = zaddCmd.ScoreMember(score, fmt.Sprintf("%d", marker.MarkerID))
		}
		c.Do(ctx, zaddCmd.Build())

		// Execute the transaction
		if err := c.Do(ctx, c.B().Exec().Build()).Error(); err != nil {
			log.Printf("Transaction failed: %v", err)
			return err
		}
		return nil
	})

	log.Printf("%d markers were randomly selected and added to Redis ranking.", numMarkers)
}
