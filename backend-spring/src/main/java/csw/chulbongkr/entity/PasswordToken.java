package csw.chulbongkr.entity;


import jakarta.persistence.*;
import lombok.Data;
import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "PasswordTokens")
@Data
public class PasswordToken {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "TokenID")
    private Integer tokenID;

    @Column(name = "Token", nullable = false)
    private String token;

    @Column(name = "Email", nullable = false, unique = true)
    private String email;

    @Column(name = "Verified", nullable = false)
    private Boolean verified = false;

    @Column(name = "ExpiresAt", nullable = false)
    private LocalDateTime expiresAt;

    @CreationTimestamp
    @Column(name = "CreatedAt", updatable = false)
    private LocalDateTime createdAt;
}