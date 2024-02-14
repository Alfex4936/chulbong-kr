package middlewares

import (
	"chulbong-kr/database"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks for a valid opaque token in the Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

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

	c.Locals("email", email) // Store email in locals for use in the handler
	return c.Next()
}
