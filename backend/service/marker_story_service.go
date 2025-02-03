package service

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	insertStoryQuery = "INSERT INTO Stories (MarkerID, UserID, Caption, PhotoURL, Blurhash, Address, ExpiresAt) VALUES (?, ?, ?, ?, ?, ?, ?)"

	selectStoriesQuery = `
SELECT 
	s.StoryID, 
	s.MarkerID, 
	s.UserID, 
	s.Caption, 
	s.PhotoURL, 
	s.CreatedAt, 
	s.ExpiresAt, 
	s.Address, 
	u.Username,
	COALESCE(SUM(CASE WHEN r.ReactionType = 'thumbsup' THEN 1 ELSE 0 END), 0) AS ThumbsUp,
	COALESCE(SUM(CASE WHEN r.ReactionType = 'thumbsdown' THEN 1 ELSE 0 END), 0) AS ThumbsDown
FROM Stories s
JOIN Users u ON s.UserID = u.UserID
LEFT JOIN Reactions r ON s.StoryID = r.StoryID
WHERE s.MarkerID = ? AND s.ExpiresAt > ?
GROUP BY s.StoryID
ORDER BY s.CreatedAt DESC
    `

	selectUserIdFromStoriesQuery = "SELECT MarkerID, UserID FROM Stories WHERE StoryID = ?"

	selectPhotoFromStoriesQuery = "SELECT PhotoURL FROM Stories WHERE StoryID = ?"

	deleteStoryByIdQuery = "DELETE FROM Stories WHERE StoryID = ?"

	addReactionToStoryQuery = `
        INSERT INTO Reactions (StoryID, UserID, ReactionType)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE ReactionType = ?
    `

	deleteReactionFromStoryQuery = "DELETE FROM Reactions WHERE StoryID = ? AND UserID = ?"
	reportStoryQuery             = "INSERT INTO StoryReports (StoryID, UserID, Reason) VALUES (?, ?, ?)"

	checkExistingStoryQuery = `
        SELECT StoryID FROM Stories 
        WHERE MarkerID = ? AND UserID = ? AND ExpiresAt > ?
        LIMIT 1
    `

	selectAllStoriesQuery = `
        SELECT s.StoryID, s.MarkerID, s.UserID, s.Caption, s.PhotoURL, s.Blurhash, s.Address, s.CreatedAt, s.ExpiresAt, u.Username
        FROM Stories s
        JOIN Users u ON s.UserID = u.UserID
        WHERE s.ExpiresAt > ?
        ORDER BY s.CreatedAt DESC
        LIMIT ? OFFSET ?
    `

	getMarkerIDFromStoryIDQuery = "SELECT MarkerID FROM Stories WHERE StoryID = ?"
	checkExistingStoryIdQuery   = "SELECT EXISTS(SELECT 1 FROM Stories WHERE StoryID = ?)"

	getMarkerAddressQuery = "SELECT Address FROM Markers WHERE MarkerID = ?"
)

type StoryService struct {
	DB        *sqlx.DB
	S3Service *S3Service
	Redis     *RedisService
	Logger    *zap.Logger
}

func NewMarkerStoryService(
	db *sqlx.DB,
	s3 *S3Service,
	redis *RedisService,
	logger *zap.Logger,

) *StoryService {
	return &StoryService{
		DB:        db,
		Redis:     redis,
		S3Service: s3,
		Logger:    logger,
	}
}

func (s *StoryService) AddStory(markerID int, userID int, caption string, photo *multipart.FileHeader) (*dto.StoryResponse, error) {
	// Check Marker existence and get Address
	var address string
	err := s.DB.Get(&address, getMarkerAddressQuery, markerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("marker does not exist")
		}
		return nil, err
	}

	// Check if user already has active story
	var existingStoryID int
	err = s.DB.Get(&existingStoryID, checkExistingStoryQuery, markerID, userID, time.Now())
	if err == nil {
		return nil, ErrAlreadyStoryPost
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Handle file reading
	file, err := photo.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	// TODO: goroutine?
	blurhashString := util.EncodeBlurHashImage(img, 6, 5)

	// Upload the photo
	folder := fmt.Sprintf("stories/%d", markerID)
	photoURL, _, err := s.S3Service.UploadFileToS3(folder, photo, true)
	if err != nil {
		return nil, err
	}

	// Begin transaction and insert story
	tx, txErr := s.DB.Beginx()
	if txErr != nil {
		return nil, txErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else if commitErr := tx.Commit(); commitErr != nil {
			err = commitErr
		}
	}()

	expiresAt := time.Now().Add(24 * time.Hour)
	res, err := tx.Exec(insertStoryQuery, markerID, userID, caption, photoURL, blurhashString, address, expiresAt)
	if err != nil {
		return nil, err
	}

	storyID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var username string
	err = s.DB.Get(&username, getUsernameByIdQuery, userID)
	if err != nil {
		return nil, err
	}

	// Invalidate caches
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))
	s.Redis.ResetAllCache("stories:all:*")

	return &dto.StoryResponse{
		StoryID:   int(storyID),
		MarkerID:  markerID,
		UserID:    userID,
		Username:  username,
		Caption:   caption,
		PhotoURL:  photoURL,
		Blurhash:  &blurhashString,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Address:   address,
	}, nil
}

func (s *StoryService) GetAllStories(page int, pageSize int) ([]dto.StoryResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("stories:all:page:%d", page)
	var stories []dto.StoryResponse
	err := s.Redis.GetCacheEntry(cacheKey, &stories)
	if err == nil {
		return stories, nil
	}

	offset := (page - 1) * pageSize

	stories = []dto.StoryResponse{}
	err = s.DB.Select(&stories, selectAllStoriesQuery, time.Now(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.Redis.SetCacheEntry(cacheKey, stories, time.Minute*10) // Cache for 10 minutes

	return stories, nil
}

func (s *StoryService) GetStories(markerID int, offset int, pageSize int) ([]dto.StoryResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("stories:%d:offset:%d", markerID, offset)
	var stories []dto.StoryResponse
	err := s.Redis.GetCacheEntry(cacheKey, &stories)
	if err == nil {
		return stories, nil
	}

	// Fetch stories with pagination
	args := []interface{}{markerID, time.Now().UTC()}

	err = s.DB.Select(&stories, selectStoriesQuery, args...)
	if err != nil {
		return nil, err
	}

	// Cache the result with expiration based on the earliest ExpiresAt
	if len(stories) > 0 {
		earliestExpiresAt := stories[0].ExpiresAt
		for _, story := range stories {
			if story.ExpiresAt.Before(earliestExpiresAt) {
				earliestExpiresAt = story.ExpiresAt
			}
		}
		duration := time.Until(earliestExpiresAt)
		if duration > 0 {
			s.Redis.SetCacheEntry(cacheKey, stories, duration)
		}
	} else {
		// Cache empty result for a short duration to prevent cache stampede
		s.Redis.SetCacheEntry(cacheKey, stories, time.Minute*5)
	}

	return stories, nil
}

func (s *StoryService) DeleteStory(markerID int, storyID int, userID int, userRole string) error {
	// Begin a transaction
	tx, txErr := s.DB.Beginx()
	if txErr != nil {
		return txErr
	}
	var err error
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	// Step 1: Check if the story exists and get details
	var dbMarkerID int
	var ownerID int
	err = tx.QueryRow(selectUserIdFromStoriesQuery, storyID).Scan(&dbMarkerID, &ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStoryNotFound // Story does not exist
		}
		return err // Some other error occurred
	}

	// Step 2: Verify that the markerID matches
	if dbMarkerID != markerID {
		return ErrStoryNotFound // The story does not belong to the specified marker
	}

	// Step 3: Check user permissions
	if userRole != "admin" && ownerID != userID {
		return ErrUnauthorized
	}

	// Get the photo URL to delete from S3
	var photoURL string
	err = tx.Get(&photoURL, selectPhotoFromStoriesQuery, storyID)
	if err != nil {
		return err
	}

	// Delete the story from the database
	_, err = tx.Exec(deleteStoryByIdQuery, storyID)
	if err != nil {
		return err
	}

	// Commit the transaction before deleting from S3
	if commitErr := tx.Commit(); commitErr != nil {
		s.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
		return commitErr
	}

	// Delete the photo from S3
	if deleteErr := s.S3Service.DeleteDataFromS3(photoURL); deleteErr != nil {
		s.Logger.Error("Failed to delete photo from S3", zap.Error(deleteErr))
		// Depending on requirements, we might return an error or just log it
	}

	// Invalidate cache
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))
	s.Redis.ResetAllCache("stories:all:*")

	return nil
}

func (s *StoryService) AddReaction(storyID int, userID int, reactionType string) error {
	// Begin a transaction
	tx, txErr := s.DB.Beginx()
	if txErr != nil {
		return txErr
	}
	var err error
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	// Insert or update the reaction
	_, err = tx.Exec(addReactionToStoryQuery, storyID, userID, reactionType, reactionType)
	if err != nil {
		return err
	}

	// Fetch the markerID using storyID
	var markerID int
	err = tx.Get(&markerID, getMarkerIDFromStoryIDQuery, storyID)
	if err != nil {
		return err
	}

	// Invalidate cache for the marker's stories
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))

	return nil
}

func (s *StoryService) RemoveReaction(storyID int, userID int) error {
	// Begin a transaction
	tx, txErr := s.DB.Beginx()
	if txErr != nil {
		return txErr
	}
	var err error
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	_, err = tx.Exec(deleteReactionFromStoryQuery, storyID, userID)
	if err != nil {
		return err
	}

	// Fetch the markerID using storyID
	var markerID int
	err = tx.Get(&markerID, getMarkerIDFromStoryIDQuery, storyID)
	if err != nil {
		return err
	}

	// Invalidate cache for the marker's stories
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))

	return nil
}

func (s *StoryService) ReportStory(storyID int, userID int, reason string) error {
	// Check if the story exists
	var exists bool
	err := s.DB.Get(&exists, checkExistingStoryIdQuery, storyID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrStoryNotFound
	}

	// Insert the report
	_, err = s.DB.Exec(reportStoryQuery, storyID, userID, reason)
	if err != nil {
		// Handle duplicate report error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("you have already reported this story")
		}
		return err
	}

	return nil
}
