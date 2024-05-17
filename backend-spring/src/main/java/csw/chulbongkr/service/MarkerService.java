package csw.chulbongkr.service;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.repository.marker.MarkerRepository;
import csw.chulbongkr.service.local.FileDownloadService;
import csw.chulbongkr.service.local.ImageProcessorService;
import csw.chulbongkr.util.CoordinatesConverter;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Optional;

@RequiredArgsConstructor
@Service
public class MarkerService {
    private final MarkerRepository markerRepository;
    private final KakaoApiService kakaoApiService;
    private final FileDownloadService fileDownloadService;
    private final ImageProcessorService imageProcessorService;

    public List<MarkerDTO.MarkerSimple> getAllMarkers() {
        return markerRepository.findAllSimplifiedMarkers();
    }

    public Optional<String> saveOfflineMap(double lat, double lng) throws IOException {
        if (!CoordinatesConverter.IsInSouthKorea(lat, lng)) {
            return Optional.empty();
        }

        // 1. Get address in Korean
        String address = kakaoApiService.fetchAddress(lat, lng);

        // 2. Convert into WCONGNAMUL
        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);

        // 3. Get static map image
        String kakaoURL = String.format("%s&MX=%f&MY=%f", kakaoApiService.getStaticImageURL(), xy.latitude(), xy.longitude());
        String baseImageFilePath = fileDownloadService.downloadKakaoBaseImage(kakaoURL);

        // 4. Load all close markers nearby
        var nearbyMarkers = findClosestNMarkersWithinDistance(lat, lng, 700, 30, 0);
        if (nearbyMarkers.isEmpty()) {
            return Optional.empty();
        }

        List<CoordinatesConverter.XYCoordinate> convertedMarkers = nearbyMarkers.stream()
                .map(marker -> CoordinatesConverter.convertWGS84ToWCONGNAMUL(marker.latitude(), marker.longitude()))
                .toList();

        // 5. Place markers on the image
        String resultImagePath = imageProcessorService.placeMarkersOnImage(baseImageFilePath, convertedMarkers, xy.latitude(), xy.longitude());
        Files.delete(Path.of(baseImageFilePath));

        // 6. Save the image in PDF
        String pdfPath = imageProcessorService.generateMapPDF(resultImagePath, address);

        return Optional.of(pdfPath);

    }

    public List<MarkerDTO.MarkerWithDistance> findClosestNMarkersWithinDistance(double lat, double lng, double distance, int pageSize, int offset) {
        return markerRepository.findMarkersWithinDistance(lat, lng, distance, pageSize, offset);
    }
}
