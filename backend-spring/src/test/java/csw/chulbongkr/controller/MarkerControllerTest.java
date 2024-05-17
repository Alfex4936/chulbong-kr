package csw.chulbongkr.controller;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.service.MarkerService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
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

@WebMvcTest(MarkerController.class)
class MarkerControllerTest {

    @MockBean
    private MarkerService markerService;

    private MockMvc mockMvc;

    @BeforeEach
    void setUp() {
        mockMvc = MockMvcBuilders.standaloneSetup(new MarkerController(markerService)).build();
    }

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
}