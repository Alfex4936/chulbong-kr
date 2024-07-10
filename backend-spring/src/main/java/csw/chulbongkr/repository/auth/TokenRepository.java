package csw.chulbongkr.repository.auth;

import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;

@Repository
public class TokenRepository {
    private final JdbcTemplate jdbcTemplate;

    public TokenRepository(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    public void saveOrUpdateToken(int userId, String token, LocalDateTime expiresAt) {
        String query = """
            INSERT INTO opaque_tokens (user_id, token, expires_at)
            VALUES (?, ?, ?)
            ON DUPLICATE KEY UPDATE token = VALUES(token), expires_at = VALUES(expires_at), updated_at = NOW()
            """;
        jdbcTemplate.update(query, userId, token, expiresAt);
    }
}
