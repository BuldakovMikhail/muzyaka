package dao

import "src/internal/models"

type Musician struct {
	ID          uint64 `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (Musician) TableName() string {
	return "musicians"
}

type MusicianPhotos struct {
	PhotoSrc   string `gorm:"column:photo_src"`
	MusicianId uint64 `gorm:"column:musician_id"`
}

func (MusicianPhotos) TableName() string {
	return "musicians_photos"
}

func ToPostgresMusician(musician *models.Musician) *Musician {
	return &Musician{
		ID:          musician.Id,
		Name:        musician.Name,
		Description: musician.Description,
	}
}

func ToPostgresMusicianPhotos(musician *models.Musician) []*MusicianPhotos {
	var res []*MusicianPhotos
	for _, v := range musician.Photos {
		res = append(res, &MusicianPhotos{
			PhotoSrc:   v,
			MusicianId: musician.Id,
		})
	}

	return res
}

func ToModelMusician(musician *Musician, photos []*MusicianPhotos) *models.Musician {
	var res []string
	for _, v := range photos {
		res = append(res, v)
	}

	return &models.Musician{
		Id:          musician.ID,
		Name:        musician.Name,
		Photos:      res,
		Description: musician.Description,
	}
}
