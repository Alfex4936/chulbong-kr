package csw.chulbongkr.controller;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.service.MarkerService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/markers")
public class MarkerController {
    private final MarkerService markerService;

    @GetMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<List<MarkerDTO.MarkerSimple>> getAllMarkers() {
        return ResponseEntity.ok(markerService.getAllMarkers());
    }
}
