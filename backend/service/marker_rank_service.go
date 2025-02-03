package service

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/axiomhq/hyperloglog"
	"github.com/jmoiron/sqlx"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/redis/rueidis"
	"github.com/zeebo/xxh3"
	"go.uber.org/zap"
)

const (
	RankUpdateTime = 3 * time.Minute
	MinClickRank   = 5

	getTopMarkersQuery = `
SELECT 
	MarkerID, 
	ST_X(Location) AS Latitude,
	ST_Y(Location) AS Longitude,
	Address
FROM 
	Markers
WHERE MarkerID IN (?)
ORDER BY FIELD(MarkerID, ?)`
)

type MarkerRankService struct {
	DB            *sqlx.DB
	Redis         *RedisService
	MarkerService *MarkerManageService
	Logger        *zap.Logger
}

func NewMarkerRankService(db *sqlx.DB, redis *RedisService, markerService *MarkerManageService, logger *zap.Logger) *MarkerRankService {
	return &MarkerRankService{
		DB:            db,
		Redis:         redis,
		MarkerService: markerService,
		Logger:        logger,
	}
}

// 클릭 이벤트를 저장할 임시 저장소
var (
	ClickEventBuffer = csmap.Create(
		csmap.WithShardCount[int, int](64),
		csmap.WithCustomHasher[int, int](func(key int) uint64 {
			// Convert int to a byte slice
			bs := make([]byte, 8)
			binary.LittleEndian.PutUint64(bs, uint64(key))
			return xxh3.Hash(bs)
		}),
	)

	SketchedLocations = csmap.Create(
		csmap.WithShardCount[string, *hyperloglog.Sketch](64),
		csmap.WithCustomHasher[string, *hyperloglog.Sketch](func(key string) uint64 {
			return xxh3.HashString(key)
		}),
	)
)

// 클릭 이벤트를 버퍼에 추가하는 함수
func (s *MarkerRankService) BufferClickEvent(markerID int) {
	// 현재 클릭 수 조회
	// 마커 ID가 존재하지 않으면 클릭 수를 1로 설정
	ClickEventBuffer.SetIfAbsent(markerID, 1)

	actual, ok := ClickEventBuffer.Load(markerID)
	if !ok {
		return
	}

	// 마커 ID가 존재하면 클릭 수를 1 증가
	newClicks := actual + 1
	ClickEventBuffer.Store(markerID, newClicks)
}

func (s *MarkerRankService) SaveUniqueVisitor(markerID string, uniqueUser string) {
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

func (s *MarkerRankService) GetUniqueVisitorCount(markerID string) int {
	sketch, ok := SketchedLocations.Load(markerID)
	if !ok {
		return 0
	}
	return int(sketch.Estimate())
}

func (s *MarkerRankService) GetAllUniqueVisitorCounts() map[string]int {
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
func (s *MarkerRankService) ProcessClickEventsBatch() {
	// 일정 시간 간격으로 배치 처리 실행
	ticker := time.NewTicker(RankUpdateTime)
	defer ticker.Stop() // 함수가 반환될 때 ticker를 정지

	for range ticker.C {
		s.IncrementMarkerClicks(ClickEventBuffer)
		// 처리 후 버퍼 초기화
		ClickEventBuffer.Clear()
	}
}

// 마커 방문 시 클릭 수를 파이프라인을 사용하여 증가
func (s *MarkerRankService) IncrementMarkerClicks(markerClicks *csmap.CsMap[int, int]) {
	ctx := context.Background()

	markerClicks.Range(func(markerID int, clicks int) bool {
		scoreIncrement := float64(clicks)
		// Build and execute the ZINCRBY command for each marker
		zIncrCmd := s.Redis.Core.Client.B().Zincrby().Key("marker_clicks").Increment(scoreIncrement).Member(fmt.Sprintf("%d", markerID)).Build()
		if err := s.Redis.Core.Client.Do(ctx, zIncrCmd).Error(); err != nil {
			s.Logger.Error("Error incrementing clicks", zap.Error(err), zap.Int("markerID", markerID))
		} else {
			// If successful, delete the marker from the map
			markerClicks.Delete(markerID)
		}
		return true // Continue iterating
	})
}

// 상위 N개 마커 랭킹 조회
func (s *MarkerRankService) GetTopMarkers(limit int) []dto.MarkerSimpleWithAddr {
	if limit < 3 {
		limit = 5
	}
	// Sorted Set에서 점수(클릭 수)가 높은 순으로 마커 ID 조회
	ctx := context.Background()

	// Convert minClickRank to string and prepare for the ZRangeByScore command
	minScore := strconv.Itoa(MinClickRank + 1) // "+1" to adjust for exclusive minimum

	// Use ZREVRANGEBYSCORE to get marker IDs in descending order based on score
	markerScores, err := s.Redis.Core.Client.Do(ctx, s.Redis.Core.Client.B().Zrevrangebyscore().
		Key("marker_clicks").
		Max("+inf").
		Min(minScore).
		Withscores().
		Limit(0, int64(limit)).
		Build()).AsZScores()

	if err != nil {
		s.Logger.Error("Error retrieving top markers", zap.Error(err))
		return nil
	}

	if len(markerScores) == 0 {
		return []dto.MarkerSimpleWithAddr{} // Early return if no markers are found.
	}

	// Collect all marker IDs from the sorted set result for a batch database query.
	markerIDs := make([]string, len(markerScores))
	for i, markerScore := range markerScores {
		markerIDs[i] = markerScore.Member // Directly use string ID to avoid unnecessary conversions.
		// log.Printf("🤣 Marker id: %s and score: %f", markerScore.Member, markerScore.Score)
	}

	// Prepare an SQL query using IN clause with sqlx.In
	query, args, err := sqlx.In(getTopMarkersQuery, markerIDs, markerIDs)
	if err != nil {
		s.Logger.Error("Error preparing query", zap.Error(err))
		return nil
	}

	// sqlx.In returns queries with the `?` bindvar, must rebind it for our specific database.
	query = s.DB.Rebind(query)

	markerRanks := make([]dto.MarkerSimpleWithAddr, 0, len(markerIDs))
	err = s.DB.Select(&markerRanks, query, args...) // args here includes markerIDs for both IN and ORDER BY clauses.
	if err != nil {
		s.Logger.Error("Error retrieving markers from DB", zap.Error(err))
		return nil
	}

	return markerRanks
}

func (s *MarkerRankService) RemoveMarkerClick(markerID int) error {
	ctx := context.Background()

	// Convert markerID to string because Redis sorted set members are strings
	member := strconv.Itoa(markerID)

	// Remove the marker from the "marker_clicks" sorted set
	err := s.Redis.Core.Client.Do(ctx, s.Redis.Core.Client.B().Zrem().Key("marker_clicks").Member(member).Build()).Error()
	if err != nil {
		s.Logger.Error("Error removing marker click", zap.Error(err))
		return err
	}

	return nil
}

// Admin
func (s *MarkerRankService) ResetAndRandomizeClickRanking() {
	ctx := context.Background()

	// Check if the "marker_clicks" sorted set already has members
	cardResp, err := s.Redis.Core.Client.Do(ctx, s.Redis.Core.Client.B().Zcard().Key("marker_clicks").Build()).AsInt64()
	if err != nil {
		return
	}
	if cardResp > 1 {
		s.Logger.Error("marker_clicks already has members. Skipping reset and randomization.")
		return
	}

	markers, err := s.MarkerService.GetAllMarkers()
	if err != nil {
		s.Logger.Error("Error fetching markers", zap.Error(err))
		return
	}

	// Filter markers with HasPhoto == true
	var filteredMarkers []dto.MarkerSimple
	for _, marker := range markers {
		if marker.HasPhoto {
			filteredMarkers = append(filteredMarkers, marker)
		}
	}

	// Ensure the slice has markers, and if not, there's nothing more to do
	if len(filteredMarkers) == 0 {
		return
	}

	// Randomly pick up to marker IDs
	rand.Shuffle(len(filteredMarkers), func(i, j int) {
		filteredMarkers[i], filteredMarkers[j] = filteredMarkers[j], filteredMarkers[i]
	})

	numMarkers := rand.IntN(100) + 30 // Random number between 30 and 130

	selectedMarkers := filteredMarkers[:numMarkers]

	// atomic
	s.Redis.Core.Client.Dedicated(func(c rueidis.DedicatedClient) error {
		// Start a transaction
		c.Do(ctx, c.B().Multi().Build())

		// Delete the existing "marker_clicks" sorted set
		c.Do(ctx, c.B().Del().Key("marker_clicks").Build())

		// Re-populate "marker_clicks" with the selected markers
		zaddCmd := c.B().Zadd().Key("marker_clicks").ScoreMember()
		for _, marker := range selectedMarkers {
			score := float64(30 + rand.IntN(10))
			zaddCmd = zaddCmd.ScoreMember(score, strconv.Itoa(marker.MarkerID))
		}
		c.Do(ctx, zaddCmd.Build())

		// Execute the transaction
		if err := c.Do(ctx, c.B().Exec().Build()).Error(); err != nil {
			s.Logger.Error("Transaction failed", zap.Error(err))
			return err
		}
		return nil
	})

	s.Logger.Info("Click ranking reset and randomization completed.", zap.Int("numMarkers", numMarkers))
}
