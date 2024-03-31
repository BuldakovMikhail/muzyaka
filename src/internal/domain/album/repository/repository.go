package repository

import "src/internal/models"

type AlbumRepository interface {
	GetAlbum(id uint64) (*models.Album, error)
	UpdateAlbum(album *models.Album) error
	AddAlbum(album *models.Album) (uint64, error)
	DeleteAlbum(id uint64) error
	AddTrackToAlbum(albumId uint64, track *models.Track) (uint64, error)
	DeleteTrackFromAlbum(albumId uint64, trackId uint64) error
	GetAllTracksForAlbum(albumId uint64) ([]*models.Track, error)
}
