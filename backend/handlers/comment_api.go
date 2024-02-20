package handlers

import (
	"chulbong-kr/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PostComment creates a new comment
func PostComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	markerID, _ := strconv.Atoi(c.FormValue("markerId"))
	commentText := c.FormValue("commentText")

	comment, err := services.CreateComment(c.Context(), markerID, userID, commentText)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}
	return c.Status(fiber.StatusOK).JSON(comment)
}

// UpdateComment edits an existing comment
func UpdateComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	commentID, _ := strconv.Atoi(c.Params("commentId"))
	commentText := c.FormValue("commentText")

	comment, err := services.UpdateComment(c.Context(), commentID, userID, commentText)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}
	return c.Status(fiber.StatusOK).JSON(comment)
}

// DeleteComment removes a comment
func DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	commentID, _ := strconv.Atoi(c.Params("commentId"))

	err := services.DeleteComment(c.Context(), commentID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comment"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment deleted successfully"})
}
