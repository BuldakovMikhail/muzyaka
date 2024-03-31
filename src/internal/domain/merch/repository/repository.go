package repository

import "src/internal/models"

type MerchRepository interface {
	GetMerch(id uint64) (*models.Merch, error)
	UpdateMerch(album *models.Merch) error
	AddMerch(album *models.Merch) (uint64, error)
	DeleteMerch(id uint64) error
}
