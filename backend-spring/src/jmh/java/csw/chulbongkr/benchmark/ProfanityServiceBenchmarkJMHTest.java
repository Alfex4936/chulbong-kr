package csw.chulbongkr.benchmark;


import csw.chulbongkr.ChulbongKrApplication;
import csw.chulbongkr.service.ProfanityService;
import org.junit.jupiter.api.Tag;
import org.openjdk.jmh.annotations.*;
import org.openjdk.jmh.infra.Blackhole;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.ConfigurableApplicationContext;

import java.util.List;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Tag("skipInCI")
@SpringBootTest(classes = ChulbongKrApplication.class)
@State(Scope.Benchmark)
public class ProfanityServiceBenchmarkJMHTest {

    @Autowired
    private ProfanityService profanityService;

    private ConfigurableApplicationContext context;

    @Setup(Level.Trial)
    public void setUp() throws Exception {
        context = SpringApplication.run(ChulbongKrApplication.class);
        profanityService = context.getBean(ProfanityService.class);
        profanityService.init();
    }

    @TearDown(Level.Trial)
    public void tearDown() throws Exception {
        context.close();
    }

    @Benchmark
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    @Warmup(iterations = 2, time = 3, timeUnit = TimeUnit.SECONDS)
    @Measurement(iterations = 10, time = 5, timeUnit = TimeUnit.SECONDS)
    public void benchmarkContainsProfanity(Blackhole blackhole) {
        int randomLength = new Random().nextInt(1_000_000) + 1_000_000; // Random length between 1,000,000 and 2,000,000
        String largeText = generateLargeRandomKoreanText(randomLength); // Generate random length text

        long startTime = System.nanoTime();
        boolean containsProfanity = profanityService.containsProfanity(largeText);
        long endTime = System.nanoTime();

        blackhole.consume(containsProfanity);
        blackhole.consume(endTime - startTime);
    }

    @Benchmark
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    @Warmup(iterations = 2, time = 3, timeUnit = TimeUnit.SECONDS)
    @Measurement(iterations = 10, time = 5, timeUnit = TimeUnit.SECONDS)
    public void benchmarkStringContains(Blackhole blackhole) {
        List<String> badWords = profanityService.getBadWords();
        int randomLength = new Random().nextInt(1_000_000) + 1_000_000; // Random length between 1,000,000 and 2,000,000
        String largeText = generateLargeRandomKoreanText(randomLength); // Generate random length text

        long startTime = System.nanoTime();
        boolean containsProfanity = false;
        for (String badWord : badWords) {
            if (largeText.contains(badWord)) {
                containsProfanity = true;
                break;
            }
        }
        long endTime = System.nanoTime();

        blackhole.consume(containsProfanity);
        blackhole.consume(endTime - startTime);
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

    public static void main(String[] args) throws Exception {
        org.openjdk.jmh.Main.main(args);
    }
}
