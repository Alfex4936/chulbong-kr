package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"log"
	"mime/multipart"
	"strconv"
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
		adminGroup.Get("/fetch", listUpdatedMarkersHandler)
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
			_ = services.DeleteDataFromS3(unreferencedURL)
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

func listUpdatedMarkersHandler(c *fiber.Ctx) error {
	postSwitch := c.Query("post", "n")
	currentDateString := c.Query("date", time.Now().Format("2006-01-02"))

	currentDate, _ := time.Parse("2006-1-2", currentDateString)

	markers, err := services.FetchLatestMarkers(currentDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if postSwitch == "y" {
		for _, m := range markers {
			latitude, err := strconv.ParseFloat(string(m.Latitude), 64)
			if err != nil {
				continue
			}

			longitude, err := strconv.ParseFloat(string(m.Longitude), 64)
			if err != nil {
				continue
			}

			if fErr := services.CheckMarkerValidity(latitude, longitude, ""); fErr != nil {
				log.Printf("➡️ Skipping: %s\n", fErr.Message)
				continue
			}

			userId := c.Locals("userID").(int)

			latitudeForm := []string{string(m.Latitude)}
			longitudeForm := []string{string(m.Longitude)}

			// Create the form with the initial value map containing the latitude and longitude.
			form := &multipart.Form{
				Value: map[string][]string{
					"latitude":  latitudeForm,
					"longitude": longitudeForm,
				},
				File: nil, // No file uploads are being handled
			}

			marker, err := services.CreateMarkerWithPhotos(&dto.MarkerRequest{
				Latitude:    latitude,
				Longitude:   longitude,
				Description: "",
			}, userId, form)
			if err != nil {
				log.Println("Error creating marker:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating marker"})
			}

			newMarkerID := marker.MarkerID

			if newMarkerID == 0 {
				log.Println("Error creating marker with 0:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating marker"})
			}

			// Now, prepare the request for setting facilities
			if m.ChulbongCount < 1 || m.PyeongCount < 1 {
				continue
				// return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "No facilities to add"})
			}

			if err := services.SetMarkerFacilities(newMarkerID, []dto.FacilityQuantity{
				{FacilityID: 1, Quantity: m.ChulbongCount},
				{FacilityID: 2, Quantity: m.PyeongCount},
			}); err != nil {
				continue
				// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set facilities for marker"})
			}

		}
	}

	return c.JSON(markers)
}
