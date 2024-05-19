package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *MarkerHandler) HandleGetMarkerRanking(c *fiber.Ctx) error {
	ranking := h.MarkerFacadeService.GetTopMarkers(50) // []dto.MarkerRank { MarkerID (string), Clicks (int) }

	return c.JSON(ranking)
}

func (h *MarkerHandler) HandleGetUniqueVisitorCount(c *fiber.Ctx) error {
	markerID := c.Query("markerId")
	if markerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	count := h.MarkerFacadeService.GetUniqueVisitorCount(markerID)

	return c.JSON(fiber.Map{"markerId": markerID, "visitors": count})
}

func (h *MarkerHandler) HandleGetAllUniqueVisitorCount(c *fiber.Ctx) error {
	count := h.MarkerFacadeService.GetAllUniqueVisitorCounts()
	return c.JSON(count)
}
