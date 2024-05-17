package csw.chulbongkr.service.local;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.when;

class FileDownloadServiceTest {

    @Mock
    private RestTemplate restTemplate;

    @Mock
    private FileCleanupService fileCleanupService;

    @InjectMocks
    private FileDownloadService fileDownloadService;

    @BeforeEach
    void setUp() throws IOException {
        MockitoAnnotations.openMocks(this);
        Path tempDir = Files.createTempDirectory("test");
        when(fileCleanupService.getTempDir()).thenReturn(tempDir);
    }

    @Test
    void testDownloadKakaoBaseImage_success() throws IOException {
        byte[] imageData = "fake image data".getBytes();
        Resource resource = new ByteArrayResource(imageData) {
            @Override
            public String getFilename() {
                return "test.png";
            }

            @Override
            public InputStream getInputStream() throws IOException {
                return super.getInputStream();
            }
        };

        when(restTemplate.getForEntity(anyString(), eq(Resource.class))).thenReturn(ResponseEntity.ok(resource));

        String resultPath = fileDownloadService.downloadKakaoBaseImage("http://example.com/fake-image.png");
        assertNotNull(resultPath);
        final Path path = Path.of(resultPath);
        assertTrue(Files.exists(path));
        assertEquals("fake image data", Files.readString(path));

        // Clean up
        Files.deleteIfExists(path);
    }

    @Test
    void testDownloadKakaoBaseImage_resourceNotExists() {
        when(restTemplate.getForEntity(anyString(), eq(Resource.class))).thenReturn(ResponseEntity.of(Optional.empty()));

        IOException exception = assertThrows(IOException.class, () -> {
            fileDownloadService.downloadKakaoBaseImage("http://example.com/fake-image.png");
        });

        assertEquals("Resource does not exist: http://example.com/fake-image.png", exception.getMessage());
    }

    @Test
    void testDownloadKakaoBaseImage_nullResource() {
        when(restTemplate.getForEntity(anyString(), eq(Resource.class))).thenReturn(ResponseEntity.ok(null));

        IOException exception = assertThrows(IOException.class, () -> {
            fileDownloadService.downloadKakaoBaseImage("http://example.com/fake-image.png");
        });

        assertEquals("Resource does not exist: http://example.com/fake-image.png", exception.getMessage());
    }
}