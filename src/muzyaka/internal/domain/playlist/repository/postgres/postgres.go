package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/domain/playlist/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type playlistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) repository.PlaylistRepository {
	return &playlistRepository{db: db}
}

func (p playlistRepository) GetAllTracks(playlistId uint64) ([]*models.TrackMeta, error) {
	var relations []*dao.PlaylistTrack
	tx := p.db.Limit(dao.MaxLimit).Find(&relations, "playlist_id = ?", playlistId)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	var ids []uint64
	for _, v := range relations {
		ids = append(ids, v.TrackId)
	}

	var pgTracks []*dao.TrackMeta
	tx = p.db.Find(&pgTracks, ids)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	var tracks []*models.TrackMeta
	for _, v := range pgTracks {
		var pgGenre dao.Genre
		tx := p.db.Where("id = ?", v.GenreRefer).Take(&pgGenre)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table playlist)")
		}

		track := &models.TrackMeta{
			Id:     v.ID,
			Source: v.Source,
			Name:   v.Name,
			Genre:  pgGenre.Name,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (p playlistRepository) GetPlaylist(id uint64) (*models.Playlist, error) {
	var playlist dao.Playlist

	tx := p.db.Where("id = ?", id).Take(&playlist)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return dao.ToModelPlaylist(&playlist), nil
}

func (p playlistRepository) UpdatePlaylist(playlist *models.Playlist) error {
	pgPlaylist := dao.ToPostgresPlaylist(playlist)

	tx := p.db.Omit("id").Updates(&pgPlaylist)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return nil
}

func (p playlistRepository) AddPlaylist(playlist *models.Playlist) (uint64, error) {
	pgPlaylist := dao.ToPostgresPlaylist(playlist)

	tx := p.db.Create(&pgPlaylist)

	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	playlist.Id = pgPlaylist.ID

	return pgPlaylist.ID, nil
}

func (p playlistRepository) DeletePlaylist(id uint64) error {
	tx := p.db.Delete(&dao.Playlist{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return nil
}

func (p playlistRepository) AddTrackToPlaylist(playlistId uint64, trackId uint64) error {
	tx := p.db.Create(&dao.PlaylistTrack{
		TrackId:    trackId,
		PlaylistId: playlistId,
	})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return nil
}

func (p playlistRepository) DeleteTrackFromPlaylist(playlistId uint64, trackId uint64) error {
	tx := p.db.Delete(&dao.PlaylistTrack{
		TrackId:    trackId,
		PlaylistId: playlistId,
	})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return nil
}
