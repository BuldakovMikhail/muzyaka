package postgres

import (
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
	"src/internal/models/dao"
)

type trackRepository struct {
	db *gorm.DB
}

func (t trackRepository) GetTrack(id uint64) (*models.Track, error) {
	var track dao.Track

	tx := t.db.Where("id = ?", id).Take(&track)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	var genre dao.Genre
	tx = t.db.Where("id = ?", track.GenreRefer).Take(&genre)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	return dao.ToModelTrack(&track, &genre), nil
}

func (t trackRepository) UpdateTrack(track *models.Track) error {
	var pgGenre dao.Genre
	tx := t.db.Where("name = ?", track.Genre).Take(&pgGenre)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return models.ErrInvalidGenre
	} else if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table track)")
	}

	pgTrack := dao.ToPostgresTrack(track, pgGenre.ID)

	err := t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("id").Updates(&pgTrack).Error; err != nil {
			return err
		}

		eventID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

		if err := tx.Create(&dao.Outbox{
			ID:         0,
			EventId:    eventID,
			TrackId:    pgTrack.ID,
			Source:     pgTrack.Source,
			Name:       pgTrack.Name,
			GenreRefer: pgTrack.GenreRefer,
			Type:       dao.TypeUpdate,
			Sent:       false,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "database error (table track)")
	}

	return nil
}
