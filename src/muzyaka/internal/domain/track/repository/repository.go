package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type TrackRepository interface {
	GetTrack(id uint64) (*models.TrackMeta, error)
	UpdateTrack(track *models.TrackMeta) error
	DeleteTrack(trackId uint64) error
}
