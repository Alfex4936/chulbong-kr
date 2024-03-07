package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateMarkerWithPhotosHandler(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// Check if latitude and longitude are provided
	latValues, latExists := form.Value["latitude"]
	longValues, longExists := form.Value["longitude"]
	if !latExists || !longExists || len(latValues[0]) == 0 || len(longValues[0]) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Latitude and longitude are required"})
	}

	// Convert latitude and longitude to float64
	latitude, err := strconv.ParseFloat(latValues[0], 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}
	longitude, err := strconv.ParseFloat(longValues[0], 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	// Location Must Be Inside South Korea
	yes := utils.IsInSouthKorea(latitude, longitude)
	if !yes {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Operation not allowed within South Korea."})
	}

	// Checking if a Marker is Nearby
	yes, _ = services.IsMarkerNearby(latitude, longitude, 7)
	if yes {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "There is a marker already nearby."})
	}

	// Set default description if it's empty or not provided
	description := "설명 없음" // Default description
	if descValues, exists := form.Value["description"]; exists && len(descValues[0]) > 0 {
		description = descValues[0]
	}

	userId := c.Locals("userID").(int)
	username := c.Locals("username").(string)

	// Construct the marker object from form values
	markerDto := dto.MarkerRequest{
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
	}

	marker, err := services.CreateMarkerWithPhotos(&markerDto, userId, form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	marker.Username = username
	marker.UserID = userId

	return c.Status(fiber.StatusCreated).JSON(marker)
}

// // CreateMarker handler
// func CreateMarkerHandler(c *fiber.Ctx) error {
// 	// assert the photo has been uploaded first.
// 	var markerDto dto.MarkerRequest
// 	userId := c.Locals("userID").(int)
// 	username := c.Locals("username").(string)

// 	if err := c.BodyParser(&markerDto); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	marker, err := services.CreateMarker(&markerDto, userId)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// Start a transaction
// 	tx, err := database.DB.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// Insert photo
// 	photo := &models.Photo{
// 		MarkerID:   marker.MarkerID,
// 		PhotoURL:   markerDto.PhotoURL,
// 		UploadedAt: time.Now(),
// 	}
// 	photoQuery := `INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt) VALUES (?, ?, NOW())`
// 	_, err = tx.Exec(photoQuery, photo.MarkerID, photo.PhotoURL)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	// Commit transaction
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}

// 	// Map the models.Marker to dto.MarkerResponse
// 	response := dto.MarkerResponse{
// 		MarkerID:    marker.MarkerID,
// 		Latitude:    marker.Latitude,
// 		Longitude:   marker.Longitude,
// 		Description: marker.Description,
// 		Username:    username,
// 		PhotoURL:    markerDto.PhotoURL,
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(response)
// }

func GetAllMarkersHandler(c *fiber.Ctx) error {
	markersWithPhotos, err := services.GetAllMarkers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(markersWithPhotos)
}

// GetMarker handler
func GetMarker(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("markerID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	marker, err := services.GetMarker(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marker not found"})
	}
	return c.JSON(marker)
}

// UpdateMarker updates an existing marker
func UpdateMarker(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("markerID"))
	markerWithPhoto, _ := services.GetMarker(id)

	if err := c.BodyParser(markerWithPhoto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.UpdateMarker(&markerWithPhoto.Marker); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(markerWithPhoto)
}

// DeleteMarkerHandler handles the HTTP request to delete a marker.
func DeleteMarkerHandler(c *fiber.Ctx) error {
	// Auth
	userID := c.Locals("userID").(int)

	// Get MarkerID from the URL parameter
	markerIDParam := c.Params("markerID")
	markerID, err := strconv.Atoi(markerIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid marker ID"})
	}

	// Call the service function to delete the marker, now passing userID as well
	err = services.DeleteMarker(userID, markerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete marker: " + err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UploadMarkerPhotoToS3Handler to upload a file to S3
func UploadMarkerPhotoToS3Handler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not get uploaded file"})
	}

	fileURL, err := services.UploadFileToS3(file)
	if err != nil {
		// Interpret the error message to set the appropriate HTTP status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"url": fileURL})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

// Find Close Markers godoc
//
// @Summary		Find close markers
// @Description	This endpoint retrieves markers that are close to a specified location within a given distance.
// @Description	It requires latitude, longitude, distance, and the number of markers (N) to return.
// @Description	If no markers are found within the specified distance, it returns a "No markers found" message.
// @Description	Returns a list of markers that meet the criteria. (maximum 3km distance allowed)
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
func FindCloseMarkersHandler(c *fiber.Ctx) error {
	var params dto.QueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	if params.Distance > 3000 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Distance cannot be greater than 3,000m (3km)"})
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
