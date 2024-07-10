package csw.chulbongkr.controller;

import csw.chulbongkr.controller.marker.MarkerController;
import csw.chulbongkr.dto.KakaoDTO;
import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.service.KakaoApiService;
import csw.chulbongkr.service.MarkerService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.hamcrest.Matchers.containsString;
import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@AutoConfigureMockMvc
@SpringBootTest
class MarkerControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private MarkerService markerService;

    @MockBean
    private KakaoApiService kakaoApiService;

    @Test
    void testGetAllMarkers() throws Exception {
        List<MarkerDTO.MarkerSimple> markerList = Collections.singletonList(
                new MarkerDTO.MarkerSimple(1, 37.529903839012064, 127.04447892740619)
        );
        when(markerService.getAllMarkers()).thenReturn(markerList);

        mockMvc.perform(get("/api/v1/markers")
                        .accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
                .andExpect(jsonPath("$[0].markerId").value(1))
                .andExpect(jsonPath("$[0].latitude").value(37.529903839012064))
                .andExpect(jsonPath("$[0].longitude").value(127.04447892740619));

        verify(markerService, times(1)).getAllMarkers();
    }

    @Test
    void testDownloadOfflineMap() throws Exception {
        String pdfPath = "kpullup-" + UUID.randomUUID() + ".pdf";
        Path filePath = Paths.get(pdfPath);
        File pdfFile = filePath.toFile();
        pdfFile.createNewFile(); // Create an empty file for testing

        when(markerService.saveOfflineMap(37.529903839012064, 127.04447892740619))
                .thenReturn(Optional.of(pdfPath));

        mockMvc.perform(get("/api/v1/markers/save-offline")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619")
                        .accept(MediaType.APPLICATION_PDF))
                .andExpect(status().isOk())
                .andExpect(header().string(HttpHeaders.CONTENT_DISPOSITION, containsString("attachment; filename=\"kpullup-")))
                .andExpect(content().contentType(MediaType.APPLICATION_PDF_VALUE));

        verify(markerService, times(1)).saveOfflineMap(37.529903839012064, 127.04447892740619);

        // Clean up
        pdfFile.delete();
    }

    @Test
    void testDownloadOfflineMap_NotFound() throws Exception {
        when(markerService.saveOfflineMap(37.529903839012064, 127.04447892740619))
                .thenReturn(Optional.empty());

        mockMvc.perform(get("/api/v1/markers/save-offline")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619")
                        .accept(MediaType.APPLICATION_PDF))
                .andExpect(status().isNotFound());

        verify(markerService, times(1)).saveOfflineMap(37.529903839012064, 127.04447892740619);
    }

    @Test
    void testDownloadOfflineMap_Exception() throws Exception {
        when(markerService.saveOfflineMap(37.529903839012064, 127.04447892740619))
                .thenThrow(new IOException("Error generating file"));

        mockMvc.perform(get("/api/v1/markers/save-offline")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619")
                        .accept(MediaType.APPLICATION_PDF))
                .andExpect(status().isInternalServerError())
                .andExpect(content().string(containsString("Error generating file")));

        verify(markerService, times(1)).saveOfflineMap(37.529903839012064, 127.04447892740619);
    }

    @Test
    void testFindCloseMarkers_ValidRequest_DefaultParams() throws Exception {
        List<MarkerDTO.MarkerWithDistance> markers = Collections.singletonList(
                new MarkerDTO.MarkerWithDistance(1, 37.529903839012064, 127.04447892740619, "Description", 100.0, null)
        );
        when(markerService.findClosestNMarkersWithinDistance(37.529903839012064, 127.04447892740619, 1000, 5, 0))
                .thenReturn(markers);

        mockMvc.perform(get("/api/v1/markers/close")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
                .andExpect(jsonPath("$.markers[0].markerId").value(1))
                .andExpect(jsonPath("$.currentPage").value(1))
                .andExpect(jsonPath("$.totalPages").value(1))
                .andExpect(jsonPath("$.totalMarkers").value(1));

        verify(markerService, times(1)).findClosestNMarkersWithinDistance(37.529903839012064, 127.04447892740619, 1000, 5, 0);
    }

    @Test
    void testFindCloseMarkers_ValidRequest_SpecifiedParams() throws Exception {
        List<MarkerDTO.MarkerWithDistance> markers = Collections.singletonList(
                new MarkerDTO.MarkerWithDistance(1, 37.529903839012064, 127.04447892740619, "Description", 100.0, null)
        );
        when(markerService.findClosestNMarkersWithinDistance(37.529903839012064, 127.04447892740619, 2000, 5, 0))
                .thenReturn(markers);

        mockMvc.perform(get("/api/v1/markers/close")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619")
                        .param("distance", "2000")
                        .param("pageSize", "0")
                        .param("page", "0"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
                .andExpect(jsonPath("$.markers[0].markerId").value(1))
                .andExpect(jsonPath("$.currentPage").value(1))
                .andExpect(jsonPath("$.totalPages").value(1))
                .andExpect(jsonPath("$.totalMarkers").value(1));

        verify(markerService, times(1)).findClosestNMarkersWithinDistance(37.529903839012064, 127.04447892740619, 2000, 5, 0);
    }

    @Test
    void testFindCloseMarkers_ExceedMaxDistance() throws Exception {
        mockMvc.perform(get("/api/v1/markers/close")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619")
                        .param("distance", "15000"))
                .andExpect(status().isForbidden())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
                .andExpect(jsonPath("$.error").value("Distance cannot be greater than 10,000m (10km)"));

        verify(markerService, never()).findClosestNMarkersWithinDistance(anyDouble(), anyDouble(), anyInt(), anyInt(), anyInt());
    }

    @Test
    void testGetWeather_Success() throws Exception {
        KakaoDTO.Weather.WeatherRequest weatherRequest = new KakaoDTO.Weather.WeatherRequest(
                "20", "Clear", "icon.png", "50%", "0", "0"
        );
        when(kakaoApiService.fetchWeather(37.529903839012064, 127.04447892740619))
                .thenReturn(weatherRequest);

        mockMvc.perform(get("/api/v1/markers/weather")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
                .andExpect(jsonPath("$.temperature").value("20"))
                .andExpect(jsonPath("$.desc").value("Clear"))
                .andExpect(jsonPath("$.iconImage").value("icon.png"))
                .andExpect(jsonPath("$.humidity").value("50%"))
                .andExpect(jsonPath("$.rainfall").value("0"))
                .andExpect(jsonPath("$.snowfall").value("0"));

        verify(kakaoApiService, times(1)).fetchWeather(37.529903839012064, 127.04447892740619);
    }

    @Test
    void testGetWeather_Exception() throws Exception {
        when(kakaoApiService.fetchWeather(37.529903839012064, 127.04447892740619))
                .thenThrow(new RuntimeException("Failed to fetch weather data"));

        mockMvc.perform(get("/api/v1/markers/weather")
                        .param("latitude", "37.529903839012064")
                        .param("longitude", "127.04447892740619"))
                .andExpect(status().isInternalServerError())
                .andExpect(content().string(containsString("Failed to fetch weather data")));

        verify(kakaoApiService, times(1)).fetchWeather(37.529903839012064, 127.04447892740619);
    }
}