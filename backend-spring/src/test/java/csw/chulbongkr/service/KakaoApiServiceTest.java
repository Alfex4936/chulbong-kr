package csw.chulbongkr.service;

import csw.chulbongkr.config.custom.KakaoConfig;
import csw.chulbongkr.dto.KakaoDTO;
import csw.chulbongkr.util.CoordinatesConverter;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.web.client.RestTemplate;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.when;

class KakaoApiServiceTest {

    @Mock
    private KakaoConfig kakaoConfig;

    @Mock
    private RestTemplate restTemplate;

    @InjectMocks
    private KakaoApiService kakaoApiService;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);

        when(kakaoConfig.getAddressInfo()).thenReturn("http://fake-address-url");
        when(kakaoConfig.getWeatherUrl()).thenReturn("http://fake-weather-url");
        when(kakaoConfig.getWeatherIconUrl()).thenReturn("http://fake-icon-url/%s.png");
    }

    @Test
    void testFetchAddress_Success() {
        // Given
        double lat = 37.5665;
        double lng = 126.9780;
        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);
        String requestURL = "http://fake-address-url&x=" + xy.latitude() + "&y=" + xy.longitude();

        KakaoDTO.KakaoMarkerData markerData = new KakaoDTO.KakaoMarkerData(
                null,
                0,
                0,
                null,
                null,
                null,
                null,
                new KakaoDTO.AreaDocument("1", "New Address", null, "Building", null, null, null, null, null)
        );
        when(restTemplate.getForObject(requestURL, KakaoDTO.KakaoMarkerData.class)).thenReturn(markerData);

        // When
        String address = kakaoApiService.fetchAddress(lat, lng);

        // Then
        assertEquals("New Address, Building", address);
    }

    @Test
    void testFetchAddress_NoAddress() {
        // Given
        double lat = 37.5665;
        double lng = 126.9780;
        String requestURL = "http://fake-address-url&x=37.5665&y=126.9780";

        when(restTemplate.getForObject(requestURL, KakaoDTO.KakaoMarkerData.class)).thenReturn(null);

        // When
        String address = kakaoApiService.fetchAddress(lat, lng);

        // Then
        assertEquals("대한민국 철봉 지도", address);
    }

    @Test
    void testFetchWeather_Success() {
        // Given
        double lat = 37.5665;
        double lng = 126.9780;
        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);
        String requestURL = "http://fake-weather-url&x=" + xy.latitude() + "&y=" + xy.longitude();

        HttpHeaders headers = new HttpHeaders();
        headers.set("Referer", requestURL);
        HttpEntity<String> entity = new HttpEntity<>(headers);

        KakaoDTO.Weather.WeatherInfo currentWeather = new KakaoDTO.Weather.WeatherInfo("01", "OK", "01", "25.0", "Sunny", "0", "0", "0");
        KakaoDTO.Weather weather = new KakaoDTO.Weather(new KakaoDTO.Weather.Codes("OK", null, null), new KakaoDTO.Weather.WeatherInfos(currentWeather, null));

        ResponseEntity<KakaoDTO.Weather> responseEntity = ResponseEntity.ok(weather);
        when(restTemplate.exchange(requestURL, HttpMethod.GET, entity, KakaoDTO.Weather.class)).thenReturn(responseEntity);

        // When
        KakaoDTO.Weather.WeatherRequest weatherRequest = kakaoApiService.fetchWeather(lat, lng);

        // Then
        assertEquals("25.0", weatherRequest.temperature());
        assertEquals("Sunny", weatherRequest.desc());
        assertEquals("http://fake-icon-url/01.png", weatherRequest.iconImage());
        assertEquals("0", weatherRequest.humidity());
        assertEquals("0", weatherRequest.rainfall());
        assertEquals("0", weatherRequest.snowfall());
    }

    @Test
    void testFetchWeather_Failure() {
        // Given
        double lat = 37.5665;
        double lng = 126.9780;
        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);
        String requestURL = "http://fake-weather-url&x=" + xy.latitude() + "&y=" + xy.longitude();

        HttpHeaders headers = new HttpHeaders();
        headers.set("Referer", requestURL);
        HttpEntity<String> entity = new HttpEntity<>(headers);

        // Ensure the response entity is null to trigger the runtime exception
        when(restTemplate.exchange(requestURL, HttpMethod.GET, entity, KakaoDTO.Weather.class)).thenReturn(null);

        // When & Then
        RuntimeException exception = assertThrows(RuntimeException.class, () -> {
            kakaoApiService.fetchWeather(lat, lng);
        });
        assertEquals("Failed to fetch weather data", exception.getMessage());
    }

    @Test
    void testGetStaticImageURL() {
        // Given
        String expectedUrl = "http://static-map-url";
        when(kakaoConfig.getStaticMap()).thenReturn(expectedUrl);

        // When
        String actualUrl = kakaoApiService.getStaticImageURL();

        // Then
        assertEquals(expectedUrl, actualUrl);
    }

    @Test
    void testFetchAddress_OldAddress() {
        // Given
        double lat = 37.5665;
        double lng = 126.9780;

        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);
        String requestURL = "http://fake-address-url&x=" + xy.latitude() + "&y=" + xy.longitude();

        KakaoDTO.KakaoMarkerData markerData = new KakaoDTO.KakaoMarkerData(
                new KakaoDTO.AreaDocument("1", "Old Address", null, "Old Building", null, null, null, null, null),
                0,
                0,
                null,
                null,
                null,
                null,
                null
        );
        when(restTemplate.getForObject(requestURL, KakaoDTO.KakaoMarkerData.class)).thenReturn(markerData);

        // When
        String address = kakaoApiService.fetchAddress(lat, lng);

        // Then
        assertEquals("Old Address, Old Building", address);
    }
}