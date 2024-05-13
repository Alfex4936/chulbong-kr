package csw.chulbongkr.util;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class CoordinatesConverterTest {

    @Test
    public void testCalculateSameDistance() {
        // Given
        double lat1 = 37.5665;
        double long1 = 126.9780;
        double lat2 = 37.5665;
        double long2 = 126.9780;

        // When
        double result = CoordinatesConverter.calculateDistanceApproximately(lat1, long1, lat2, long2);

        // Then
        assertEquals(0, result);
    }

    @Test
    public void testCalculateCloseDistances() {
        // Given
        double lat1 = 37.293536;
        double long1 = 127.061558;
        double lat2 = 37.29355;
        double long2 = 127.061476;

        // When
        double result = CoordinatesConverter.calculateDistanceApproximately(lat1, long1, lat2, long2); // 7.427212904540012

        // Then
        assertEquals(7, result, 1.0);
    }

    @Test
    public void testTransformWGS84ToKoreaTM() {
        // Given
        double lat = 37.5478543870196;
        double lon = 129.105530908533;

        // When
        var result = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lon);

        // Then
        assertEquals(965186.0, result.x);
        assertEquals(1129749.0, result.y);
    }

    @Test
    public void testTransformKoreaTmToWGS84() {
        // Given
        double lat = 965186.0;
        double lon = 1129749.0;

        // When
        var result = CoordinatesConverter.convertWCONGNAMULToWGS84(lat, lon);

        // Then
        assertEquals(37.5478543870196, result.x, 0.001);
        assertEquals(129.105530908533, result.y, 0.001);
    }
}
