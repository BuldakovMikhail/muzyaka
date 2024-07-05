package postgres

import (
	"bytes"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	repository2 "src/internal/domain/merch/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type merchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) repository2.MerchRepository {
	return &merchRepository{db: db}
}

func (m *merchRepository) GetMerchByPartName(name string, offset int, limit int) ([]*models.Merch, error) {
	var merch []*dao.Merch

	tx := m.db.
		Offset(offset).
		Limit(limit).
		Where("name LIKE ?", "%"+name+"%").
		Order("name").
		Find(&merch)
	if err := tx.Error; err != nil {
		return nil, errors.Wrap(err, "database error (table track)")
	}

	var modelMerch []*models.Merch

	for _, v := range merch {
		var photos []*dao.MerchPhotos
		tx = m.db.Where("merch_id = ?", v.ID).Limit(dao.MaxLimit).Find(&photos)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table track)")
		}
		modelMerch = append(modelMerch, dao.ToModelMerch(v, photos))
	}

	return modelMerch, nil
}

func (m *merchRepository) IsMerchOwned(merchId uint64, musicianId uint64) (bool, error) {
	var merch dao.Merch
	tx := m.db.Where("id = ?", merchId).Take(&merch)
	if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table album)")
	}

	return musicianId == merch.MusicianID, nil
}

func (m *merchRepository) GetMerch(id uint64) (*models.Merch, error) {
	var merch dao.Merch
	var merchPhotos []*dao.MerchPhotos

	tx := m.db.Where("id = ?", id).Take(&merch)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	tx = m.db.Where("merch_id = ?", id).Limit(dao.MaxLimit).Find(&merchPhotos)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	return dao.ToModelMerch(&merch, merchPhotos), nil
}

func (m *merchRepository) GetMusicianForMerch(merchId uint64) (uint64, error) {
	var merch dao.Merch

	tx := m.db.Where("id = ?", merchId).Take(&merch)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table merch)")
	}

	return merch.MusicianID, nil
}

func (m *merchRepository) GetAllMerchForMusician(musicianId uint64) ([]*models.Merch, error) {
	var merch []*dao.Merch

	tx := m.db.Where("musician_id = ?", musicianId).Limit(dao.MaxLimit).Find(&merch)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	var res []*models.Merch
	for _, v := range merch {
		var merchPhotos []*dao.MerchPhotos
		tx = m.db.Where("merch_id = ?", v.ID).Limit(dao.MaxLimit).Find(&merchPhotos)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table merch)")
		}

		res = append(res, dao.ToModelMerch(v, merchPhotos))
	}

	return res, nil
}

func (m *merchRepository) UpdateMerch(merch *models.Merch) error {
	pgMerch := dao.ToPostgresMerch(merch, 0)
	pgMerchPhotos := dao.ToPostgresMerchPhotos(merch)

	var existingFiles []*dao.MerchPhotos
	tx := m.db.Limit(dao.MaxLimit).Find(&existingFiles, "merch_id = ?", pgMerch.ID)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table merch)")
	}

	var filesToDelete []*dao.MerchPhotos
	var filesToAdd []*dao.MerchPhotos

	for _, v := range pgMerchPhotos {
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
		for _, vi := range pgMerchPhotos {
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
		if err := tx.Omit("id", "musician_id").Updates(pgMerch).Error; err != nil {
			return err
		}

		for _, v := range filesToDelete {
			if err := tx.
				Where("photo_file = ? AND merch_id = ?", v.PhotoFile, v.MerchId).
				Delete(&dao.MerchPhotos{}).Error; err != nil {
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
		return errors.Wrap(err, "database error (table merch)")
	}

	return nil
}

func (m *merchRepository) AddMerch(merch *models.Merch, musicianId uint64) (uint64, error) {
	pgMerch := dao.ToPostgresMerch(merch, musicianId)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMerch).Error; err != nil {
			return err
		}
		temp := merch
		temp.Id = pgMerch.ID
		pgMerchPhotos := dao.ToPostgresMerchPhotos(temp)
		if len(pgMerchPhotos) != 0 {
			if err := tx.Create(&pgMerchPhotos).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		merch.Id = 0
		return 0, errors.Wrap(err, "database error (table merch)")
	}

	merch.Id = pgMerch.ID
	return pgMerch.ID, nil
}

func (m *merchRepository) DeleteMerch(id uint64) error {
	tx := m.db.Delete(&dao.Merch{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table merch)")
	}

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}
