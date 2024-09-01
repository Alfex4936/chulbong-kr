package service

import "errors"

var (
	ErrFileUpload     = errors.New("an error during file upload")
	ErrNoFiles        = errors.New("upload at least one picture to prove")
	ErrMarkerNotExist = errors.New("check if a marker exists")
	ErrNoPhotos       = errors.New("upload at least one photo")
)
