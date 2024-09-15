package dto

import "src/internal/models"

type Genres struct {
	Genres []string `json:"genres"`
}

type TrackMeta struct {
	Id    uint64  `json:"id"`
	Name  string  `json:"name"`
	Genre *string `json:"genre"`
}

type TrackMetaWithoutId struct {
	Name  string  `json:"name"`
	Genre *string `json:"genre"`
}

type TrackObjectWithoutId struct {
	TrackMetaWithoutId
	Payload []byte `json:"payload"`
}

type TrackObject struct {
	TrackMeta
	Payload []byte `json:"payload"`
}

type TrackObjectWithSource struct {
	TrackMeta
	Source  string `json:"source"`
	Payload []byte `json:"payload"`
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
		Id:    m.Id,
		Name:  m.Name,
		Genre: genre,
	}
}

func ToModelTrackObjectWithoutId(t *TrackObjectWithoutId, id uint64, source string) *models.TrackObject {
	var genre string
	if t.Genre == nil {
		genre = ""
	} else {
		genre = *t.Genre
	}
	return &models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     id,
			Source: source,
			Name:   t.Name,
			Genre:  genre,
		},
		Payload: t.Payload,
	}
}

func ToDtoTrackObjectWithSource(t *models.TrackObject) *TrackObjectWithSource {
	var genre *string
	if t.Genre == "" {
		genre = nil
	} else {
		genre = &t.Genre
	}
	return &TrackObjectWithSource{
		TrackMeta: TrackMeta{
			Id:    t.Id,
			Name:  t.Name,
			Genre: genre,
		},
		Source:  t.Source,
		Payload: t.Payload,
	}
}
