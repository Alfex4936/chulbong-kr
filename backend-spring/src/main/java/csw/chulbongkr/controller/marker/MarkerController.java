package csw.chulbongkr.controller.marker;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.service.KakaoApiService;
import csw.chulbongkr.service.MarkerService;
import lombok.RequiredArgsConstructor;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/markers")
public class MarkerController {
    private final MarkerService markerService;
    private final KakaoApiService kakaoApiService;

    @GetMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<List<MarkerDTO.MarkerSimple>> getAllMarkers() {
        return ResponseEntity.ok(markerService.getAllMarkers());
    }

    @GetMapping(value = "/close", produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Map<String, Object>> findCloseMarkers(
            // @formatter:off
            @RequestParam double latitude,
            @RequestParam double longitude,
            @RequestParam(defaultValue = "1000") int distance,
            @RequestParam(defaultValue = "5") int pageSize,
            @RequestParam(defaultValue = "1") int page) {
        // @formatter:on
        if (distance > 10000) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", "Distance cannot be greater than 10,000m (10km)"));
        }

        if (page < 1) {
            page = 1;
        }

        if (pageSize < 1) {
            pageSize = 5;
        }

        int offset = (page - 1) * pageSize;

        List<MarkerDTO.MarkerWithDistance> markers = markerService.findClosestNMarkersWithinDistance(latitude, longitude, distance, pageSize, offset);
        int total = markers.size();

        int totalPages = total / pageSize;
        if (total % pageSize != 0) {
            totalPages++;
        }

        if (page > totalPages) {
            page = totalPages;
        }

        Map<String, Object> response = new HashMap<>();
        response.put("markers", markers);
        response.put("currentPage", page);
        response.put("totalPages", totalPages);
        response.put("totalMarkers", total);

        return ResponseEntity.ok(response);
    }

    @GetMapping(value = "/weather", produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<?> getWeather(@RequestParam double latitude, @RequestParam double longitude) {
        try {
            return ResponseEntity.ok(kakaoApiService.fetchWeather(latitude, longitude));
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Failed to fetch weather data");
        }
    }

    @GetMapping(value = "/save-offline", produces = MediaType.APPLICATION_PDF_VALUE)
    public ResponseEntity<?> downloadOfflineMap(@RequestParam double latitude, @RequestParam double longitude) {
        try {
            Optional<String> pdfPathOpt = markerService.saveOfflineMap(latitude, longitude);
            if (pdfPathOpt.isEmpty()) {
                return ResponseEntity.notFound().build();
            }

            String pdfPath = pdfPathOpt.get();
            Path filePath = Paths.get(pdfPath);
            Resource resource = new FileSystemResource(filePath);

            String contentDisposition = "attachment; filename=\"" + resource.getFilename() + "\"";
            return ResponseEntity.ok().contentType(MediaType.APPLICATION_PDF).header(HttpHeaders.CONTENT_DISPOSITION, contentDisposition).body(resource);
        } catch (IOException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Error generating file");
        }
    }
}
