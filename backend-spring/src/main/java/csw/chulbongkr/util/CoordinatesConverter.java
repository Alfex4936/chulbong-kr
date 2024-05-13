package csw.chulbongkr.util;

public class CoordinatesConverter {
    // Constants related to the WGS84 ellipsoid.
    private static final double A_WGS84 = 6378137; // Semi-major axis.
    private static final double FLATTENING_FACTOR = 0.0033528106647474805;

    // Constants for Korea TM projection.
    private static final double K0 = 1; // Scale factor.
    private static final double DX = 500000; // False Easting.
    private static final double DY = 200000; // False Northing.
    private static final double LAT0 = 38; // Latitude of origin.
    private static final double LON0 = 127; // Longitude of origin.
    private static final double RADIANS_PER_DEGREE = Math.PI / 180;
    private static final double SCALE_FACTOR = 2.5;
    private static final double DEGREES_TO_RADIANS = Math.PI / 180;

    public static double calculateDistanceApproximately(double lat1, double long1, double lat2, double long2) {
        double lat1Rad = lat1 * RADIANS_PER_DEGREE;
        double lat2Rad = lat2 * RADIANS_PER_DEGREE;
        double deltaLat = (lat2 - lat1) * RADIANS_PER_DEGREE;
        double deltaLong = (long2 - long1) * RADIANS_PER_DEGREE;
        double x = deltaLong * Math.cos((lat1Rad + lat2Rad) / 2);
        return Math.sqrt(x * x + deltaLat * deltaLat) * A_WGS84;
    }

    public static XYCoordinate convertWGS84ToWCONGNAMUL(double lat, double lon) {
        double[] tmCoords = transformWGS84ToKoreaTM(lat, lon);
        double x = Math.round(tmCoords[0] * SCALE_FACTOR);
        double y = Math.round(tmCoords[1] * SCALE_FACTOR);
        return new XYCoordinate(x, y);
    }

    private static double[] transformWGS84ToKoreaTM(double lat, double lon) {
        double latRad = lat * RADIANS_PER_DEGREE;
        double lonRad = lon * RADIANS_PER_DEGREE;
        double lRad = CoordinatesConverter.LAT0 * RADIANS_PER_DEGREE;
        double mRad = CoordinatesConverter.LON0 * RADIANS_PER_DEGREE;

        double w = 1 / CoordinatesConverter.FLATTENING_FACTOR;

        double z = CoordinatesConverter.A_WGS84 * (w - 1) / w;
        double G = 1 - (z * z) / (CoordinatesConverter.A_WGS84 * CoordinatesConverter.A_WGS84);
        w = (CoordinatesConverter.A_WGS84 * CoordinatesConverter.A_WGS84 - z * z) / (z * z);
        z = (CoordinatesConverter.A_WGS84 - z) / (CoordinatesConverter.A_WGS84 + z);

        final double z2 = z * z;
        final double z3 = z2 * z;
        final double z4 = z2 * z2;
        final double z5 = z4 * z;

        double E = CoordinatesConverter.A_WGS84 * (1 - z + 5 * (z2 - z3) / 4 + 81 * (z4 - z5) / 64);
        double I = 3 * CoordinatesConverter.A_WGS84 * (z - z2 + 7 * (z3 - z4) / 8 + 55 * z5 / 64) / 2;
        double J = 15 * CoordinatesConverter.A_WGS84 * (z2 - z3 + 3 * (z4 - z5) / 4) / 16;
        double L = 35 * CoordinatesConverter.A_WGS84 * (z3 - z4 + 11 * z5 / 16) / 48;
        double M = 315 * CoordinatesConverter.A_WGS84 * (z4 - z5) / 512;

        double D = lonRad - mRad;
        double u = E * lRad - I * Math.sin(2 * lRad) + J * Math.sin(4 * lRad) - L * Math.sin(6 * lRad) + M * Math.sin(8 * lRad);
        z = u * CoordinatesConverter.K0;
        final double sinLat = Math.sin(latRad);
        final double cosLat = Math.cos(latRad);
        double t = sinLat / cosLat;
        G = CoordinatesConverter.A_WGS84 / Math.sqrt(1 - G * sinLat * sinLat);

        u = E * latRad - I * Math.sin(2 * latRad) + J * Math.sin(4 * latRad) - L * Math.sin(6 * latRad) + M * Math.sin(8 * latRad);
        double o = u * CoordinatesConverter.K0;

        E = G * sinLat * cosLat * CoordinatesConverter.K0 / 2;
        I = G * sinLat * Math.pow(cosLat, 3) * CoordinatesConverter.K0 * (5 - t * t + 9 * w + 4 * w * w) / 24;
        J = G * sinLat * Math.pow(cosLat, 5) * CoordinatesConverter.K0 * (61 - 58 * t * t + t * t * t * t + 270 * w - 330 * t * t * w + 445 * w * w + 324 * w * w * w - 680 * t * t * w * w + 88 * w * w * w * w - 600 * t * t * w * w * w - 192 * t * t * w * w * w * w) / 720;
        double H = G * sinLat * Math.pow(cosLat, 7) * CoordinatesConverter.K0 * (1385 - 3111 * t * t + 543 * t * t * t * t - t * t * t * t * t * t) / 40320;
        o += D * D * E + D * D * D * D * I + D * D * D * D * D * D * J + D * D * D * D * D * D * D * D * H;
        double y = o - z + CoordinatesConverter.DX;

        o = G * cosLat * CoordinatesConverter.K0;
        z = G * Math.pow(cosLat, 3) * CoordinatesConverter.K0 * (1 - t * t + w) / 6;
        w = G * Math.pow(cosLat, 5) * CoordinatesConverter.K0 * (5 - 18 * t * t + t * t * t * t + 14 * w - 58 * t * t * w + 13 * w * w + 4 * w * w * w - 64 * t * t * w * w - 25 * t * t * w * w * w) / 120;
        u = G * Math.pow(cosLat, 7) * CoordinatesConverter.K0 * (61 - 479 * t * t + 179 * t * t * t * t - t * t * t * t * t * t) / 5040;
        double x = CoordinatesConverter.DY + D * o + D * D * D * z + D * D * D * D * D * w + D * D * D * D * D * D * D * u;

        return new double[]{x, y};
    }

    public static class XYCoordinate {
        double x;
        double y;

        public XYCoordinate(double x, double y) {
            this.x = x;
            this.y = y;
        }

        public double latitude() {
            return x;
        }

        public double longitude() {
            return y;
        }

        @Override
        public String toString() {
            return "WCONGNAMULCoord{" + "x=" + x + ", y=" + y + '}';
        }
    }

    public static XYCoordinate convertWCONGNAMULToWGS84(double lat, double lon) {
        double[] tmCoords = transformKoreaTMToWGS84(lat/2.5, lon/2.5);
        return new XYCoordinate(tmCoords[0], tmCoords[1]);
    }

    private static double[] transformKoreaTMToWGS84(double x, double y) {
        double u = CoordinatesConverter.FLATTENING_FACTOR;

        double w = Math.atan(1) / 45; // Conversion factor from degrees to radians
        double o = CoordinatesConverter.LAT0 * w;
        double D = CoordinatesConverter.LON0 * w;
        u = 1 / u;

        double B = CoordinatesConverter.A_WGS84 * (u - 1) / u;
        double z = (CoordinatesConverter.A_WGS84 * CoordinatesConverter.A_WGS84 - B * B) / (CoordinatesConverter.A_WGS84 * CoordinatesConverter.A_WGS84);
        u = (CoordinatesConverter.A_WGS84 * CoordinatesConverter.A_WGS84 - B * B) / (B * B);
        B = (CoordinatesConverter.A_WGS84 - B) / (CoordinatesConverter.A_WGS84 + B);

        // Precompute powers of B
        double B2 = B * B;
        double B3 = B2 * B;
        double B4 = B3 * B;
        double B5 = B4 * B;

        double G = CoordinatesConverter.A_WGS84 * (1 - B + 5 * (B2 - B3) / 4 + 81 * (B4 - B5) / 64);
        double E = 3 * CoordinatesConverter.A_WGS84 * (B - B2 + 7 * (B3 - B4) / 8 + 55 * B5 / 64) / 2;
        double I = 15 * CoordinatesConverter.A_WGS84 * (B2 - B3 + 3 * (B4 - B5) / 4) / 16;
        double J = 35 * CoordinatesConverter.A_WGS84 * (B3 - B4 + 11 * B5 / 16) / 48;
        double L = 315 * CoordinatesConverter.A_WGS84 * (B4 - B5) / 512;

        o = G * o - E * Math.sin(2 * o) + I * Math.sin(4 * o) - J * Math.sin(6 * o) + L * Math.sin(8 * o);
        // o *= CoordinatesConverter.K0;
        o = y + o - CoordinatesConverter.DX;
        double M = o / CoordinatesConverter.K0;
        double H = CoordinatesConverter.A_WGS84 * (1 - z) / Math.pow(Math.sqrt(1 - z * Math.pow(Math.sin(0), 2)), 3);
        o = M / H;
        for (int i = 0; i < 5; i++) {
            B = G * o - E * Math.sin(2 * o) + I * Math.sin(4 * o) - J * Math.sin(6 * o) + L * Math.sin(8 * o);
            H = CoordinatesConverter.A_WGS84 * (1 - z) / Math.pow(Math.sqrt(1 - z * Math.pow(Math.sin(o), 2)), 3);
            o += (M - B) / H;
        }
        H = CoordinatesConverter.A_WGS84 * (1 - z) / Math.pow(Math.sqrt(1 - z * Math.pow(Math.sin(o), 2)), 3);
        G = CoordinatesConverter.A_WGS84 / Math.sqrt(1 - z * Math.pow(Math.sin(o), 2));
        B = Math.sin(o);
        z = Math.cos(o);
        E = B / z;
        u *= z * z;
        double A = x - CoordinatesConverter.DY;
        B = E / (2 * H * G * Math.pow(CoordinatesConverter.K0, 2));
        I = E * (5 + 3 * E * E + u - 4 * u * u - 9 * E * E * u) / (24 * H * G * G * G * Math.pow(CoordinatesConverter.K0, 4));
        J = E * (61 + 90 * E * E + 46 * u + 45 * E * E * E * E - 252 * E * E * u - 3 * u * u + 100 * u * u * u - 66 * E * E * u * u - 90 * E * E * E * E * u + 88 * u * u * u * u + 225 * E * E * E * E * u * u + 84 * E * E * u * u * u - 192 * E * E * u * u * u * u) / (720 * H * G * G * G * G * G * Math.pow(CoordinatesConverter.K0, 6));
        H = E * (1385 + 3633 * E * E + 4095 * E * E * E * E + 1575 * E * E * E * E * E * E) / (40320 * H * G * G * G * G * G * G * G * Math.pow(CoordinatesConverter.K0, 8));
        o = o - Math.pow(A, 2) * B + Math.pow(A, 4) * I - Math.pow(A, 6) * J + Math.pow(A, 8) * H;
        B = 1 / (G * z * CoordinatesConverter.K0);
        H = (1 + 2 * E * E + u) / (6 * G * G * G * z * z * z * Math.pow(CoordinatesConverter.K0, 3));
        u = (5 + 6 * u + 28 * E * E - 3 * u * u + 8 * E * E * u + 24 * E * E * E * E - 4 * u * u * u + 4 * E * E * u * u + 24 * E * E * u * u * u) / (120 * G * G * G * G * G * z * z * z * z * z * Math.pow(CoordinatesConverter.K0, 5));
        z = (61 + 662 * E * E + 1320 * E * E * E * E + 720 * E * E * E * E * E * E) / (5040 * G * G * G * G * G * G * G * z * z * z * z * z * z * z * Math.pow(CoordinatesConverter.K0, 7));
        A = A * B - Math.pow(A, 3) * H + Math.pow(A, 5) * u - Math.pow(A, 7) * z;
        D += A;

        return new double[]{o / w, D / w}; // LATITUDE, LONGITUDE
    }
}
