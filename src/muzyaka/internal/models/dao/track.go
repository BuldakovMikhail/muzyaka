package dao

import "src/internal/models"

type Genre struct {
	ID   *uint64 `gorm:"column:id"`
	Name string  `gorm:"column:name"`
}

func (Genre) TableName() string {
	return "genres"
}

type TrackMeta struct {
	ID         uint64  `gorm:"column:id"`
	Source     string  `gorm:"column:source"`
	Name       string  `gorm:"column:name"`
	GenreRefer *uint64 `gorm:"column:genre"`
	AlbumID    uint64  `gorm:"column:album_id"`
}

func (TrackMeta) TableName() string {
	return "tracks"
}

func ToPostgresTrack(e *models.TrackMeta, genreRefer *uint64, albumId uint64) *TrackMeta {
	var refer *uint64
	if genreRefer == nil {
		refer = nil
	} else {
		refer = genreRefer
	}

	return &TrackMeta{
		ID:         e.Id,
		Source:     e.Source,
		Name:       e.Name,
		GenreRefer: refer,
		AlbumID:    albumId,
	}
}

func ToModelTrack(track *TrackMeta, genre *Genre) *models.TrackMeta {
	return &models.TrackMeta{
		Id:     track.ID,
		Source: track.Source,
		Name:   track.Name,
		Genre:  genre.Name,
	}
}
