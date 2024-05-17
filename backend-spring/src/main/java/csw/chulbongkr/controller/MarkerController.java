package csw.chulbongkr.controller;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.service.MarkerService;
import lombok.RequiredArgsConstructor;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.core.io.UrlResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.net.MalformedURLException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Optional;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/markers")
public class MarkerController {
    private final MarkerService markerService;

    @GetMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<List<MarkerDTO.MarkerSimple>> getAllMarkers() {
        return ResponseEntity.ok(markerService.getAllMarkers());
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
            return ResponseEntity.ok()
                    .contentType(MediaType.APPLICATION_PDF)
                    .header(HttpHeaders.CONTENT_DISPOSITION, contentDisposition)
                    .body(resource);
        } catch (IOException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error generating file");
        }
    }
}
