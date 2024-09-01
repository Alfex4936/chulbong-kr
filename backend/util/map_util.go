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
	k0               float64 = 1      // Scale factor.
	dx               float64 = 500000 // False Easting.
	dy               float64 = 200000 // False Northing.
	lat0             float64 = 38     // Latitude of origin.
	lon0             float64 = 127    // Longitude of origin.
	radiansPerDegree float64 = math.Pi / 180
	scaleFactor      float64 = 2.5
)

var wConst = math.Atan(1) / 45 // Precomputed constant value

type MapUtil struct {
	TimeZoneFinder tzf.F
}

func NewMapUtil(finder tzf.F) *MapUtil {
	return &MapUtil{
		TimeZoneFinder: finder,
	}
}

// Haversine formula
func CalculateDistanceApproximately(lat1, long1, lat2, long2 float64) float64 {
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

func (t *MapUtil) IsInSouthKoreaPrecisely(lat, lng float64) bool {
	// Get timezone name for the coordinates
	return t.TimeZoneFinder.GetTimezoneName(lng, lat) == KoreaTimeZone
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
	return WCONGNAMULCoord{X: math.Round(x * scaleFactor), Y: math.Round(y * scaleFactor)}
}

// transformWGS84ToKoreaTM optimizes the coordinate conversion calculation.
func transformWGS84ToKoreaTM(d, e, h, f, c, l, m, lat, lon float64) (float64, float64) {
	latRad := lat * radiansPerDegree
	lonRad := lon * radiansPerDegree
	lRad := l * radiansPerDegree
	mRad := m * radiansPerDegree

	sinLat := math.Sin(latRad)
	cosLat := math.Cos(latRad)

	w := 1 / e
	if e > 1 {
		w = e
	}

	z := d * (w - 1) / w
	zSquared := z * z
	dSquared := d * d
	G := 1 - zSquared/dSquared
	w = (dSquared - zSquared) / zSquared
	z = (d - z) / (d + z)

	// Precompute powers of z
	z2 := z * z
	z3 := z2 * z
	z4 := z3 * z
	z5 := z4 * z

	E := d * (1 - z + 5*(z2-z3)/4 + 81*(z4-z5)/64)
	I := 3 * d * (z - z2 + 7*(z3-z4)/8 + 55*z5/64) / 2
	J := 15 * d * (z2 - z3 + 3*(z4-z5)/4) / 16
	L := 35 * d * (z3 - z4 + 11*z5/16) / 48
	M := 315 * d * (z4 - z5) / 512

	D := lonRad - mRad
	u := E*lRad - I*math.Sin(2*lRad) + J*math.Sin(4*lRad) - L*math.Sin(6*lRad) + M*math.Sin(8*lRad)
	z = u * c
	t := sinLat / cosLat
	G = d / math.Sqrt(1-G*sinLat*sinLat)

	u = E*latRad - I*math.Sin(2*latRad) + J*math.Sin(4*latRad) - L*math.Sin(6*latRad) + M*math.Sin(8*latRad)
	o := u * c

	E = G * sinLat * cosLat * c * 0.5                                                                                                                                                       // G * sinLat * cosLat * c / 2
	I = G * sinLat * math.Pow(cosLat, 3) * c * (5 - t*t + 9*w + 4*w*w) * (1.0 / 24)                                                                                                         // G * sinLat * cosLat^3 * c * (5 - t^2 + 9w + 4w^2) / 24
	J = G * sinLat * math.Pow(cosLat, 5) * c * (61 - 58*t*t + t*t*t*t + 270*w - 330*t*t*w + 445*w*w + 324*w*w*w - 680*t*t*w*w + 88*w*w*w*w - 600*t*t*w*w*w - 192*t*t*w*w*w*w) * (1.0 / 720) // G * sinLat * cosLat^5 * c * (61 - 58t^2 + t^4 + 270w - 330t^2w + 445w^2 + 324w^3 - 680t^2w^2 + 88w^4 - 600t^2w^3 - 192t^2w^4) / 720
	H := G * sinLat * math.Pow(cosLat, 7) * c * (1385 - 3111*t*t + 543*t*t*t*t - t*t*t*t*t*t) * (1.0 / 40320)                                                                               // G * sinLat * cosLat^7 * c * (1385 - 3111t^2 + 543t^4 - t^6) / 40320
	o += D*D*E + D*D*D*I + D*D*D*D*D*J + D*D*D*D*D*D*D*H
	y := o - z + h

	o = G * cosLat * c
	z = G * math.Pow(cosLat, 3) * c * (1 - t*t + w) * (1.0 / 6)                                                                             // G * cosLat^3 * c * (1 - t^2 + w) / 6
	w = G * math.Pow(cosLat, 5) * c * (5 - 18*t*t + t*t*t*t + 14*w - 58*t*t*w + 13*w*w + 4*w*w*w - 64*t*t*w*w - 25*t*t*w*w*w) * (1.0 / 120) // G * cosLat^5 * c * (5 - 18t^2 + t^4 + 14w - 58t^2w + 13w^2 + 4w^3 - 64t^2w^2 - 25t^2w^3) / 120
	u = G * math.Pow(cosLat, 7) * c * (61 - 479*t*t + 179*t*t*t*t - t*t*t*t*t*t) * (1.0 / 5040)                                             // G * cosLat^7 * c * (61 - 479t^2 + 179t^4 - t^6) / 5040
	x := f + D*o + D*D*D*z + D*D*D*D*D*w + D*D*D*D*D*D*D*u

	return x, y
}

// ConvertWCONGToWGS84 translates WCONGNAMUL coordinates to WGS84.
func ConvertWCONGToWGS84(x, y float64) (float64, float64) {
	return transformKoreaTMToWGS84(aWGS84, flatteningFactor, dx, dy, k0, lat0, lon0, x/2.5, y/2.5)
}

// transformKoreaTMToWGS84 transforms coordinates from Korea TM to WGS84.
func transformKoreaTMToWGS84(d, e, h, f, c, l, m, x, y float64) (float64, float64) {
	u := e
	if u > 1 {
		u = 1 / u
	}
	w := wConst // Conversion factor from degrees to radians
	o := l * w
	D := m * w
	u = 1 / u
	B := d * (u - 1) / u
	z := (d*d - B*B) / (d * d)
	u = (d*d - B*B) / (B * B)
	B = (d - B) / (d + B)

	G := d * (1 - B + 5*(B*B-B*B*B)/4 + 81*(B*B*B*B-B*B*B*B*B)/64)
	E := 3 * d * (B - B*B + 7*(B*B*B-B*B*B*B)/8 + 55*B*B*B*B*B/64) / 2
	I := 15 * d * (B*B - B*B*B + 3*(B*B*B*B-B*B*B*B*B)/4) / 16
	J := 35 * d * (B*B*B - B*B*B*B + 11*B*B*B*B*B/16) / 48
	L := 315 * d * (B*B*B*B - B*B*B*B*B) / 512

	o = G*o - E*math.Sin(2*o) + I*math.Sin(4*o) - J*math.Sin(6*o) + L*math.Sin(8*o)
	o *= c
	o = y + o - h
	M := o / c
	H := d * (1 - z) / math.Pow(math.Sqrt(1-z*math.Pow(math.Sin(0), 2)), 3)
	o = M / H
	for i := 0; i < 5; i++ {
		B = G*o - E*math.Sin(2*o) + I*math.Sin(4*o) - J*math.Sin(6*o) + L*math.Sin(8*o)
		H = d * (1 - z) / math.Pow(math.Sqrt(1-z*math.Pow(math.Sin(o), 2)), 3)
		o += (M - B) / H
	}
	H = d * (1 - z) / math.Pow(math.Sqrt(1-z*math.Pow(math.Sin(o), 2)), 3)
	G = d / math.Sqrt(1-z*math.Pow(math.Sin(o), 2))
	B = math.Sin(o)
	z = math.Cos(o)
	E = B / z
	u *= z * z
	A := x - f
	B = E / (2 * H * G * math.Pow(c, 2))
	I = E * (5 + 3*E*E + u - 4*u*u - 9*E*E*u) / (24 * H * G * G * G * math.Pow(c, 4))
	J = E * (61 + 90*E*E + 46*u + 45*E*E*E*E - 252*E*E*u - 3*u*u + 100*u*u*u - 66*E*E*u*u - 90*E*E*E*E*u + 88*u*u*u*u + 225*E*E*E*E*u*u + 84*E*E*u*u*u - 192*E*E*u*u*u*u) / (720 * H * G * G * G * G * G * math.Pow(c, 6))
	H = E * (1385 + 3633*E*E + 4095*E*E*E*E + 1575*E*E*E*E*E*E) / (40320 * H * G * G * G * G * G * G * G * math.Pow(c, 8))
	o = o - math.Pow(A, 2)*B + math.Pow(A, 4)*I - math.Pow(A, 6)*J + math.Pow(A, 8)*H
	B = 1 / (G * z * c)
	H = (1 + 2*E*E + u) / (6 * G * G * G * z * z * z * math.Pow(c, 3))
	u = (5 + 6*u + 28*E*E - 3*u*u + 8*E*E*u + 24*E*E*E*E - 4*u*u*u + 4*E*E*u*u + 24*E*E*u*u*u) / (120 * G * G * G * G * G * z * z * z * z * z * math.Pow(c, 5))
	z = (61 + 662*E*E + 1320*E*E*E*E + 720*E*E*E*E*E*E) / (5040 * G * G * G * G * G * G * G * z * z * z * z * z * z * z * math.Pow(c, 7))
	A = A*B - math.Pow(A, 3)*H + math.Pow(A, 5)*u - math.Pow(A, 7)*z
	D += A

	return o / w, D / w // LATITUDE, LONGITUDE
}
