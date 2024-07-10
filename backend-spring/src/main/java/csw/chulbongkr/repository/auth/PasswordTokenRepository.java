package csw.chulbongkr.repository.auth;

import csw.chulbongkr.entity.PasswordToken;
import org.springframework.data.jpa.repository.JpaRepository;

public interface PasswordTokenRepository extends JpaRepository<PasswordToken, Integer> {
    PasswordToken findByToken(String token);
    void deleteByToken(String token);
}
