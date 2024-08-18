package csw.chulbongkr.service.search;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.entity.lucene.MarkerSearch;
import csw.chulbongkr.repository.marker.MarkerRepository;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

import static csw.chulbongkr.util.LuceneUtil.extractInitialConsonants;
import static csw.chulbongkr.util.LuceneUtil.standardizeProvince;

@Slf4j
@Service
public class LuceneIndexerService {
    @Autowired
    private LuceneService luceneService;

    @Autowired
    private MarkerRepository markerRepository;

    @PostConstruct
    public void indexData() throws IOException {
//        saveMarker();
        log.info("❤️❤️❤️ Indexing data...");

        ObjectMapper mapper = new ObjectMapper();
        ClassPathResource resource = new ClassPathResource("lucene/markers.json");

        try (InputStream inputStream = resource.getInputStream()) {
            List<MarkerSearch> markers = mapper.readValue(inputStream, new TypeReference<>() {});

            int batchSize = 1000;

            // Use virtual threads for concurrency
            try (ExecutorService executor = Executors.newVirtualThreadPerTaskExecutor()) {
                // ((ThreadPoolExecutor) executor).setKeepAliveTime(10, TimeUnit.SECONDS);
                for (int i = 0; i < markers.size(); i += batchSize) {
                    int end = Math.min(i + batchSize, markers.size());
                    List<MarkerSearch> batch = markers.subList(i, end);

                    executor.submit(() -> {
                        try {
                            for (MarkerSearch marker : batch) {
                                String[] addressParts = marker.getAddress().split(" ");
                                if (addressParts.length > 1) {
                                    marker.setProvince(standardizeProvince(addressParts[0]));
                                    marker.setCity(addressParts[1]);
                                }
                                marker.setFullAddress(marker.getAddress());
                                marker.setInitialConsonants(extractInitialConsonants(marker.getAddress()));
                            }
                            luceneService.indexMarkerBatch(batch);
                        } catch (IOException e) {
                            log.error("Error indexing batch: {}", e.getMessage());
                        }
                    });
                }
            } // Virtual threads will be automatically cleaned up here
        }

        log.info("❤️❤️❤️ Indexing done!");
    }

    private void saveMarker() throws IOException {
        log.info("❤️❤️❤️ Saving data...");
        // Fetch the markers from the repository
        List<MarkerDTO.MarkerSimpleWithAddr> markers = markerRepository.findAllSimplifiedMarkersWithAddress();

        // Convert the list of markers to JSON
        ObjectMapper mapper = new ObjectMapper();
        String jsonContent = mapper.writeValueAsString(markers);

        // Specify the path to the resource (under src/main/resources)
        Path path = Paths.get("src/main/resources/lucene/markers.json");

        // Ensure the directory exists
        Files.createDirectories(path.getParent());

        // Write the JSON content to the file
        Files.writeString(path, jsonContent);
    }

    public void ensureInitialized() {
        // This method is intentionally left empty. Its sole purpose is to ensure the bean is initialized.
    }
}