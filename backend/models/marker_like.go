package models

import "time"

// MarkerLike corresponds to the MarkerLikes table in the database
type MarkerLike struct {
	LikeID   int       `json:"likeId" db:"LikeID"`
	MarkerID int       `json:"markerId" db:"MarkerID"`
	UserID   int       `json:"userId" db:"UserID"`
	LikedAt  time.Time `json:"likedAt" db:"LikedAt"`
}
