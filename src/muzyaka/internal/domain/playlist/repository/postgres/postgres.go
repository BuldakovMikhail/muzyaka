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

func (p playlistRepository) IsPlaylistOwned(playlistId uint64, userId uint64) (bool, error) {
	var playlist dao.Playlist

	tx := p.db.Where("id = ?", playlistId).Take(&playlist)
	if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return playlist.UserID == userId, nil
}

func (p playlistRepository) GetUserForPlaylist(playlistId uint64) (uint64, error) {
	var playlist dao.Playlist

	tx := p.db.Where("id = ?", playlistId).Take(&playlist)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return playlist.UserID, nil
}

func (p playlistRepository) GetAllTracks(playlistId uint64) ([]uint64, error) {
	var relations []*dao.PlaylistTrack
	tx := p.db.Limit(dao.MaxLimit).Find(&relations, "playlist_id = ?", playlistId)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table playlist)")
	}

	var ids []uint64
	for _, v := range relations {
		ids = append(ids, v.TrackId)
	}

	return ids, nil
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
	pgPlaylist := dao.ToPostgresPlaylist(playlist, 0)

	tx := p.db.Omit("id", "user_id").Updates(&pgPlaylist)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	return nil
}

func (p playlistRepository) AddPlaylist(playlist *models.Playlist, userId uint64) (uint64, error) {
	pgPlaylist := dao.ToPostgresPlaylist(playlist, userId)

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

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
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
	tx := p.db.Delete(&dao.PlaylistTrack{},
		"track_id = ? AND playlist_id = ?", trackId, playlistId)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table playlist)")
	}

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}
