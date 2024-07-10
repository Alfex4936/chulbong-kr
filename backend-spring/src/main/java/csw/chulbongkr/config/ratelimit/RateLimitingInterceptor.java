package csw.chulbongkr.config.ratelimit;

import csw.chulbongkr.util.IdentityUtil;
import lombok.RequiredArgsConstructor;

import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@RequiredArgsConstructor
@Component
public class RateLimitingInterceptor implements HandlerInterceptor {
    private final TokenBucketService tokenBucketService;
    private final int SC_TOO_MANY_REQUESTS = 429;

    @Override
    public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response,
            @NonNull Object handler) throws Exception {
        String userId = IdentityUtil.getClientIp(request);
        String endpoint = request.getRequestURI();

        if (userId == null) {
            response.sendError(HttpServletResponse.SC_BAD_REQUEST, "Missing user ID");
            return false;
        }

        if (!tokenBucketService.allowRequest(userId, endpoint)) {
            response.sendError(SC_TOO_MANY_REQUESTS, "Rate limit exceeded");
            return false;
        }

        return true;
    }
}