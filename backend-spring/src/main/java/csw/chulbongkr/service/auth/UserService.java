package csw.chulbongkr.service.auth;

import csw.chulbongkr.entity.User;
import csw.chulbongkr.repository.auth.UserRepository;
import csw.chulbongkr.util.TokenUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;

@RequiredArgsConstructor
@Service
public class UserService {
    private final UserRepository userRepository;
    private final BCryptPasswordEncoder passwordEncoder;
    private final TokenUtil tokenUtil;

    public Optional<User> getUserById(Integer id) {
        return userRepository.findById(id);
    }

    public Optional<User> getUserByEmail(String email) {
        return userRepository.findByEmail(email);
    }
    public User registerUser(String username, String email, String password, String provider, String providerID) {
        String hashedPassword = passwordEncoder.encode(password);
        User user = new User();
        user.setUsername(username);
        user.setEmail(email);
        user.setPasswordHash(hashedPassword);
        user.setProvider(provider);
        user.setProviderID(providerID);
        user.setRole("user");
        return userRepository.save(user);
    }

    public Optional<User> loginUser(String email, String password) {
        Optional<User> userOptional = userRepository.findByEmail(email);
        if (userOptional.isPresent()) {
            User user = userOptional.get();
            if (passwordEncoder.matches(password, user.getPasswordHash())) {
                return Optional.of(user);
            }
        }
        return Optional.empty();
    }

    public String generatePasswordResetToken(String email) {
        Optional<User> userOptional = userRepository.findByEmail(email);
        if (userOptional.isPresent()) {
            String token = tokenUtil.generateOpaqueToken(16);
            userRepository.savePasswordResetToken(userOptional.get().getUserID(), token, LocalDateTime.now().plusDays(1));
            return token;
        }
        return null;
    }

    public boolean resetPassword(String token, String newPassword) {
        Optional<Integer> userIDOptional = userRepository.findUserIDByResetToken(token);
        if (userIDOptional.isPresent()) {
            Integer userID = userIDOptional.get();
            String hashedPassword = passwordEncoder.encode(newPassword);
            userRepository.updatePassword(userID, hashedPassword);
            userRepository.deleteResetToken(token);
            return true;
        }
        return false;
    }
}
