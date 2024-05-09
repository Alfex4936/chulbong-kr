package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SchedulerService struct {
	DB                *sqlx.DB
	TokenService      *TokenService
	S3Service         *S3Service
	MarkerRankService *MarkerRankService
	RedisService      *RedisService
	cron              *cron.Cron
}

func NewSchedulerService(
	db *sqlx.DB, tokenService *TokenService, s3Service *S3Service, rankService *MarkerRankService, redisService *RedisService,
) *SchedulerService {
	return &SchedulerService{
		DB:                db,
		TokenService:      tokenService,
		S3Service:         s3Service,
		MarkerRankService: rankService,
		RedisService:      redisService,
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
	}
}

func RegisterSchedulerLifecycle(lifecycle fx.Lifecycle, scheduler *SchedulerService, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			scheduler.cron.Start()
			scheduler.RunAllCrons()
			logger.Info("Scheduler! words loaded successfully")
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})
}

func (s *SchedulerService) RunAllCrons() {
	s.CronCleanUpToken()
	s.CronCleanUpPasswordTokens()
	// s.CronResetClickRanking()
	s.CronOrphanedPhotosCleanup()
	s.CronCleanUpOldDirs()
	s.CronProcessClickEventsBatch(RankUpdateTime)
}

// Schedule a new job with a cron specification.
func (s *SchedulerService) Schedule(spec string, cmd func()) (cron.EntryID, error) {
	job := cron.FuncJob(cmd)
	return s.cron.AddJob(spec, job)
}

func (s *SchedulerService) CronCleanUpPasswordTokens() {
	_, err := s.Schedule("@daily", func() {
		if err := s.TokenService.DeleteExpiredPasswordTokens(); err != nil {
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

func (s *SchedulerService) CronResetClickRanking() {
	_, err := s.Schedule("0 2 * * 1", func() { // 2 AM every Monday

		// handlers.CacheMutex.Lock()
		// handlers.MarkersLocalCache = nil
		// handlers.CacheMutex.Unlock()

		if SketchedLocations != nil {
			SketchedLocations.Clear() // unique visitor 도 초기화
		}

		if err := s.RedisService.ResetCache("marker_clicks"); err != nil {
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

func (s *SchedulerService) CronProcessClickEventsBatch(interval time.Duration) {
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

	_, err := s.Schedule(spec, func() {
		s.MarkerRankService.IncrementMarkerClicks(ClickEventBuffer)
		// 처리 후 버퍼 초기화
		// clickEventBuffer.Clear()
	})
	if err != nil {
		fmt.Printf("Error setting up cron job: %v\n", err)
		return
	}
}

func (s *SchedulerService) CronCleanUpToken() {
	_, err := s.Schedule("@daily", func() {
		if err := s.TokenService.DeleteExpiredTokens(); err != nil {
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
func (s *SchedulerService) CronOrphanedPhotosCleanup() {
	_, err := s.Schedule("@daily", func() {
		if err := s.deleteOrphanedPhotos(); err != nil {
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

func (s *SchedulerService) CronNotificationCleanup() {
	// Schedules the cleanup job to run daily at midnight.
	_, err := s.Schedule("@daily", func() {
		if err := s.cleanUpViewedNotifications(); err != nil {
			fmt.Printf("Error cleaning up viewed notifcations: %v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling the notification cleanup job: %v\n", err)
		return
	}
}

// CronCleanUpOldDirs periodically checks and removes directories older than maxAge.
func (s *SchedulerService) CronCleanUpOldDirs() {
	tempDir := os.TempDir()
	maxAge := 2 * time.Minute

	_, err := s.Schedule("*/10 * * * *", func() { // every 10 minutes
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
func (s *SchedulerService) deleteOrphanedPhotos() error {
	// Find photos with no corresponding marker.
	orphanedPhotosQuery := `
	SELECT PhotoID, PhotoURL FROM Photos
	LEFT JOIN Markers ON Photos.MarkerID = Markers.MarkerID
	WHERE Markers.MarkerID IS NULL
	`
	rows, err := s.DB.Query(orphanedPhotosQuery)
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
	tx, err := s.DB.BeginTx(context.Background(), nil)
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
		if err := s.S3Service.DeleteDataFromS3(photoURL); err != nil {
			// Log the error but do not stop the process for other photos.
			fmt.Printf("failed to delete photo URL %s from S3: %v\n", photoURL, err)
		}
	}

	return nil
}

func (s *SchedulerService) cleanUpViewedNotifications() error {
	// how long to retain viewed notifications
	retentionDays := 7
	query := `DELETE FROM Notifications WHERE Viewed = TRUE AND UpdatedAt < NOW() - INTERVAL ? DAY`
	_, err := s.DB.Exec(query, retentionDays)
	if err != nil {
		return err
	} else {
		fmt.Println("Viewed notifications cleanup executed successfully")
	}

	return nil
}
