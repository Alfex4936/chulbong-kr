package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/model"

	"github.com/gofiber/fiber/v2"
)

type MarkerHandler struct {
	MarkerFacadeService *facade.MarkerFacadeService

	AuthMiddleware *middleware.AuthMiddleware
}

// NewMarkerHandler creates a new MarkerHandler with dependencies injected
func NewMarkerHandler(authMiddleware *middleware.AuthMiddleware, facade *facade.MarkerFacadeService,
) *MarkerHandler {
	return &MarkerHandler{
		MarkerFacadeService: facade,

		AuthMiddleware: authMiddleware,
	}
}

func RegisterMarkerRoutes(api fiber.Router, handler *MarkerHandler, authMiddleware *middleware.AuthMiddleware) {
	api.Get("/markers", handler.HandleGetAllMarkersLocal)
	api.Get("/markers/new", handler.HandleGetAllNewMarkers)

	api.Get("/markers/:markerId/details", authMiddleware.VerifySoft, handler.HandleGetMarker)
	api.Get("/markers/:markerID/facilities", handler.HandleGetFacilities)
	api.Get("/markers/close", handler.HandleFindCloseMarkers)
	api.Get("/markers/ranking", handler.HandleGetMarkerRanking)
	api.Get("/markers/unique-ranking", handler.HandleGetUniqueVisitorCount)
	api.Get("/markers/unique-ranking/all", handler.HandleGetAllUniqueVisitorCount)
	api.Get("/markers/area-ranking", handler.HandleGetCurrentAreaMarkerRanking)
	api.Get("/markers/convert", handler.HandleConvertWGS84ToWCONGNAMUL)
	api.Get("/markers/location-check", handler.HandleIsInSouthKorea)
	api.Get("/markers/weather", handler.HandleGetWeatherByWGS84)

	api.Get("/markers/save-offline", authMiddleware.Verify, handler.HandleSaveOfflineMap2)

	api.Post("/markers/upload", authMiddleware.CheckAdmin, handler.HandleUploadMarkerPhotoToS3)

	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(authMiddleware.Verify)

		markerGroup.Get("/my", handler.HandleGetUserMarkers)
		markerGroup.Get("/:markerID/dislike-status", handler.HandleCheckDislikeStatus)
		// markerGroup.Get("/:markerId", handlers.GetMarker)

		markerGroup.Post("/new", handler.HandleCreateMarkerWithPhotos)
		markerGroup.Post("/facilities", handler.HandleSetMarkerFacilities)
		markerGroup.Post("/:markerID/dislike", handler.HandleLeaveDislike)
		markerGroup.Post("/:markerID/favorites", handler.HandleAddFavorite)

		markerGroup.Put("/:markerID", handler.HandleUpdateMarker)

		markerGroup.Delete("/:markerID", handler.HandleDeleteMarker)
		markerGroup.Delete("/:markerID/dislike", handler.HandleUndoDislike)
		markerGroup.Delete("/:markerID/favorites", handler.HandleRemoveFavorite)
	}

}

// HandleGetAllMarkers handles the HTTP request to get all markers
func (h *MarkerHandler) HandleGetAllMarkers(c *fiber.Ctx) error {
	markers, err := h.MarkerFacadeService.GetAllMarkers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch markers",
		})
	}
	return c.JSON(markers)
}

func (h *MarkerHandler) HandleGetAllMarkersLocal(c *fiber.Ctx) error {
	// Check the Referer header and redirect if it matches the specific URL pattern
	// if !strings.HasSuffix(c.Get("Referer"), ".k-pullup.com") || c.Get("Referer") != "https://www.k-pullup.com/" {
	// 	return c.Redirect("https://k-pullup.com", fiber.StatusFound) // Use HTTP 302 for standard redirection
	// }

	cached := h.MarkerFacadeService.GetMarkerCache()

	c.Set("Content-type", "application/json")

	if cached != nil {
		// If cache is not empty, directly return the cached binary data as JSON
		c.Append("X-Cache", "hit")
		return c.Send(cached)
	}

	// Fetch markers if cache is empty
	markers, err := h.MarkerFacadeService.GetAllMarkers() // []dto.MarkerSimple, err
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Marshal the markers to JSON for caching and response
	markersJSON, err := json.Marshal(markers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode markers"})
	}

	// Update cache
	h.MarkerFacadeService.SetMarkerCache(markersJSON)

	return c.Send(markersJSON)
}

// HandleGetAllNewMarkers handles requests to fetch a paginated list of newly added markers.
func (h *MarkerHandler) HandleGetAllNewMarkers(c *fiber.Ctx) error {
	// Extract page and pageSize from query parameters. Provide default values if not specified.
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not specified
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10")) // Default to 10 markers per page if not specified
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
	}

	// Call the service to get markers
	markers, err := h.MarkerFacadeService.GetAllNewMarkers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch markers"})
	}

	return c.JSON(markers)
}

// HandleGetMarker handler
func (h *MarkerHandler) HandleGetMarker(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	userID, userOK := c.Locals("userID").(int)
	chulbong, _ := c.Locals("chulbong").(bool)

	marker, err := h.MarkerFacadeService.GetMarker(markerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marker not found: " + err.Error()})
	}

	if userOK {
		// Checking dislikes and favorites only if user is authenticated
		marker.Disliked, _ = h.MarkerFacadeService.CheckUserDislike(userID, markerID)
		marker.Favorited, _ = h.MarkerFacadeService.CheckUserFavorite(userID, markerID)

		// Check ownership. If marker.UserID is nil, chulbong remains as set earlier.
		if !chulbong && marker.UserID != nil {
			marker.IsChulbong = *marker.UserID == userID
		} else {
			marker.IsChulbong = chulbong
		}
	}

	go h.MarkerFacadeService.BufferClickEvent(markerID)
	go h.MarkerFacadeService.SaveUniqueVisitor(c.Params("markerId"), c)
	return c.JSON(marker)
}

// ADMIN
func (h *MarkerHandler) HandleGetAllMarkersWithAddr(c *fiber.Ctx) error {
	markersWithPhotos, err := h.MarkerFacadeService.GetAllMarkersWithAddr()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(markersWithPhotos)
}

func (h *MarkerHandler) HandleCreateMarkerWithPhotos(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// Check if latitude and longitude are provided
	latitude, longitude, err := GetLatLngFromForm(form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse latitude/longitude"})
	}

	description := GetDescriptionFromForm(form)

	// check first
	if fErr := h.MarkerFacadeService.CheckMarkerValidity(latitude, longitude, description); fErr != nil {
		return c.Status(fErr.Code).JSON(fiber.Map{"error": fErr.Message})
	}

	// no errors
	userID := c.Locals("userID").(int)

	marker, err := h.MarkerFacadeService.CreateMarkerWithPhotos(&dto.MarkerRequest{
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
	}, userID, form)
	if err != nil {
		if strings.Contains(err.Error(), "an error during file") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "an error during file upload"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error happened, try again later"})
	}

	return c.Status(fiber.StatusCreated).JSON(marker)
}

// UpdateMarker updates an existing marker
func (h *MarkerHandler) HandleUpdateMarker(c *fiber.Ctx) error {
	markerID, _ := strconv.Atoi(c.Params("markerID"))
	description := c.FormValue("description")

	if err := h.MarkerFacadeService.UpdateMarkerDescriptionOnly(markerID, description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"description": description})
}

// DeleteMarkerHandler handles the HTTP request to delete a marker.
func (h *MarkerHandler) HandleDeleteMarker(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)
	userRole := c.Locals("role").(string)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to delete the marker, now passing userID as well
	err = h.MarkerFacadeService.DeleteMarker(userID, markerID, userRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete marker"})
	}

	go h.MarkerFacadeService.RemoveMarkerClick(markerID)
	// go services.ResetCache(services.ALL_MARKERS_KEY)

	h.MarkerFacadeService.SetMarkerCache(nil)

	go h.MarkerFacadeService.ResetRedisCache(fmt.Sprintf("facilities:%d", markerID))
	go h.MarkerFacadeService.ResetAllRedisCache(fmt.Sprintf("userMarkers:%d:page:*", userID))

	return c.SendStatus(fiber.StatusOK)
}

// UploadMarkerPhotoToS3Handler to upload a file to S3
func (h *MarkerHandler) HandleUploadMarkerPhotoToS3(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	markerIDstr, markerIDExists := form.Value["markerId"]
	if !markerIDExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	markerID, err := strconv.Atoi(markerIDstr[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	files := form.File["photos"]

	urls, err := h.MarkerFacadeService.UploadMarkerPhotoToS3(markerID, files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload photos"})
	}

	return c.JSON(fiber.Map{"urls": urls})
}

// / DISLIKE
func (h *MarkerHandler) HandleLeaveDislike(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to leave a dislike, passing userID and markerID
	err = h.MarkerFacadeService.LeaveDislike(userID, markerID)
	if err != nil {
		// Handle specific error cases here, for example, a duplicate dislike
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to leave dislike: " + err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *MarkerHandler) HandleUndoDislike(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to undo a dislike
	err = h.MarkerFacadeService.UndoDislike(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to undo dislike: " + err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *MarkerHandler) HandleGetUserMarkers(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not authenticated"})
	}

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1 // default to first page
	}
	pageSize := 4

	// Construct a unique cache key using userID and page
	cacheKey := fmt.Sprintf("userMarkers:%d:page:%d", userID, page)
	var cachedResponse fiber.Map

	// Attempt to retrieve from cache
	err = h.MarkerFacadeService.GetRedisCache(cacheKey, &cachedResponse)

	if err == nil {
		// Cache hit, return cached response
		return c.JSON(cachedResponse)
	}

	markersWithPhotos, total, err := h.MarkerFacadeService.GetAllMarkersByUserWithPagination(userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "check if you have added markers"})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Prepare the response
	response := dto.UserMarkers{
		MarkersWithPhotos: markersWithPhotos,
		CurrentPage:       page,
		TotalPages:        totalPages,
		TotalMarkers:      total,
	}

	// Cache the response for future requests
	go h.MarkerFacadeService.SetRedisCache(cacheKey, response, 30*time.Minute)

	// Return the response
	return c.JSON(response)
}

// CheckDislikeStatus handler
func (h *MarkerHandler) HandleCheckDislikeStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	disliked, err := h.MarkerFacadeService.CheckUserDislike(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking dislike status"})
	}

	return c.JSON(fiber.Map{"disliked": disliked})
}

// HandleAddFavorite adds a new favorite marker for the user.
func (h *MarkerHandler) HandleAddFavorite(c *fiber.Ctx) error {
	userData, err := h.MarkerFacadeService.GetUserFromContext(c)
	if err != nil {
		return err // fiber err
	}

	// Extracting marker ID from request parameters or body
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid marker ID",
		})
	}

	err = h.MarkerFacadeService.AddFavorite(userData.UserID, markerID)
	if err != nil {
		// Respond differently based on the type of error
		if err.Error() == "maximum number of favorites reached" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.MarkerFacadeService.ResetFavCache(userData.Username, userData.UserID)

	// Successfully added the favorite
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Favorite added successfully",
	})
}

func (h *MarkerHandler) HandleRemoveFavorite(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not found"})
	}
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username not found"})
	}

	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	err = h.MarkerFacadeService.RemoveFavorite(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go h.MarkerFacadeService.ResetFavCache(username, userID)

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content is appropriate for a DELETE success with no response body
}

// GetFacilitiesHandler handles requests to get facilities by marker ID.
func (h *MarkerHandler) HandleGetFacilities(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	// Attempt to retrieve from cache first
	var facilities []model.Facility
	cacheKey := fmt.Sprintf("facilities:%d", markerID)
	cacheErr := h.MarkerFacadeService.GetRedisCache(cacheKey, &facilities)
	if cacheErr == nil && facilities != nil {
		c.Append("X-Cache", "hit")
		// Cache hit, return cached facilities
		return c.JSON(facilities)

	}

	facilities, err = h.MarkerFacadeService.GetFacilitiesByMarkerID(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve facilities"})
	}

	// Cache the result for future requests
	go h.MarkerFacadeService.SetRedisCache(cacheKey, facilities, 1*time.Hour)

	return c.JSON(facilities)
}

func (h *MarkerHandler) HandleSetMarkerFacilities(c *fiber.Ctx) error {
	req := new(dto.FacilityRequest)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if err := h.MarkerFacadeService.SetMarkerFacilities(req.MarkerID, req.Facilities); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set facilities for marker"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateMarkersAddressesHandler handles the request to update all markers' addresses.
func (h *MarkerHandler) updateMarkersAddressesHandler(c *fiber.Ctx) error {
	updatedMarkers, err := h.MarkerFacadeService.UpdateMarkersAddresses()
	if err != nil {
		// Log the error and return a generic error message to the client
		fmt.Printf("Error updating marker addresses: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update marker addresses",
		})
	}

	return c.JSON(fiber.Map{
		"message":        "Successfully updated marker addresses",
		"updatedMarkers": updatedMarkers,
	})
}

// helpers

func GetLatLngFromForm(form *multipart.Form) (float64, float64, error) {
	latStr, latOk := form.Value["latitude"]
	longStr, longOk := form.Value["longitude"]
	if !latOk || !longOk || len(latStr[0]) == 0 || len(longStr[0]) == 0 {
		return 0, 0, fmt.Errorf("latitude and longitude are required")
	}

	latitude, err := strconv.ParseFloat(latStr[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude")
	}

	longitude, err := strconv.ParseFloat(longStr[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude")
	}

	return latitude, longitude, nil
}

func GetDescriptionFromForm(form *multipart.Form) string {
	if descValues, exists := form.Value["description"]; exists && len(descValues[0]) > 0 {
		return descValues[0]
	}
	return ""
}

func GetMarkerIDFromForm(form *multipart.Form) string {
	if descValues, exists := form.Value["markerId"]; exists && len(descValues[0]) > 0 {
		return descValues[0]
	}
	return ""
}
