package handlers

import (
	"chulbong-kr/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetMarkerRanking(c *fiber.Ctx) error {
	currentTime := time.Now().Format("20060102")

	ranking := services.EstimateMarkerPopularity(currentTime)

	return c.JSON(ranking)
}
