package csw.chulbongkr.controller.auth;

import csw.chulbongkr.dto.UserDTO;
import csw.chulbongkr.entity.User;
import csw.chulbongkr.service.auth.AuthService;
import csw.chulbongkr.service.auth.TokenService;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;
    private final TokenService tokenService;

    @PostMapping("/login")
    public ResponseEntity<UserDTO.LoginResponse> handleLogin(@RequestBody UserDTO.LoginRequest request, HttpServletResponse response) {
        try {
            User user = authService.login(request.email(), request.password());
            String token = tokenService.generateAndSaveToken(user.getUserID(), user.getEmail());

            UserDTO.LoginResponse loginResponse = new UserDTO.LoginResponse(user, token);

            Cookie cookie = new Cookie("Authorization", token);
            cookie.setHttpOnly(true);
            cookie.setSecure(true);
            cookie.setPath("/");
            cookie.setMaxAge((int) tokenService.getTokenExpirationTime().toSeconds());
            response.addCookie(cookie);

            return ResponseEntity.ok(loginResponse);
        } catch (Exception e) {
            return ResponseEntity.status(401).body(new UserDTO.LoginResponse(null, "Invalid email or password"));
        }
    }
}