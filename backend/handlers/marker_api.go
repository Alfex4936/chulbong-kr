package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/middlewares"
	"chulbong-kr/models"
	"chulbong-kr/protos"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/proto"
)

var (
	// cache to store encoded marker data
	MarkersLocalCache []byte // 400 kb is fine here
	CacheMutex        sync.RWMutex
)

// RegisterMarkerRoutes sets up the routes for markers handling within the application.
func RegisterMarkerRoutes(api fiber.Router) {
	// Marker routes
	// api.Get("/markers2", getAllMarkersHandler)
	// api.Get("/markers2", handlers.GetAllMarkersProtoHandler)
	api.Get("/markers", getAllMarkersLocalHandler)
	api.Get("/markers/new", getAllNewMarkersHandler)

	// api.Get("/markers-addr", middlewares.AdminOnly, handlers.GetAllMarkersWithAddrHandler)
	// api.Post("/markers-addr", middlewares.AdminOnly, handlers.UpdateMarkersAddressesHandler)
	// api.Get("/markers-db", middlewares.AdminOnly, handlers.GetMarkersClosebyAdmin)

	api.Get("/markers/:markerId/details", middlewares.AuthSoftMiddleware, getMarker)
	api.Get("/markers/:markerID/facilities", getFacilitiesHandler)
	api.Get("/markers/close", findCloseMarkersHandler)
	api.Get("/markers/ranking", getMarkerRankingHandler)
	api.Get("/markers/unique-ranking", getUniqueVisitorCountHandler)
	api.Get("/markers/unique-ranking/all", getAllUniqueVisitorCountHandler)
	api.Get("/markers/area-ranking", getCurrentAreaMarkerRankingHandler)
	api.Get("/markers/convert", convertWGS84ToWCONGNAMULHandler)
	api.Get("/markers/location-check", isInSouthKoreaHandler)
	api.Get("/markers/weather", getWeatherByWGS84Handler)

	api.Get("/markers/save-offline", middlewares.AuthMiddleware, saveOfflineMap2Handler)

	api.Post("/markers/upload", middlewares.AdminOnly, uploadMarkerPhotoToS3Handler)

	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(middlewares.AuthMiddleware)

		markerGroup.Get("/my", getUserMarkersHandler)
		markerGroup.Get("/:markerID/dislike-status", checkDislikeStatus)
		// markerGroup.Get("/:markerId", handlers.GetMarker)

		markerGroup.Post("/new", createMarkerWithPhotosHandler)
		markerGroup.Post("/facilities", setMarkerFacilitiesHandler)
		markerGroup.Post("/:markerID/dislike", leaveDislikeHandler)
		markerGroup.Post("/:markerID/favorites", addFavoriteHandler)

		markerGroup.Put("/:markerID", updateMarker)

		markerGroup.Delete("/:markerID", deleteMarkerHandler)
		markerGroup.Delete("/:markerID/dislike", undoDislikeHandler)
		markerGroup.Delete("/:markerID/favorites", removeFavoriteHandler)
	}
}

func createMarkerWithPhotosHandler(c *fiber.Ctx) error {
	// go services.ResetCache(services.ALL_MARKERS_KEY)
	CacheMutex.Lock()
	MarkersLocalCache = nil
	CacheMutex.Unlock()

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

	// Location Must Be Inside South Korea
	if !utils.IsInSouthKoreaPrecisely(latitude, longitude) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Operations are only allowed within South Korea."})
	}

	// Checking if there's a marker close to the latitude and longitude
	if nearby, _ := services.IsMarkerNearby(latitude, longitude, 10); nearby {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "There is a marker already nearby."})
	}

	// Set default description if it's empty or not provided
	description := GetDescriptionFromForm(form)
	if containsBadWord, _ := utils.CheckForBadWords(description); containsBadWord {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment contains inappropriate content."})
	}

	userId := c.Locals("userID").(int)

	marker, err := services.CreateMarkerWithPhotos(&dto.MarkerRequest{
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
	}, userId, form)
	if err != nil {
		if strings.Contains(err.Error(), "an error during file") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "an error during file upload"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error happened, try again later"})
	}

	go services.ResetAllCache(fmt.Sprintf("userMarkers:%d:page:*", userId))

	return c.Status(fiber.StatusCreated).JSON(marker)
}

func getAllMarkersHandler(c *fiber.Ctx) error {
	cachedMarkers, cacheErr := services.GetCacheEntry[[]dto.MarkerSimple](services.ALL_MARKERS_KEY)
	if cacheErr == nil && cachedMarkers != nil {
		// Cache hit
		c.Append("X-Cache", "hit")
		return c.JSON(cachedMarkers)
	}

	markers, err := services.GetAllMarkers() // []dto.MarkerSimple, err
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go services.SetCacheEntry(services.ALL_MARKERS_KEY, markers, 180*time.Minute)

	return c.JSON(markers)
}

func getAllMarkersLocalHandler(c *fiber.Ctx) error {
	CacheMutex.RLock()
	cached := MarkersLocalCache
	CacheMutex.RUnlock()

	c.Set("Content-type", "application/json")

	if cached != nil {
		// If cache is not empty, directly return the cached binary data as JSON
		c.Append("X-Cache", "hit")
		return c.Send(cached)
	}

	// Fetch markers if cache is empty
	markers, err := services.GetAllMarkers() // []dto.MarkerSimple, err
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Marshal the markers to JSON for caching and response
	markersJSON, err := json.Marshal(markers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode markers"})
	}

	// Update cache
	CacheMutex.Lock()
	MarkersLocalCache = markersJSON
	CacheMutex.Unlock()

	return c.Send(markersJSON)
}

func getAllMarkersProtoHandler(c *fiber.Ctx) error {
	markers, err := services.GetAllMarkersProto()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	markerList := &protos.MarkerList{
		Markers: markers,
	}

	data, err := proto.Marshal(markerList)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	c.Type("application/protobuf")
	return c.Send(data)
}

// GetAllNewMarkersHandler handles requests to fetch a paginated list of newly added markers.
func getAllNewMarkersHandler(c *fiber.Ctx) error {
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
	markers, err := services.GetAllNewMarkers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch markers"})
	}

	return c.JSON(markers)
}

// ADMIN
func getAllMarkersWithAddrHandler(c *fiber.Ctx) error {
	markersWithPhotos, err := services.GetAllMarkersWithAddr()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(markersWithPhotos)
}

// GetMarker handler
func getMarker(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	userID, userOK := c.Locals("userID").(int)
	chulbong, _ := c.Locals("chulbong").(bool)

	marker, err := services.GetMarker(markerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marker not found: " + err.Error()})
	}

	if userOK {
		// Checking dislikes and favorites only if user is authenticated
		marker.Disliked, _ = services.CheckUserDislike(userID, markerID)
		marker.Favorited, _ = services.CheckUserFavorite(userID, markerID)

		// Check ownership. If marker.UserID is nil, chulbong remains as set earlier.
		if !chulbong && marker.UserID != nil {
			marker.IsChulbong = *marker.UserID == userID
		} else {
			marker.IsChulbong = chulbong
		}
	}

	go services.BufferClickEvent(markerID)
	go services.SaveUniqueVisitor(c.Params("markerId"), utils.GetUserIP(c))
	return c.JSON(marker)
}

// UpdateMarker updates an existing marker
func updateMarker(c *fiber.Ctx) error {
	markerID, _ := strconv.Atoi(c.Params("markerID"))
	description := c.FormValue("description")

	if err := services.UpdateMarkerDescriptionOnly(markerID, description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"description": description})
}

// DeleteMarkerHandler handles the HTTP request to delete a marker.
func deleteMarkerHandler(c *fiber.Ctx) error {
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
	err = services.DeleteMarker(userID, markerID, userRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete marker"})
	}

	go services.RemoveMarkerClick(markerID)
	// go services.ResetCache(services.ALL_MARKERS_KEY)

	CacheMutex.Lock()
	MarkersLocalCache = nil
	CacheMutex.Unlock()

	go services.ResetCache(fmt.Sprintf("facilities:%d", markerID))
	go services.ResetAllCache(fmt.Sprintf("userMarkers:%d:page:*", userID))

	return c.SendStatus(fiber.StatusOK)
}

// UploadMarkerPhotoToS3Handler to upload a file to S3
func uploadMarkerPhotoToS3Handler(c *fiber.Ctx) error {
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

	urls, err := services.UploadMarkerPhotoToS3(markerID, files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload photos"})
	}

	return c.JSON(fiber.Map{"urls": urls})
}

// / DISLIKE
func leaveDislikeHandler(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to leave a dislike, passing userID and markerID
	err = services.LeaveDislike(userID, markerID)
	if err != nil {
		// Handle specific error cases here, for example, a duplicate dislike
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to leave dislike: " + err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func undoDislikeHandler(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to undo a dislike
	err = services.UndoDislike(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to undo dislike: " + err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func getUserMarkersHandler(c *fiber.Ctx) error {
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
	cachedResponse, _ = services.GetCacheEntry[fiber.Map](cacheKey)

	if cachedResponse != nil {
		// Cache hit, return cached response
		return c.JSON(cachedResponse)
	}

	markersWithPhotos, total, err := services.GetAllMarkersByUserWithPagination(userID, page, pageSize)
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
	go services.SetCacheEntry(cacheKey, response, 30*time.Minute)

	// Return the response
	return c.JSON(response)
}

// CheckDislikeStatus handler
func checkDislikeStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	disliked, err := services.CheckUserDislike(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking dislike status"})
	}

	return c.JSON(fiber.Map{"disliked": disliked})
}

// AddFavoriteHandler adds a new favorite marker for the user.
func addFavoriteHandler(c *fiber.Ctx) error {
	userData, err := services.GetUserFromContext(c)
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

	err = services.AddFavorite(userData.UserID, markerID)
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

	userFavKey := fmt.Sprintf("%s:%d:%s", services.USER_FAV_KEY, userData.UserID, userData.Username)
	services.ResetCache(userFavKey)

	// Successfully added the favorite
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Favorite added successfully",
	})
}

func removeFavoriteHandler(c *fiber.Ctx) error {
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

	err = services.RemoveFavorite(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	userFavKey := fmt.Sprintf("%s:%d:%s", services.USER_FAV_KEY, userID, username)
	go services.ResetCache(userFavKey)

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content is appropriate for a DELETE success with no response body
}

// GetFacilitiesHandler handles requests to get facilities by marker ID.
func getFacilitiesHandler(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	// Attempt to retrieve from cache first
	var facilities []models.Facility
	cacheKey := fmt.Sprintf("facilities:%d", markerID)
	cachedFacilities, cacheErr := services.GetCacheEntry[[]models.Facility](cacheKey)
	if cacheErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to access cache"})
	}

	if cachedFacilities != nil {
		c.Append("X-Cache", "hit")
		// Cache hit, return cached facilities
		return c.JSON(cachedFacilities)
	}

	facilities, err = services.GetFacilitiesByMarkerID(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve facilities"})
	}

	// Cache the result for future requests
	go services.SetCacheEntry(cacheKey, facilities, 1*time.Hour)

	return c.JSON(facilities)
}

func setMarkerFacilitiesHandler(c *fiber.Ctx) error {
	var req dto.FacilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if err := services.SetMarkerFacilities(req.MarkerID, req.Facilities); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set facilities for marker"})
	}

	services.ResetCache(fmt.Sprintf("facilities:%d", req.MarkerID))

	return c.SendStatus(fiber.StatusOK)
}

// UpdateMarkersAddressesHandler handles the request to update all markers' addresses.
func updateMarkersAddressesHandler(c *fiber.Ctx) error {
	updatedMarkers, err := services.UpdateMarkersAddresses()
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
