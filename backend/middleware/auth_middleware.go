package middleware

import (
	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"

	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	getUserRoleQuery = "SELECT Username, Email, Role FROM Users WHERE UserID = ?"
)

type AuthMiddleware struct {
	AuthService *service.AuthService
	Config      *config.AppConfig
	TokenUtil   *util.TokenUtil

	getUserRoleStmt *sql.Stmt
}

func NewAuthMiddleware(authService *service.AuthService, config *config.AppConfig, token *util.TokenUtil) *AuthMiddleware {
	getUserRoleStmt, _ := authService.DB.Prepare(getUserRoleQuery)

	return &AuthMiddleware{
		AuthService: authService, Config: config, TokenUtil: token,

		getUserRoleStmt: getUserRoleStmt,
	}
}

// Verify checks for a valid opaque token in the Authorization header
func (m *AuthMiddleware) Verify(c *fiber.Ctx) error {
	// check for the cookie
	jwtCookie := c.Cookies(m.Config.LoginTokenCookie)
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	// ExpiresAt > CURRENT_TIMESTAMP doesn't work well
	userID, expiresAt, err := m.AuthService.VerifyOpaqueToken(jwtCookie)

	// Token is invalid or expired, delete the cookie
	if err != nil {
		cookie := m.TokenUtil.ClearLoginCookie()
		c.Cookie(&cookie)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	if time.Now().After(expiresAt) {
		cookie := m.TokenUtil.ClearLoginCookie()
		c.Cookie(&cookie)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	}

	// Fetch UserID and Username based on Email
	var username, email, role string
	err = m.getUserRoleStmt.QueryRow(userID).Scan(&username, &email, &role)

	if err != nil {
		cookie := m.TokenUtil.ClearLoginCookie()
		c.Cookie(&cookie)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error fetching user details"})
	}

	// Store UserID, Username and Email in locals for use in subsequent handlers
	c.Locals("userID", userID)
	c.Locals("username", username)
	c.Locals("email", email)
	c.Locals("role", role)

	// log.Printf("[DEBUG] Authenticated. %s", email)
	return c.Next()
}

// CheckAdmin checks admin permission
func (m *AuthMiddleware) CheckAdmin(c *fiber.Ctx) error {
	// Check for the cookie
	jwtCookie := c.Cookies(m.Config.LoginTokenCookie)
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	userDetails, err := m.AuthService.FetchUserDetails(jwtCookie, false)

	// Check for no rows found or other errors
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	// Check if the token has expired
	if time.Now().After(userDetails.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	}

	// Check if the user role is not admin
	if userDetails.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	c.Locals("userID", userDetails.UserID)
	c.Locals("isAdmin", true)

	// Proceed to the next handler if the user is an admin
	return c.Next()
}

// VerifySoft checks for a valid opaque token in the Authorization header (no error returns)
func (m *AuthMiddleware) VerifySoft(c *fiber.Ctx) error {
	// check for the cookie
	jwtCookie := c.Cookies(m.Config.LoginTokenCookie)
	if jwtCookie == "" {
		return c.Next()
	}
	userDetails, err := m.AuthService.FetchUserDetails(jwtCookie, true)

	// Adjust the error check to specifically look for no rows found, indicating an invalid or expired token.
	if err == sql.ErrNoRows {
		return c.Next()
	} else if err != nil {
		return c.Next()
	}

	if time.Now().After(userDetails.ExpiresAt) {
		return c.Next()
	}

	// Store in locals for use in subsequent handlers
	c.Locals("userID", userDetails.UserID)
	c.Locals("username", userDetails.Username)
	c.Locals("email", userDetails.Email)
	c.Locals("chulbong", userDetails.Role == "admin")

	return c.Next()
}
