package csw.chulbongkr.repository.auth;

import csw.chulbongkr.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Integer> {
    Optional<User> findByUsername(String username);
    Optional<User> findByEmail(String email);

    @Transactional
    @Modifying
    @Query("UPDATE User u SET u.passwordHash = :passwordHash WHERE u.userID = :userID")
    void updatePassword(Integer userID, String passwordHash);

    @Transactional
    @Modifying
    @Query(value = "INSERT INTO PasswordResetTokens (UserID, Token, ExpiresAt) VALUES (:userID, :token, :expiresAt) ON DUPLICATE KEY UPDATE token = VALUES(token), expiresAt = VALUES(expiresAt)", nativeQuery = true)
    void savePasswordResetToken(@Param("userID") Integer userID, @Param("token") String token, @Param("expiresAt") LocalDateTime expiresAt);

    @Query("SELECT p.userID FROM PasswordResetToken p WHERE p.token = :token AND p.expiresAt > CURRENT_TIMESTAMP")
    Optional<Integer> findUserIDByResetToken(String token);

    @Transactional
    @Modifying
    @Query("DELETE FROM PasswordResetToken p WHERE p.token = :token")
    void deleteResetToken(String token);
}
