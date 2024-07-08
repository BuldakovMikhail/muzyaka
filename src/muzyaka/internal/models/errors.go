package models

import "errors"

var (
	ErrNotFound        = errors.New("item is not found")
	ErrAlredyExists    = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidGenre    = errors.New("invalid genre")
	ErrInvalidToken    = errors.New("invalid token")
	ErrEmptyAlbum      = errors.New("album cannot be empty")

	ErrAccessDenied   = errors.New("access denied")
	ErrInvalidContext = errors.New("error in context parsing")

	ErrNothingToDelete  = errors.New("nothing to delete")
	ErrInvalidParameter = errors.New("error in query parameters")

	ErrInvalidPayload    = errors.New("error, invalid payload")
	ErrInvalidFileFormat = errors.New("error, invalid file format")
)
