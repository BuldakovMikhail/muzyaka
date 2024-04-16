package postgres

import (
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

func (m *merchRepository) GetMerch(id uint64) (*models.Merch, error) {
	var merch dao.Merch
	var merchPhotos []*dao.MerchPhotos

	tx := m.db.Where("id = ?", id).Take(&merch)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	tx = m.db.Where("merch_id = ?", id).Find(&merchPhotos)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	return dao.ToModelMerch(&merch, merchPhotos), nil
}

func (m *merchRepository) UpdateMerch(merch *models.Merch) error {
	pgMerch := dao.ToPostgresMerch(merch)
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
			if vi.PhotosSrc == v.PhotosSrc {
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
			if vi.PhotosSrc == v.PhotosSrc {
				flag = true
				break
			}
		}
		if !flag {
			filesToDelete = append(filesToDelete, v)
		}
	}

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := m.db.Omit("id").Updates(pgMerch).Error; err != nil {
			return err
		}

		for _, v := range filesToDelete {
			if err := m.db.
				Where("photo_src = ? AND merch_id = ?", v.PhotosSrc, v.MerchId).
				Delete(&dao.MerchPhotos{}).Error; err != nil {
				return err
			}
		}
		if filesToAdd != nil {
			if err := m.db.Create(&filesToAdd).Error; err != nil {
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

func (m *merchRepository) AddMerch(merch *models.Merch) (uint64, error) {
	pgMerch := dao.ToPostgresMerch(merch)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMerch).Error; err != nil {
			return err
		}
		merch.Id = pgMerch.ID
		pgMerchPhotos := dao.ToPostgresMerchPhotos(merch)

		if err := tx.Create(&pgMerchPhotos).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		merch.Id = 0
		return 0, errors.Wrap(err, "database error (table merch)")
	}

	return pgMerch.ID, nil
}

func (m *merchRepository) DeleteMerch(id uint64) error {
	tx := m.db.Delete(&dao.Merch{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table merch)")
	}

	return nil
}
