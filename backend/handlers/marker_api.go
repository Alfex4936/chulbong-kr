package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// QueryParamsExample handler
func QueryParamsExample(c *fiber.Ctx) error {
	// Capture query parameters
	query1 := c.Query("query")
	query2 := c.Query("query2")

	// You can also provide a default value if a query parameter is missing
	query3 := c.Query("query3", "default value")

	return c.JSON(fiber.Map{
		"query1": query1,
		"query2": query2,
		"query3": query3,
	})
}

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
	yes := services.IsInSouthKorea(latitude, longitude)
	if !yes {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Operation not allowed within South Korea."})
	}

	// Checking if a Marker is Nearby
	yes, _ = services.IsMarkerNearby(latitude, longitude)
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

func FindCloseMarkersHandler(c *fiber.Ctx) error {
	var params dto.QueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	// Find nearby markers within the specified distance
	markers, err := services.FindClosestNMarkersWithinDistance(params.Latitude, params.Longitude, params.Distance, params.N)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(markers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No markers found within the specified distance"})
	}

	// Return the found markers
	return c.JSON(markers)
}
