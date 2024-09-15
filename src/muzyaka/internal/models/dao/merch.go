package dao

import "src/internal/models"

type Merch struct {
	ID         uint64 `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Desc       string `gorm:"column:description"`
	Link       string `gorm:"column:link"`
	MusicianID uint64 `gorm:"column:musician_id"`
}

func (Merch) TableName() string {
	return "merch"
}

type MerchPhotos struct {
	ID        uint64 `gorm:"column:id"`
	MerchId   uint64 `gorm:"column:merch_id"`
	PhotoFile []byte `gorm:"column:photo_file"`
}

func (MerchPhotos) TableName() string {
	return "merch_photos"
}

func ToPostgresMerch(e *models.Merch, musicianId uint64) *Merch {
	return &Merch{
		ID:         e.Id,
		Name:       e.Name,
		Desc:       e.Description,
		Link:       e.OrderUrl,
		MusicianID: musicianId,
	}
}

func ToPostgresMerchPhotos(e *models.Merch) []*MerchPhotos {
	var merchPhotos []*MerchPhotos

	for _, v := range e.PhotoFiles {
		merchPhotos = append(
			merchPhotos,
			&MerchPhotos{
				MerchId:   e.Id,
				PhotoFile: v,
			},
		)
	}

	return merchPhotos
}

func ToModelMerch(e *Merch, mp []*MerchPhotos) *models.Merch {
	var photos [][]byte

	for _, v := range mp {
		photos = append(photos, v.PhotoFile)
	}

	return &models.Merch{
		Id:          e.ID,
		Name:        e.Name,
		PhotoFiles:  photos,
		Description: e.Desc,
		OrderUrl:    e.Link,
	}
}
