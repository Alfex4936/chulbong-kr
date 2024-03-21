package handlers

import (
	"chulbong-kr/services"

	"github.com/gofiber/fiber/v2"
)

func GetMarkerRankingHandler(c *fiber.Ctx) error {
	ranking := services.GetTopMarkers(10) // []dto.MarkerRank { MarkerID (string), Clicks (int) }

	return c.JSON(ranking)
}
