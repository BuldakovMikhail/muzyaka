package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type AlbumRepository interface {
	GetAlbum(id uint64) (*models.Album, error)
	UpdateAlbum(album *models.Album) error
	AddAlbumWithTracks(album *models.Album, tracks []*models.Track) (uint64, error)
	DeleteAlbum(id uint64) error
	AddTrackToAlbum(albumId uint64, track *models.Track) (uint64, error)
	DeleteTrackFromAlbum(albumId uint64, trackId uint64) error
	GetAllTracksForAlbum(albumId uint64) ([]*models.Track, error)
}
