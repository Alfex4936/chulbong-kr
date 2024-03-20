package services

import (
	"chulbong-kr/dto"
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// AddMarkerVistior
func AddMarkerVisitor(markerID int, user string) {
	// 마커 방문의 고유성을 검증하기 위해 세션 ID를 사용
	key := fmt.Sprintf("visit:marker:%d:ip:%s", markerID, user)

	// 이전에 같은 세션에서 이미 방문 기록이 있는지 확인
	exists, err := RedisStore.Conn().Exists(context.Background(), key).Result()
	if err != nil || exists > 0 {
		// 이미 방문 기록이 있거나 오류 발생 시, 추가 작업을 수행하지 않음
		return
	}

	// 마커 방문 기록
	hllKey := fmt.Sprintf("hll:marker:%s:%d", time.Now().Format("20060102"), markerID)
	_, err = RedisStore.Conn().PFAdd(context.Background(), hllKey, user).Result()
	if err != nil {
		log.Printf("Error adding visitor to HyperLogLog: %v", err)
		return
	}

	// 방문 기록의 고유성을 위한 키 설정, 짧은 만료 시간 설정
	RedisStore.Conn().Set(context.Background(), key, 1, time.Minute*10).Result()

	// 해당 키에 대해 3일(259200초) 후 만료되도록 설정
	_, err = RedisStore.Conn().Expire(context.Background(), hllKey, 259200*time.Second).Result()
	if err != nil {
		log.Printf("Error setting expiration for key %s: %v", hllKey, err)
	}
}

// 마커 인기도 예측 서비스
// 특정 시간 동안의 마커 ID 기반으로 키를 생성하여 PFCount 수행
// 결과 정렬 후 반환
// 예시: 현재 시간 2023-10-05 일 때, 마커 인기도 예측 서비스 실행 시 실행 결과
//
//	[{MarkerID: 1, Count: 5}
func EstimateMarkerPopularity(date string) []dto.MarkerPopularity {
	markerPopularityList := make([]dto.MarkerPopularity, 0)
	var cursor uint64

	// 특정 시간 동안의 마커 ID 기반으로 키를 생성하여 PFCount 수행
	matchPattern := fmt.Sprintf("hll:marker:%s:*", date)

	for {
		var keys []string
		var err error
		// NOTICE: KEYS 에서 SCAN 으로 변경
		// 성능과 확장성 향상: 큰 데이터 세트를 처리할 때 SCAN을 사용하면 Redis의 성능 저하를 방지하고, 응답성을 유지
		// 서버 부하 감소: SCAN은 서버에 부담을 주지 않으면서 데이터를 점진적으로 처리
		keys, cursor, err = RedisStore.Conn().Scan(context.Background(), cursor, matchPattern, 10).Result()
		if err != nil {
			log.Printf("Error scanning keys: %v", err)
			break
		}

		for _, key := range keys {
			// 키에서 마커 ID 추출 (파싱)
			splitKey := strings.Split(key, ":")
			markerID := splitKey[len(splitKey)-1]
			count, err := RedisStore.Conn().PFCount(context.Background(), key).Result()
			if err != nil {
				log.Printf("Error counting HyperLogLog: %v", err)
				continue
			}
			markerPopularityList = append(markerPopularityList, dto.MarkerPopularity{MarkerID: markerID, Count: count})
		}

		// 모든 키를 스캔했으면 종료
		if cursor == 0 {
			break
		}
	}

	// 결과 정렬
	sort.Slice(markerPopularityList, func(i, j int) bool {
		return markerPopularityList[i].Count > markerPopularityList[j].Count
	})
	return markerPopularityList
}
