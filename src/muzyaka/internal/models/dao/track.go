package dao

import "src/internal/models"

type Genre struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (Genre) TableName() string {
	return "genres"
}

type Track struct {
	ID         uint64 `gorm:"column:id"`
	Source     string `gorm:"column:source"`
	Name       string `gorm:"column:name"`
	GenreRefer uint64 `gorm:"column:genre"`
}

func (Track) TableName() string {
	return "tracks"
}

func ToPostgresTrack(e *models.Track, genreRefer uint64) *Track {
	return &Track{
		ID:         e.Id,
		Source:     e.Source,
		Name:       e.Name,
		GenreRefer: genreRefer,
	}
}

func ToModelTrack(track *Track, genre *Genre) *models.Track {
	return &models.Track{
		Id:     track.ID,
		Source: track.Source,
		Name:   track.Name,
		Genre:  genre.Name,
	}
}
