package csw.chulbongkr.service;


import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.List;
import java.util.Random;

import static org.junit.jupiter.api.Assertions.*;

@Tag("skipInCI")
@SpringBootTest
class ProfanityServiceBenchmarkTest {

    @Autowired
    private ProfanityService profanityService;

    @BeforeEach
    void setUp() throws Exception {
        profanityService.init();
    }

    /*
    20240516 (DAT)

    Average input text length: 1368812.2
    Average duration: 225100.0 ns (0.2251 ms, 2.251E-4 s)

    My DAT
    Average input text length: 1424660.8
    Average duration: 210940.0 ns (0.21094 ms, 2.1094E-4 s)

    Average input text length: 1663676.6
    Average duration: 239500.0 ns (0.2395 ms, 2.395E-4 s)
     */
    @Test
    void benchmarkContainsProfanity() {
        int iterations = 10;
        long totalDurationNs = 0;
        int totalLength = 0;

        for (int i = 0; i < iterations; i++) {
            int randomLength = new Random().nextInt(1_000_000) + 1_000_000; // Random length between 1,000,000 and 2,000,000
            String largeText = generateLargeRandomKoreanText(randomLength); // Generate random length text
            totalLength += randomLength;

            long startTime = System.nanoTime();
            boolean containsProfanity = profanityService.containsProfanity(largeText);
            long endTime = System.nanoTime();

            long durationNs = endTime - startTime; // Duration in nanoseconds
            totalDurationNs += durationNs; // Accumulate total duration

            double durationMs = durationNs / 1_000_000.0; // Convert to milliseconds
            double durationS = durationNs / 1_000_000_000.0; // Convert to seconds

            System.out.println("Run " + (i + 1) + ": " + randomLength + " texts - Profanity check duration: " + durationNs + " ns (" + durationMs + " ms, " + durationS + " s)");
            assertTrue(containsProfanity, "Profanity check failed on large Korean text");
        }

        double averageDurationNs = totalDurationNs / (double) iterations;
        double averageDurationMs = averageDurationNs / 1_000_000.0;
        double averageDurationS = averageDurationNs / 1_000_000_000.0;
        double averageLength = totalLength / (double) iterations;

        System.out.println("Average input text length: " + averageLength);
        System.out.println("Average duration: " + averageDurationNs + " ns (" + averageDurationMs + " ms, " + averageDurationS + " s)");

        // Assert that the average duration is less than 10 milliseconds (10,000,000 nanoseconds)
        assertTrue(averageDurationNs < 10_000_000, "Average profanity check should be faster than 10ms");
    }

    /*
    20240516 (String.contains)

    Average input text length: 1592083.3
    Average duration: 5.794606E7 ns (57.94606 ms, 0.05794606 s)
     */
    @Test
    void benchmarkStringContains() {
        int iterations = 10;
        long totalDurationNs = 0;
        int totalLength = 0;
        List<String> badWords = profanityService.getBadWords();

        for (int i = 0; i < iterations; i++) {
            int randomLength = new Random().nextInt(1_000_000) + 1_000_000; // Random length between 1,000,000 and 2,000,000
            String largeText = generateLargeRandomKoreanText(randomLength); // Generate random length text
            totalLength += randomLength;

            long startTime = System.nanoTime();
            boolean containsProfanity = false;
            for (String badWord : badWords) {
                if (largeText.contains(badWord)) {
                    containsProfanity = true;
                    break;
                }
            }
            long endTime = System.nanoTime();

            long durationNs = endTime - startTime; // Duration in nanoseconds
            totalDurationNs += durationNs; // Accumulate total duration

            double durationMs = durationNs / 1_000_000.0; // Convert to milliseconds
            double durationS = durationNs / 1_000_000_000.0; // Convert to seconds

            System.out.println("Run " + (i + 1) + ": " + randomLength + " texts - String.contains check duration: " + durationNs + " ns (" + durationMs + " ms, " + durationS + " s)");
            assertTrue(containsProfanity, "String.contains check failed on large Korean text");
        }

        double averageDurationNs = totalDurationNs / (double) iterations;
        double averageDurationMs = averageDurationNs / 1_000_000.0;
        double averageDurationS = averageDurationNs / 1_000_000_000.0;
        double averageLength = totalLength / (double) iterations;

        System.out.println("Average input text length: " + averageLength);
        System.out.println("Average duration: " + averageDurationNs + " ns (" + averageDurationMs + " ms, " + averageDurationS + " s)");

        // Assert that the average duration is less than 10000 milliseconds
        assertTrue(averageDurationS < 10, "Average String.contains check should be faster than 10s");
    }

    private String generateLargeRandomKoreanText(int length) {
        Random random = new Random();
        StringBuilder sb = new StringBuilder(length);

        String badWord = "시발"; // Korean bad word for testing
        int insertPosition = random.nextInt(length - badWord.length());

        for (int i = 0; i < length; i++) {
            if (i == insertPosition) {
                sb.append(badWord);
                i += badWord.length() - 1;
            } else {
                char c = (char) (random.nextInt(0xD7A3 - 0xAC00 + 1) + 0xAC00); // Random Korean character
                sb.append(c);
            }
        }

        return sb.toString();
    }

    private String generateLargeRandomText(int length) {
        Random random = new Random();
        StringBuilder sb = new StringBuilder(length);

        String badWord = "시발"; // Korean bad word for testing
        int insertPosition = random.nextInt(length - badWord.length());

        for (int i = 0; i < length; i++) {
            if (i == insertPosition) {
                sb.append(badWord);
                i += badWord.length() - 1;
            } else {
                char c = (char) (random.nextInt(26) + 'a'); // Random lower case letter
                sb.append(c);
            }
        }

        return sb.toString();
    }
}