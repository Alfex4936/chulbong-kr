package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PostCommentHandler creates a new comment
func PostCommentHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	var req dto.CommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	comment, err := services.CreateComment(req.MarkerID, userID, req.CommentText)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}
	return c.Status(fiber.StatusOK).JSON(comment)
}

func UpdateCommentHandler(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("commentId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	// Extract userID and newCommentText from the request
	userID := c.Locals("userID").(int)
	var request struct {
		CommentText string `json:"commentText"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Call the service function to update the comment
	if err := services.UpdateComment(commentID, userID, request.CommentText); err != nil {
		if err.Error() == "comment not found or not owned by user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle other potential errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment updated successfully"})
}

func RemoveCommentHandler(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("commentId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	userID := c.Locals("userID").(int)

	err = services.RemoveComment(commentID, userID)
	if err != nil {
		if err.Error() == "comment not found or already deleted" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove comment"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment removed successfully"})
}

func LoadCommentsHandler(c *fiber.Ctx) error {
	markerID, err := c.ParamsInt("markerId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid marker ID provided",
		})
	}

	// Call service function to load comments for the marker
	comments, err := services.LoadCommentsForMarker(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(comments)
}
