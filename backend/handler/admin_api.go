package handler

import (
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	AdminFacade *facade.AdminFacadeService
	UserService *service.UserService

	TokenUtil *util.TokenUtil
}

// NewAdminHandler creates a new AdminHandler with dependencies injected
func NewAdminHandler(
	admin *facade.AdminFacadeService,
	user *service.UserService,
	tutil *util.TokenUtil,
) *AdminHandler {
	return &AdminHandler{
		AdminFacade: admin,
		UserService: user,
		TokenUtil:   tutil,
	}
}

// RegisterAdminRoutes sets up the routes for admin handling within the application.
func RegisterAdminRoutes(api fiber.Router, handler *AdminHandler, authMiddleware *middleware.AuthMiddleware) {
	api.Post("/chat/ban/:markerID/:userID", authMiddleware.CheckAdmin, handler.HandleBanUser)

	adminGroup := api.Group("/admin")
	{
		adminGroup.Use(authMiddleware.CheckAdmin)
		adminGroup.Get("/dead", handler.HandleListUnreferencedS3Objects)
		adminGroup.Get("/fetch", handler.HandleListUpdatedMarkers)
	}
}

func (h *AdminHandler) HandleListUnreferencedS3Objects(c *fiber.Ctx) error {
	killSwitch := c.Query("kill", "n")

	dbURLs, err := h.AdminFacade.FetchAllPhotoURLsFromDB()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "fetching URLs from database"})
	}

	s3Keys, err := h.AdminFacade.ListAllObjectsInS3()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "fetching keys from S3"})
	}

	unreferenced := h.AdminFacade.FindUnreferencedS3Objects(dbURLs, s3Keys)

	if killSwitch == "y" {
		for _, unreferencedURL := range unreferenced {
			_ = h.AdminFacade.DeleteDataFromS3(unreferencedURL)
		}
	}

	return c.JSON(unreferenced)
}

func (h *AdminHandler) HandleBanUser(c *fiber.Ctx) error {
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
	err := h.AdminFacade.BanUser(markerID, userID, duration)
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

func (h *AdminHandler) HandleListUpdatedMarkers(c *fiber.Ctx) error {
	postSwitch := c.Query("post", "n")
	currentDateString := c.Query("date", time.Now().Format("2006-01-02"))

	currentDate, _ := time.Parse("2006-1-2", currentDateString)

	markers, err := h.AdminFacade.FetchLatestMarkers(currentDate)
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

			if fErr := h.AdminFacade.CheckMarkerValidity(latitude, longitude, ""); fErr != nil {
				log.Printf("➡️ Skipping: %s\n", fErr.Message)
				continue
			}

			userID := c.Locals("userID").(int)

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

			marker, err := h.AdminFacade.CreateMarkerWithPhotos(&dto.MarkerRequest{
				Latitude:    latitude,
				Longitude:   longitude,
				Description: "",
			}, userID, form)
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

			if err := h.AdminFacade.SetMarkerFacilities(newMarkerID, []dto.FacilityQuantity{
				{FacilityID: 1, Quantity: m.ChulbongCount},
				{FacilityID: 2, Quantity: m.PyeongCount},
			}); err != nil {
				continue
				// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set facilities for marker"})
			}

		}
	}

	h.AdminFacade.ResetMarkerCache()

	return c.JSON(markers)
}
