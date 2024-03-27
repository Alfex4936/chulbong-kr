package services

import (
	"chulbong-kr/database"
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

func CronCleanUpPasswordTokens() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		fmt.Println("Running cleanup job...")
		if err := DeleteExpiredPasswordTokens(); err != nil {
			// Log the error
			fmt.Printf("Error deleting expired tokens: %v\n", err)
		} else {
			fmt.Println("Expired tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		fmt.Printf("Error scheduling the token cleanup job: %v\n", err)
		return
	}
	c.Start()

	// Optionally, keep the scheduler running indefinitely
	// select {}
}

func CronResetClickRanking() {
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() { // 2 AM every day
		if err := ResetCache("marker_clicks"); err != nil {
			// Log the error
			fmt.Printf("Error reseting marker clicks: %v\n", err)
		} else {
			fmt.Println("Marker ranking cleanup executed successfully")
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling the marker ranking cleanup job: %v\n", err)
		return
	}
	c.Start()
}

func CronCleanUpToken() {
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		if err := DeleteExpiredTokens(); err != nil {
			// Log the error
			fmt.Printf("Error deleting expired tokens: %v\n", err)
		} else {
			fmt.Println("Expired tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		fmt.Printf("Error scheduling the token cleanup job: %v\n", err)
		return
	}
	c.Start()
}

// StartOrphanedPhotosCleanupCron starts the cron job for cleaning up orphaned photos.
func StartOrphanedPhotosCleanupCron() {
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		if err := deleteOrphanedPhotos(); err != nil {
			fmt.Printf("Error cleaning up orphaned photos: %v\n", err)
		} else {
			fmt.Println("Orphaned photos cleanup executed successfully")
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling the orphaned photos cleanup job: %v\n", err)
		return
	}
	c.Start()
}

// deleteOrphanedPhotos checks for photos without a corresponding marker and deletes them.
func deleteOrphanedPhotos() error {
	// Find photos with no corresponding marker.
	orphanedPhotosQuery := `
	SELECT PhotoID, PhotoURL FROM Photos
	LEFT JOIN Markers ON Photos.MarkerID = Markers.MarkerID
	WHERE Markers.MarkerID IS NULL
	`
	rows, err := database.DB.Query(orphanedPhotosQuery)
	if err != nil {
		return fmt.Errorf("querying orphaned photos: %w", err)
	}
	defer rows.Close()

	// Prepare to delete photos from the database and S3.
	deletePhotoQuery := "DELETE FROM Photos WHERE PhotoID = ?"
	var photoIDsToDelete []int
	var photoURLsToDelete []string

	for rows.Next() {
		var photoID int
		var photoURL string
		if err := rows.Scan(&photoID, &photoURL); err != nil {
			return fmt.Errorf("scanning orphaned photos: %w", err)
		}
		photoIDsToDelete = append(photoIDsToDelete, photoID)
		photoURLsToDelete = append(photoURLsToDelete, photoURL)
	}

	// Begin a transaction for batch deletion.
	tx, err := database.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	// Delete orphaned photos from the database.
	for _, photoID := range photoIDsToDelete {
		if _, err := tx.Exec(deletePhotoQuery, photoID); err != nil {
			tx.Rollback()
			return fmt.Errorf("deleting photo ID %d: %w", photoID, err)
		}
	}

	// Commit the database transaction.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	// Delete orphaned photos from S3.
	for _, photoURL := range photoURLsToDelete {
		if err := DeleteDataFromS3(photoURL); err != nil {
			// Log the error but do not stop the process for other photos.
			fmt.Printf("failed to delete photo URL %s from S3: %v\n", photoURL, err)
		}
	}

	return nil
}
