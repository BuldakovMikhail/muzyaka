package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type MerchRepository interface {
	GetMerch(id uint64) (*models.Merch, error)
	UpdateMerch(merch *models.Merch) error
	AddMerch(merch *models.Merch) (uint64, error)
	DeleteMerch(id uint64) error
}
