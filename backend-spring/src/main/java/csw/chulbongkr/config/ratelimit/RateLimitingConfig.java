package csw.chulbongkr.config.ratelimit;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.Scheduled;

@RequiredArgsConstructor
@Configuration
@EnableScheduling
public class RateLimitingConfig {
    private final TokenBucketService tokenBucketService;

    // Refill tokens every 3 seconds
    @Scheduled(fixedRate = 3000)
    public void refillTokens() {
        tokenBucketService.refillAllBuckets();
    }
}
