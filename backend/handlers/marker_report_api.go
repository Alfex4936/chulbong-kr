package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RegisterReportRoutes sets up the routes for report handling within the application.
func RegisterReportRoutes(api fiber.Router) {
	reportGroup := api.Group("/reports")
	{
		reportGroup.Get("/all", getAllReportsHandler)
		reportGroup.Get("/marker/:markerID", getMarkerReportsHandler)
		reportGroup.Post("", middlewares.AuthSoftMiddleware, createReportHandler)
	}
}

// GetAllReportsHandler retrieves all reports for all markers, grouped by MarkerID.
func getAllReportsHandler(c *fiber.Ctx) error {
	reports, err := services.GetAllReports()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports"})
	}

	if len(reports) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No reports found"})
	}

	// Group reports by MarkerID
	groupedReports := make(map[int]dto.MarkerReports)
	for _, report := range reports {
		groupedReports[report.MarkerID] = dto.MarkerReports{
			Reports: append(groupedReports[report.MarkerID].Reports, report),
		}
	}

	// Create response structure
	response := dto.ReportsResponse{
		TotalReports: len(reports),
		Markers:      groupedReports,
	}

	return c.JSON(response)
}

func getMarkerReportsHandler(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	reports, err := services.GetAllReportsBy(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get reports"})
	}
	return c.JSON(reports)
}

func createReportHandler(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse form"})
	}

	// Check if latitude and longitude are provided
	// if user didn't change, frontend must send original point
	latitude, longitude, err := GetLatLngFromForm(form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse latitude and longitude"})
	}

	// Location Must Be Inside South Korea
	if !utils.IsInSouthKoreaPrecisely(latitude, longitude) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Operations only allowed within South Korea."})
	}

	description := GetDescriptionFromForm(form)
	if containsBadWord, _ := utils.CheckForBadWords(description); containsBadWord {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment contains inappropriate content."})
	}

	markerIDstr := GetMarkerIDFromForm(form)
	if markerIDstr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Inappropriate markerId."})
	}
	markerID, err := strconv.Atoi(markerIDstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	userID, _ := c.Locals("userID").(int) // userID will be 0 if not logged in

	err = services.CreateReport(&dto.MarkerReportRequest{
		MarkerID:    markerID,
		UserID:      userID,
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
	}, form)
	if err != nil {
		if strings.Contains(err.Error(), "an error during file") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "an error during file upload"})
		} else if strings.Contains(err.Error(), "no files file") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "upload at least one picture to prove"})
		} else if strings.Contains(err.Error(), "Error 1452 (23000)") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "check if a marker exists"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create report"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "report created successfully"})
}
