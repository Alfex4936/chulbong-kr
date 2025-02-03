package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

const WEATHER_MINUTES = 15 * time.Minute

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

	if params.Distance > 50000 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Distance cannot be greater than 15,000m (15km)"})
	}

	// Set default page to 1 if not specified
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 4
	}

	offset := (params.Page - 1) * params.PageSize

	// Generate a cache key based on the query parameters
	cacheKey := fmt.Sprintf("close_markers:%f:%f:%d:%d:%d", params.Latitude, params.Longitude, params.Distance, params.Page, params.PageSize)

	// Attempt to fetch from cache
	cachedData, err := h.CacheService.GetCloseMarkersCache(cacheKey)
	if err == nil && len(cachedData) > 0 {
		// Cache hit, return the cached data
		c.Append("X-Cache", "hit")
		return c.Send(cachedData)
	}

	// Cache miss: Find nearby markers within the specified distance and page
	markers, total, err := h.MarkerFacadeService.FindClosestNMarkersWithinDistance(params.Latitude, params.Longitude, params.Distance, params.PageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers"})
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

	// Prepare the response data
	response := dto.MarkersClose{
		Markers:      markers,
		CurrentPage:  params.Page,
		TotalPages:   totalPages,
		TotalMarkers: total,
	}

	// Marshal the response for caching
	responseJSON, err := sonic.Marshal(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode response"})
	}

	// Cache the response for future use
	go h.CacheService.SetCloseMarkersCache(cacheKey, responseJSON, 10*time.Minute)

	// Return the response to the client
	return c.Send(responseJSON)
}

func (h *MarkerHandler) HandleGetCurrentAreaMarkerRanking(c *fiber.Ctx) error {
	limitParam := c.Query("limit", "10") // Default limit
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit"})
	}

	// "current area"
	const currentAreaDistance = 10000 // Meters

	markers, err := h.MarkerFacadeService.FindRankedMarkersInCurrentArea(lat, lng, currentAreaDistance, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve markers"})
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

	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate a short cache key by hashing the lat/long combination
	weather, cacheErr := h.CacheService.GetWcongCache(lat, lng)
	if cacheErr == nil && weather != nil {
		c.Append("X-Cache", "hit")
		// Cache hit, return cached weather (10mins)
		return c.JSON(weather)
	}

	result, err := h.MarkerFacadeService.FetchWeatherFromAddress(lat, lng)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch weather from address"})
	}

	// Cache the result for future requests
	go h.CacheService.SetWcongCache(lat, lng, result)

	return c.JSON(result)
}

func (h *MarkerHandler) HandleConvertWGS84ToWCONGNAMUL(c *fiber.Ctx) error {
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := util.ConvertWGS84ToWCONGNAMUL(lat, lng)

	return c.JSON(result)
}

func (h *MarkerHandler) HandleIsInSouthKorea(c *fiber.Ctx) error {
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := h.MarkerFacadeService.IsInSouthKoreaPrecisely(lat, lng)

	return c.JSON(fiber.Map{"result": result})
}

// DEPRECATED: Use version 2
func (h *MarkerHandler) HandleSaveOfflineMap(c *fiber.Ctx) error {
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pdf, err := h.MarkerFacadeService.SaveOfflineMap(lat, lng)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF: " + err.Error()})
	}

	return c.Download(pdf)
}

func (h *MarkerHandler) HandleTestDynamic(c *fiber.Ctx) error {
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sParam := c.Query("scale")
	wParam := c.Query("width")
	hParam := c.Query("height")
	scale, err := strconv.ParseFloat(sParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}
	width, err := strconv.ParseInt(wParam, 0, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}
	height, err := strconv.ParseInt(hParam, 0, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	h.MarkerFacadeService.TestDynamic(lat, lng, scale, width, height)
	return c.SendString("Dynamic API test")
}

func (h *MarkerHandler) HandleSaveOfflineMap2(c *fiber.Ctx) error {
	lat, lng, err := GetLatLong(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pdf, _, err := h.MarkerFacadeService.SaveOfflineMap2(lat, lng)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create a PDF"})
	}
	if pdf == "" {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "no content for this location"})
	}

	// Use Fiber's SendFile method
	// err = c.SendFile(pdf, true) // 'true' to enable compression
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send file"})
	// }
	// return nil
	return c.Download(pdf) // sendfile systemcall
}
