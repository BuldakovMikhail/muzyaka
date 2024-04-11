package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
	"src/internal/models/dao"
)

type merchRepository struct {
	db *gorm.DB
}

func (m *merchRepository) GetMerch(id uint64) (*models.Merch, error) {
	var merch dao.Merch
	var merchPhotos []*dao.MerchPhotos

	tx := m.db.Where("id = ?", id).Take(&merch)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	tx = m.db.Where("merch_id = ?", id).Take(&merchPhotos)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table merch)")
	}

	return dao.ToModelMerch(&merch, merchPhotos), nil
}

func (m *merchRepository) UpdateMerch(merch *models.Merch) error {
	pgMerch := dao.ToPostgresMerch(merch)
	pgMerchPhotos := dao.ToPostgresMerchPhotos(merch)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := m.db.Omit("id").Updates(pgMerch).Error; err != nil {
			return err
		}
		if err := m.db.Delete(&dao.MerchPhotos{}, "merch_id = ?", pgMerch.ID).Error; err != nil {
			return err
		}
		if err := m.db.Create(pgMerchPhotos).Error; err != nil {
			return err
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

		pgMerchPhotos := dao.ToPostgresMerchPhotos(merch)
		if err := tx.Create(&pgMerchPhotos).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
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

	return nil
}
