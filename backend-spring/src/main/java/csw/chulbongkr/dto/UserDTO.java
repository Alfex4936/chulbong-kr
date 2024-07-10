package csw.chulbongkr.dto;

import csw.chulbongkr.entity.User;

public class UserDTO {
    public record UserProfile(Integer userId, String username, String email) {}

    public record LoginResponse(User user, String token) {}
    public record LoginRequest(String email, String password) {}
}
