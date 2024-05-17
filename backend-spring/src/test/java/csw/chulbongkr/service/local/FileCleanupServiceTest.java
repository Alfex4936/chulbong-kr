package csw.chulbongkr.service.local;


import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.MockitoAnnotations;

import java.io.IOException;
import java.nio.file.*;
import java.nio.file.attribute.BasicFileAttributes;
import java.nio.file.attribute.FileTime;
import java.time.Instant;
import java.time.temporal.ChronoUnit;

import static org.junit.jupiter.api.Assertions.*;

class FileCleanupServiceTest {

    private FileCleanupService fileCleanupService;
    private Path tempDir;

    @BeforeEach
    void setUp() throws IOException {
        MockitoAnnotations.openMocks(this);
        fileCleanupService = new FileCleanupService();
        tempDir = fileCleanupService.getTempDir();
    }


    @Test
    void testTempDirCreation() {
        assertNotNull(fileCleanupService.getTempDir());
        assertTrue(Files.exists(fileCleanupService.getTempDir()));
    }

    @Test
    void testCleanupOldFiles() throws IOException {
        Path tempDir = fileCleanupService.getTempDir();
        Path oldFile = Files.createTempFile(tempDir, "oldFile", ".tmp");
        Path newFile = Files.createTempFile(tempDir, "newFile", ".tmp");

        // Set the creation time of oldFile to 20 minutes ago
        Files.setAttribute(oldFile, "basic:creationTime", FileTime.from(Instant.now().minus(20, ChronoUnit.MINUTES)));

        // Verify both files exist
        assertTrue(Files.exists(oldFile));
        assertTrue(Files.exists(newFile));

        // Run the cleanup method
        fileCleanupService.cleanupOldFiles();

        // Verify oldFile is deleted and newFile still exists
        assertFalse(Files.exists(oldFile));
        assertTrue(Files.exists(newFile));

        // Clean up
        Files.deleteIfExists(newFile);
    }

    @Test
    void testCleanupOnShutdown() throws IOException {
        Path tempDir = fileCleanupService.getTempDir();
        Path file1 = Files.createTempFile(tempDir, "file1", ".tmp");
        Path file2 = Files.createTempFile(tempDir, "file2", ".tmp");

        // Verify both files exist
        assertTrue(Files.exists(file1));
        assertTrue(Files.exists(file2));

        // Run the cleanupOnShutdown method
        fileCleanupService.cleanupOnShutdown();

        // Verify both files are deleted
        assertFalse(Files.exists(file1));
        assertFalse(Files.exists(file2));
    }

    @AfterEach
    void tearDown() throws IOException {
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(tempDir)) {
            for (Path file : stream) {
                Files.deleteIfExists(file);
            }
        }
        Files.deleteIfExists(tempDir);
    }
}