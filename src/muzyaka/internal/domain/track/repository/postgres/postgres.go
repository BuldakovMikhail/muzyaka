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

func (t trackRepository) GetGenres() ([]string, error) {
	var genres []dao.Genre
	tx := t.db.Find(&genres)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table genre)")
	}

	var genresNames []string
	for _, v := range genres {
		genresNames = append(genresNames, v.Name)
	}

	return genresNames, nil
}

func (t trackRepository) GetTracksByPartName(name string, offset int, limit int) ([]*models.TrackMeta, error) {
	var tracks []*dao.TrackMeta

	tx := t.db.
		Offset(offset).
		Limit(limit).
		Where("name LIKE ?", "%"+name+"%").
		Order("name").
		Find(&tracks)
	if err := tx.Error; err != nil {
		return nil, errors.Wrap(err, "database error (table track)")
	}

	var modelTracks []*models.TrackMeta

	for _, v := range tracks {
		var genre dao.Genre
		tx = t.db.Where("id = ?", v.GenreRefer).Limit(1).Find(&genre)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table track)")
		}
		modelTracks = append(modelTracks, dao.ToModelTrack(v, &genre))
	}

	return modelTracks, nil
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
