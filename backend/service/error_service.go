package service

import "errors"

var (
	ErrFileUpload = errors.New("an error during file upload")
	ErrNoFiles    = errors.New("no files provided")
	ErrNoPhotos   = errors.New("upload at least one photo")

	// Comment
	ErrMarkerNotFound     = errors.New("marker not found")
	ErrMaxCommentsReached = errors.New("user has reached the maximum number of comments")

	// Report
	ErrBeginTransaction   = errors.New("could not begin transaction")
	ErrInsertReport       = errors.New("failed to insert report")
	ErrLastInsertID       = errors.New("failed to get last insert ID")
	ErrInsertReportPhoto  = errors.New("failed to insert report photo")
	ErrCommitTransaction  = errors.New("could not commit transaction")
	ErrMarkerDoesNotExist = errors.New("marker does not exist")
)
