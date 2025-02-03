package handler

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/gofiber/fiber/v2"
)

// MarkerHandler.go
func (h *MarkerHandler) HandleAddStory(c *fiber.Ctx) error {
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	userID := c.Locals("userID").(int)

	// Parse the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	caption := ""
	if captions, ok := form.Value["caption"]; ok && len(captions) > 0 {
		caption = captions[0]
		if len(caption) > 30 {
			caption = caption[:30]
		}
	}

	caption, _ = h.MarkerFacadeService.BadWordUtil.ReplaceBadWords(caption)

	// Get the photo
	var photo *multipart.FileHeader
	if photos, ok := form.File["photo"]; ok && len(photos) > 0 {
		photo = photos[0]
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Photo is required"})
	}

	// Call the service to add the story
	storyResponse, err := h.MarkerFacadeService.StoryService.AddStory(markerID, userID, caption, photo)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyStoryPost) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Story already posted"}) // 409
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add story"})
	}

	return c.Status(fiber.StatusCreated).JSON(storyResponse)
}

func (h *MarkerHandler) HandleGetStories(c *fiber.Ctx) error {
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "30"))
	if err != nil || pageSize < 1 {
		pageSize = 30
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Call the service to get stories
	stories, err := h.MarkerFacadeService.StoryService.GetStories(markerID, offset, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get stories"})
	}

	return c.JSON(stories)
}

func (h *MarkerHandler) HandleGetAllStories(c *fiber.Ctx) error {
	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeParam := c.Query("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Call the service to get all stories
	stories, err := h.MarkerFacadeService.StoryService.GetAllStories(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get stories" + err.Error()})
	}

	return c.JSON(stories)
}

func (h *MarkerHandler) HandleDeleteStory(c *fiber.Ctx) error {
	markerIDParam := c.Params("markerID")
	storyIDParam := c.Params("storyID")

	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}
	storyID, err := strconv.Atoi(storyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid story ID"})
	}

	userRole := c.Locals("role").(string)
	userID := c.Locals("userID").(int)

	// Call the service to delete the story
	err = h.MarkerFacadeService.StoryService.DeleteStory(markerID, storyID, userID, userRole)
	if err != nil {
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to delete this story"})
		} else if err == service.ErrStoryNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Story not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete story"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Story deleted"})
}

func (h *MarkerHandler) HandleAddReaction(c *fiber.Ctx) error {
	storyIDParam := c.Params("storyID")
	storyID, err := strconv.Atoi(storyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid story ID"})
	}

	userID := c.Locals("userID").(int)

	// Parse the request body to get the reaction type
	var reactionRequest dto.ReactionRequest
	if err := c.BodyParser(&reactionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if reactionRequest.ReactionType != "thumbsup" && reactionRequest.ReactionType != "thumbsdown" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid reaction type"})
	}

	// Call the service to add the reaction
	err = h.MarkerFacadeService.StoryService.AddReaction(storyID, userID, reactionRequest.ReactionType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add reaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Reaction added"})
}

func (h *MarkerHandler) HandleRemoveReaction(c *fiber.Ctx) error {
	storyIDParam := c.Params("storyID")
	storyID, err := strconv.Atoi(storyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid story ID"})
	}

	userID := c.Locals("userID").(int)

	// Call the service to remove the reaction
	err = h.MarkerFacadeService.StoryService.RemoveReaction(storyID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove reaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Reaction removed"})
}

func (h *MarkerHandler) HandleReportStory(c *fiber.Ctx) error {
	storyIDParam := c.Params("storyID")
	storyID, err := strconv.Atoi(storyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid story ID"})
	}

	userID := c.Locals("userID").(int)

	// Get reason from request body
	var reportRequest struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&reportRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(reportRequest.Reason) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reason is too long"})
	}

	// Call the service to report the story
	err = h.MarkerFacadeService.StoryService.ReportStory(storyID, userID, reportRequest.Reason)
	if err != nil {
		if errors.Is(err, service.ErrStoryNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Story not found"})
		}
		// Handle duplicate report error
		if strings.Contains(err.Error(), "duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You have already reported this story"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to report story"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Story reported"})
}
