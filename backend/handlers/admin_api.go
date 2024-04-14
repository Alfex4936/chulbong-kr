package handlers

import (
	"chulbong-kr/services"

	"github.com/gofiber/fiber/v2"
)

func ListUnreferencedS3ObjectsHandler(c *fiber.Ctx) error {
	killSwitch := c.Query("kill", "n")

	dbURLs, err := services.FetchAllPhotoURLsFromDB()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "fetching URLs from database"})
	}

	s3Keys, err := services.ListAllObjectsInS3()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "fetching keys from S3"})
	}

	unreferenced := services.FindUnreferencedS3Objects(dbURLs, s3Keys)

	if killSwitch == "y" {
		for _, unreferencedURL := range unreferenced {
			services.DeleteDataFromS3(unreferencedURL)
		}
	}

	return c.JSON(unreferenced)
}
