package handlers

import (
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RegisterAdminRoutes sets up the routes for admin handling within the application.
func RegisterAdminRoutes(api fiber.Router) {
	api.Post("/chat/ban/:markerID/:userID", middlewares.AdminOnly, banUserHandler)

	adminGroup := api.Group("/admin")
	{
		adminGroup.Use(middlewares.AdminOnly)
		adminGroup.Get("/dead", listUnreferencedS3ObjectsHandler)
	}
}

func listUnreferencedS3ObjectsHandler(c *fiber.Ctx) error {
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

func banUserHandler(c *fiber.Ctx) error {
	// Extract markerID and userID from the path parameters
	markerID := c.Params("markerID")
	userID := c.Params("userID")

	// assert duration is sent in the request body as JSON
	var requestBody struct {
		DurationInMinutes int `json:"duration"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		requestBody = struct {
			DurationInMinutes int `json:"duration"`
		}{
			DurationInMinutes: 5, // default 5 minutes banned
		}
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": "Invalid request format",
		// })
	}

	if requestBody.DurationInMinutes < 1 {
		requestBody.DurationInMinutes = 5
	} else if requestBody.DurationInMinutes > 15 {
		requestBody.DurationInMinutes = 10 // max 10 minutes
	}

	// Convert duration to time.Duration
	duration := time.Duration(requestBody.DurationInMinutes) * time.Minute

	// Call the BanUser method on the manager instance
	err := services.WsRoomManager.BanUser(markerID, userID, duration)
	if err != nil {
		// Log the error or handle it as needed
		// log.Printf("Error banning user %s from room %s: %v", userID, markerID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to ban user",
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "User successfully banned",
		"time":    duration,
	})
}
