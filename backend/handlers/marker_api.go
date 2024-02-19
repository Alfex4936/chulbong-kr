package handlers

import (
	"chulbong-kr/dto"
	"chulbong-kr/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetExample handler
func GetExample(c *fiber.Ctx) error {
	return c.SendString("GET request example")
}

// PutExample handler
func PutExample(c *fiber.Ctx) error {
	return c.SendString("PUT request example")
}

// DynamicRouteExample handler
func DynamicRouteExample(c *fiber.Ctx) error {
	// Capture string and id from the path
	stringParam := c.Params("string")
	idParam := c.Params("id")

	// Optionally, convert idParam to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID must be a number",
		})
	}

	return c.JSON(fiber.Map{
		"string": stringParam,
		"id":     id,
	})
}

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
	id, err := strconv.Atoi(c.Params("id"))
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
	id, _ := strconv.Atoi(c.Params("id"))
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
