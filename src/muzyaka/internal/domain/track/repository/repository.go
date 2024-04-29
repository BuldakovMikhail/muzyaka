package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type TrackRepository interface {
	GetTrack(id uint64) (*models.TrackMeta, error)
	UpdateTrackOutbox(track *models.TrackMeta) error
}
