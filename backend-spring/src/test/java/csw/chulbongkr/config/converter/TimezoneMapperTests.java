package csw.chulbongkr.config.converter;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class TimezoneMapperTests {

    @Test
    public void testKoreaTZ() {
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(37.5665, 126.9780));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(33.45049302403202, 126.57055468146439));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(35.1581232984585, 129.1598440928477));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(36.08502506194445, 129.55140108962055));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(33.51412972779723, 126.97244569597137));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(37.2426, 131.8597));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(37.98488937628463, 124.68608584402796));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(34.077014440034155, 125.11863713970902));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(33.40308591727227, 125.33012986877029));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(37.48996936541193, 123.90063886609349));
        assertEquals("Asia/Seoul", TimezoneMapper.latLngToTimezoneString(38.31298701550817, 127.13494497492005));
    }

    @Test
    public void testNorthKoreaTZ() {
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(37.69080614508024, 125.34030764910038));
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(38.02931441647555, 124.7248411833974));
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(37.70707171696747, 125.69373359560583));
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(39.040122308158885, 125.75997459218848));
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(39.040122308158885, 125.75997459218848));
        assertEquals("Asia/Pyongyang", TimezoneMapper.latLngToTimezoneString(37.79250413112327, 126.65242762559188));
    }

    @Test
    public void testForeignTZ() {
        assertEquals("Asia/Tokyo", TimezoneMapper.latLngToTimezoneString(32.124463344828854, 125.18301360832207));
        assertEquals("Asia/Shanghai", TimezoneMapper.latLngToTimezoneString(37.265628007634014, 122.89450076989124));

        assertEquals("America/New_York", TimezoneMapper.latLngToTimezoneString(40.7128, -74.0060));
        assertEquals("America/Los_Angeles", TimezoneMapper.latLngToTimezoneString(34.0522, -118.2437));
        assertEquals("America/Chicago", TimezoneMapper.latLngToTimezoneString(41.8781, -87.6298));
        assertEquals("America/Denver", TimezoneMapper.latLngToTimezoneString(39.7392, -104.9903));
        assertEquals("America/Phoenix", TimezoneMapper.latLngToTimezoneString(33.4484, -112.0740));
        assertEquals("America/Anchorage", TimezoneMapper.latLngToTimezoneString(61.2181, -149.9003));
        assertEquals("Pacific/Honolulu", TimezoneMapper.latLngToTimezoneString(21.3069, -157.8583));
    }
}
