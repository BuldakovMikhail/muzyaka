package postgres

import (
	"gorm.io/gorm"
	"src/internal/models"
)

type playlistRepository struct {
	db *gorm.DB
}

func (p playlistRepository) GetPlaylist(id uint64) (*models.Playlist, error) {
	//TODO implement me
	panic("implement me")
}

func (p playlistRepository) UpdatePlaylist(playlist *models.Playlist) error {
	//TODO implement me
	panic("implement me")
}

func (p playlistRepository) AddPlaylist(playlist *models.Playlist) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (p playlistRepository) DeletePlaylist(id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (p playlistRepository) AddTrackToPlaylist(playlistId uint64, trackId uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (p playlistRepository) DeleteTrackFromPlaylist(playlistId uint64, trackId uint64) error {
	//TODO implement me
	panic("implement me")
}
