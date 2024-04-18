package services

import (
	"chulbong-kr/database"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func RunAllCrons() {
	CronCleanUpToken()
	CronCleanUpPasswordTokens()
	CronResetClickRanking()
	CronOrphanedPhotosCleanup()
	CronCleanUpOldDirs()
	CronProcessClickEventsBatch(RANK_UPDATE_TIME)
}

// CronService holds a reference to a cron scheduler and its related setup.
type CronService struct {
	cron *cron.Cron
}

var (
	instance *CronService
	once     sync.Once
)

// GetCronService returns a singleton instance of CronService, creating it if not already present.
func GetCronService() *CronService {
	once.Do(func() {
		c := cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		))
		instance = &CronService{
			cron: c,
		}
		instance.cron.Start()
	})
	return instance
}

// Schedule a new job with a cron specification.
func (cs *CronService) Schedule(spec string, cmd func()) (cron.EntryID, error) {
	job := cron.FuncJob(cmd)
	return cs.cron.AddJob(spec, job)
}

func CronCleanUpPasswordTokens() {
	c := GetCronService()
	_, err := c.Schedule("@daily", func() {
		if err := DeleteExpiredPasswordTokens(); err != nil {
			// Log the error
			fmt.Printf("Error deleting expired tokens: %v\n", err)
		} else {
			fmt.Println("Expired password tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		fmt.Printf("Error scheduling the token cleanup job: %v\n", err)
		return
	}

}

func CronResetClickRanking() {
	c := GetCronService()
	_, err := c.Schedule("0 2 * * 1", func() { // 2 AM every Monday

		// handlers.CacheMutex.Lock()
		// handlers.MarkersLocalCache = nil
		// handlers.CacheMutex.Unlock()

		if SketchedLocations != nil {
			SketchedLocations.Clear() // unique visitor 도 초기화
		}

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
}

func CronProcessClickEventsBatch(interval time.Duration) {
	c := GetCronService()
	var spec string

	switch {
	case interval < time.Hour:
		// Minutes interval
		minutes := int(interval.Minutes())
		spec = fmt.Sprintf("*/%d * * * *", minutes)
	case interval >= time.Hour && interval < 24*time.Hour:
		// Hourly interval
		hours := int(interval.Hours())
		spec = fmt.Sprintf("0 */%d * * *", hours)
	default:
		// Default to every 10 minutes if the interval is oddly long or unspecified
		spec = "*/10 * * * *"
	}
	// spec = "*/1 * * * *"

	_, err := c.Schedule(spec, func() {
		IncrementMarkerClicks(clickEventBuffer)
		// 처리 후 버퍼 초기화
		// clickEventBuffer.Clear()
	})
	if err != nil {
		fmt.Printf("Error setting up cron job: %v\n", err)
		return
	}
}

func CronCleanUpToken() {
	c := GetCronService()
	_, err := c.Schedule("@daily", func() {
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
}

// CronOrphanedPhotosCleanup starts the cron job for cleaning up orphaned photos.
func CronOrphanedPhotosCleanup() {
	c := GetCronService()
	_, err := c.Schedule("@daily", func() {
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
}

func CronNotificationCleanup() {
	c := GetCronService()
	// Schedules the cleanup job to run daily at midnight.
	_, err := c.Schedule("@daily", func() {
		if err := cleanUpViewedNotifications(); err != nil {
			fmt.Printf("Error cleaning up viewed notifcations: %v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling the notification cleanup job: %v\n", err)
		return
	}
}

// CronCleanUpOldDirs periodically checks and removes directories older than maxAge.
func CronCleanUpOldDirs() {
	c := GetCronService()
	tempDir := os.TempDir()
	maxAge := 2 * time.Minute

	_, err := c.Schedule("*/10 * * * *", func() { // every 10 minutes
		if err := cleanTempDir(tempDir, maxAge); err != nil {
			fmt.Printf("Error cleaning up orphaned photos: %v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling the orphaned photos cleanup job: %v\n", err)
		return
	}
}

// -----HELPER

// cleanTempDir removes temp directories that are older than the maxAge.
func cleanTempDir(dir string, maxAge time.Duration) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list directories in %s: %w", dir, err)
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "chulbongkr-") {
			dirPath := filepath.Join(dir, file.Name())
			fileInfo, err := file.Info()
			if err != nil {
				// log.Printf("Failed to get info for directory %s: %v", dirPath, err)
				continue
			}

			if now.Sub(fileInfo.ModTime()) > maxAge {
				os.RemoveAll(dirPath)
				//if err := os.RemoveAll(dirPath); err != nil {
				//	log.Printf("Failed to delete old directory %s: %v", dirPath, err)
				//} else {
				//	log.Printf("Deleted old directory %s", dirPath)
				//}
			}
		}
	}
	return nil
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

func cleanUpViewedNotifications() error {
	// how long to retain viewed notifications
	retentionDays := 7
	query := `DELETE FROM Notifications WHERE Viewed = TRUE AND UpdatedAt < NOW() - INTERVAL ? DAY`
	_, err := database.DB.Exec(query, retentionDays)
	if err != nil {
		return err
	} else {
		fmt.Println("Viewed notifications cleanup executed successfully")
	}

	return nil
}
