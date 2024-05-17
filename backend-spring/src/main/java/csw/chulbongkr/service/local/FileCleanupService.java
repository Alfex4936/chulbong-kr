package csw.chulbongkr.service.local;

import jakarta.annotation.PreDestroy;
import lombok.Getter;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.nio.file.*;
import java.nio.file.attribute.BasicFileAttributes;
import java.time.Instant;
import java.time.temporal.ChronoUnit;

@Getter
@Service
public class FileCleanupService {

    private final Path tempDir;

    public FileCleanupService() throws IOException {
        this.tempDir = Files.createTempDirectory("kpullup");
    }

    @Scheduled(fixedRate = 15 * 60 * 1000) // Run every 15 minutes
    public void cleanupOldFiles() {
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(tempDir)) {
            for (Path file : stream) {
                try {
                    BasicFileAttributes attrs = Files.readAttributes(file, BasicFileAttributes.class);
                    if (attrs.creationTime().toInstant().isBefore(Instant.now().minus(15, ChronoUnit.MINUTES))) {
                        Files.delete(file);
                    }
                } catch (IOException e) {
                }
            }
        } catch (IOException e) {
        }
    }

    @PreDestroy
    public void cleanupOnShutdown() {
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(tempDir)) {
            for (Path file : stream) {
                try {
                    Files.delete(file);
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
