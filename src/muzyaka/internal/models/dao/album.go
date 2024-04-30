package dao

import "src/internal/models"

type Album struct {
	ID    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Cover []byte `gorm:"column:cover"`
	Type  string `gorm:"column:type"`
}

func (Album) TableName() string {
	return "albums"
}

func ToPostgresAlbum(e *models.Album) *Album {
	return &Album{
		ID:    e.Id,
		Name:  e.Name,
		Cover: e.Cover,
		Type:  e.Type,
	}
}

func ToModelAlbum(e *Album) *models.Album {
	return &models.Album{
		Id:    e.ID,
		Name:  e.Name,
		Cover: e.Cover,
		Type:  e.Type,
	}
}
