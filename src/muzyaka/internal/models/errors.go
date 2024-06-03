package models

import "errors"

var (
	ErrNotFound        = errors.New("item is not found")
	ErrAlredyExists    = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidGenre    = errors.New("invalid genre")
	ErrInvalidToken    = errors.New("invalid token")

	ErrAccessDenied   = errors.New("access denied")
	ErrInvalidContext = errors.New("error in context parsing")

	ErrNothingToDelete = errors.New("nothing to delete")
)
