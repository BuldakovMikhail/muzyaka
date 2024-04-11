package postgres

import (
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

	var trackAuthors []*dao.Author
	tx = t.db.Model(&dao.TrackAuthor{}).
		Select("authors.id", "authors.name").
		Limit(dao.MaxLimit).
		Where("track_id = ?", id).
		Joins("left join authors on authors.id = track_authors.author_id").
		Scan(&trackAuthors)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table track)")
	}

	return dao.ToModelTrack(&track, &genre, trackAuthors), nil
}

func (t trackRepository) UpdateTrack(track *models.Track) error {
	//TODO implement me
	panic("implement me")
}
