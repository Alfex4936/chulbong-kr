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

const (
	getOrphanedPhotosQuery = `
SELECT PhotoID, PhotoURL FROM Photos
LEFT JOIN Markers ON Photos.MarkerID = Markers.MarkerID
WHERE Markers.MarkerID IS NULL`

	deletePhotoByPhotoIdQuery = "DELETE FROM Photos WHERE PhotoID = ?"

	getOrphanedPhotosByReportQuery = `
SELECT PhotoID, PhotoURL FROM ReportPhotos
LEFT JOIN Reports ON ReportPhotos.ReportID = Reports.ReportID
WHERE Reports.ReportID IS NULL`

	deleteViewedNotificationsQuery = "DELETE FROM Notifications WHERE Viewed = TRUE AND UpdatedAt < NOW() - INTERVAL ? DAY"
)

type SchedulerService struct {
	DB                  *sqlx.DB
	TokenService        *TokenService
	S3Service           *S3Service
	MarkerRankService   *MarkerRankService
	MarkerManageService *MarkerManageService
	RedisService        *RedisService
	SmtpService         *SmtpService
	ReportService       *ReportService
	BleveSearchService  *BleveSearchService
	cron                *cron.Cron
	adminEmail          string

	GetOrphanedPhotosStmt         *sqlx.Stmt
	GetOrphanedPhotosByReportStmt *sqlx.Stmt
	DeleteViewedNotificationsStmt *sqlx.Stmt
}

func NewSchedulerService(
	db *sqlx.DB, tokenService *TokenService,
	s3Service *S3Service, rankService *MarkerRankService,
	markerService *MarkerManageService, redisService *RedisService,
	smtpService *SmtpService, reportService *ReportService,
	bleveService *BleveSearchService,

) *SchedulerService {
	// Prepare query parameters
	GetOrphanedPhotosStmt, _ := db.Preparex(getOrphanedPhotosQuery)
	GetOrphanedPhotosByReportStmt, _ := db.Preparex(getOrphanedPhotosByReportQuery)
	DeleteViewedNotificationsStmt, _ := db.Preparex(deleteViewedNotificationsQuery)

	return &SchedulerService{
		DB:                  db,
		TokenService:        tokenService,
		S3Service:           s3Service,
		MarkerRankService:   rankService,
		MarkerManageService: markerService,
		RedisService:        redisService,
		SmtpService:         smtpService,
		ReportService:       reportService,
		BleveSearchService:  bleveService,
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
		adminEmail:                    os.Getenv("SMTP_PRIVATE_EMAIL"),
		GetOrphanedPhotosStmt:         GetOrphanedPhotosStmt,
		GetOrphanedPhotosByReportStmt: GetOrphanedPhotosByReportStmt,
		DeleteViewedNotificationsStmt: DeleteViewedNotificationsStmt,
	}
}

func RegisterSchedulerLifecycle(lifecycle fx.Lifecycle, scheduler *SchedulerService, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			scheduler.cron.Start()
			scheduler.RunAllCrons(logger)
			logger.Info("Scheduler! words loaded successfully")
			return nil
		},
		OnStop: func(context.Context) error {
			scheduler.GetOrphanedPhotosStmt.Close()
			scheduler.GetOrphanedPhotosByReportStmt.Close()
			scheduler.DeleteViewedNotificationsStmt.Close()
			return nil
		},
	})
}

func (s *SchedulerService) RunAllCrons(logger *zap.Logger) {
	s.CronCleanUpToken(logger)
	s.CronUpdateRSS(logger)
	s.CronCleanUpPasswordTokens(logger)
	s.CronResetClickRanking(logger)
	s.CronOrphanedPhotosCleanup(logger)
	s.CronCleanUpOldDirs(logger)
	s.CronProcessClickEventsBatch(RankUpdateTime, logger)
	s.CronSendPendingReportsEmail(logger)
	s.CronCheckMarkerIndex(logger)

	// reports, err := s.ReportService.GetPendingReports()
	// if err != nil {
	// 	// Log the error
	// 	fmt.Printf("Error fetching pending reports: %v\n", err)
	// 	return
	// }

	// if len(reports) > 0 {
	// 	if err := s.SmtpService.SendPendingReportsEmail(s.adminEmail, reports); err != nil {
	// 		// Log the error
	// 		fmt.Printf("Error sending pending reports email: %v\n", err)
	// 	} else {
	// 		fmt.Println("Pending reports email sent successfully")
	// 	}
	// }
}

// Schedule a new job with a cron specification.
func (s *SchedulerService) Schedule(spec string, cmd func()) (cron.EntryID, error) {
	job := cron.FuncJob(cmd)
	return s.cron.AddJob(spec, job)
}

func (s *SchedulerService) CronCleanUpPasswordTokens(logger *zap.Logger) {
	_, err := s.Schedule("@daily", func() {
		if err := s.TokenService.DeleteExpiredPasswordTokens(); err != nil {
			// Log the error
			logger.Error("Error deleting expired tokens", zap.Error(err))
		} else {
			logger.Info("Expired password tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		logger.Error("Error scheduling the token cleanup job", zap.Error(err))
		return
	}

}

func (s *SchedulerService) CronResetClickRanking(logger *zap.Logger) {
	_, err := s.Schedule("0 2 * * 1", func() { // 2 AM every Monday

		// handlers.CacheMutex.Lock()
		// handlers.MarkersLocalCache = nil
		// handlers.CacheMutex.Unlock()

		if SketchedLocations != nil {
			SketchedLocations.Clear() // unique visitor 도 초기화
		}

		// if err := s.RedisService.ResetCache("marker_clicks"); err != nil {
		// 	// Log the error
		// 	fmt.Printf("Error reseting marker clicks: %v\n", err)
		// } else {
		// 	fmt.Println("Marker ranking cleanup executed successfully")
		// }
	})
	if err != nil {
		logger.Error("Error scheduling the marker ranking cleanup job", zap.Error(err))
		return
	}
}

func (s *SchedulerService) CronSendPendingReportsEmail(logger *zap.Logger) {
	// Convert 12 PM KST to UTC (3 AM UTC)
	_, err := s.Schedule("0 3 * * *", func() {
		reports, err := s.ReportService.GetPendingReports()
		if err != nil {
			// Log the error
			logger.Error("Error fetching pending reports", zap.Error(err))
			return
		}

		if len(reports) > 0 {
			if err := s.SmtpService.SendPendingReportsEmail(s.adminEmail, reports); err != nil {
				// Log the error
				logger.Error("Error sending pending reports email", zap.Error(err))
			} else {
				logger.Info("Pending reports email sent successfully")
			}
		}
	})

	if err != nil {
		logger.Error("Error scheduling the pending reports email job", zap.Error(err))
		return
	}
}

func (s *SchedulerService) CronProcessClickEventsBatch(interval time.Duration, logger *zap.Logger) {
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
		logger.Error("Error setting up cron job: %v\n", zap.Error(err))
		return
	}
}

func (s *SchedulerService) CronCleanUpToken(logger *zap.Logger) {
	_, err := s.Schedule("@daily", func() {
		if err := s.TokenService.DeleteExpiredTokens(); err != nil {
			// Log the error
			logger.Error("Error deleting expired tokens", zap.Error(err))
		} else {
			logger.Info("Expired tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		logger.Error("Error scheduling the token cleanup job", zap.Error(err))
		return
	}
}

func (s *SchedulerService) CronUpdateRSS(logger *zap.Logger) {
	_, err := s.Schedule("@daily", func() {
		rss, _ := s.MarkerManageService.GenerateRSS()

		err := saveRSSToFile(rss, "marker_rss.xml")
		if err != nil {
			logger.Error("Error saving RSS to file", zap.Error(err))
		}

	})
	if err != nil {
		logger.Error("Error scheduling the token cleanup job", zap.Error(err))
		return
	}
}

// CronOrphanedPhotosCleanup starts the cron job for cleaning up orphaned photos.
func (s *SchedulerService) CronOrphanedPhotosCleanup(logger *zap.Logger) {
	_, err := s.Schedule("@daily", func() {
		if err := s.deleteOrphanedPhotos(); err != nil {
			logger.Error("Error cleaning up orphaned photos", zap.Error(err))
		} else {
			logger.Info("Orphaned photos cleanup executed successfully")
		}

		if err := s.deleteOrphanedReportPhotos(); err != nil {
			logger.Error("Error cleaning up orphaned report photos", zap.Error(err))
		} else {
			logger.Info("Orphaned report photos cleanup executed successfully")
		}
	})

	if err != nil {
		logger.Error("Error scheduling the orphaned photos cleanup job", zap.Error(err))
		return
	}
}

func (s *SchedulerService) CronNotificationCleanup(logger *zap.Logger) {
	// Schedules the cleanup job to run daily at midnight.
	_, err := s.Schedule("@daily", func() {
		if err := s.cleanUpViewedNotifications(); err != nil {
			logger.Error("Error cleaning up viewed notifcations", zap.Error(err))
		}
	})
	if err != nil {
		logger.Error("Error scheduling the notification cleanup job", zap.Error(err))
		return
	}
}

// CronCleanUpOldDirs periodically checks and removes directories older than maxAge.
func (s *SchedulerService) CronCleanUpOldDirs(logger *zap.Logger) {
	tempDir := os.TempDir()
	maxAge := 2 * time.Minute

	_, err := s.Schedule("*/15 * * * *", func() { // every 15 minutes
		if err := cleanTempDir(tempDir, maxAge); err != nil {
			logger.Error("Error cleaning up orphaned photos: %v\n", zap.Error(err))
		}
	})
	if err != nil {
		logger.Error("Error scheduling the orphaned photos cleanup job: %v\n", zap.Error(err))
		return
	}
}

// CronCheckMarkerIndex periodically checks and removes indexes of bleve.
func (s *SchedulerService) CronCheckMarkerIndex(logger *zap.Logger) {
	_, err := s.Schedule("0 17 * * *", func() { // UTC 17pm
		if err := s.BleveSearchService.CheckIndexes(); err != nil {
			logger.Error("Error cleaning up orphaned photos: %v\n", zap.Error(err))
		} else {
			s.BleveSearchService.InvalidateCache()
		}
	})
	if err != nil {
		logger.Error("Error scheduling the orphaned photos cleanup job: %v\n", zap.Error(err))
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
	rows, err := s.GetOrphanedPhotosStmt.Query()
	if err != nil {
		return fmt.Errorf("querying orphaned photos: %w", err)
	}
	defer rows.Close()

	// Prepare to delete photos from the database and S3.
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
		if _, err := tx.Exec(deletePhotoByPhotoIdQuery, photoID); err != nil {
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

func (s *SchedulerService) deleteOrphanedReportPhotos() error {
	// Find photos with no corresponding marker.
	rows, err := s.GetOrphanedPhotosByReportStmt.Query()
	if err != nil {
		return fmt.Errorf("querying orphaned photos: %w", err)
	}
	defer rows.Close()

	// Prepare to delete photos from the database and S3.
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
		if _, err := tx.Exec(deleteReportPhotosQuery, photoID); err != nil {
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
	_, err := s.DB.Exec(deleteViewedNotificationsQuery, retentionDays)
	if err != nil {
		return err
	} else {
		fmt.Println("Viewed notifications cleanup executed successfully")
	}

	return nil
}
