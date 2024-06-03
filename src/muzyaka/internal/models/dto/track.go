package dto

import "src/internal/models"

type TrackMeta struct {
	Id     uint64  `json:"id"`
	Source string  `json:"source"`
	Name   string  `json:"name"`
	Genre  *string `json:"genre"`
}

type TrackMetaWithoutId struct {
	Source string  `json:"source"`
	Name   string  `json:"name"`
	Genre  *string `json:"genre"`
}

type TrackObjectWithoutId struct {
	TrackMetaWithoutId
	Payload     []byte `json:"payload"`
	PayloadSize int64  `json:"payload_size"`
}

type TrackObject struct {
	TrackMeta
	Payload     []byte `json:"payload"`
	PayloadSize int64  `json:"payload_size"`
}

type TracksMetaCollection struct {
	Tracks []*TrackMeta `json:"tracks"`
}

func ToDtoTrackMeta(m *models.TrackMeta) *TrackMeta {
	var genre *string
	if m.Genre == "" {
		genre = nil
	} else {
		genre = &m.Genre
	}

	return &TrackMeta{
		Id:     m.Id,
		Source: m.Source,
		Name:   m.Name,
		Genre:  genre,
	}
}

func ToModelTrackMeta(t *TrackMeta) *models.TrackMeta {
	var genre string
	if t.Genre == nil {
		genre = ""
	} else {
		genre = *t.Genre
	}

	return &models.TrackMeta{
		Id:     t.Id,
		Source: t.Source,
		Name:   t.Name,
		Genre:  genre,
	}
}

func ToModelTrackObject(t *TrackObject) *models.TrackObject {
	var genre string
	if t.Genre == nil {
		genre = ""
	} else {
		genre = *t.Genre
	}
	return &models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     t.Id,
			Source: t.Source,
			Name:   t.Name,
			Genre:  genre,
		},
		Payload:     t.Payload,
		PayloadSize: t.PayloadSize,
	}
}

func ToModelTrackObjectWithoutId(t *TrackObjectWithoutId, id uint64) *models.TrackObject {
	var genre string
	if t.Genre == nil {
		genre = ""
	} else {
		genre = *t.Genre
	}
	return &models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     id,
			Source: t.Source,
			Name:   t.Name,
			Genre:  genre,
		},
		Payload:     t.Payload,
		PayloadSize: t.PayloadSize,
	}
}

func ToDtoTrackObject(t *models.TrackObject) *TrackObject {
	var genre *string
	if t.Genre == "" {
		genre = nil
	} else {
		genre = &t.Genre
	}
	return &TrackObject{
		TrackMeta: TrackMeta{
			Id:     t.Id,
			Source: t.Source,
			Name:   t.Name,
			Genre:  genre,
		},
		Payload:     t.Payload,
		PayloadSize: t.PayloadSize,
	}
}
