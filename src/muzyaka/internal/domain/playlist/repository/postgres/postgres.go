package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
	"src/internal/models/dao"
)

type playlistRepository struct {
	db *gorm.DB
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
