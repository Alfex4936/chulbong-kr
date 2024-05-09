package service

import (
	"fmt"
	"log"
	"mime/multipart"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/jmoiron/sqlx"
)

type ReportService struct {
	DB              *sqlx.DB
	S3Service       *S3Service
	LocationService *MarkerLocationService
}

func NewReportService(db *sqlx.DB, s3Service *S3Service, location *MarkerLocationService) *ReportService {
	return &ReportService{
		DB:              db,
		S3Service:       s3Service,
		LocationService: location,
	}
}

// GetAllReports retrieves reports for all markers from the database.
func (s *ReportService) GetAllReports() ([]dto.MarkerReportResponse, error) {
	const query = `
    SELECT r.ReportID, r.MarkerID, r.UserID, ST_X(r.Location) AS Latitude, ST_Y(r.Location) AS Longitude,
    ST_X(r.NewLocation) AS NewLatitude, ST_Y(r.NewLocation) AS NewLongitude,
    r.Description, r.CreatedAt, r.Status, p.PhotoURL
    FROM Reports r
    LEFT JOIN ReportPhotos p ON r.ReportID = p.ReportID
    ORDER BY r.CreatedAt DESC
    `
	rows, err := s.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("error querying reports: %w", err)
	}
	defer rows.Close()

	reportMap := make(map[int]*dto.MarkerReportResponse)
	for rows.Next() {
		var (
			r   dto.MarkerReportResponse
			url string
		)
		if err := rows.Scan(&r.ReportID, &r.MarkerID, &r.UserID, &r.Latitude, &r.Longitude,
			&r.NewLatitude, &r.NewLongitude, &r.Description, &r.CreatedAt, &r.Status, &url); err != nil {
			return nil, err
		}
		if report, exists := reportMap[r.ReportID]; exists {
			report.PhotoURLs = append(report.PhotoURLs, url)
		} else {
			r.PhotoURLs = []string{url}
			reportMap[r.ReportID] = &r
		}
	}

	// Convert map to slice
	reports := make([]dto.MarkerReportResponse, 0, len(reportMap))
	for _, report := range reportMap {
		reports = append(reports, *report)
	}

	return reports, nil
}

func (s *ReportService) GetAllReportsBy(markerID int) ([]dto.MarkerReportResponse, error) {
	const query = `
    SELECT r.ReportID, r.MarkerID, r.UserID, ST_X(r.Location) AS Latitude, ST_Y(r.Location) AS Longitude,
    ST_X(r.NewLocation) AS NewLatitude, ST_Y(r.NewLocation) AS NewLongitude,
    r.Description, r.CreatedAt, r.Status, p.PhotoURL
    FROM Reports r
    LEFT JOIN ReportPhotos p ON r.ReportID = p.ReportID
    WHERE r.MarkerID = ?
    ORDER BY r.CreatedAt DESC
    `
	rows, err := s.DB.Queryx(query, markerID)
	if err != nil {
		return nil, fmt.Errorf("error querying reports by marker ID: %w", err)
	}
	defer rows.Close()

	reportMap := make(map[int]*dto.MarkerReportResponse)
	for rows.Next() {
		var (
			r   dto.MarkerReportResponse
			url string
		)
		if err := rows.Scan(&r.ReportID, &r.MarkerID, &r.UserID, &r.Latitude, &r.Longitude,
			&r.NewLatitude, &r.NewLongitude, &r.Description, &r.CreatedAt, &r.Status, &url); err != nil {
			return nil, err
		}
		if report, exists := reportMap[r.ReportID]; exists {
			report.PhotoURLs = append(report.PhotoURLs, url)
		} else {
			r.PhotoURLs = []string{url}
			reportMap[r.ReportID] = &r
		}
	}

	// Convert map to slice
	reports := make([]dto.MarkerReportResponse, 0, len(reportMap))
	for _, report := range reportMap {
		reports = append(reports, *report)
	}

	return reports, nil
}

// CreateReport handles the logic for creating a report and uploading photos related to that report.
func (s *ReportService) CreateReport(report *dto.MarkerReportRequest, form *multipart.Form) error {
	// Begin a transaction for database operations
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback() // Ensure the transaction is rolled back in case of error

	// Insert the main report record
	const reportQuery = `INSERT INTO Reports (MarkerID, UserID, Location, NewLocation, Description) VALUES (?, ?, ST_PointFromText(?, 4326), ST_PointFromText(?, 4326), ?)`
	res, err := tx.Exec(reportQuery, report.MarkerID, report.UserID, fmt.Sprintf("POINT(%f %f)", report.Latitude, report.Longitude), fmt.Sprintf("POINT(%f %f)", report.NewLatitude, report.NewLongitude), report.Description)
	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	// Get the last inserted ID for the report
	reportID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	// Process file uploads from the multipart form
	files := form.File["photos"]
	if len(files) == 0 {
		return fmt.Errorf("no files to process")
	}

	const photoQuery = `INSERT INTO ReportPhotos (ReportID, PhotoURL) VALUES (?, ?)`
	var wg sync.WaitGroup
	errorChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			fileURL, err := s.S3Service.UploadFileToS3("reports", file)
			if err != nil {
				errorChan <- fmt.Errorf("failed to upload file to S3: %w", err)
				return
			}
			if _, err := tx.Exec(photoQuery, reportID, fileURL); err != nil {
				errorChan <- fmt.Errorf("failed to execute database operation: %w", err)
				return
			}
		}(file)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Check for errors in the error channel
	for err := range errorChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	// Commit the transaction after all operations succeed
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (s *ReportService) ApproveReport(reportID, userID int) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Use a derived table to avoid Error 1093
	// SQL query tries to update a table (Reports) and simultaneously select from the same table within a subquery.
	const approveQuery = `
    UPDATE Reports
    SET Status = 'APPROVED'
    WHERE ReportID = ?
    AND (
        MarkerID IN (
            SELECT MarkerID
            FROM (
                SELECT MarkerID
                FROM Markers
                WHERE UserID = ? AND MarkerID IN (
                    SELECT MarkerID
                    FROM Reports
                    WHERE ReportID = ?
                )
            ) AS subquery1
        )
        OR EXISTS (
            SELECT 1
            FROM (
                SELECT UserID
                FROM Users
                WHERE UserID = ? AND Role = 'admin'
            ) AS subquery2
        )
    )
    `
	res, err := tx.Exec(approveQuery, reportID, userID, reportID, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error approving report: %w", err)
	}

	if count, _ := res.RowsAffected(); count == 0 {
		tx.Rollback()
		return fmt.Errorf("no report updated, either report does not exist or user is not the owner")
	}

	if err := s.UpdateMarkerWithReportDetailsTx(tx, reportID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// s.UpdateDbLocation(markerID)
	return nil
}

func (s *ReportService) DenyReport(reportID, userID int) error {
	const denyQuery = `
	UPDATE Reports
	SET Status = 'DENIED'
	WHERE ReportID = ?
	AND (
		MarkerID IN (
			SELECT MarkerID
			FROM Markers
			WHERE UserID = ? AND MarkerID = (
				SELECT MarkerID
				FROM Reports
				WHERE ReportID = ?
			)
		)
		OR EXISTS (
			SELECT 1
			FROM Users
			WHERE UserID = ? AND Role = 'admin'
		)
	)
	`
	res, err := s.DB.Exec(denyQuery, reportID, userID, reportID, userID)
	if err != nil {
		return fmt.Errorf("error denying report: %w", err)
	}
	// Check if the row was actually updated
	if count, _ := res.RowsAffected(); count == 0 {
		return fmt.Errorf("no report updated, either report does not exist or user is not the owner")
	}
	return nil
}

func (s *ReportService) UpdateMarkerWithReportDetailsTx(tx *sqlx.Tx, reportID int) error {
	// Update marker details based on the approved report
	// CASE update statement to conditionally update the description only if the Reports.Description is not an empty string
	const updateMarkerQuery = `
	UPDATE Markers 
	JOIN Reports ON Markers.MarkerID = Reports.MarkerID
	SET Markers.Location = COALESCE(Reports.NewLocation, Markers.Location),
		Markers.Description = CASE WHEN Reports.Description != '' THEN COALESCE(Reports.Description, Markers.Description) ELSE Markers.Description END,
		Markers.UpdatedAt = CURRENT_TIMESTAMP
	WHERE Reports.ReportID = ? AND Reports.Status = 'APPROVED'	
    `

	if _, err := tx.Exec(updateMarkerQuery, reportID); err != nil {
		return fmt.Errorf("error updating marker with report details: %w", err)
	}

	// Move approved report photos to the Photos table
	const insertPhotosQuery = `
    INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt)
    SELECT r.MarkerID, rp.PhotoURL, rp.UploadedAt
    FROM ReportPhotos rp
    JOIN Reports r ON rp.ReportID = r.ReportID
    WHERE r.ReportID = ? AND r.Status = 'APPROVED'
    `
	if _, err := tx.Exec(insertPhotosQuery, reportID); err != nil {
		return fmt.Errorf("error transferring report photos to marker photos: %w", err)
	}

	// Delete the moved photos from ReportPhotos
	const deleteReportPhotosQuery = `
    DELETE FROM ReportPhotos
    WHERE ReportID = ?
    `
	if _, err := tx.Exec(deleteReportPhotosQuery, reportID); err != nil {
		return fmt.Errorf("error deleting report photos after moving: %w", err)
	}

	return nil
}

func (s *ReportService) DeleteReport(reportID, userID, markerID int) error {
	// Start a transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if the user is authorized to delete the report
	const authQuery = `
    SELECT COUNT(*) FROM Reports
    JOIN Markers ON Reports.MarkerID = Markers.MarkerID
    WHERE Reports.ReportID = ? AND Markers.MarkerID = ? AND (Markers.UserID = ? OR EXISTS (
        SELECT 1 FROM Users WHERE UserID = ? AND Role = 'admin'
    ))
    `
	var count int
	if err := tx.Get(&count, authQuery, reportID, markerID, userID, userID); err != nil {
		return fmt.Errorf("error checking authorization: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("unauthorized to delete this report")
	}

	// Delete report photos first
	var photoURLs []string
	const selectPhotosQuery = `
    SELECT PhotoURL FROM ReportPhotos WHERE ReportID = ?
    `
	if err := tx.Select(&photoURLs, selectPhotosQuery, reportID); err != nil {
		return fmt.Errorf("error deleting report photos: %w", err)
	}

	// Delete photos from database
	const deletePhotosQuery = `DELETE FROM ReportPhotos WHERE ReportID = ?`
	if _, err := tx.Exec(deletePhotosQuery, markerID); err != nil {
		tx.Rollback()
		return fmt.Errorf("deleting photos: %w", err)
	}

	// Then delete the report
	const deleteReportQuery = `
    DELETE FROM Reports WHERE ReportID = ?
    `
	if _, err := tx.Exec(deleteReportQuery, reportID); err != nil {
		return fmt.Errorf("error deleting report: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	// Delete photos from S3 in a goroutine
	go func(photoURLs []string) {
		for _, photoURL := range photoURLs {
			if err := s.S3Service.DeleteDataFromS3(photoURL); err != nil {
				log.Printf("Failed to delete photo from S3: %s, error: %v", photoURL, err)
			}
		}
	}(photoURLs)

	return nil
}

func (s *ReportService) UpdateDbLocation(markerID int64, latitude, longitude float64) {
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

			address, err = s.LocationService.FacilityService.FetchAddressFromMap(latitude, longitude)
			if err == nil && address != "" {
				break // Success, exit the retry loop
			}

			log.Printf("Attempt %d failed to fetch address for marker %d: %v", attempt, markerID, err)
		}

		if err != nil || address == "" {
			address2, _ := s.LocationService.FacilityService.FetchRegionFromAPI(latitude, longitude)
			if address2 == "-2" {
				// delete the marker 북한 or 일본
				_, err = s.LocationService.DB.Exec("DELETE FROM Markers WHERE MarkerID = ?", markerID)
				if err != nil {
					log.Printf("Failed to update address for marker %d: %v", markerID, err)
				}
				return // no need to insert in failures
			}

			var errMsg string
			if err != nil {
				errMsg = fmt.Sprintf("Final attempt failed to fetch address for marker %d: %v", markerID, err)
			} else {
				errMsg = fmt.Sprintf("No address found for marker %d after %d attempts", markerID, maxRetries)
			}

			water, _ := s.LocationService.FacilityService.FetchRegionWaterInfo(latitude, longitude)
			if water {
				errMsg = fmt.Sprintf("The marker (%d) might be above on water", markerID)
			}

			url := fmt.Sprintf("%sd=%d&la=%f&lo=%f", s.LocationService.Config.ClientURL, markerID, latitude, longitude)

			logFailureStmt := "INSERT INTO MarkerAddressFailures (MarkerID, ErrorMessage, URL) VALUES (?, ?, ?)"
			if _, logErr := s.LocationService.DB.Exec(logFailureStmt, markerID, errMsg, url); logErr != nil {
				log.Printf("Failed to log address fetch failure for marker %d: %v", markerID, logErr)
			}
			return
		}

		// Update the marker's address in the database after successful fetch
		_, err = s.LocationService.DB.Exec("UPDATE Markers SET Address = ? WHERE MarkerID = ?", address, markerID)
		if err != nil {
			log.Printf("Failed to update address for marker %d: %v", markerID, err)
		}

		// userIDstr := strconv.Itoa(userID)
		// updateMsg := fmt.Sprintf("새로운 철봉이 [ %s ]에 등록되었습니다!", address)
		// metadata := notification.NotificationMarkerMetadata{
		// 	MarkerID:  markerID,
		// 	Latitude:  latitude,
		// 	Longitude: longitude,
		// 	Address:   address,
		// }

		// rawMetadata, _ := json.Marshal(metadata)
		// PostNotification(userIDstr, "NewMarker", "k-pullup!", updateMsg, rawMetadata)

		// TODO: update when frontend updates
		// if address != "" {
		// 	updateMsg := fmt.Sprintf("새로운 철봉이 등록되었습니다! %s", address)
		// 	PublishMarkerUpdate(updateMsg)
		// }

	}(markerID, latitude, longitude)
}
