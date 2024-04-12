package postgres

import (
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/domain/album/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type albumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) repository.AlbumRepository {
	return &albumRepository{db: db}
}

func (ar *albumRepository) GetAlbum(id uint64) (*models.Album, error) {
	var album dao.Album

	tx := ar.db.Where("id = ?", id).Take(&album)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table album)")
	}
	return dao.ToModelAlbum(&album), nil
}

func (ar *albumRepository) UpdateAlbum(album *models.Album) error {
	pgAlbum := dao.ToPostgresAlbum(album)
	tx := ar.db.Omit("id").Updates(pgAlbum)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
	}

	return nil
}

func (ar *albumRepository) AddAlbumWithTracks(album *models.Album, tracks []*models.Track) (uint64, error) {
	pgAlbum := dao.ToPostgresAlbum(album)
	var pgTracks []*dao.Track

	for _, v := range tracks {
		var pgGenre dao.Genre
		tx := ar.db.Where("name = ?", v.Genre).Take(&pgGenre)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return 0, models.ErrInvalidGenre
		} else if tx.Error != nil {
			return 0, errors.Wrap(tx.Error, "database error (table album)")
		}

		pgTracks = append(pgTracks, &dao.Track{
			ID:         0,
			Source:     v.Source,
			Name:       v.Name,
			GenreRefer: pgGenre.ID,
		})
	}
	err := ar.db.Transaction(func(tx *gorm.DB) error {
		if err := ar.db.Create(pgAlbum).Error; err != nil {
			return err
		}

		if err := ar.db.Create(&pgTracks).Error; err != nil {
			return err
		}

		for _, v := range pgTracks {
			if err := ar.db.Create(&dao.AlbumTrack{
				AlbumId: pgAlbum.ID,
				TrackId: v.ID,
			}).Error; err != nil {
				return err
			}

			eventID, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}

			if err := ar.db.Create(&dao.Outbox{
				ID:         0,
				EventId:    eventID,
				TrackId:    v.ID,
				Source:     v.Source,
				Name:       v.Name,
				GenreRefer: v.GenreRefer,
				Type:       dao.TypeAdd,
				Sent:       false,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "database error (table album)")
	}
	album.Id = pgAlbum.ID
	return pgAlbum.ID, nil
}

func (ar *albumRepository) DeleteAlbum(id uint64) error {

	err := ar.db.Transaction(func(tx *gorm.DB) error {
		var relations []*dao.AlbumTrack
		if err := ar.db.Limit(dao.MaxLimit).Find(&relations, "album_id = ?", id).Error; err != nil {
			return err
		}

		for _, v := range relations {
			if err := ar.db.Delete(&dao.Track{}, v.TrackId).Error; err != nil {
				return err
			}

			eventID, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}

			if err := ar.db.Create(&dao.Outbox{
				ID:      0,
				EventId: eventID,
				TrackId: v.TrackId,
				Type:    dao.TypeDelete,
				Sent:    false,
			}).Error; err != nil {
				return err
			}
		}

		if err := ar.db.Delete(&dao.Album{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "database error (table album)")
	}

	return nil
}

func (ar *albumRepository) AddTrackToAlbum(albumId uint64, track *models.Track) (uint64, error) {
	var pgGenre dao.Genre
	tx := ar.db.Where("name = ?", track.Genre).Take(&pgGenre)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return 0, models.ErrInvalidGenre
	} else if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table album)")
	}

	pgTrack := dao.ToPostgresTrack(track, pgGenre.ID)

	err := ar.db.Transaction(func(tx *gorm.DB) error {
		// Add track to tracks table
		if err := tx.Create(&pgTrack).Error; err != nil {
			return err
		}
		// Set album-track relation
		if err := tx.Create(&dao.AlbumTrack{AlbumId: albumId, TrackId: pgTrack.ID}).Error; err != nil {
			return err
		}

		// Add event to outbox
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
			Type:       dao.TypeAdd,
			Sent:       false,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "database error (table album)")
	}

	return pgTrack.ID, nil
}

func (ar *albumRepository) DeleteTrackFromAlbum(albumId uint64, trackId uint64) error {
	err := ar.db.Transaction(func(tx *gorm.DB) error {
		if err := ar.db.Delete(dao.Track{}, trackId).Error; err != nil {
			return err
		}

		eventID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

		if err := ar.db.Create(&dao.Outbox{
			ID:      0,
			EventId: eventID,
			TrackId: trackId,
			Type:    dao.TypeDelete,
			Sent:    false,
		}).Error; err != nil {
			return err
		}

		res := ar.db.Limit(1).Find(&dao.AlbumTrack{}, "album_id = ?", albumId)
		if res.Error != nil {
			return res.Error
		}

		// Delete album if no tracks left
		if res.RowsAffected <= 0 {
			if err := ar.db.Delete(&dao.Album{}, albumId).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "database error (table album)")
	}

	return nil
}

func (ar *albumRepository) GetAllTracksForAlbum(albumId uint64) ([]*models.Track, error) {
	var relations []*dao.AlbumTrack
	tx := ar.db.Limit(dao.MaxLimit).Find(&relations, "album_id = ?", albumId)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table album)")
	}

	var ids []uint64
	for _, v := range relations {
		ids = append(ids, v.TrackId)
	}

	var pgTracks []*dao.Track
	tx = ar.db.Find(&pgTracks, ids)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table album)")
	}

	var tracks []*models.Track
	for _, v := range pgTracks {
		var pgGenre dao.Genre
		tx := ar.db.Where("id = ?", v.GenreRefer).Take(&pgGenre)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table album)")
		}

		track := &models.Track{
			Id:     v.ID,
			Source: v.Source,
			Name:   v.Name,
			Genre:  pgGenre.Name,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
