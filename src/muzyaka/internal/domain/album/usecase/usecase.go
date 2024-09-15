package usecase

import (
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"src/internal/domain/album/repository"
	repository2 "src/internal/domain/track/repository"
	"src/internal/models"
)

type AlbumUseCase interface {
	GetAlbum(id uint64) (*models.Album, error)
	UpdateAlbum(album *models.Album) error
	AddAlbumWithTracks(album *models.Album, tracks []*models.TrackObject, musicianId uint64) (uint64, error)
	DeleteAlbum(id uint64) error
	AddTrack(albumId uint64, track *models.TrackObject) (uint64, error)
	DeleteTrack(trackId uint64) error
	GetAllTracks(albumId uint64) ([]*models.TrackMeta, error)

	IsAlbumOwned(albumId uint64, musicianId uint64) (bool, error)
	GetAlbumIdForTrack(trackId uint64) (uint64, error)
	GetAllAlbumsForMusician(musicianId uint64) ([]*models.Album, error)
}

type usecase struct {
	albumRep   repository.AlbumRepository
	storageRep repository2.TrackStorage
	trackRep   repository2.TrackRepository
}

func NewAlbumUseCase(albumRepository repository.AlbumRepository,
	storage repository2.TrackStorage,
	trackRepository repository2.TrackRepository) AlbumUseCase {
	return &usecase{albumRep: albumRepository, storageRep: storage, trackRep: trackRepository}
}

func (u *usecase) GetAllAlbumsForMusician(musicianId uint64) ([]*models.Album, error) {
	albums, err := u.albumRep.GetAllAlbumsForMusician(musicianId)
	if err != nil {
		return nil, errors.Wrap(err, "track.usecase.GetAllAlbumsForMusician error while get")
	}

	return albums, nil
}

func (u *usecase) GetAlbumIdForTrack(trackId uint64) (uint64, error) {
	id, err := u.albumRep.GetAlbumId(trackId)
	if err != nil {
		return 0, errors.Wrap(err, "track.usecase.GetAlbumIdForTrack error while get")
	}

	return id, nil
}

func (u *usecase) IsAlbumOwned(albumId uint64, musicianId uint64) (bool, error) {
	res, err := u.albumRep.IsAlbumOwned(albumId, musicianId)

	if err != nil {
		return false, errors.Wrap(err, "album.usecase.IsAlbumOwned error while get")
	}

	return res, nil
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

func (u *usecase) AddAlbumWithTracks(album *models.Album, tracks []*models.TrackObject, musicianId uint64) (uint64, error) {
	var tracksMeta []*models.TrackMeta
	for _, v := range tracks {
		if len(v.Payload) == 0 {
			return 0, models.ErrInvalidPayload
		}

		newSource, err := uuid.GenerateUUID()
		if err != nil {
			return 0, errors.Wrap(err, "album.usecase.AddAlbum error in UUID gen")
		}

		v.Source = newSource

		err = u.storageRep.UploadObject(v)
		if err != nil {
			return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
		}

		tracksMeta = append(tracksMeta, v.ExtractMeta())
	}
	id, err := u.albumRep.AddAlbumWithTracksOutbox(album, tracksMeta, musicianId)

	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddAlbum error while add")
	}

	return id, nil
}

func (u *usecase) DeleteAlbum(id uint64) error {
	tracks, err := u.albumRep.GetAllTracksForAlbum(id)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteAlbumOutbox error while delete")
	}

	err = u.albumRep.DeleteAlbumOutbox(id)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteAlbumOutbox error while delete")
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
	if len(track.Payload) == 0 {
		return 0, models.ErrInvalidPayload
	}

	newSource, err := uuid.GenerateUUID()
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddTrack error in UUID gen")
	}
	track.Source = newSource

	err = u.storageRep.UploadObject(track)
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddTrack error while add")
	}

	id, err := u.albumRep.AddTrackToAlbumOutbox(albumId, track.ExtractMeta())
	if err != nil {
		return 0, errors.Wrap(err, "album.usecase.AddTrack error while add")
	}

	return id, nil
}

func (u *usecase) DeleteTrack(trackId uint64) error {
	trackMeta, err := u.trackRep.GetTrack(trackId)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while get")
	}

	err = u.albumRep.DeleteTrackFromAlbumOutbox(trackId)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while delete")
	}

	err = u.storageRep.DeleteObject(trackMeta)
	if err != nil {
		return errors.Wrap(err, "album.usecase.DeleteTrack error while add")
	}

	return nil
}

func (u *usecase) GetAllTracks(albumId uint64) ([]*models.TrackMeta, error) {
	tracks, err := u.albumRep.GetAllTracksForAlbum(albumId)

	if err != nil {
		return nil, errors.Wrap(err, "album.usecase.GetAllTracks error while get")
	}

	return tracks, err
}
