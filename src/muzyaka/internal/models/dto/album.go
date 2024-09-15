package dto

import (
	"src/internal/models"
)

type Album struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	CoverFile []byte `json:"cover_file"`
	Type      string `json:"type"`
}

type AlbumWithoutId struct {
	Name      string `json:"name"`
	CoverFile []byte `json:"cover_file"`
	Type      string `json:"type"`
}

type AlbumsCollection struct {
	Albums []*Album `json:"albums"`
}

type AlbumWithTracks struct {
	AlbumWithoutId
	Tracks []*TrackObjectWithoutId `json:"tracks"`
}

type CreateAlbumResponse struct {
	Id uint64 `json:"id"`
}

type CreateTrackResponse struct {
	Id uint64 `json:"id"`
}

func ToDtoAlbum(a *models.Album) *Album {
	return &Album{
		Id:        a.Id,
		Name:      a.Name,
		CoverFile: a.CoverFile,
		Type:      a.Type,
	}
}

func ToModelAlbum(a *Album) *models.Album {
	return &models.Album{
		Id:        a.Id,
		Name:      a.Name,
		CoverFile: a.CoverFile,
		Type:      a.Type,
	}
}

func ToModelAlbumWithId(id uint64, a *AlbumWithoutId) *models.Album {
	return &models.Album{
		Id:        id,
		Name:      a.Name,
		CoverFile: a.CoverFile,
		Type:      a.Type,
	}
}
