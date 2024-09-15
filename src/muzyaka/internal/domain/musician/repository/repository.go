package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type MusicianRepository interface {
	GetMusician(id uint64) (*models.Musician, error)
	UpdateMusician(musician *models.Musician) error
	AddMusician(musician *models.Musician) (uint64, error)
	DeleteMusician(id uint64) error
	GetMusicianIdForUser(userId uint64) (uint64, error)
}
