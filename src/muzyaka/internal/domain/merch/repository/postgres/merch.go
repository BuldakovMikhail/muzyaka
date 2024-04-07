package postgres

import (
	"gorm.io/gorm"
	"src/internal/models"
)

type merchRepository struct {
	db *gorm.DB
}

func (m *merchRepository) GetMerch(id uint64) (*models.Merch, error) {
	//TODO implement me
	panic("implement me")
}

func (m *merchRepository) UpdateMerch(album *models.Merch) error {
	//TODO implement me
	panic("implement me")
}

func (m *merchRepository) AddMerch(album *models.Merch) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *merchRepository) DeleteMerch(id uint64) error {
	//TODO implement me
	panic("implement me")
}
