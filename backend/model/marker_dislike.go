package model

import "time"

// MarkerDislike is a dislike
type MarkerDislike struct {
	DislikeID  int       `db:"DislikeID"`
	MarkerID   int       `db:"MarkerID"`
	UserID     int       `db:"UserID"`
	DislikedAt time.Time `db:"DislikedAt"`
}
