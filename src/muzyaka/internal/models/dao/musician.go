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
	ID         uint64 `gorm:"column:id"`
	PhotoFile  []byte `gorm:"column:photo_file"`
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
	for _, v := range musician.PhotoFiles {
		res = append(res, &MusicianPhotos{
			PhotoFile:  v,
			MusicianId: musician.Id,
		})
	}

	return res
}

func ToModelMusician(musician *Musician, photos []*MusicianPhotos) *models.Musician {
	var res [][]byte
	for _, v := range photos {
		res = append(res, v.PhotoFile)
	}

	return &models.Musician{
		Id:          musician.ID,
		Name:        musician.Name,
		PhotoFiles:  res,
		Description: musician.Description,
	}
}
