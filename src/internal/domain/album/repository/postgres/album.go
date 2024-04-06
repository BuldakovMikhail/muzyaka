package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
)

type Album struct {
	ID    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Cover string `gorm:"column:cover"`
	Type  string `gorm:"column:type"`
}

func (Album) TableName() string {
	return "album"
}

type Outbox struct {
	ID      uint64 `gorm:"column:id"`
	EventId uint64 `gorm:"column:event_id"`
	OrderId uint64 `gorm:"column:order_id"`
	Type    string `gorm:"column:type"`
	Sent    bool   `gorm:"column:sent"`
}

func (Outbox) TableName() string {
	return "outbox"
}

type albumRepository struct {
	db *gorm.DB
}

func toPostgresAlbum(e *models.Album) *Album {
	return &Album{
		ID:    e.Id,
		Name:  e.Name,
		Cover: e.Cover,
		Type:  e.Type,
	}
}

func toModelAlbum(e *Album) *models.Album {
	return &models.Album{
		Id:    e.ID,
		Name:  e.Name,
		Cover: e.Cover,
		Type:  e.Type,
	}
}

func (ar *albumRepository) GetAlbum(id uint64) (*models.Album, error) {
	var album Album

	tx := ar.db.Where("id = ?", id).Take(&album)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table entry)")
	}
	return toModelAlbum(&album), nil
}

func (ar *albumRepository) UpdateAlbum(album *models.Album) error {
	//TODO implement me
	panic("implement me")
}

func (ar *albumRepository) AddAlbum(album *models.Album) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (ar *albumRepository) DeleteAlbum(id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (ar *albumRepository) AddTrackToAlbum(albumId uint64, track *models.Track) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (ar *albumRepository) DeleteTrackFromAlbum(albumId uint64, trackId uint64) error {
	//TODO implement me
	panic("implement me")
}

func (ar *albumRepository) GetAllTracksForAlbum(albumId uint64) ([]*models.Track, error) {
	//TODO implement me
	panic("implement me")
}
