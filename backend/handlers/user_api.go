package handlers

import (
	"chulbong-kr/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

// DeleteUserHandler deletes the currently authenticated user
func DeleteUserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}

	log.Printf("[DEBUG][HANDLER] Deleting user %v", userID)

	if err := services.DeleteUserWithRelatedData(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 for successful deletion with no content in response
}
