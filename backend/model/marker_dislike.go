package model

import "time"

// MarkerDislike is a dislike
type MarkerDislike struct {
	DislikedAt time.Time `db:"DislikedAt"`
	DislikeID  int       `db:"DislikeID"`
	MarkerID   int       `db:"MarkerID"`
	UserID     int       `db:"UserID"`
}
