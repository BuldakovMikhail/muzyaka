package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/playlist/repository"
	repository2 "src/internal/domain/track/repository"
	"src/internal/models"
)

type PlaylistUseCase interface {
	UpdatedPlaylist(playlist *models.Playlist) error
	AddPlaylist(playlist *models.Playlist, userId uint64) (uint64, error)
	DeletePlaylist(id uint64) error
	GetPlaylist(id uint64) (*models.Playlist, error)
	AddTrack(playlistId uint64, trackId uint64) error
	DeleteTrack(playlistId uint64, trackId uint64) error
	GetAllTracks(playlistId uint64) ([]*models.TrackMeta, error)
	GetUserForPlaylist(playlistId uint64) (uint64, error)

	IsPlaylistOwned(playlistId uint64, userId uint64) (bool, error)
	GetAllPlaylistsForUser(userId uint64) ([]*models.Playlist, error)
}

type usecase struct {
	playlistRep repository.PlaylistRepository
	trackRep    repository2.TrackRepository
}

// TODO: добавить ограничение на количество загружаемых сущностей

func NewPlaylistUseCase(rep repository.PlaylistRepository, trackRep repository2.TrackRepository) PlaylistUseCase {
	return &usecase{playlistRep: rep, trackRep: trackRep}
}

func (u *usecase) GetAllPlaylistsForUser(userId uint64) ([]*models.Playlist, error) {
	playlists, err := u.playlistRep.GetAllPlaylistsForUser(userId)
	if err != nil {
		return nil, errors.Wrap(err, "playlist.usecase.GetAllPlaylistsForUser error while get")
	}

	return playlists, nil
}

func (u *usecase) IsPlaylistOwned(playlistId uint64, userId uint64) (bool, error) {
	isAllowed, err := u.playlistRep.IsPlaylistOwned(playlistId, userId)
	if err != nil {
		return false, errors.Wrap(err, "playlist.usecase.IsPlaylistOwned error while get")
	}

	return isAllowed, nil
}

func (u *usecase) GetUserForPlaylist(playlistId uint64) (uint64, error) {
	userId, err := u.playlistRep.GetUserForPlaylist(playlistId)
	if err != nil {
		return 0, errors.Wrap(err, "playlist.usecase.GetUserForPlaylist error while get")
	}

	return userId, nil
}

func (u *usecase) GetAllTracks(playlistId uint64) ([]*models.TrackMeta, error) {
	trackIds, err := u.playlistRep.GetAllTracks(playlistId)
	if err != nil {
		return nil, errors.Wrap(err, "playlist.usecase.GetAllTracks error while get")
	}

	var tracks []*models.TrackMeta
	for _, v := range trackIds {
		track, err := u.trackRep.GetTrack(v)
		if err != nil {
			return nil, errors.Wrap(err, "playlist.usecase.GetAllTracks error while get")
		}

		tracks = append(tracks, track)
	}

	return tracks, err
}

func (u *usecase) UpdatedPlaylist(playlist *models.Playlist) error {
	err := u.playlistRep.UpdatePlaylist(playlist)

	if err != nil {
		return errors.Wrap(err, "playlist.usecase.UpdatedPlaylist error while update")
	}

	return nil
}

func (u *usecase) AddPlaylist(playlist *models.Playlist, userId uint64) (uint64, error) {
	id, err := u.playlistRep.AddPlaylist(playlist, userId)

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

func (u *usecase) AddTrack(playlistId uint64, trackId uint64) error {
	err := u.playlistRep.AddTrackToPlaylist(playlistId, trackId)

	if err != nil {
		return errors.Wrap(err, "playlist.usecase.AddTrack error while add")
	}

	return nil
}

func (u *usecase) DeleteTrack(playlistId uint64, trackId uint64) error {
	err := u.playlistRep.DeleteTrackFromPlaylist(playlistId, trackId)
	if err != nil {
		return errors.Wrap(err, "playlist.usecase.DeleteTrack error while delete")
	}

	return nil
}
