package service

import "errors"

var (
	ErrFileUpload = errors.New("an error during file upload")
	ErrNoFiles    = errors.New("upload at least one picture to prove")
	ErrNoPhotos   = errors.New("upload at least one photo")

	// Comment
	ErrMarkerNotFound     = errors.New("marker not found")
	ErrMaxCommentsReached = errors.New("user has reached the maximum number of comments")
)
