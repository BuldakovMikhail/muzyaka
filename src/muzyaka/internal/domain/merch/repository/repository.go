package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type MerchRepository interface {
	GetMerch(id uint64) (*models.Merch, error)
	GetAllMerchForMusician(musicianId uint64) ([]*models.Merch, error)
	UpdateMerch(merch *models.Merch) error
	AddMerch(merch *models.Merch, musicianId uint64) (uint64, error)
	DeleteMerch(id uint64) error
	GetMusicianForMerch(merchId uint64) (uint64, error)

	IsMerchOwned(merchId uint64, musicianId uint64) (bool, error)
}
