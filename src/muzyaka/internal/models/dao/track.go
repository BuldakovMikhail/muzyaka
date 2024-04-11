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

type Author struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (Author) TableName() string {
	return "authors"
}

type TrackAuthor struct {
	AuthorId uint64 `gorm:"column:author_id"`
	TrackId  uint64 `gorm:"column:track_id"`
}

func (TrackAuthor) TableName() string {
	return "track_authors"
}

func ToPostgresTrack(e *models.Track, genreRefer uint64) *Track {
	return &Track{
		ID:         e.Id,
		Source:     e.Source,
		Name:       e.Name,
		GenreRefer: genreRefer,
	}
}

func ToModelTrack(track *Track, genre *Genre, authors []*Author) *models.Track {
	var authorsNames []string

	for _, v := range authors {
		authorsNames = append(authorsNames, v.Name)
	}

	return &models.Track{
		Id:      track.ID,
		Source:  track.Source,
		Authors: authorsNames,
		Name:    track.Name,
		Genre:   genre.Name,
	}
}
