package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
	"src/internal/models/dao"
)

type musicianRepository struct {
	db *gorm.DB
}

func (m musicianRepository) GetMusician(id uint64) (*models.Musician, error) {
	var musician dao.Musician
	var musicianPhotots []*dao.MusicianPhotos

	tx := m.db.Where("id = ?", id).Take(&musician)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table musician)")
	}

	tx = m.db.Where("merch_id = ?", id).Take(&musicianPhotots)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table musician)")
	}

	return dao.ToModelMusician(&musician, musicianPhotots), nil
}

func (m musicianRepository) UpdateMusician(musician *models.Musician) error {
	pgMusician := dao.ToPostgresMusician(musician)
	pgMusicianPhotos := dao.ToPostgresMusicianPhotos(musician)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := m.db.Omit("id").Updates(pgMusician).Error; err != nil {
			return err
		}
		if err := m.db.Delete(&dao.MusicianPhotos{}, "musician_id = ?", pgMusician.ID).Error; err != nil {
			return err
		}
		if err := m.db.Create(pgMusicianPhotos).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "database error (table musician)")
	}

	return nil
}

func (m musicianRepository) AddMusician(musician *models.Musician) (uint64, error) {
	pgMusician := dao.ToPostgresMusician(musician)
	pgMusicianPhotos := dao.ToPostgresMusicianPhotos(musician)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMusician).Error; err != nil {
			return err
		}

		if err := tx.Create(&pgMusicianPhotos).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "database error (table musician)")
	}
	musician.Id = pgMusician.ID
	return pgMusician.ID, nil
}

func (m musicianRepository) DeleteMusician(id uint64) error {
	tx := m.db.Delete(&dao.Musician{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table musician)")
	}

	return nil
}
