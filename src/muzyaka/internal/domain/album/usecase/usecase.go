package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/album/repository"
	"src/internal/models"
)

type AlbumUseCase interface {
	GetAlbum(id uint64) (*models.Album, error)
	UpdateAlbum(album *models.Album) error
	AddAlbumWithTracks(album *models.Album, tracks []*models.TrackObject) (uint64, error)
	DeleteAlbum(id uint64) error
	AddTrack(albumId uint64, track *models.TrackObject) (uint64, error)
	DeleteTrack(albumId uint64, track *models.TrackMeta) error
	GetAllTracks(albumId uint64) ([]*models.TrackMeta, error)
}

type usecase struct {
	albumRep   repository.AlbumRepository
	storageRep repository.TrackStorage
}

func NewAlbumUseCase(albumRepository repository.AlbumRepository, storage repository.TrackStorage) AlbumUseCase {
	return &usecase{albumRep: albumRepository, storageRep: storage}
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

func (u *usecase) AddAlbumWithTracks(album *models.Album, tracks []*models.TrackObject) (uint64, error) {
	var tracksMeta []*models.TrackMeta
	for _, v := range tracks {
		err := u.storageRep.UploadObject(v)
		if err != nil {
			return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
		}

		tracksMeta = append(tracksMeta, v.ExtractMeta())
	}

	id, err := u.albumRep.AddAlbumWithTracks(album, tracksMeta)

	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
	}

	return id, nil
}

func (u *usecase) DeleteAlbum(id uint64) error {
	tracks, err := u.albumRep.GetAllTracksForAlbum(id)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteAlbum error while delete")
	}

	err = u.albumRep.DeleteAlbum(id)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteAlbum error while delete")
	}

	for _, v := range tracks {
		err = u.storageRep.DeleteObject(v)
		if err != nil {
			return errors.Wrap(err, "album.usecase.AddAlbum error while add")
		}
	}

	return nil
}

func (u *usecase) AddTrack(albumId uint64, track *models.TrackObject) (uint64, error) {
	err := u.storageRep.UploadObject(track)
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
	}

	id, err := u.albumRep.AddTrackToAlbum(albumId, track.ExtractMeta())
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddTrack error while add")
	}

	return id, nil
}

func (u *usecase) DeleteTrack(album_id uint64, track *models.TrackMeta) error {

	err := u.albumRep.DeleteTrackFromAlbum(album_id, track)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while delete")
	}

	err = u.storageRep.DeleteObject(track)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while add")
	}

	return nil
}

// TODO: разделить на данные и метаданные, через другой сервис могу извлекать данные.
// Грузить в принципе можно и тут
// Есть сервис терков, можно возвращать байтики из него
// Модель могу засплитить спокойно.

func (u *usecase) GetAllTracks(albumId uint64) ([]*models.TrackMeta, error) {
	tracks, err := u.albumRep.GetAllTracksForAlbum(albumId)

	if err != nil {
		return nil, errors.Wrap(err, "album.usecase.GetAllTracks error while get")
	}

	return tracks, err
}
