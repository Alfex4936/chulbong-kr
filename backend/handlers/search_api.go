package handlers

import (
	"chulbong-kr/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RegisterSearchRoutes sets up the routes for search handling within the application.
func RegisterSearchRoutes(api fiber.Router) {
	searchGroup := api.Group("/search")
	{
		searchGroup.Get("/marker", searchMarkerAddressHandler)
	}
}

// Handler for searching marker addresses
func searchMarkerAddressHandler(c *fiber.Ctx) error {
	term := c.Query("term")
	term = strings.TrimSpace(term)
	if term == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search term is required",
		})
	}

	// Call the service function
	response, err := services.SearchMarkerAddress(term)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
