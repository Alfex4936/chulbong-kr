package csw.chulbongkr.util;

import jakarta.servlet.http.HttpServletRequest;

public class IdentityUtil {
    public static String getClientIp(HttpServletRequest request) {
        if (request == null) {
            return "";
        }

        String clientIp = request.getHeader("Fly-Client-IP");
        if (clientIp == null || clientIp.isEmpty()) {
            clientIp = request.getHeader("Fly-Client-Ip");
        }
        if (clientIp == null || clientIp.isEmpty()) {
            clientIp = request.getHeader("X-Forwarded-For");
        }
        if (clientIp == null || clientIp.isEmpty()) {
            clientIp = request.getHeader("X-Real-IP");
        }
        if (clientIp == null || clientIp.isEmpty()) {
            clientIp = request.getRemoteAddr();
        }
        return clientIp;
    }
}
