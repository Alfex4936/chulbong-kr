package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/alphadose/haxmap"

	"chulbong-kr/database"
	"chulbong-kr/dto"
)

// 클릭 이벤트를 저장할 임시 저장소
var clickEventBuffer = haxmap.New[int, int]()

const RANK_UPDATE_TIME = 3 * time.Minute

// 클릭 이벤트를 버퍼에 추가하는 함수
func BufferClickEvent(markerID int) {
	// 현재 클릭 수 조회
	val, ok := clickEventBuffer.Get(markerID)
	if !ok {
		// 마커 ID가 존재하지 않으면 클릭 수를 1로 설정
		clickEventBuffer.Set(markerID, 1)
	} else {
		// 마커 ID가 존재하면 클릭 수를 1 증가
		newClicks := val + 1
		clickEventBuffer.Set(markerID, newClicks)
	}
}

// 정해진 시간 간격마다 클릭 이벤트 배치 처리를 실행하는 함수
func ProcessClickEventsBatch() {
	// 일정 시간 간격으로 배치 처리 실행
	ticker := time.NewTicker(RANK_UPDATE_TIME)
	defer ticker.Stop() // 함수가 반환될 때 ticker를 정지

	for range ticker.C {
		IncrementMarkerClicks(clickEventBuffer)
		// 처리 후 버퍼 초기화
		clickEventBuffer = haxmap.New[int, int]()
	}
}

// 마커 방문 시 클릭 수를 파이프라인을 사용하여 증가
func IncrementMarkerClicks(markerClicks *haxmap.Map[int, int]) {
	pipe := RedisStore.Conn().TxPipeline()

	clickEventBuffer.ForEach(func(markerID int, clicks int) bool {
		// map에서 가져온 클릭 수만큼 점수 증가
		scoreIncrement := float64(clicks)
		pipe.ZIncrBy(context.Background(), "marker_clicks", scoreIncrement, fmt.Sprintf("%d", markerID))
		return true // return `true` to continue iteration and `false` to break iteration
	})

	// Execute all commands in the pipeline
	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error incrementing marker clicks: %v", err)
	}
}

// TODO: 내 주변 상위 마커 N 개 랭킹 조회도 만들기
// 상위 N개 마커 랭킹 조회
func GetTopMarkers(limit int) []dto.MarkerSimpleWithAddr {
	if limit < 3 {
		limit = 5
	}
	// Sorted Set에서 점수(클릭 수)가 높은 순으로 마커 ID 조회
	markerScores, err := RedisStore.Conn().ZRevRangeWithScores(context.Background(), "marker_clicks", 0, int64(limit-1)).Result()
	if err != nil {
		log.Printf("Error retrieving top markers: %v", err)
		return nil
	}

	var markerIDs []int // 조회한 마커 ID를 저장할 슬라이스
	for _, markerScore := range markerScores {
		// Redis에서 조회한 마커 ID를 정수로 변환
		markerIDStr, _ := markerScore.Member.(string)
		markerID, err := strconv.Atoi(markerIDStr)
		if err != nil {
			log.Printf("Error converting marker ID: %v", err)
			continue
		}

		markerIDs = append(markerIDs, markerID)
	}

	// 데이터베이스에서 마커의 상세 정보 조회
	markerRanks := make([]dto.MarkerSimpleWithAddr, 0)
	const markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude,
        Address
    FROM 
        Markers
	WHERE MarkerID = ?`

	for _, markerID := range markerIDs {
		var marker dto.MarkerSimpleWithAddr
		err := database.DB.Get(&marker, markerQuery, markerID)
		if err != nil {
			log.Printf("Error retrieving marker: %v", err)
			continue
		}

		markerRanks = append(markerRanks, marker)
	}

	return markerRanks
}

func RemoveMarkerClick(markerID int) error {
	ctx := context.Background()

	// Convert markerID to string because Redis sorted set members are strings
	member := fmt.Sprintf("%d", markerID)

	// Remove the marker from the "marker_clicks" sorted set
	_, err := RedisStore.Conn().ZRem(ctx, "marker_clicks", member).Result()
	if err != nil {
		log.Printf("Error removing marker click: %v", err)
		return err
	}

	return nil
}
