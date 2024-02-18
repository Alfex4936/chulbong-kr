package services

// Go convention to keep test files alongside the files they're testing, typically in the same package.

import (
	"math"
	"testing"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := distance(tt.lat1, tt.long1, tt.lat2, tt.long2)

			// Log the distance calculated for this test case
			t.Logf("Calculated distance for %q: %v meters", tt.name, result)

			if math.Abs(result-tt.expectedResult) > 1 { // Allowing a margin of error of 1 meter
				t.Errorf("distance() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
