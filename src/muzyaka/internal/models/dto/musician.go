package dto

import (
	"src/internal/models"
)

type Musician struct {
	Id          uint64   `json:"id"`
	Name        string   `json:"musician_name"`
	PhotoFiles  [][]byte `json:"photo_files"`
	Description string   `json:"description"`
}

type MusicianWithoutId struct {
	Name        string   `json:"musician_name"`
	PhotoFiles  [][]byte `json:"photo_files"`
	Description string   `json:"description"`
}

type CreateMusicianResponse struct {
	Id uint64 `json:"id"`
}

func ToModelMusician(musician *Musician) *models.Musician {
	return &models.Musician{
		Id:          musician.Id,
		Name:        musician.Name,
		PhotoFiles:  musician.PhotoFiles,
		Description: musician.Description,
	}
}

func ToModelMusicianWithoutId(musician *MusicianWithoutId, id uint64) *models.Musician {
	return &models.Musician{
		Id:          id,
		Name:        musician.Name,
		PhotoFiles:  musician.PhotoFiles,
		Description: musician.Description,
	}
}

func ToDtoMusician(musician *models.Musician) *Musician {
	return &Musician{
		Id:          musician.Id,
		Name:        musician.Name,
		PhotoFiles:  musician.PhotoFiles,
		Description: musician.Description,
	}
}
