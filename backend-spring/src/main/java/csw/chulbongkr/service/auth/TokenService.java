package csw.chulbongkr.service.auth;

import csw.chulbongkr.config.security.JwtTokenProvider;
import csw.chulbongkr.entity.PasswordToken;
import csw.chulbongkr.repository.auth.PasswordTokenRepository;
import csw.chulbongkr.repository.auth.TokenRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import javax.crypto.SecretKey;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class TokenService {

    private final PasswordTokenRepository passwordTokenRepository;
    private final JwtTokenProvider jwtTokenProvider;

    public String generateAndSaveToken(Integer userId, String email) {
        String token = jwtTokenProvider.generateToken(userId);

        PasswordToken passwordToken = new PasswordToken();
        passwordToken.setToken(token);
        passwordToken.setEmail(email);
        passwordToken.setExpiresAt(LocalDateTime.now().plus(jwtTokenProvider.getTokenExpirationTime()));
        passwordToken.setVerified(false);

        passwordTokenRepository.save(passwordToken);

        return token;
    }

    public boolean validateToken(String token) {
        return jwtTokenProvider.validateToken(token);
    }

    public boolean verifyToken(String token) {
        PasswordToken passwordToken = passwordTokenRepository.findByToken(token);

        if (passwordToken == null || passwordToken.getExpiresAt().isBefore(LocalDateTime.now())) {
            return false;
        }

        passwordToken.setVerified(true);
        passwordTokenRepository.save(passwordToken);

        return true;
    }

    public void deleteToken(String token) {
        passwordTokenRepository.deleteByToken(token);
    }

    public Duration getTokenExpirationTime() {
        return jwtTokenProvider.getTokenExpirationTime();
    }
}
