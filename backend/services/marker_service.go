package services

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"strings"

	"chulbong-kr/database"
	"chulbong-kr/dto"
	"chulbong-kr/models"

	"github.com/jmoiron/sqlx"
)

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

	// After successfully creating the marker, process and upload the files
	// Process file uploads from the multipart form
	files := form.File["photos"]
	for _, file := range files {
		fileURL, err := UploadFileToS3(file)
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

func GetAllMarkersByUserWithPagination(userID, page, pageSize int) ([]models.MarkerWithPhotos, int, error) {
	offset := (page - 1) * pageSize

	// Query to select markers created by a specific user with LIMIT and OFFSET for pagination
	markerQuery := `
SELECT 
    M.MarkerID, 
    M.UserID, 
    ST_X(M.Location) AS Latitude, 
    ST_Y(M.Location) AS Longitude, 
    M.Description, 
    U.Username, 
    M.CreatedAt, 
    M.UpdatedAt, 
    IFNULL(D.DislikeCount, 0) AS DislikeCount
FROM 
    Markers M
INNER JOIN 
    Users U ON M.UserID = U.UserID
LEFT JOIN 
    (
        SELECT 
            MarkerID, 
            COUNT(DislikeID) AS DislikeCount
        FROM 
            MarkerDislikes
        GROUP BY 
            MarkerID
    ) D ON M.MarkerID = D.MarkerID
WHERE 
    M.UserID = ?
ORDER BY 
    M.CreatedAt DESC
LIMIT ? OFFSET ?
`

	var markersWithUsernames []dto.MarkerWithDislike
	err := database.DB.Select(&markersWithUsernames, markerQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Fetch all photos at once
	photoQuery := `SELECT * FROM Photos WHERE MarkerID IN (?)`
	query, args, err := sqlx.In(photoQuery, getMarkerIDs(markersWithUsernames))
	if err != nil {
		return nil, 0, err
	}
	var allPhotos []models.Photo
	err = database.DB.Select(&allPhotos, database.DB.Rebind(query), args...)
	if err != nil {
		return nil, 0, err
	}

	// Map photos to their markers
	photoMap := make(map[int][]models.Photo) // markerID to photos
	for _, photo := range allPhotos {
		photoMap[photo.MarkerID] = append(photoMap[photo.MarkerID], photo)
	}

	// Assemble the final structure
	markersWithPhotos := make([]models.MarkerWithPhotos, 0)
	for _, marker := range markersWithUsernames {
		markersWithPhotos = append(markersWithPhotos, models.MarkerWithPhotos{
			Marker:       marker.Marker,
			Photos:       photoMap[marker.MarkerID],
			Username:     marker.Username,
			DislikeCount: marker.DislikeCount,
		})
	}

	// Query to get the total count of markers for the user
	countQuery := `SELECT COUNT(DISTINCT Markers.MarkerID) FROM Markers WHERE Markers.UserID = ?`
	var total int
	err = database.DB.Get(&total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return markersWithPhotos, total, nil
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
  IFNULL(D.DislikeCount, 0) AS DislikeCount
FROM Markers M
LEFT JOIN Users U ON M.UserID = U.UserID
LEFT JOIN (
  SELECT
    MarkerID,
    COUNT(DislikeID) AS DislikeCount
  FROM MarkerDislikes
  GROUP BY MarkerID
) D ON M.MarkerID = D.MarkerID
WHERE M.MarkerID = ?`

	var markersWithUsernames struct {
		models.Marker
		Username     string `db:"Username"`
		DislikeCount int    `db:"DislikeCount"`
	}
	err := database.DB.Get(&markersWithUsernames, query, markerID)
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
		Marker:       markersWithUsernames.Marker,
		Photos:       photoMap[markersWithUsernames.MarkerID],
		Username:     markersWithUsernames.Username,
		DislikeCount: markersWithUsernames.DislikeCount,
	}

	PublishMarkerUpdate(fmt.Sprintf("user: %s", markersWithPhotos.Username))

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

// LeaveDislike user's dislike for a marker
func LeaveDislike(userID int, markerID int) error {
	_, err := database.DB.Exec(
		"INSERT INTO MarkerDislikes (MarkerID, UserID) VALUES (?, ?) ON DUPLICATE KEY UPDATE DislikedAt=VALUES(DislikedAt)",
		markerID, userID,
	)
	if err != nil {
		return fmt.Errorf("inserting dislike: %w", err)
	}
	return nil
}

// UndoDislike nudo user's dislike for a marker
func UndoDislike(userID int, markerID int) error {
	result, err := database.DB.Exec(
		"DELETE FROM MarkerDislikes WHERE UserID = ? AND MarkerID = ?",
		userID, markerID,
	)
	if err != nil {
		return fmt.Errorf("deleting dislike: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no dislike found to undo")
	}

	return nil
}

// This service function checks if the given user has disliked the specified marker.
func CheckUserDislike(userID, markerID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM MarkerDislikes WHERE UserID = ? AND MarkerID = ?)"
	err := database.DB.Get(&exists, query, userID, markerID)
	return exists, err
}

func CheckUserFavorite(userID, markerID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Favorites WHERE UserID = ? AND MarkerID = ?)"
	err := database.DB.Get(&exists, query, userID, markerID)
	return exists, err
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

	// After successfully deleting from the database, delete photos from S3
	for _, photoURL := range photoURLs {
		if err := DeleteDataFromS3(photoURL); err != nil {
			return fmt.Errorf("deleting photo from S3: %w", err)
		}
	}

	return nil
}

// meters_per_degree = 40075000 / 360 / 1000
// IsMarkerNearby checks if there's a marker within n meters of the given latitude and longitude
func IsMarkerNearby(lat, long float64, bufferDistanceMeters int) (bool, error) {
	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	query := `
SELECT EXISTS (
    SELECT 1 
    FROM Markers
    WHERE ST_Within(Location, ST_Buffer(ST_GeomFromText(?, 4326), ?))
) AS Nearby;
`

	// Execute the query
	var nearby bool
	err := database.DB.Get(&nearby, query, point, bufferDistanceMeters)
	if err != nil {
		return false, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	return nearby, nil
}

// FindClosestNMarkersWithinDistance
func FindClosestNMarkersWithinDistance(lat, long float64, distance, pageSize, offset int) ([]dto.MarkerWithDistance, int, error) {
	point := fmt.Sprintf("POINT(%f %f)", lat, long)

	// Query to find markers within N meters and calculate total
	query := `
SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude, Description, ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) AS distance
FROM Markers
WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?
ORDER BY distance ASC
LIMIT ? OFFSET ?`

	var markers []dto.MarkerWithDistance
	err := database.DB.Select(&markers, query, point, point, distance, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking for nearby markers: %w", err)
	}

	// Query to get total count of markers within distance
	countQuery := `
SELECT COUNT(*)
FROM Markers
WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?`

	var total int
	err = database.DB.Get(&total, countQuery, point, distance)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total markers count: %w", err)
	}

	return markers, total, nil
}

// ------------ FAVORITES ------------

// AddFavoriteHandler adds a new favorite marker for the user.
func AddFavorite(userID, markerID int) error {
	// First, count the existing favorites for the user
	var count int
	err := database.DB.QueryRowx("SELECT COUNT(*) FROM Favorites WHERE UserID = ?", userID).Scan(&count)
	if err != nil {
		return err
	}

	// Check if the user already has 10 favorites
	if count >= 10 {
		return fmt.Errorf("maximum number of favorites reached")
	}

	// If not, insert the new favorite
	_, err = database.DB.Exec("INSERT INTO Favorites (UserID, MarkerID) VALUES (?, ?)", userID, markerID)
	if err != nil {
		// Convert error to string and check if it contains the MySQL error code for duplicate entry
		if strings.Contains(err.Error(), "1062") {
			return fmt.Errorf("you have already favorited this marker")
		}
		return fmt.Errorf("failed to add favorite: %w", err)
	}

	return nil
}

func RemoveFavorite(userID, markerID int) error {
	// Delete the specified favorite for the user
	_, err := database.DB.Exec("DELETE FROM Favorites WHERE UserID = ? AND MarkerID = ?", userID, markerID)
	if err != nil {
		return fmt.Errorf("error removing favorite: %w", err)
	}
	return nil
}

// Helper function to extract marker IDs
func getMarkerIDs(markers []dto.MarkerWithDislike) []interface{} {
	ids := make([]interface{}, len(markers))
	for i, marker := range markers {
		ids[i] = marker.MarkerID
	}
	return ids
}
