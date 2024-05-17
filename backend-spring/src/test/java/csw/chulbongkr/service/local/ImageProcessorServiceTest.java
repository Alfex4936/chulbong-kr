package csw.chulbongkr.service.local;


import csw.chulbongkr.util.CoordinatesConverter;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.core.io.ClassPathResource;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

class ImageProcessorServiceTest {

    @Mock
    private FileCleanupService fileCleanupService;

    @InjectMocks
    private TestImageProcessorService  imageProcessorService;

    @BeforeEach
    void setUp() throws IOException {
        MockitoAnnotations.openMocks(this);
        when(fileCleanupService.getTempDir()).thenReturn(Files.createTempDirectory("test"));

        imageProcessorService.init();
    }

    @Test
    void testPlaceMarkersOnImage() throws IOException {
        List<CoordinatesConverter.XYCoordinate> markers = List.of(
                new CoordinatesConverter.XYCoordinate(37.1, 127.1),
                new CoordinatesConverter.XYCoordinate(37.2, 127.2)
        );

        String resultPath = imageProcessorService.placeMarkersOnImage("src/test/resources/test_base_image.png", markers, 127.0, 37.0);
        assertNotNull(resultPath);
        assertTrue(new File(resultPath).exists());

        // Clean up
        Files.deleteIfExists(Path.of(resultPath));
    }

    @Test
    void testGenerateMapPDF() throws IOException {
        String pdfPath = imageProcessorService.generateMapPDF("src/test/resources/test_base_image.png", "Test Title");
        assertNotNull(pdfPath);
        assertTrue(new File(pdfPath).exists());

        // Clean up
        Files.deleteIfExists(Path.of(pdfPath));
    }

    @Test
    void testInit() throws IOException {
        imageProcessorService.init();
        assertNotNull(imageProcessorService.getMarkerIcon());
        assertNotNull(imageProcessorService.getNanumFont());
    }

    @Test
    void testLoadWebP() throws IOException {
        BufferedImage image = imageProcessorService.loadWebP("map_marker.webp");
        assertNotNull(image);
    }

    @Test
    void testLoadNanumFont() throws IOException {
        File fontFile = imageProcessorService.loadNanumFont("fonts/nanum.ttf");
        assertNotNull(fontFile);
        assertTrue(fontFile.exists());
    }

    private BufferedImage loadTestImage() throws IOException {
        return ImageIO.read(new ClassPathResource("test_base_image.png").getFile());
    }

    private File loadTestFont() throws IOException {
        return new ClassPathResource("fonts/nanum.ttf").getFile();
    }

    // Create a subclass for testing
    static class TestImageProcessorService extends ImageProcessorService {

        public TestImageProcessorService(FileCleanupService fileCleanupService) {
            super(fileCleanupService);
        }

        @Override
        protected BufferedImage loadWebP(String resourcePath) throws IOException {
            return ImageIO.read(new ClassPathResource(resourcePath).getFile());
        }

        @Override
        protected File loadNanumFont(String resourcePath) throws IOException {
            return new ClassPathResource(resourcePath).getFile();
        }
    }
}