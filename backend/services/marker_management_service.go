package services

import (
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"
)

var CLIENT_ADDR = os.Getenv("CLIENT_ADDR")

func CreateMarkerWithPhotos(markerDto *dto.MarkerRequest, userID int, form *multipart.Form) (*dto.MarkerResponse, error) {
	// Begin a transaction for database operations
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}

	// Insert the marker into the database
	res, err := tx.Exec(
		"INSERT INTO Markers (UserID, Location, Description, CreatedAt, UpdatedAt) VALUES (?, ST_PointFromText(?, 4326), ?, NOW(), NOW())",
		userID, fmt.Sprintf("POINT(%f %f)", markerDto.Latitude, markerDto.Longitude), markerDto.Description,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	markerID, _ := res.LastInsertId()
	folder := fmt.Sprintf("markers/%d", markerID)

	// After successfully creating the marker, process and upload the files
	// Process file uploads from the multipart form
	files := form.File["photos"]
	for _, file := range files {
		fileURL, err := UploadFileToS3(folder, file)
		if err != nil {
			fmt.Printf("Failed to upload file to S3: %v\n", err)
			continue // Skip this file and continue with the next
		}
		// Associate each photo with the marker in the database
		if _, err := tx.Exec("INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt) VALUES (?, ?, NOW())", markerID, fileURL); err != nil {
			tx.Rollback()

			// Attempt to delete the uploaded file from S3
			if delErr := DeleteDataFromS3(fileURL); delErr != nil {
				fmt.Printf("Also failed to delete the file from S3: %v\n", delErr)
			}
			return nil, err
		}
	}

	// Commit the transaction after all operations succeed
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	go func(markerID int64, latitude, longitude float64) {
		maxRetries := 3
		retryDelay := 5 * time.Second // Delay between retries

		var address string
		var err error

		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				log.Printf("Retrying to fetch address for marker %d, attempt %d", markerID, attempt)
				time.Sleep(retryDelay) // Wait before retrying
			}

			address, err = FetchAddressFromAPI(latitude, longitude)
			if err == nil && address != "-1" {
				break // Success, exit the retry loop
			}

			log.Printf("Attempt %d failed to fetch address for marker %d: %v", attempt, markerID, err)
		}

		if err != nil || address == "" {
			var errMsg string
			if err != nil {
				errMsg = fmt.Sprintf("Final attempt failed to fetch address for marker %d: %v", markerID, err)
			} else {
				errMsg = fmt.Sprintf("No address found for marker %d after %d attempts", markerID, maxRetries)
			}

			url := fmt.Sprintf("%sd=%d&la=%f&lo=%f", CLIENT_ADDR, markerID, markerDto.Latitude, markerDto.Longitude)

			logFailureStmt := "INSERT INTO MarkerAddressFailures (MarkerID, ErrorMessage, URL) VALUES (?, ?, ?)"
			if _, logErr := database.DB.Exec(logFailureStmt, markerID, errMsg, url); logErr != nil {
				log.Printf("Failed to log address fetch failure for marker %d: %v", markerID, logErr)
			}
			return
		}

		// Update the marker's address in the database after successful fetch
		_, err = database.DB.Exec("UPDATE Markers SET Address = ? WHERE MarkerID = ?", address, markerID)
		if err != nil {
			log.Printf("Failed to update address for marker %d: %v", markerID, err)
		}
		if address != "" {
			updateMsg := fmt.Sprintf("새로운 철봉이 등록되었습니다! %s", address)
			PublishMarkerUpdate(updateMsg)
		}

	}(markerID, markerDto.Latitude, markerDto.Longitude)

	// Construct and return the response
	return &dto.MarkerResponse{
		MarkerID:    int(markerID),
		Latitude:    markerDto.Latitude,
		Longitude:   markerDto.Longitude,
		Description: markerDto.Description,
		// PhotoURLs:   photoURLs,
	}, nil
}

// GetAllMarkers now returns a simplified list of markers
func GetAllMarkers() ([]dto.MarkerSimple, error) {
	// Simplified query to fetch only the marker IDs, latitudes, and longitudes
	const markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude
    FROM 
        Markers;`

	var markers []dto.MarkerSimple
	err := database.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	return markers, nil
}

// GetAllMarkersWithAddr fetches all markers and returns only those with an address not found or empty.
func GetAllMarkersWithAddr() ([]dto.MarkerSimpleWithAddr, error) {
	const markerQuery = `SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude FROM Markers;`

	var markers []dto.MarkerSimpleWithAddr
	err := database.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	var filteredMarkers []dto.MarkerSimpleWithAddr
	for i := range markers {
		// func FetchAddressFromAPI(latitude float64, longitude float64) (string, error)
		address, err := FetchAddressFromAPI(markers[i].Latitude, markers[i].Longitude)
		if err != nil {
			continue
		}

		// Update the address and include only if "Address not found" or empty
		if address == "Address not found" || address == "" {
			markers[i].Address = address
			filteredMarkers = append(filteredMarkers, markers[i])
		}
	}

	return filteredMarkers, nil
}

func GetAllMarkersByUserWithPagination(userID, page, pageSize int) ([]dto.MarkerSimpleWithDescrption, int, error) {
	offset := (page - 1) * pageSize

	// Query to select markers created by a specific user with LIMIT and OFFSET for pagination
	markerQuery := `
SELECT 
    M.MarkerID,
    ST_X(M.Location) AS Latitude, 
    ST_Y(M.Location) AS Longitude, 
    M.Description,
    M.CreatedAt,
	M.Address
FROM 
    Markers M
WHERE 
    M.UserID = ?
ORDER BY 
    M.CreatedAt DESC
LIMIT ? OFFSET ?
`

	markersWithDescription := make([]dto.MarkerSimpleWithDescrption, 0)
	err := database.DB.Select(&markersWithDescription, markerQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Query to get the total count of markers for the user
	countQuery := `SELECT COUNT(DISTINCT Markers.MarkerID) FROM Markers WHERE Markers.UserID = ?`
	var total int
	err = database.DB.Get(&total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return markersWithDescription, total, nil
}

// GetMarker retrieves a single marker and its associated photo by the marker's ID
func GetMarker(markerID int) (*models.MarkerWithPhotos, error) {
	const query = `
	SELECT
	M.MarkerID,
	M.UserID,
	ST_X(M.Location) AS Latitude,
	ST_Y(M.Location) AS Longitude,
	M.Description,
	COALESCE(U.Username, '탈퇴한 사용자') AS Username,
	M.CreatedAt,
	M.UpdatedAt,
	M.Address,
	COALESCE(D.DislikeCount, 0) AS DislikeCount,
	COALESCE(F.FavoriteCount, 0) AS FavoriteCount
  FROM Markers M
  LEFT JOIN Users U ON M.UserID = U.UserID
  LEFT JOIN (
	SELECT
	  MarkerID,
	  COUNT(*) AS DislikeCount
	FROM MarkerDislikes
	WHERE MarkerID = ?
	GROUP BY MarkerID
  ) D ON M.MarkerID = D.MarkerID
  LEFT JOIN (
	SELECT
	  MarkerID,
	  COUNT(*) AS FavoriteCount
	FROM Favorites
	WHERE MarkerID = ?
	GROUP BY MarkerID
  ) F ON M.MarkerID = F.MarkerID
  WHERE M.MarkerID = ?`

	var markersWithUsernames struct {
		models.Marker
		Username      string `db:"Username"`
		DislikeCount  int    `db:"DislikeCount"`
		FavoriteCount int    `db:"FavoriteCount"`
	}
	err := database.DB.Get(&markersWithUsernames, query, markerID, markerID, markerID)
	if err != nil {
		return nil, err
	}

	// Fetch all photos at once
	const photoQuery = `SELECT * FROM Photos`
	var allPhotos []models.Photo
	err = database.DB.Select(&allPhotos, photoQuery)
	if err != nil {
		return nil, err
	}

	// Map photos to their markers
	photoMap := make(map[int][]models.Photo) // markerID to photos
	for _, photo := range allPhotos {
		photoMap[photo.MarkerID] = append(photoMap[photo.MarkerID], photo)
	}

	// Assemble the final structure
	markersWithPhotos := models.MarkerWithPhotos{
		Marker:        markersWithUsernames.Marker,
		Photos:        photoMap[markersWithUsernames.MarkerID],
		Username:      markersWithUsernames.Username,
		DislikeCount:  markersWithUsernames.DislikeCount,
		FavoriteCount: markersWithUsernames.FavoriteCount,
	}

	// PublishMarkerUpdate(fmt.Sprintf("user: %s", markersWithPhotos.Username))

	return &markersWithPhotos, nil
}

// UpdateMarker updates an existing marker's latitude, longitude, and description
func UpdateMarker(marker *models.Marker) error {
	const query = `UPDATE Markers SET Latitude = ?, Longitude = ?, Description = ?, UpdatedAt = NOW() 
                   WHERE MarkerID = ?`
	_, err := database.DB.Exec(query, marker.Latitude, marker.Longitude, marker.Description, marker.MarkerID)
	return err
}

func UpdateMarkerDescriptionOnly(markerID int, description string) error {
	const query = `UPDATE Markers SET Description = ?, UpdatedAt = NOW() 
                   WHERE MarkerID = ?`
	_, err := database.DB.Exec(query, description, markerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no marker found with markerID %d", markerID)
		}
		return fmt.Errorf("error updating a marker: %w", err)
	}

	return nil
}

// DeleteMarker deletes a marker and its associated photos from the database and S3.
func DeleteMarker(userID, markerID int) error {
	// Start a transaction
	tx, err := database.DB.Beginx()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	// Check if the marker belongs to the user
	var ownerID int
	const checkOwnerQuery = `SELECT UserID FROM Markers WHERE MarkerID = ?`
	err = tx.Get(&ownerID, checkOwnerQuery, markerID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("checking marker ownership: %w", err)
	}
	if ownerID != userID {
		tx.Rollback()
		return fmt.Errorf("user %d is not authorized to delete marker %d", userID, markerID)
	}

	// fetch photo URLs associated with the marker
	var photoURLs []string
	const selectPhotosQuery = `SELECT PhotoURL FROM Photos WHERE MarkerID = ?`
	err = tx.Select(&photoURLs, selectPhotosQuery, markerID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("fetching photos: %w", err)
	}

	// Delete associated photos from the database
	const deletePhotosQuery = `DELETE FROM Photos WHERE MarkerID = ?`
	if _, err := tx.Exec(deletePhotosQuery, markerID); err != nil {
		tx.Rollback()
		return fmt.Errorf("deleting photos: %w", err)
	}

	// Now delete the marker
	const deleteMarkerQuery = `DELETE FROM Markers WHERE MarkerID = ?`
	if _, err := tx.Exec(deleteMarkerQuery, markerID); err != nil {
		tx.Rollback()
		return fmt.Errorf("deleting marker: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	// After successfully deleting from the database, delete photos from S3 in a goroutine
	go func(photoURLs []string) {
		for _, photoURL := range photoURLs {
			// Attempt to delete the photo from S3
			if err := DeleteDataFromS3(photoURL); err != nil {
				// Log the error or handle it according to your application's requirements
				log.Printf("Failed to delete photo from S3: %s, error: %v", photoURL, err)
			}
		}
	}(photoURLs)

	return nil
}

func UploadMarkerPhotoToS3(markerID int, files []*multipart.FileHeader) ([]string, error) {
	// Begin a transaction for database operations
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	folder := fmt.Sprintf("markers/%d", markerID)

	picUrls := make([]string, 0)
	// Process file uploads from the multipart form
	for _, file := range files {
		fileURL, err := UploadFileToS3(folder, file)
		if err != nil {
			fmt.Printf("Failed to upload file to S3: %v\n", err)
			continue // Skip this file and continue with the next
		}
		picUrls = append(picUrls, fileURL)
		// Associate each photo with the marker in the database
		if _, err := tx.Exec("INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt) VALUES (?, ?, NOW())", markerID, fileURL); err != nil {
			// Attempt to delete the uploaded file from S3
			if delErr := DeleteDataFromS3(fileURL); delErr != nil {
				fmt.Printf("Also failed to delete the file from S3: %v\n", delErr)
			}
			return nil, err
		}
	}

	// Update Marker's UpdatedAt field
	if _, err := tx.Exec("UPDATE Markers SET UpdatedAt = NOW() WHERE MarkerID = ?", markerID); err != nil {
		return nil, err
	}

	// If no errors, commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return picUrls, nil
}
