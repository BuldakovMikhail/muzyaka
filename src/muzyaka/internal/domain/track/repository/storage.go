package repository

import "src/internal/models"

//go:generate mockgen -source=storage.go -destination=mocks/storage.go

type TrackStorage interface {
	UploadObject(track *models.TrackObject) error
	LoadObject(track *models.TrackMeta) (*models.TrackObject, error)
	DeleteObject(track *models.TrackMeta) error
}
