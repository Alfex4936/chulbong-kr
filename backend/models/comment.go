package models

import "time"

// Comment corresponds to the Comments table in the database
type Comment struct {
	CommentID   int        `json:"commentId" db:"CommentID"`
	MarkerID    int        `json:"markerId" db:"MarkerID"`
	UserID      int        `json:"userId" db:"UserID"`
	CommentText string     `json:"commentText" db:"CommentText"`
	PostedAt    time.Time  `json:"postedAt" db:"PostedAt"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"UpdatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"DeletedAt"`
}
