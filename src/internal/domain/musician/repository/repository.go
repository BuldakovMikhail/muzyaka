package repository

import "src/internal/models"

type MusicianRepository interface {
	GetMusician(id uint64) (*models.Musician, error)
	UpdateMusician(album *models.Musician) error
	AddMusician(album *models.Musician) (uint64, error)
	DeleteMusician(id uint64) error
}
