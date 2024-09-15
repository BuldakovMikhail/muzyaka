package dao

import "src/internal/models"

type Playlist struct {
	ID          uint64 `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	CoverFile   []byte `gorm:"column:cover_file"`
	Description string `gorm:"column:description"`
	UserID      uint64 `gorm:"column:user_id"`
}

type PlaylistTrack struct {
	TrackId    uint64 `gorm:"column:track_id"`
	PlaylistId uint64 `gorm:"column:playlist_id"`
}

func (Playlist) TableName() string {
	return "playlists"
}

func (PlaylistTrack) TableName() string {
	return "track_playlist"
}

func ToModelPlaylist(playlist *Playlist) *models.Playlist {
	return &models.Playlist{
		Id:          playlist.ID,
		Name:        playlist.Name,
		CoverFile:   playlist.CoverFile,
		Description: playlist.Description,
	}
}

func ToPostgresPlaylist(playlist *models.Playlist, userId uint64) *Playlist {
	return &Playlist{
		ID:          playlist.Id,
		Name:        playlist.Name,
		CoverFile:   playlist.CoverFile,
		Description: playlist.Description,
		UserID:      userId,
	}
}
