package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type TrackRepository interface {
	GetTrack(id uint64) (*models.Track, error)
	UpdateTrack(track *models.Track) error
}
