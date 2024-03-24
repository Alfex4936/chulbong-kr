package dto

import "time"

type CommentRequest struct {
	MarkerID    int    `json:"markerId"`
	CommentText string `json:"commentText"`
}

type CommentLoadParams struct {
	N    int `query:"n"`
	Page int `query:"page"`
}

type CommentWithUsername struct {
	CommentID   int       `json:"commentId" db:"CommentID"`
	MarkerID    int       `json:"markerId" db:"MarkerID"`
	UserID      int       `json:"userId" db:"UserID"`
	CommentText string    `json:"commentText" db:"CommentText"`
	PostedAt    time.Time `json:"postedAt" db:"PostedAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"UpdatedAt"`
	Username    string    `json:"username" db:"Username"`
}
