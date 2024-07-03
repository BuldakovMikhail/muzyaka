package dto

import "src/internal/models"

type Playlist struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	CoverFile   []byte `json:"cover_file"`
	Description string `json:"description"`
}

type PlaylistWithUser struct {
	Playlist
	UserId uint64 `json:"user_id"`
}

type PlaylistWithoutId struct {
	Name        string `json:"name"`
	CoverFile   []byte `json:"cover_file"`
	Description string `json:"description"`
}

type CreatePlaylistResponse struct {
	Id uint64 `json:"id"`
}

type AddTrackPlaylistRequest struct {
	TrackId uint64 `json:"track_id"`
}

type AddTrackPlaylistResponse struct {
	Status string `json:"status"`
}

type PlaylistsCollection struct {
	Playlists []*Playlist `json:"playlists"`
}

func ToModelPlaylist(playlist *Playlist) *models.Playlist {
	return &models.Playlist{
		Id:          playlist.Id,
		Name:        playlist.Name,
		CoverFile:   playlist.CoverFile,
		Description: playlist.Description,
	}
}

func ToModelPlaylistWithoutId(playlist *PlaylistWithoutId, id uint64) *models.Playlist {
	return &models.Playlist{
		Id:          id,
		Name:        playlist.Name,
		CoverFile:   playlist.CoverFile,
		Description: playlist.Description,
	}
}

func ToDtoPlaylistWithUser(p *models.Playlist, userId uint64) *PlaylistWithUser {
	return &PlaylistWithUser{
		Playlist: Playlist{
			Id:          p.Id,
			Name:        p.Name,
			CoverFile:   p.CoverFile,
			Description: p.Description,
		},
		UserId: userId,
	}
}

func ToDtoPlaylist(playlist *models.Playlist) *Playlist {
	return &Playlist{
		Name:        playlist.Name,
		CoverFile:   playlist.CoverFile,
		Description: playlist.Description,
		Id:          playlist.Id,
	}
}
