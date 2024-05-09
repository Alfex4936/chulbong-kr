package service

import (
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/jmoiron/sqlx"
)

type ReportService struct {
	DB        *sqlx.DB
	S3Service *S3Service
}

func NewReportService(db *sqlx.DB, s3Service *S3Service) *ReportService {
	return &ReportService{
		DB:        db,
		S3Service: s3Service,
	}
}

// GetAllReports retrieves reports for all markers from the database.
func (s *ReportService) GetAllReports() ([]dto.MarkerReportResponse, error) {
	const query = `
SELECT ReportID, MarkerID, UserID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude,
       Description, ReportImageURL, CreatedAt
FROM Reports
ORDER BY CreatedAt DESC
`
	var reports []dto.MarkerReportResponse
	if err := s.DB.Select(&reports, query); err != nil {
		return nil, fmt.Errorf("error querying reports: %w", err)
	}

	return reports, nil
}

func (s *ReportService) GetAllReportsBy(markerID int) ([]dto.MarkerReportResponse, error) {
	const query = `
SELECT ReportID, MarkerID, UserID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude,
       Description, ReportImageURL, CreatedAt
FROM Reports
WHERE MarkerID = ?
ORDER BY CreatedAt DESC
`
	reports := make([]dto.MarkerReportResponse, 0)
	if err := s.DB.Select(&reports, query, markerID); err != nil {
		return nil, fmt.Errorf("error querying reports: %w", err)
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

	const query = `INSERT INTO Reports (MarkerID, UserID, Location, Description, ReportImageURL) VALUES (?, ?, ST_PointFromText(?, 4326), ?, ?)`

	// Process file uploads from the multipart form
	files := form.File["photos"]
	if len(files) == 0 {
		return fmt.Errorf("no files to process")
	}

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
			point := fmt.Sprintf("POINT(%f %f)", report.Latitude, report.Longitude)
			if _, err := tx.Exec(query, report.MarkerID, report.UserID, point, report.Description, fileURL); err != nil {
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
