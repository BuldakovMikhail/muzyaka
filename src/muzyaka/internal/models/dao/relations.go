package dao

type AlbumTrack struct {
	AlbumId uint64 `gorm:"column:album_id"`
	TrackId uint64 `gorm:"column:track_id"`
}

func (AlbumTrack) TableName() string {
	return "albums_tracks"
}
