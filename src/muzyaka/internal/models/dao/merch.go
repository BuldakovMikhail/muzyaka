package dao

import "src/internal/models"

type Merch struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Desc string `gorm:"column:description"`
	Link string `gorm:"column:link"`
}

func (Merch) TableName() string {
	return "merch"
}

type MerchPhotos struct {
	MerchId   uint64 `gorm:"column:merch_id"`
	PhotosSrc string `gorm:"column:photo_src"`
}

func (MerchPhotos) TableName() string {
	return "merch_photos"
}

func ToPostgresMerch(e *models.Merch) *Merch {
	return &Merch{
		ID:   e.Id,
		Name: e.Name,
		Desc: e.Description,
		Link: e.OrderUrl,
	}
}

func ToPostgresMerchPhotos(e *models.Merch) []*MerchPhotos {
	var merchPhotos []*MerchPhotos

	for _, v := range e.Photos {
		merchPhotos = append(
			merchPhotos,
			&MerchPhotos{
				MerchId:   e.Id,
				PhotosSrc: v,
			},
		)
	}

	return merchPhotos
}

func ToModelMerch(e *Merch, mp []*MerchPhotos) *models.Merch {
	var photos []string

	for _, v := range mp {
		photos = append(photos, v.PhotosSrc)
	}

	return &models.Merch{
		Id:          e.ID,
		Name:        e.Name,
		Photos:      photos,
		Description: e.Desc,
		OrderUrl:    e.Link,
	}
}
