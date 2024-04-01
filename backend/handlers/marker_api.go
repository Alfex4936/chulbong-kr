package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateMarkerWithPhotosHandler(c *fiber.Ctx) error {
	go services.ResetCache(services.ALL_MARKERS_KEY)

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// Check if latitude and longitude are provided
	latitude, longitude, err := getLatLong(form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Location Must Be Inside South Korea
	if !utils.IsInSouthKorea(latitude, longitude) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Operations only allowed within South Korea."})
	}

	// Checking if there's a marker close to the latitude and longitude
	if nearby, _ := services.IsMarkerNearby(latitude, longitude, 10); nearby {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "There is a marker already nearby."})
	}

	// Set default description if it's empty or not provided
	description := getDescription(form)
	if containsBadWord, _ := utils.CheckForBadWords(description); containsBadWord {
		return c.Status(fiber.StatusBadRequest).SendString("Comment contains inappropriate content.")
	}

	userId := c.Locals("userID").(int)

	marker, err := services.CreateMarkerWithPhotos(&dto.MarkerRequest{
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
	}, userId, form)
	if err != nil {
		if strings.Contains(err.Error(), "an error during file") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(marker)
}

func GetAllMarkersHandler(c *fiber.Ctx) error {
	cachedMarkers, cacheErr := services.GetCacheEntry[[]dto.MarkerSimple](services.ALL_MARKERS_KEY)
	if cacheErr == nil && cachedMarkers != nil {
		// Cache hit
		c.Append("X-Cache", "hit")
		return c.JSON(cachedMarkers)
	}

	markersWithPhotos, err := services.GetAllMarkers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go services.SetCacheEntry(services.ALL_MARKERS_KEY, markersWithPhotos, 30*time.Minute)

	return c.JSON(markersWithPhotos)
}

// ADMIN
func GetAllMarkersWithAddrHandler(c *fiber.Ctx) error {
	markersWithPhotos, err := services.GetAllMarkersWithAddr()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(markersWithPhotos)
}

// GetMarker handler
func GetMarker(c *fiber.Ctx) error {
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
	return c.JSON(marker)
}

// UpdateMarker updates an existing marker
func UpdateMarker(c *fiber.Ctx) error {
	markerID, _ := strconv.Atoi(c.Params("markerID"))
	description := c.FormValue("description")

	if err := services.UpdateMarkerDescriptionOnly(markerID, description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"description": description})
}

// DeleteMarkerHandler handles the HTTP request to delete a marker.
func DeleteMarkerHandler(c *fiber.Ctx) error {
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
	go services.ResetCache(services.ALL_MARKERS_KEY)

	return c.SendStatus(fiber.StatusOK)
}

// UploadMarkerPhotoToS3Handler to upload a file to S3
func UploadMarkerPhotoToS3Handler(c *fiber.Ctx) error {
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

// DeleteObjectFromS3Handler handles requests to delete objects from S3.
func DeleteObjectFromS3Handler(c *fiber.Ctx) error {
	var requestBody struct {
		ObjectURL string `json:"objectUrl"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body is not valid"})
	}

	// Ensure the object URL is not empty
	if requestBody.ObjectURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Object URL is required"})
	}

	// Call the service function to delete the object from S3
	if err := services.DeleteDataFromS3(requestBody.ObjectURL); err != nil {
		// Determine if the error should be a 404 not found or a 500 internal server error
		if err.Error() == "object not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Object not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete object from S3"})
	}

	// Return a success response
	return c.SendStatus(fiber.StatusNoContent)
}

// / DISLIKE
func LeaveDislikeHandler(c *fiber.Ctx) error {
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

func UndoDislikeHandler(c *fiber.Ctx) error {
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

func GetUserMarkersHandler(c *fiber.Ctx) error {
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

	markersWithPhotos, total, err := services.GetAllMarkersByUserWithPagination(userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "check if you have added markers"})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Return the filtered markers with pagination info
	return c.JSON(fiber.Map{
		"markers":      markersWithPhotos,
		"currentPage":  page,
		"totalPages":   totalPages,
		"totalMarkers": total,
	})
}

// CheckDislikeStatus handler
func CheckDislikeStatus(c *fiber.Ctx) error {
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
func AddFavoriteHandler(c *fiber.Ctx) error {
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

func RemoveFavoriteHandler(c *fiber.Ctx) error {
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
func GetFacilitiesHandler(c *fiber.Ctx) error {
	markerID, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Marker ID"})
	}

	facilities, err := services.GetFacilitiesByMarkerID(markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve facilities"})
	}

	return c.JSON(facilities)
}

func SetMarkerFacilitiesHandler(c *fiber.Ctx) error {
	var req dto.FacilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if err := services.SetMarkerFacilities(req.MarkerID, req.Facilities); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set facilities for marker"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateMarkersAddressesHandler handles the request to update all markers' addresses.
func UpdateMarkersAddressesHandler(c *fiber.Ctx) error {
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

func getLatLong(form *multipart.Form) (float64, float64, error) {
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

func getDescription(form *multipart.Form) string {
	if descValues, exists := form.Value["description"]; exists && len(descValues[0]) > 0 {
		return descValues[0]
	}
	return ""
}
