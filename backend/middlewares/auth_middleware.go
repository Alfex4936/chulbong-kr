package middlewares

import (
	"chulbong-kr/database"
	"chulbong-kr/util"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks for a valid opaque token in the Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	// check for the cookie
	jwtCookie := c.Cookies(util.LOGIN_TOKEN_COOKIE)
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}
	token := jwtCookie

	// ExpiresAt > CURRENT_TIMESTAMP doesn't work well
	query := `SELECT UserID, ExpiresAt FROM OpaqueTokens WHERE OpaqueToken = ?`
	var userID int
	var expiresAt time.Time
	err := database.DB.QueryRow(query, token).Scan(&userID, &expiresAt)

	// Token is invalid or expired, delete the cookie
	if err == sql.ErrNoRows || time.Now().After(expiresAt) {
		cookie := util.ClearLoginCookie()
		c.Cookie(&cookie)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	// Fetch UserID and Username based on Email
	userQuery := `SELECT Username, Email, Role FROM Users WHERE UserID = ?`
	var username, email, role string
	err = database.DB.QueryRow(userQuery, userID).Scan(&username, &email, &role)
	if err != nil {
		cookie := util.ClearLoginCookie()
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

// AdminOnly checks admin permission
func AdminOnly(c *fiber.Ctx) error {
	// Check for the cookie
	jwtCookie := c.Cookies(util.LOGIN_TOKEN_COOKIE)
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	// query to also fetch the user's role
	query := `
    SELECT u.UserID, u.Role, ot.ExpiresAt 
    FROM OpaqueTokens ot
    JOIN Users u ON ot.UserID = u.UserID
    WHERE ot.OpaqueToken = ?`

	var userID int
	var role string
	var expiresAt time.Time
	err := database.DB.QueryRow(query, jwtCookie).Scan(&userID, &role, &expiresAt)

	// Check for no rows found or other errors
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	// Check if the token has expired
	if time.Now().After(expiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	}

	// Check if the user role is not admin
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	c.Locals("isAdmin", true)

	// Proceed to the next handler if the user is an admin
	return c.Next()
}

// AuthSoftMiddleware checks for a valid opaque token in the Authorization header (no error returns)
func AuthSoftMiddleware(c *fiber.Ctx) error {
	// check for the cookie
	jwtCookie := c.Cookies(util.LOGIN_TOKEN_COOKIE)
	if jwtCookie == "" {
		return c.Next()
	}
	token := jwtCookie

	query := `SELECT UserID, ExpiresAt FROM OpaqueTokens WHERE OpaqueToken = ?`
	var userID int
	var expiresAt time.Time
	err := database.DB.QueryRow(query, token).Scan(&userID, &expiresAt)

	// Adjust the error check to specifically look for no rows found, indicating an invalid or expired token.
	if err == sql.ErrNoRows {
		return c.Next()
	} else if err != nil {
		return c.Next()
	}

	if time.Now().After(expiresAt) {
		return c.Next()
	}

	// Fetch based on Email
	userQuery := `SELECT Username, Email, Role FROM Users WHERE UserID = ?`
	var username string
	var email string
	var chulbong string
	err = database.DB.QueryRow(userQuery, userID).Scan(&username, &email, &chulbong)
	if err != nil {
		return c.Next()
	}

	// Store in locals for use in subsequent handlers
	c.Locals("userID", userID)
	c.Locals("username", username)
	c.Locals("email", email)
	c.Locals("chulbong", chulbong == "admin")

	return c.Next()
}
