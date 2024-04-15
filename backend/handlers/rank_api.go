package handlers

import (
	"chulbong-kr/services"

	"github.com/gofiber/fiber/v2"
)

func getMarkerRankingHandler(c *fiber.Ctx) error {
	ranking := services.GetTopMarkers(10) // []dto.MarkerRank { MarkerID (string), Clicks (int) }

	return c.JSON(ranking)
}

func getUniqueVisitorCountHandler(c *fiber.Ctx) error {
	markerID := c.Query("markerId")
	if markerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	count := services.GetUniqueVisitorCount(markerID)

	return c.JSON(fiber.Map{"markerId": markerID, "visitors": count})
}

func getAllUniqueVisitorCountHandler(c *fiber.Ctx) error {
	count := services.GetAllUniqueVisitorCounts()
	return c.JSON(count)
}
