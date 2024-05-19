package util

// Go convention to keep test files alongside the files they're testing, typically in the same package.

import (
	"math"
	"testing"
)

// Test

// GP2TM performs the coordinate conversion calculation.
func GP2TM2(d, e, h, f, c, l, m, lat, lon float64) (float64, float64) {
	a := lat
	b := lon
	B := f
	w := e
	if w > 1 {
		w = 1 / w
	}
	A := math.Atan(1) / 45
	o := a * A
	D := b * A
	u := l * A
	A *= m
	w = 1 / w
	z := d * (w - 1) / w
	G := (math.Pow(d, 2) - math.Pow(z, 2)) / math.Pow(d, 2)
	w = (math.Pow(d, 2) - math.Pow(z, 2)) / math.Pow(z, 2)
	z = (d - z) / (d + z)
	E := d * (1 - z + 5*(math.Pow(z, 2)-math.Pow(z, 3))/4 + 81*(math.Pow(z, 4)-math.Pow(z, 5))/64)
	I := 3 * d * (z - math.Pow(z, 2) + 7*(math.Pow(z, 3)-math.Pow(z, 4))/8 + 55*math.Pow(z, 5)/64) / 2
	J := 15 * d * (math.Pow(z, 2) - math.Pow(z, 3) + 3*(math.Pow(z, 4)-math.Pow(z, 5))/4) / 16
	L := 35 * d * (math.Pow(z, 3) - math.Pow(z, 4) + 11*math.Pow(z, 5)/16) / 48
	M := 315 * d * (math.Pow(z, 4) - math.Pow(z, 5)) / 512
	D -= A
	u = E*u - I*math.Sin(2*u) + J*math.Sin(4*u) - L*math.Sin(6*u) + M*math.Sin(8*u)
	z = u * c
	H := math.Sin(o)
	u = math.Cos(o)
	A = H / u
	w *= math.Pow(u, 2)
	G = d / math.Sqrt(1-G*math.Pow(math.Sin(o), 2))
	o = E*o - I*math.Sin(2*o) + J*math.Sin(4*o) - L*math.Sin(6*o) + M*math.Sin(8*o)
	o *= c
	E = G * H * u * c / 2
	I = G * H * math.Pow(u, 3) * c * (5 - math.Pow(A, 2) + 9*w + 4*math.Pow(w, 2)) / 24
	J = G * H * math.Pow(u, 5) * c * (61 - 58*math.Pow(A, 2) + math.Pow(A, 4) + 270*w - 330*math.Pow(A, 2)*w + 445*math.Pow(w, 2) + 324*math.Pow(w, 3) - 680*math.Pow(A, 2)*math.Pow(w, 2) + 88*math.Pow(w, 4) - 600*math.Pow(A, 2)*math.Pow(w, 3) - 192*math.Pow(A, 2)*math.Pow(w, 4)) / 720
	H = G * H * math.Pow(u, 7) * c * (1385 - 3111*math.Pow(A, 2) + 543*math.Pow(A, 4) - math.Pow(A, 6)) / 40320
	o = o + math.Pow(D, 2)*E + math.Pow(D, 4)*I + math.Pow(D, 6)*J + math.Pow(D, 8)*H
	y := o - z + h
	o = G * u * c
	z = G * math.Pow(u, 3) * c * (1 - math.Pow(A, 2) + w) / 6
	w = G * math.Pow(u, 5) * c * (5 - 18*math.Pow(A, 2) + math.Pow(A, 4) + 14*w - 58*math.Pow(A, 2)*w + 13*math.Pow(w, 2) + 4*math.Pow(w, 3) - 64*math.Pow(A, 2)*math.Pow(w, 2) - 25*math.Pow(A, 2)*math.Pow(w, 3)) / 120
	u = G * math.Pow(u, 7) * c * (61 - 479*math.Pow(A, 2) + 179*math.Pow(A, 4) - math.Pow(A, 6)) / 5040
	x := B + D*o + math.Pow(D, 3)*z + math.Pow(D, 5)*w + math.Pow(D, 7)*u

	return x, y
}

func TestWCONGNAMUL(t *testing.T) {
	tests := []struct {
		name                 string
		lat, long            float64 // Starting point
		expectedX, expectedY float64 // Expected ending point
	}{
		{
			name:      "Test 1",
			lat:       37.248098895147216,
			long:      126.99116337285824,
			expectedX: 498040.0,
			expectedY: 1041367.0,
		},
		{
			name:      "Test Not korea 1",
			lat:       33.248098895147216,
			long:      126.99116337285824,
			expectedX: 497941.0,
			expectedY: -68085.0,
		},
		{
			name:      "Test 3",
			lat:       35.73294563400083,
			long:      127.37264182214031,
			expectedX: 584279.0,
			expectedY: 621193.0,
		},
		{
			name:      "Test 4",
			lat:       35.7328,
			long:      127.37264182214031,
			expectedX: 584280.0,
			expectedY: 621153.0,
		},
		{
			name:      "Test 5",
			lat:       34.248098895147216,
			long:      126.96666,
			expectedX: 492322.0,
			expectedY: 209211.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertWGS84ToWCONGNAMUL(tt.lat, tt.long)

			t.Logf("Calculated WCONGNAMUL: %f, %f", result.X, result.Y)
			if math.Abs(result.X-tt.expectedX) > 0.001 || math.Abs(result.Y-tt.expectedY) > 0.001 { // Allowing a margin of error of 1
				t.Errorf("ConvertWGS84ToWCONGNAMUL(%v, %v) = (%v, %v), want (%v, %v)",
					tt.lat, tt.long, result.X, result.Y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

func TestWCONGNAMULToWGS84(t *testing.T) {
	tests := []struct {
		name                 string
		lat, long            float64 // Starting point
		expectedX, expectedY float64 // Expected ending point
	}{
		{
			name:      "Test 1",
			lat:       498040.0,
			long:      1041367.0,
			expectedX: 37.248098895147216,
			expectedY: 126.99116337285824,
		},
		{
			name:      "Test Not korea 1",
			lat:       497941.0,
			long:      -68085.0,
			expectedX: 33.248098895147216,
			expectedY: 126.99116337285824,
		},
		{
			name:      "Test 3",
			lat:       584279.0,
			long:      621193.0,
			expectedX: 35.73294563400083,
			expectedY: 127.37264182214031,
		},
		{
			name:      "Test 4",
			lat:       584280.0,
			long:      621153.0,
			expectedX: 35.7328,
			expectedY: 127.37264182214031,
		},
		{
			name:      "Test 5",
			lat:       492322.0,
			long:      209211.0,
			expectedX: 34.248098895147216,
			expectedY: 126.96666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			X, Y := ConvertWCONGToWGS84(tt.lat, tt.long)

			t.Logf("Calculated WGS84: %f, %f", X, Y)
			if math.Abs(X-tt.expectedX) > 0.01 || math.Abs(Y-tt.expectedY) > 0.00001 { // Allowing a margin of error
				t.Errorf("ConvertWCONGToWGS84(%v, %v) = (%v, %v), want (%v, %v)",
					tt.lat, tt.long, X, Y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// Test cases for distance function
func TestDistance(t *testing.T) {
	tests := []struct {
		name           string
		lat1, long1    float64 // starting point
		lat2, long2    float64 // ending point
		expectedResult float64 // expected distance in meters
	}{
		{
			name:           "Same location",
			lat1:           40.748817,
			long1:          -73.985428,
			lat2:           40.748817,
			long2:          -73.985428,
			expectedResult: 0,
		},
		{
			name:           "Nearby point", // actual 24.879m distance
			lat1:           40.748817,
			long1:          -73.985428,
			lat2:           40.7486,
			long2:          -73.9855,
			expectedResult: 24, // Roughly 24 meters apart
		},
		{
			name:           "Very close distance", // actual 4.1749m distance
			lat1:           33.450701,
			long1:          126.570667,
			lat2:           33.450701,
			long2:          126.570712, // Approximately 5 meters away in longitude
			expectedResult: 5,          // Expecting the result to be close to 5 meters
		},
		{
			name:           "Very close distance 2", // actual 2.010888m distance
			lat1:           37.8580352854713,
			long1:          126.80789827370542,
			lat2:           37.85803307018021,
			long2:          126.80792100630472,
			expectedResult: 2, // Expecting the result to be close to 5 meters
		},
		{
			name:           "Very close distance 3", // actual 7.418885708137619m distance
			lat1:           37.293536,
			long1:          127.061558, // gwangyo
			lat2:           37.29355,
			long2:          127.061476,
			expectedResult: 7, // Expecting the result to be close to 7 meters
		},
		{
			name:           "Very close distance 4", // actual 45.015298479052085m distance
			lat1:           37.267622323623456,
			long1:          127.08362857620969,
			lat2:           37.267885604618314,
			long2:          127.0840150071249,
			expectedResult: 45, // Expecting the result to be close to 45 meters
		},
		{
			name:           "Very close distance 5",
			lat1:           33.202504,
			long1:          126.28944700000001,
			lat2:           33.201874152046415,
			long2:          126.29117918946883,
			expectedResult: 175, // Expecting the result to be close to 45 meters
		},
		{
			name:           "Very close distance 6",
			lat1:           33.23596,
			long1:          126.56218600000001,
			lat2:           33.23596,
			long2:          126.56208600000001,
			expectedResult: 0, // Expecting the result to be close to 45 meters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDistanceApproximately(tt.lat1, tt.long1, tt.lat2, tt.long2)

			// Log the distance calculated for this test case
			t.Logf("Calculated distance for %q: %v meters", tt.name, result)

			if math.Abs(result-tt.expectedResult) > 1 { // Allowing a margin of error of 1 meter
				t.Errorf("distance() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestIsInSouthKorea(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		want      bool
	}{
		{"서울", 37.5665, 126.9780, true},
		{"제주 카카오오름", 33.45049302403202, 126.57055468146439, true},
		{"해운대", 35.1581232984585, 129.1598440928477, true},  // 해운대 해수욕장
		{"포항", 36.08502506194445, 129.55140108962055, true}, // 포항
		{"세종", 36.481550006080006, 127.28920084353089, true},
		{"제주도", 33.4890, 126.4983, true},
		{"우도", 33.51412972779723, 126.97244569597137, true},
		{"마라도", 33.11294701534852, 126.2662987980748, true},
		{"독도", 37.2426, 131.8597, true},
		{"울릉도", 37.4845, 130.9057, true},
		{"차귀도", 33.311273820042125, 126.14345298508049, true}, // 차귀도- 제주특별자치도 제주시 한경면 고산리
		{"대강리", 38.61453830741445, 128.35799152766955, true},  // northernmost point, 강원특별자치도 고성군 현내면
		{"백령도", 37.96909906079667, 124.609983839757, true},    // westernmost point, 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"백령도2", 37.98488937628463, 124.68608584402796, true}, // 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"철원", 38.31374456713513, 127.13423745903036, true},   // 강원특별자치도 철원군 철원읍 가단리 52
		{"거제도", 34.54419719852532, 128.43864110479205, true},  // 거제도
		{"광도", 34.269977354595504, 127.53055654653483, true},
		{"가거도", 34.077014440034155, 125.11863713970902, true},
		// false
		{"이어도", 32.124463344828854, 125.18301360832207, false}, // southernmost point, 이어도. cannot build stuff.
		{"Los Angeles", 34.0522, -118.2437, false},
		{"Tokyo", 35.6895, 139.6917, false},
		{"Beijing", 39.9042, 116.4074, false},
		{"Uni Island", 34.707351308730146, 129.43478825264333, false},
		{"Uni Island2", 34.43217756058352, 129.33997781093186, false},
		{"Uni Island3", 34.636217082470296, 129.4828167691493, false},
		{"Uni Island4", 34.29666974505072, 129.3871993238883, false},
		{"Uni Island5", 34.0854739629158, 129.2154168085643, false},
		{"Fukuoka 1 (Japan)", 33.784029222960406, 130.53443527389945, false},
		{"Fukuoka 2 (Japan)", 34.296085822281455, 130.93051474444093, false},
		{"Fukuoka 3 (Japan)", 32.69461329871054, 128.79495039442563, false},
		{"Fukuoka 4 (Japan)", 32.95445481630956, 129.09330313600782, false},
		{"Fukuoka 5 (Japan)", 33.53700218298737, 130.3983824405139, false},
		{"Shimanae (Japan)", 35.03719336610837, 132.4915325911786, false},
		{"Okinoshimajo (Japan)", 36.27042331297408, 133.24889805463428, false},
		{"Shimayama Island (Japan)", 32.683327616680096, 128.64905526405005, false},
		{"Kyoto (Japan)", 35.277030942449066, 135.4727941919809, false},
		{"Yantai (China)", 37.45460313491269, 122.43159543394779, false},
		{"평양 (N.Korea)", 39.040122308158885, 125.75997459218848, false},
		// {"Iki Island", 33.833510640897295, 129.6794423356137, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsInSouthKorea(tc.latitude, tc.longitude)
			if got != tc.want {
				t.Errorf("FAIL: %s - IsInSouthKorea(%f, %f) = %v; want %v", tc.name, tc.latitude, tc.longitude, got, tc.want)
			} else {
				// Provide clearer messages indicating the correctness of the test result
				if got {
					t.Logf("PASS: %s is correctly identified as inside South Korea.", tc.name)
				} else {
					t.Logf("PASS: %s is correctly identified as outside South Korea.", tc.name)
				}
			}
		})
	}
}
