package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/domain/track/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type trackRepository struct {
	db *gorm.DB
}

func NewTrackRepository(db *gorm.DB) repository.TrackRepository {
	return &trackRepository{db: db}
}

func (t trackRepository) DeleteTrack(trackId uint64) error {

	tx := t.db.Delete(dao.TrackMeta{}, trackId)

	if err := tx.Error; err != nil {
		return errors.Wrap(err, "database error (table album)")
	}

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}

func (t trackRepository) GetTrack(id uint64) (*models.TrackMeta, error) {
	var track dao.TrackMeta

	tx := t.db.Where("id = ?", id).Take(&track)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	var genre dao.Genre
	tx = t.db.Where("id = ?", track.GenreRefer).Limit(1).Find(&genre)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	return dao.ToModelTrack(&track, &genre), nil
}

func (t trackRepository) UpdateTrack(track *models.TrackMeta) error {
	var pgGenre dao.Genre
	tx := t.db.Where("name = ?", track.Genre).Limit(1).Find(&pgGenre)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table track)")
	}

	pgTrack := dao.ToPostgresTrack(track, pgGenre.ID, 0)

	if err := t.db.Omit("id", "album_id").Updates(&pgTrack).Error; err != nil {
		return errors.Wrap(err, "database error (table track)")
	}

	return nil
}
