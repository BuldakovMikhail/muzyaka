package postgres

import (
	"gorm.io/gorm"
	"src/internal/models"
)

type musicianRepository struct {
	db *gorm.DB
}

func (m musicianRepository) GetMusician(id uint64) (*models.Musician, error) {
	//TODO implement me
	panic("implement me")
}

func (m musicianRepository) UpdateMusician(album *models.Musician) error {
	//TODO implement me
	panic("implement me")
}

func (m musicianRepository) AddMusician(album *models.Musician) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m musicianRepository) DeleteMusician(id uint64) error {
	//TODO implement me
	panic("implement me")
}
