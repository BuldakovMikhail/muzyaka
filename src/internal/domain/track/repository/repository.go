package repository

import "src/internal/models"

type TrackRepository interface {
	GetTrack(id uint64) (*models.Track, error)
	UpdateTrack(track *models.Track) error
}
