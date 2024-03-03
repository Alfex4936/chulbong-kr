package utils

import "math"

/*
북한 제외
극동: 경상북도 울릉군의 독도(獨島)로 동경 131° 52′20“, → 131.87222222
극서: 전라남도 신안군의 소흑산도(小黑山島)로 동경 125° 04′, → 125.06666667
극북: 강원도 고성군 현내면 송현진으로 북위 38° 27′00, → 38.45000000
극남: 제주도 남제주군 마라도(馬羅島)로 북위 33° 06′00" → 33.10000000
섬 포함 우리나라의 중심점은 강원도 양구군 남면 도촌리 산48번지
북위 38도 03분 37.5초, 동경 128도 02분 2.5초 → 38.05138889, 128.03388889
섬을 제외하고 육지만을 놓고 한반도의 중심점을 계산하면 북한에 위치한 강원도 회양군 현리 인근
북위(lon): 38도 39분 00초, 동경(lat) 127도 28분 55초 → 33.10000000, 127.48194444
대한민국
도분초: 37° 34′ 8″ N, 126° 58′ 36″ E
소수점 좌표: 37.568889, 126.976667
*/
// South Korea's bounding box
const (
	SouthKoreaMinLat  = 33.0
	SouthKoreaMaxLat  = 38.615
	SouthKoreaMinLong = 124.0
	SouthKoreaMaxLong = 132.0
)

// Tsushima (Uni Island) bounding box
const (
	TsushimaMinLat  = 34.080
	TsushimaMaxLat  = 34.708
	TsushimaMinLong = 129.164396
	TsushimaMaxLong = 129.4938
)

// Nagasaki bounding box
const (
	NagasakiMinLat  = 32.75
	NagasakiMaxLat  = 34.41
	NagasakiMinLong = 128.67
	NagasakiMaxLong = 131.07
)

// Fukuoka bounding box
const (
	FukuokaMinLat  = 33.14
	FukuokaMaxLat  = 33.88
	FukuokaMinLong = 129.10
	FukuokaMaxLong = 130.64
)

// Haversine formula
func approximateDistance(lat1, long1, lat2, long2 float64) float64 {
	const R = 6370986 // Radius of the Earth in meters
	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)
	deltaLat := (lat2 - lat1) * (math.Pi / 180)
	deltaLong := (long2 - long1) * (math.Pi / 180)
	x := deltaLong * math.Cos((lat1Rad+lat2Rad)/2)
	y := deltaLat
	return math.Sqrt(x*x+y*y) * R
}

// distance calculates the distance between two geographic coordinates in meters
func distance(lat1, long1, lat2, long2 float64) float64 {
	var deltaLat = (lat2 - lat1) * (math.Pi / 180)
	var deltaLong = (long2 - long1) * (math.Pi / 180)
	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(deltaLong/2)*math.Sin(deltaLong/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return 6370986 * c // Earth radius in meters
}

// IsInSouthKorea checks if given latitude and longitude are within South Korea (roughly)
func IsInSouthKorea(lat, long float64) bool {
	// Check if within Tsushima (Uni Island) and return false if true
	if lat >= TsushimaMinLat && lat <= TsushimaMaxLat && long >= TsushimaMinLong && long <= TsushimaMaxLong {
		return false // The point is within Tsushima Island, not South Korea
	}

	// Check if within Nagasaki and return false if true
	if lat >= NagasakiMinLat && lat <= NagasakiMaxLat && long >= NagasakiMinLong && long <= NagasakiMaxLong {
		return false // The point is within Nagasaki, not South Korea
	}

	// Check if within Fukuoka and return false if true
	if lat >= FukuokaMinLat && lat <= FukuokaMaxLat && long >= FukuokaMinLong && long <= FukuokaMaxLong {
		return false // The point is within Fukuoka, not South Korea
	}

	// Check if within South Korea's bounding box
	return lat >= SouthKoreaMinLat && lat <= SouthKoreaMaxLat && long >= SouthKoreaMinLong && long <= SouthKoreaMaxLong
}
