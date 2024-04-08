package postgres

import (
	"gorm.io/gorm"
	"src/internal/models"
)

type trackRepository struct {
	db *gorm.DB
}

func (t trackRepository) GetTrack(id uint64) (*models.Track, error) {
	//TODO implement me
	panic("implement me")
}

func (t trackRepository) UpdateTrack(track *models.Track) error {
	//TODO implement me
	panic("implement me")
}
