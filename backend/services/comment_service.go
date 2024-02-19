package services

import (
	"chulbong-kr/database"
	"chulbong-kr/models"
	"context"
	"time"
)

type Comment = models.Comment

// CreateComment inserts a new comment into the database
func CreateComment(ctx context.Context, markerID, userID int, commentText string) (Comment, error) {
	comment := Comment{
		MarkerID:    markerID,
		UserID:      userID,
		CommentText: commentText,
		PostedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}
	const query = `INSERT INTO Comments (MarkerID, UserID, CommentText, PostedAt, UpdatedAt) 
				   VALUES (:MarkerID, :UserID, :CommentText, :PostedAt, :UpdatedAt) 
				   RETURNING CommentID`
	err := database.DB.QueryRowxContext(ctx, query, comment).Scan(&comment.CommentID)
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

// UpdateComment updates an existing comment
func UpdateComment(ctx context.Context, commentID, userID int, commentText string) (Comment, error) {
	comment := Comment{
		CommentID:   commentID,
		CommentText: commentText,
		UpdatedAt:   time.Now(),
	}
	const query = `UPDATE Comments 
				   SET CommentText = :CommentText, UpdatedAt = :UpdatedAt 
				   WHERE CommentID = :CommentID AND UserID = :UserID
				   RETURNING MarkerID, UserID, PostedAt`
	err := database.DB.QueryRowxContext(ctx, query, comment).Scan(&comment.MarkerID, &comment.UserID, &comment.PostedAt)
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

// DeleteComment deletes a comment from the database
func DeleteComment(ctx context.Context, commentID, userID int) error {
	const query = `DELETE FROM Comments WHERE CommentID = $1 AND UserID = $2`
	_, err := database.DB.ExecContext(ctx, query, commentID, userID)
	return err
}
