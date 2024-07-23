package csw.chulbongkr.service.search;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import csw.chulbongkr.entity.lucene.MarkerSearch;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;
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

    @PostConstruct
    public void indexData() throws IOException {
        log.info("❤️❤️❤️ Indexing data...");
        ObjectMapper mapper = new ObjectMapper();
        ClassPathResource resource = new ClassPathResource("lucene/markers.json");

        try (InputStream inputStream = resource.getInputStream()) {
            List<MarkerSearch> markers = mapper.readValue(inputStream, new TypeReference<>() {});

            int batchSize = 1000;
            int numberOfThreads = 4;

            ThreadPoolExecutor executor = (ThreadPoolExecutor) Executors.newFixedThreadPool(numberOfThreads);

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
                        log.error("Error indexing batch: " + e.getMessage());
                    }
                });
            }

            executor.shutdown();
            executor.awaitTermination(1, TimeUnit.HOURS);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }

        log.info("❤️❤️❤️ Indexing done!");
    }

    public void ensureInitialized() {
        // This method is intentionally left empty. Its sole purpose is to ensure the bean is initialized.
    }
}