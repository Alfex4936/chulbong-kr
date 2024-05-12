package handler

import (
	"strconv"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/gofiber/fiber/v2"
)

// Find Close Markers godoc
//
// @Summary		Find close markers
// @Description	This endpoint retrieves markers that are close to a specified location within a given distance.
// @Description	It requires latitude, longitude, distance, and the markers to return.
// @Description	If no markers are found within the specified distance, it returns a "No markers found" message.
// @Description	Returns a list of markers that meet the criteria. (maximum 10km distance allowed)
// @ID			find-close-markers
// @Tags		markers
// @Accept		json
// @Produce	json
// @Param		latitude	query	number	true	"Latitude of the location (float)"
// @Param		longitude	query	number	true	"Longitude of the location (float)"
// @Param		distance	query	int		true	"Search radius distance (meters)"
// @Param		N			query	int		true	"Page size"
// @Param		page			query	int		true	"Page Index number"
// @Security	ApiKeyAuth
// @Success	200	{object}	map[string]interface{}	"Markers found successfully (with distance) in pages"
// @Failure	400	{object}	map[string]interface{}	"Invalid query parameters"
// @Failure	404	{object}	map[string]interface{}	"No markers found within the specified distance"
// @Failure	500	{object}	map[string]interface{}	"Internal server error"
// @Router		/markers/close [get]
func (h *MarkerHandler) HandleFindCloseMarkers(c *fiber.Ctx) error {
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

	if params.PageSize < 1 {
		params.PageSize = 4
	}

	offset := (params.Page - 1) * params.PageSize

	// Find nearby markers within the specified distance and page
	markers, total, err := h.MarkerFacadeService.FindClosestNMarkersWithinDistance(params.Latitude, params.Longitude, params.Distance, params.PageSize, offset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate total pages
	totalPages := total / params.PageSize
	if total%params.PageSize != 0 {
		totalPages++
	}

	// Adjust the current page if the calculated offset exceeds the number of markers
	if params.Page > totalPages {
		params.Page = totalPages
	}
	if params.Page < 1 {
		params.Page = 1 // Ensure page is set to 1 if totalPages calculates to 0 (i.e., no markers found)
	}

	// Return the found markers along with pagination info
	return c.JSON(fiber.Map{
		"markers":      markers,
		"currentPage":  params.Page,
		"totalPages":   totalPages,
		"totalMarkers": total,
	})
}

func (h *MarkerHandler) HandleGetCurrentAreaMarkerRanking(c *fiber.Ctx) error {
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

	markers, err := h.MarkerFacadeService.FindRankedMarkersInCurrentArea(lat, long, currentAreaDistance, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers " + err.Error()})
	}

	if markers == nil {
		return c.JSON([]dto.MarkerWithDistance{})
	}

	return c.JSON(markers)
}

func (h *MarkerHandler) HandleGetMarkersClosebyAdmin(c *fiber.Ctx) error {
	markers, err := h.MarkerFacadeService.CheckNearbyMarkersInDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers: " + err.Error()})
	}

	return c.JSON(markers)
}

func (h *MarkerHandler) HandleGetWeatherByWGS84(c *fiber.Ctx) error {
	// Check the Referer header and redirect if it matches the specific URL pattern
	// if !strings.HasSuffix(c.Get("Referer"), ".k-pullup.com") || c.Get("Referer") != "https://www.k-pullup.com/" {
	// 	return c.Redirect("https://k-pullup.com", fiber.StatusFound) // Use HTTP 302 for standard redirection
	// }

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

	result, err := h.MarkerFacadeService.FetchWeatherFromAddress(lat, long)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch weather from address: " + err.Error()})
	}

	return c.JSON(result)
}
func (h *MarkerHandler) HandleConvertWGS84ToWCONGNAMUL(c *fiber.Ctx) error {
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

	result := util.ConvertWGS84ToWCONGNAMUL(lat, long)

	return c.JSON(result)
}

func (h *MarkerHandler) HandleIsInSouthKorea(c *fiber.Ctx) error {
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

	result := h.MarkerFacadeService.IsInSouthKoreaPrecisely(lat, long)

	return c.JSON(fiber.Map{"result": result})
}

// DEPRECATED: Use version 2
func (h *MarkerHandler) HandleSaveOfflineMap(c *fiber.Ctx) error {
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

	pdf, err := h.MarkerFacadeService.SaveOfflineMap(lat, long)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF: " + err.Error()})
	}

	return c.Download(pdf)
}

func (h *MarkerHandler) HandleSaveOfflineMap2(c *fiber.Ctx) error {
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

	pdf, err := h.MarkerFacadeService.SaveOfflineMap2(lat, long)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF"})
	}
	if pdf == "" {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "no content for this location"})
	}

	return c.Download(pdf)
}
