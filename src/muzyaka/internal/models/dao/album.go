package dao

import "src/internal/models"

type Album struct {
	ID         uint64 `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Cover      []byte `gorm:"column:cover_file"`
	Type       string `gorm:"column:type"`
	MusicianID uint64 `gorm:"column:musician_id"`
}

func (Album) TableName() string {
	return "albums"
}

func ToPostgresAlbum(e *models.Album, musicianId uint64) *Album {
	return &Album{
		ID:         e.Id,
		Name:       e.Name,
		Cover:      e.CoverFile,
		Type:       e.Type,
		MusicianID: musicianId,
	}
}

func ToModelAlbum(e *Album) *models.Album {
	return &models.Album{
		Id:        e.ID,
		Name:      e.Name,
		CoverFile: e.Cover,
		Type:      e.Type,
	}
}
