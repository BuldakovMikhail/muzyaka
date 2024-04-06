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
	album.Id = pgAlbum.ID

	return nil
}

func (ar *albumRepository) AddAlbum(album *models.Album) (uint64, error) {
	pgAlbum := dao.ToPostgresAlbum(album)
	tx := ar.db.Create(pgAlbum)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table album)")
	}

	return pgAlbum.ID, nil
}

func (ar *albumRepository) DeleteAlbum(id uint64) error {
	tx := ar.db.Delete(&dao.Album{}, id)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
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
			ID:      0,
			EventId: eventID,
			TrackId: pgTrack.ID,
			Type:    "add",
			Sent:    false,
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
	tx := ar.db.Delete(dao.Track{}, trackId)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table album)")
	}

	return nil
}

func (ar *albumRepository) GetAllTracksForAlbum(albumId uint64) ([]*models.Track, error) {
	var relations []dao.AlbumTrack
	tx := ar.db.Limit(models.MaxLimit).Find(&relations, "album_id = ?", albumId)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table album)")
	}

	var ids []uint64
	for _, v := range relations {
		ids = append(ids, v.TrackId)
	}

	var pgTracks []dao.Track
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

		// TODO: fill all gaps
		track := &models.Track{
			Id:         v.ID,
			Source:     v.Source,
			Producers:  nil,
			Authors:    nil,
			Performers: nil,
			Name:       v.Name,
			Genre:      pgGenre.Name,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
