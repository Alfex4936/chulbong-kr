package service

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	insertStoryQuery = "INSERT INTO Stories (MarkerID, UserID, Caption, PhotoURL, ExpiresAt) VALUES (?, ?, ?, ?, ?)"

	selectStoriesQuery = `
        SELECT 
            s.StoryID, 
            s.MarkerID, 
            s.UserID, 
            s.Caption, 
            s.PhotoURL, 
            s.CreatedAt, 
            s.ExpiresAt, 
            u.Username,
            COALESCE(SUM(CASE WHEN r.ReactionType = 'thumbsup' THEN 1 ELSE 0 END), 0) AS ThumbsUp,
            COALESCE(SUM(CASE WHEN r.ReactionType = 'thumbsdown' THEN 1 ELSE 0 END), 0) AS ThumbsDown
        FROM Stories s
        JOIN Users u ON s.UserID = u.UserID
        LEFT JOIN Reactions r ON s.StoryID = r.StoryID
        WHERE s.MarkerID = ? AND s.ExpiresAt > ?
        GROUP BY s.StoryID
        ORDER BY s.CreatedAt DESC
        LIMIT ? OFFSET ?
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
        SELECT s.StoryID, s.MarkerID, s.UserID, s.Caption, s.PhotoURL, s.CreatedAt, s.ExpiresAt, u.Username
        FROM Stories s
        JOIN Users u ON s.UserID = u.UserID
        WHERE s.ExpiresAt > ?
        ORDER BY s.CreatedAt DESC
        LIMIT ? OFFSET ?
    `

	getMarkerIDFromStoryIDQuery = "SELECT MarkerID FROM Stories WHERE StoryID = ?"
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
	// Check if the user already has an active story for this marker
	var existingStoryID int
	err := s.DB.Get(&existingStoryID, checkExistingStoryQuery, markerID, userID, time.Now())
	if err == nil {
		// User already has an active story
		return nil, ErrAlreadyStoryPost
	} else if err != sql.ErrNoRows {
		// An unexpected error occurred
		return nil, err
	}
	// No existing story, proceed to insert a new one

	// Begin a transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Upload the photo to S3
	folder := fmt.Sprintf("stories/%d", markerID)
	photoURL, err := s.S3Service.UploadFileToS3(folder, photo)
	if err != nil {
		return nil, err
	}

	// Set the expiration time (1 day like Instagram)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert the story into the database
	res, err := tx.Exec(insertStoryQuery, markerID, userID, caption, photoURL, expiresAt)
	if err != nil {
		return nil, err
	}

	storyID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))
	s.Redis.ResetAllCache("stories:all:*")

	// Fetch username for response
	var username string
	err = s.DB.Get(&username, getUsernameByIdQuery, userID)
	if err != nil {
		return nil, err
	}

	storyResponse := &dto.StoryResponse{
		StoryID:   int(storyID),
		MarkerID:  markerID,
		UserID:    userID,
		Username:  username,
		Caption:   caption,
		PhotoURL:  photoURL,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	return storyResponse, nil
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
	s.Redis.SetCacheEntry(cacheKey, stories, time.Minute*10) // Cache for10 minutes

	return stories, nil
}

func (s *StoryService) GetStories(markerID int, page int, pageSize int) ([]dto.StoryResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("stories:%d:page:%d", markerID, page)
	var stories []dto.StoryResponse
	err := s.Redis.GetCacheEntry(cacheKey, &stories)
	if err == nil {
		return stories, nil
	}

	offset := (page - 1) * pageSize

	stories = []dto.StoryResponse{}
	// maybe time.Now().UTC()?
	err = s.DB.Select(&stories, selectStoriesQuery, markerID, time.Now(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.Redis.SetCacheEntry(cacheKey, stories, time.Hour)

	return stories, nil
}

func (s *StoryService) DeleteStory(markerID int, storyID int, userID int, userRole string) error {
	// Step 1: Check if the story exists
	var dbMarkerID int
	var ownerID int
	err := s.DB.QueryRow(selectUserIdFromStoriesQuery, storyID).Scan(&dbMarkerID, &ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrStoryNotFound // Story does not exist
		}
		return err // Some other error occurred
	}

	// Step 2: Verify that the markerID matches
	if dbMarkerID != markerID {
		return ErrStoryNotFound // The story does not belong to the specified marker
	}

	// Step 3: Check user permissions
	// Admins can delete any story, but users can only delete their own stories
	if userRole != "admin" && ownerID != userID {
		return ErrUnauthorized
	}

	// Begin a transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

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

	// Delete the photo from S3
	err = s.S3Service.DeleteDataFromS3(photoURL)
	if err != nil {
		s.Logger.Error("Failed to delete photo from S3", zap.Error(err))
	}

	// Invalidate cache
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))
	s.Redis.ResetAllCache("stories:all:*")

	return nil
}

func (s *StoryService) AddReaction(storyID int, userID int, reactionType string) error {
	// Insert or update the reaction
	_, err := s.DB.Exec(addReactionToStoryQuery, storyID, userID, reactionType, reactionType)
	if err != nil {
		return err
	}

	// Fetch the markerID using storyID
	var markerID int
	err = s.DB.Get(&markerID, getMarkerIDFromStoryIDQuery, storyID)
	if err != nil {
		return err
	}

	// Invalidate cache for the marker's stories
	s.Redis.ResetAllCache(fmt.Sprintf("stories:%d:*", markerID))

	return nil
}

func (s *StoryService) RemoveReaction(storyID int, userID int) error {
	_, err := s.DB.Exec(deleteReactionFromStoryQuery, storyID, userID)
	if err != nil {
		return err
	}

	// Fetch the markerID using storyID
	var markerID int
	err = s.DB.Get(&markerID, getMarkerIDFromStoryIDQuery, storyID)
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
	err := s.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM Stories WHERE StoryID = ?)", storyID)
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
