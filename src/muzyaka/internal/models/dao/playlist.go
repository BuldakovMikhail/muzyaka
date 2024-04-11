package dao

import "src/internal/models"

type Playlist struct {
	ID          uint64 `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	Cover       string `gorm:"column:cover"`
	Description string `gorm:"column:description"`
}

type PlaylistTrack struct {
	TrackId    uint64 `gorm:"column:track_id"`
	PlaylistId uint64 `gorm:"column:playlist_id"`
}

func (Playlist) TableName() string {
	return "playlists"
}

func (PlaylistTrack) TableName() string {
	return "playlists_tracks"
}

func ToModelPlaylist(playlist *Playlist) *models.Playlist {
	return &models.Playlist{
		Id:          playlist.ID,
		Name:        playlist.Name,
		Cover:       playlist.Cover,
		Description: playlist.Description,
	}
}

func ToPostgresPlaylist(playlist *models.Playlist) *Playlist {
	return &Playlist{
		ID:          playlist.Id,
		Name:        playlist.Name,
		Cover:       playlist.Cover,
		Description: playlist.Description,
	}
}
