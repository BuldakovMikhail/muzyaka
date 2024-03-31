package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/album/repository"
	"src/internal/models"
)

type AlbumUseCase interface {
	GetAlbum(id uint64) (*models.Album, error)
	UpdateAlbum(album *models.Album) error
	AddAlbum(album *models.Album) (uint64, error)
	DeleteAlbum(id uint64) error
	AddTrack(album_id uint64, track *models.Track) (uint64, error)
	DeleteTrack(album_id uint64, track_id uint64) error
	GetAllTracks(album_id uint64) ([]*models.Track, error)
}

type usecase struct {
	albumRep repository.AlbumRepository
}

func (u *usecase) GetAlbum(id uint64) (*models.Album, error) {
	res, err := u.albumRep.GetAlbum(id)

	if err != nil {
		return nil, errors.Wrap(err, "album.usecase.GetAlbum error while get")
	}

	return res, nil
}

func (u *usecase) UpdateAlbum(album *models.Album) error {
	err := u.albumRep.UpdateAlbum(album)

	if err != nil {
		return errors.Wrap(err, "album.usecase.UpdateAlbum error while update")
	}

	return nil
}

func (u *usecase) AddAlbum(album *models.Album) (uint64, error) {
	id, err := u.albumRep.AddAlbum(album)

	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
	}

	return id, nil
}

func (u *usecase) DeleteAlbum(id uint64) error {
	err := u.albumRep.DeleteAlbum(id)

	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteAlbum error while delete")
	}

	return nil
}

func (u *usecase) AddTrack(album_id uint64, track *models.Track) (uint64, error) {
	id, err := u.albumRep.AddTrackToAlbum(album_id, track)
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddTrack error while add")
	}

	return id, nil
}

func (u *usecase) DeleteTrack(album_id uint64, track_id uint64) error {
	err := u.albumRep.DeleteTrackFromAlbum(album_id, track_id)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while delete")
	}

	return nil
}

func (u *usecase) GetAllTracks(albumId uint64) ([]*models.Track, error) {
	tracks, err := u.albumRep.GetAllTracksForAlbum(albumId)
	if err != nil {
		return nil, errors.Wrap(err, "album.usecase.GetAllTracks error while get")
	}

	return tracks, err
}
