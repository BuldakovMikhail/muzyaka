package postgres

import (
	"bytes"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/domain/musician/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type musicianRepository struct {
	db *gorm.DB
}

func NewMusicianRepository(db *gorm.DB) repository.MusicianRepository {
	return &musicianRepository{db: db}
}

func (m musicianRepository) GetMusician(id uint64) (*models.Musician, error) {
	var musician dao.Musician
	var musicianPhotots []*dao.MusicianPhotos

	tx := m.db.Where("id = ?", id).Take(&musician)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table musician)")
	}

	tx = m.db.Where("musician_id = ?", id).Limit(dao.MaxLimit).Find(&musicianPhotots)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table musician)")
	}

	return dao.ToModelMusician(&musician, musicianPhotots), nil
}

func (m musicianRepository) GetMusicianIdForUser(userId uint64) (uint64, error) {
	var relation dao.UserMusician
	tx := m.db.Where("user_id = ?", userId).Take(&relation)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table users_musicians)")
	}

	return relation.MusicianId, nil
}

func (m musicianRepository) UpdateMusician(musician *models.Musician) error {
	pgMusician := dao.ToPostgresMusician(musician)
	pgMusicianPhotos := dao.ToPostgresMusicianPhotos(musician)

	var existingFiles []*dao.MusicianPhotos
	tx := m.db.Limit(dao.MaxLimit).Find(&existingFiles, "musician_id = ?", pgMusician.ID)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table musician)")
	}

	var filesToDelete []*dao.MusicianPhotos
	var filesToAdd []*dao.MusicianPhotos

	for _, v := range pgMusicianPhotos {
		flag := false
		for _, vi := range existingFiles {
			if bytes.Equal(vi.PhotoFile, v.PhotoFile) {
				flag = true
				break
			}
		}

		if !flag {
			filesToAdd = append(filesToAdd, v)
		}
	}

	for _, v := range existingFiles {
		flag := false
		for _, vi := range pgMusicianPhotos {
			if bytes.Equal(vi.PhotoFile, v.PhotoFile) {
				flag = true
				break
			}
		}
		if !flag {
			filesToDelete = append(filesToDelete, v)
		}
	}

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("id").Updates(pgMusician).Error; err != nil {
			return err
		}
		for _, v := range filesToDelete {
			if err := tx.
				Where("id = ?", v.ID).
				Delete(&dao.MusicianPhotos{}).Error; err != nil {
				return err
			}
		}
		if filesToAdd != nil {
			if err := tx.Create(&filesToAdd).Error; err != nil {
				return err
			}
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

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMusician).Error; err != nil {
			return err
		}

		temp := musician
		temp.Id = pgMusician.ID
		pgMusicianPhotos := dao.ToPostgresMusicianPhotos(temp)
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

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}
