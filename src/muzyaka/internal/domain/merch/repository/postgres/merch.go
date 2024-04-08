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
	tx := m.db.Omit("id").Updates(pgMerch)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
	}

	return nil
}

func (m *merchRepository) AddMerch(merch *models.Merch) (uint64, error) {
	pgMerch := dao.ToPostgresMerch(merch)

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMerch).Error; err != nil {
			return err
		}

		pgMerchPhotos := dao.ToPostgresMerchPhotos(pgMerch.ID, merch.Photos)
		if err := tx.Create(&pgMerchPhotos).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return pgMerch.ID, errors.Wrap(err, "database error (table album)")
	}
	merch.Id = pgMerch.ID
	return pgMerch.ID, nil
}

func (m *merchRepository) DeleteMerch(id uint64) error {
	tx := m.db.Delete(&dao.Merch{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
	}

	return nil
}

func (m *merchRepository) UpdateMerchPhotos(merch *models.Merch) error {
	pgMerchPhotos := dao.ToPostgresMerchPhotos(merch.Id, merch.Photos)

	tx := m.db.Create(&pgMerchPhotos)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
	}

	return nil
}
