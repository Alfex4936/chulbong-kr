package middlewares

import (
	"chulbong-kr/database"
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks for a valid opaque token in the Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	// Split the Authorization header to extract the token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer {token}"})
	}

	token := parts[1] // The actual token part

	query := `SELECT Email, ExpiresAt FROM OpaqueTokens WHERE OpaqueToken = ?`
	var email string
	var expiresAt time.Time
	err := database.DB.QueryRow(query, token).Scan(&email, &expiresAt)

	// Adjust the error check to specifically look for no rows found, indicating an invalid or expired token.
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	if time.Now().After(expiresAt) {
		// // Token has expired, delete it
		// delQuery := `DELETE FROM OpaqueTokens WHERE OpaqueToken = ?`
		// if _, delErr := database.DB.Exec(delQuery, token); delErr != nil {
		// 	// Log the error; decide how you want to handle the failure of deleting an expired token
		// 	fmt.Println("Failed to delete expired token:", delErr)
		// }
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	}

	// Fetch UserID and Username based on Email
	userQuery := `SELECT UserID, Username FROM Users WHERE Email = ?`
	var userID int
	var username string
	err = database.DB.QueryRow(userQuery, email).Scan(&userID, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error fetching user details"})
	}

	// Store UserID, Username and Email in locals for use in subsequent handlers
	c.Locals("userID", userID)
	c.Locals("username", username)
	c.Locals("email", email)
	return c.Next()
}
