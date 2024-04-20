package util

import (
	"math"

	"github.com/Alfex4936/tzf"
)

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

const RadiusOfEarthMeters float64 = 6370986
const KoreaTimeZone = "Asia/Seoul"

const (
	// Constants related to the WGS84 ellipsoid.
	aWGS84           float64 = 6378137 // Semi-major axis.
	flatteningFactor float64 = 0.0033528106647474805

	// Constants for Korea TM projection.
	k0          float64 = 1      // Scale factor.
	dx          float64 = 500000 // False Easting.
	dy          float64 = 200000 // False Northing.
	lat0        float64 = 38     // Latitude of origin.
	lon0        float64 = 127    // Longitude of origin.
	scaleFactor float64 = 2.5
)

var TimeZoneFinder tzf.F

// Haversine formula
func approximateDistance(lat1, long1, lat2, long2 float64) float64 {
	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)
	deltaLat := (lat2 - lat1) * (math.Pi / 180)
	deltaLong := (long2 - long1) * (math.Pi / 180)
	x := deltaLong * math.Cos((lat1Rad+lat2Rad)/2)
	y := deltaLat
	return math.Sqrt(x*x+y*y) * RadiusOfEarthMeters
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

func IsInSouthKoreaPrecisely(lat, lng float64) bool {
	// Get timezone name for the coordinates
	return TimeZoneFinder.GetTimezoneName(lng, lat) == KoreaTimeZone
}

// CONVERT ----------------------------------------------------------------
// WCONGNAMULCoord represents a coordinate in the WCONGNAMUL system.
type WCONGNAMULCoord struct {
	X float64 // X coordinate
	Y float64 // Y coordinate
}

// ConvertWGS84ToWCONGNAMUL converts coordinates from WGS84 to WCONGNAMUL.
func ConvertWGS84ToWCONGNAMUL(lat, long float64) WCONGNAMULCoord {
	x, y := transformWGS84ToKoreaTM(aWGS84, flatteningFactor, dx, dy, k0, lat0, lon0, lat, long)
	// x, y := transformWGS84ToKoreaTM(aWGS84, flatteningFactor, dx, dy, k0, lat0, lon0, lat, long)
	x = math.Round(x * scaleFactor)
	y = math.Round(y * scaleFactor)
	return WCONGNAMULCoord{X: x, Y: y}
}

// transformWGS84ToKoreaTM optimizes the coordinate conversion calculation.
func transformWGS84ToKoreaTM(d, e, h, f, c, l, m, lat, lon float64) (float64, float64) {
	A := math.Pi / 180
	latRad := lat * A
	lonRad := lon * A
	lRad := l * A
	mRad := m * A

	w := 1 / e
	if e > 1 {
		w = e
	}

	z := d * (w - 1) / w
	G := 1 - (z*z)/(d*d)
	w = (d*d - z*z) / (z * z)
	z = (d - z) / (d + z)

	E := d * (1 - z + 5*(z*z-z*z*z)/4 + 81*(z*z*z*z-z*z*z*z*z)/64)
	I := 3 * d * (z - z*z + 7*(z*z*z-z*z*z*z)/8 + 55*z*z*z*z*z/64) / 2
	J := 15 * d * (z*z - z*z*z + 3*(z*z*z*z-z*z*z*z*z)/4) / 16
	L := 35 * d * (z*z*z - z*z*z*z + 11*z*z*z*z*z/16) / 48
	M := 315 * d * (z*z*z*z - z*z*z*z*z) / 512

	D := lonRad - mRad
	u := E*lRad - I*math.Sin(2*lRad) + J*math.Sin(4*lRad) - L*math.Sin(6*lRad) + M*math.Sin(8*lRad)
	z = u * c
	sinLat := math.Sin(latRad)
	cosLat := math.Cos(latRad)
	t := sinLat / cosLat
	G = d / math.Sqrt(1-G*sinLat*sinLat)

	u = E*latRad - I*math.Sin(2*latRad) + J*math.Sin(4*latRad) - L*math.Sin(6*latRad) + M*math.Sin(8*latRad)
	o := u * c

	E = G * sinLat * cosLat * c / 2
	I = G * sinLat * math.Pow(cosLat, 3) * c * (5 - t*t + 9*w + 4*w*w) / 24
	J = G * sinLat * math.Pow(cosLat, 5) * c * (61 - 58*t*t + t*t*t*t + 270*w - 330*t*t*w + 445*w*w + 324*w*w*w - 680*t*t*w*w + 88*w*w*w*w - 600*t*t*w*w*w - 192*t*t*w*w*w*w) / 720
	H := G * sinLat * math.Pow(cosLat, 7) * c * (1385 - 3111*t*t + 543*t*t*t*t - t*t*t*t*t*t) / 40320
	o += D*D*E + D*D*D*D*I + D*D*D*D*D*D*J + D*D*D*D*D*D*D*D*H
	y := o - z + h

	o = G * cosLat * c
	z = G * math.Pow(cosLat, 3) * c * (1 - t*t + w) / 6
	w = G * math.Pow(cosLat, 5) * c * (5 - 18*t*t + t*t*t*t + 14*w - 58*t*t*w + 13*w*w + 4*w*w*w - 64*t*t*w*w - 25*t*t*w*w*w) / 120
	u = G * math.Pow(cosLat, 7) * c * (61 - 479*t*t + 179*t*t*t*t - t*t*t*t*t*t) / 5040
	x := f + D*o + D*D*D*z + D*D*D*D*D*w + D*D*D*D*D*D*D*u

	return x, y
}
