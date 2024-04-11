package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type PlaylistRepository interface {
	GetPlaylist(id uint64) (*models.Playlist, error)
	UpdatePlaylist(playlist *models.Playlist) error
	AddPlaylist(playlist *models.Playlist) (uint64, error)
	DeletePlaylist(id uint64) error
	AddTrackToPlaylist(playlistId uint64, trackId uint64) error
	DeleteTrackFromPlaylist(playlistId uint64, trackId uint64) error
}
