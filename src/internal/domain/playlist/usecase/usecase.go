package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/playlist/repository"
	"src/internal/models"
)

type PlaylistUseCase interface {
	UpdatedPlaylist(playlist *models.Playlist) error
	AddPlaylist(playlist *models.Playlist) (uint64, error)
	DeletePlaylist(id uint64) error
	GetPlaylist(id uint64) (*models.Playlist, error)
	AddTrack(playlistId uint64, trackId uint64) (uint64, error)
	DeleteTrack(playlistId uint64, trackId uint64) error
}

type usecase struct {
	playlistRep repository.PlaylistRepository
}

func (u *usecase) UpdatedPlaylist(playlist *models.Playlist) error {
	err := u.playlistRep.UpdatePlaylist(playlist)

	if err != nil {
		return errors.Wrap(err, "playlist.usecase.UpdatedPlaylist error while update")
	}

	return nil
}

func (u *usecase) AddPlaylist(playlist *models.Playlist) (uint64, error) {
	id, err := u.playlistRep.AddPlaylist(playlist)

	if err != nil {
		return 0, errors.Wrap(err, "playlist.usecase.AddPlaylist error while add")
	}

	return id, nil
}

func (u *usecase) DeletePlaylist(id uint64) error {
	err := u.playlistRep.DeletePlaylist(id)

	if err != nil {
		return errors.Wrap(err, "playlist.usecase.DeletePlaylist error while delete")
	}

	return nil
}

func (u *usecase) GetPlaylist(id uint64) (*models.Playlist, error) {
	res, err := u.playlistRep.GetPlaylist(id)

	if err != nil {
		return nil, errors.Wrap(err, "playlist.usecase.GetPlaylist error while get")
	}

	return res, nil
}

func (u *usecase) AddTrack(playlistId uint64, trackId uint64) (uint64, error) {
	id, err := u.playlistRep.AddTrackToPlaylist(playlistId, trackId)

	if err != nil {
		return 0, errors.Wrap(err, "playlist.usecase.AddTrack error while add")
	}

	return id, nil
}

func (u *usecase) DeleteTrack(playlistId uint64, trackId uint64) error {
	err := u.playlistRep.DeleteTrackFromPlaylist(playlistId, trackId)
	if err != nil {
		return errors.Wrap(err, "playlist.usecase.DeleteTrack error while delete")
	}

	return nil
}
