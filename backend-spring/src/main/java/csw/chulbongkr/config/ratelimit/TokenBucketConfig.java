package csw.chulbongkr.config.ratelimit;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.HashMap;
import java.util.Map;

@Configuration
public class TokenBucketConfig {

    @Bean
    public TokenBucketService tokenBucketService() {
        Map<String, TokenBucketService.BucketConfig> bucketConfigs = new HashMap<>();
        bucketConfigs.put("/api/v1/auth(?:/.*)?", new TokenBucketService.BucketConfig(10));  // Less capacity for auth endpoints
        bucketConfigs.put("/api/v1/markers(?:/.*)?", new TokenBucketService.BucketConfig(50));
        bucketConfigs.put("/api/v1/comments(?:/.*)?", new TokenBucketService.BucketConfig(50));
        return new TokenBucketService(bucketConfigs);
    }
}