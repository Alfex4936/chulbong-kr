package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// RegisterCommentRoutes sets up the routes for comments handling within the application.
func RegisterCommentRoutes(api fiber.Router) {
	api.Get("/comments/:markerId/comments", loadCommentsHandler)

	commentGroup := api.Group("/comments")
	{
		commentGroup.Use(middlewares.AuthMiddleware)
		commentGroup.Post("", postCommentHandler)
		commentGroup.Patch("/:commentId", updateCommentHandler)
		commentGroup.Delete("/:commentId", removeCommentHandler)
	}
}

// postCommentHandler creates a new comment
func postCommentHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	var req dto.CommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	containsBadWord, _ := utils.CheckForBadWords(req.CommentText)
	if containsBadWord {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment contains inappropriate content."})
	}

	comment, err := services.CreateComment(req.MarkerID, userID, req.CommentText)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}
	return c.Status(fiber.StatusOK).JSON(comment)
}

func updateCommentHandler(c *fiber.Ctx) error {
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
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Failed to update the comment"})
		}
		// Handle other potential errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update the comment"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment updated successfully"})
}

func removeCommentHandler(c *fiber.Ctx) error {
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

func loadCommentsHandler(c *fiber.Ctx) error {
	var params dto.CommentLoadParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	markerID, err := c.ParamsInt("markerId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid marker ID provided",
		})
	}

	// default
	if params.N < 1 {
		params.N = 4
	}
	if params.Page < 1 {
		params.Page = 1
	}

	pageSize := 4 // Define page size
	offset := (params.Page - 1) * pageSize

	// Call service function to load comments for the marker
	comments, total, err := services.LoadCommentsForMarker(markerID, pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate total pages
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"comments":      comments,
		"currentPage":   params.Page,
		"totalPages":    totalPages,
		"totalComments": total,
	})
}
