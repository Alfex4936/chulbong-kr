package handler

import (
	"strings"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/service"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	SearchService *service.ZincSearchService
}

// NewSearchHandler creates a new SearchHandler with dependencies injected
func NewSearchHandler(rank *service.ZincSearchService,
) *SearchHandler {
	return &SearchHandler{
		SearchService: rank,
	}
}

// RegisterSearchRoutes sets up the routes for search handling within the application.
func RegisterSearchRoutes(api fiber.Router, handler *SearchHandler) {
	searchGroup := api.Group("/search")
	{
		searchGroup.Get("/marker", handler.HandleSearchMarkerAddress)
		// searchGroup.Post("/marker", handler.HandleInsertMarkerAddressTest)
		// searchGroup.Delete("", handler.HandleDeleteMarkerAddressTest)
	}
}

// Handler for searching marker addresses
func (h *SearchHandler) HandleSearchMarkerAddress(c *fiber.Ctx) error {
	term := c.Query("term")
	term = strings.TrimSpace(term)
	if term == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search term is required",
		})
	}

	// Call the service function
	response, err := h.SearchService.SearchMarkerAddress(term)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Handler for searching marker addresses
func (h *SearchHandler) HandleInsertMarkerAddressTest(c *fiber.Ctx) error {
	// Call the service function
	err := h.SearchService.InsertMarkerIndex(dto.MarkerIndexData{MarkerID: 9999, Address: "test address"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("success")
}
func (h *SearchHandler) HandleDeleteMarkerAddressTest(c *fiber.Ctx) error {
	markerID := c.Query("markerId")
	markerID = strings.TrimSpace(markerID)
	if markerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search term is required",
		})
	}

	// Call the service function
	err := h.SearchService.DeleteMarkerIndex(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("success")
}
