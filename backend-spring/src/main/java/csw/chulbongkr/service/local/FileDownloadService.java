package csw.chulbongkr.service.local;

import lombok.RequiredArgsConstructor;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.util.UUID;

@RequiredArgsConstructor
@Service
public class FileDownloadService {
    private final RestTemplate restTemplate;
    private final FileCleanupService fileCleanupService;

    public String downloadKakaoBaseImage(String url) throws IOException {
        ResponseEntity<Resource> response = restTemplate.getForEntity(url, Resource.class);
        Resource resource = response.getBody();

        if (resource == null || !resource.exists()) {
            throw new IOException("Resource does not exist: " + url);
        }

        String fileName = "base_map-" + UUID.randomUUID() + ".png";
        Path destPath = fileCleanupService.getTempDir().resolve(fileName);

        try (InputStream in = resource.getInputStream();
             OutputStream out = Files.newOutputStream(destPath, StandardOpenOption.CREATE, StandardOpenOption.TRUNCATE_EXISTING, StandardOpenOption.WRITE)) {
            byte[] buffer = new byte[8192];
            int bytesRead;
            while ((bytesRead = in.read(buffer)) != -1) {
                out.write(buffer, 0, bytesRead);
            }
        }

        return destPath.toString();
    }
}
