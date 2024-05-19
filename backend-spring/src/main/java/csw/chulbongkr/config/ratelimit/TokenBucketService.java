package csw.chulbongkr.config.ratelimit;

import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

@Slf4j
@Service
public class TokenBucketService {
    private final Map<String, BucketConfig> bucketConfigs;
    private final Map<String, Map<String, AtomicInteger>> buckets = new ConcurrentHashMap<>();

    public TokenBucketService(Map<String, BucketConfig> bucketConfigs) {
        this.bucketConfigs = bucketConfigs;
    }

    public boolean allowRequest(String userId, String endpoint) {
        BucketConfig config = getConfigForEndpoint(endpoint);
        Map<String, AtomicInteger> userBuckets = buckets.computeIfAbsent(userId, k -> new ConcurrentHashMap<>());
        AtomicInteger tokens = userBuckets.computeIfAbsent(endpoint, k -> new AtomicInteger(config.getCapacity()));

        synchronized (tokens) {
//            log.info("User: {}, Endpoint: {}, Tokens before request: {}", userId, endpoint, tokens.get());
            if (tokens.get() > 0) {
                tokens.decrementAndGet();
//                log.info("Tokens after request: {}", tokens.get());
                return true;
            }
//            log.warn("Rate limit exceeded for user: {} on endpoint: {}", userId, endpoint);
            return false;
        }
    }

    public void refillAllBuckets() {
        for (Map<String, AtomicInteger> userBuckets : buckets.values()) {
            for (Map.Entry<String, AtomicInteger> entry : userBuckets.entrySet()) {
                BucketConfig config = getConfigForEndpoint(entry.getKey());
                synchronized (entry.getValue()) {
                    entry.getValue().set(config.getCapacity());
//                    log.info("Refilled tokens for userBucket: {} to capacity: {}", entry.getKey(), config.getCapacity());
                }
            }
        }
    }

    private BucketConfig getConfigForEndpoint(String endpoint) {
        return bucketConfigs.entrySet().stream()
                .filter(entry -> endpoint.matches(entry.getKey()))
                .map(Map.Entry::getValue)
                .findFirst()
                .orElseThrow(() -> new IllegalArgumentException("No bucket config found for endpoint: " + endpoint));
    }

    @Getter
    public static class BucketConfig {
        private final int capacity;

        public BucketConfig(int capacity) {
            this.capacity = capacity;
        }
    }
}