package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Find Close Markers godoc
//
// @Summary		Find close markers
// @Description	This endpoint retrieves markers that are close to a specified location within a given distance.
// @Description	It requires latitude, longitude, distance, and the number of markers (N) to return.
// @Description	If no markers are found within the specified distance, it returns a "No markers found" message.
// @Description	Returns a list of markers that meet the criteria. (maximum 10km distance allowed)
// @ID			find-close-markers
// @Tags		markers
// @Accept		json
// @Produce	json
// @Param		latitude	query	number	true	"Latitude of the location (float)"
// @Param		longitude	query	number	true	"Longitude of the location (float)"
// @Param		distance	query	int		true	"Search radius distance (meters)"
// @Param		N			query	int		true	"Number of markers to return"
// @Param		page			query	int		true	"Page Index number"
// @Security	ApiKeyAuth
// @Success	200	{object}	map[string]interface{}	"Markers found successfully (with distance) in pages"
// @Failure	400	{object}	map[string]interface{}	"Invalid query parameters"
// @Failure	404	{object}	map[string]interface{}	"No markers found within the specified distance"
// @Failure	500	{object}	map[string]interface{}	"Internal server error"
// @Router		/markers/close [get]
func findCloseMarkersHandler(c *fiber.Ctx) error {
	var params dto.QueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	if params.Distance > 10000 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Distance cannot be greater than 10,000m (10km)"})
	}

	// Set default page to 1 if not specified
	if params.Page < 1 {
		params.Page = 1
	}

	pageSize := 4 // Define page size
	offset := (params.Page - 1) * pageSize

	// Find nearby markers within the specified distance and page
	markers, total, err := services.FindClosestNMarkersWithinDistance(params.Latitude, params.Longitude, params.Distance, pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// if len(markers) == 0 {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No markers found within the specified distance"})
	// }

	// Calculate total pages
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	// Return the found markers along with pagination info
	return c.JSON(fiber.Map{
		"markers":      markers,
		"currentPage":  params.Page,
		"totalPages":   totalPages,
		"totalMarkers": total,
	})
}

func getCurrentAreaMarkerRankingHandler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")
	limitParam := c.Query("limit", "5") // Default limit

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}

	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit"})
	}

	// "current area"
	const currentAreaDistance = 10000 // Meters

	markers, err := services.FindRankedMarkersInCurrentArea(lat, long, currentAreaDistance, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers"})
	}

	if markers == nil {
		return c.JSON([]dto.MarkerWithDistance{})
	}

	return c.JSON(markers)
}

func getMarkersClosebyAdmin(c *fiber.Ctx) error {
	markers, err := services.CheckNearbyMarkersInDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers: " + err.Error()})
	}

	return c.JSON(markers)
}

func getWeatherByWGS84Handler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}

	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	result, err := services.FetchWeatherFromAddress(lat, long)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch weather from address: " + err.Error()})
	}

	return c.JSON(result)
}
func convertWGS84ToWCONGNAMULHandler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}

	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	result := utils.ConvertWGS84ToWCONGNAMUL(lat, long)

	return c.JSON(result)
}

func isInSouthKoreaHandler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}

	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	result := utils.IsInSouthKoreaPrecisely(lat, long)

	return c.JSON(fiber.Map{"result": result})
}

// DEPRECATED: Use version 2
func saveOfflineMapHandler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}
	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	pdf, err := services.SaveOfflineMap(lat, long)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF: " + err.Error()})
	}

	return c.Download(pdf)
}

func saveOfflineMap2Handler(c *fiber.Ctx) error {
	latParam := c.Query("latitude")
	longParam := c.Query("longitude")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}
	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	pdf, err := services.SaveOfflineMap2(lat, long)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF"})
	}
	if pdf == "" {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "no content for this location"})
	}

	return c.Download(pdf)
}
