package models

type TrackMeta struct {
	Id     uint64
	Source string
	Name   string
	Genre  string
}

type TrackObject struct {
	TrackMeta
	Payload     []byte
	PayloadSize uint64
}

func (t *TrackObject) ExtractMeta() *TrackMeta {
	return &t.TrackMeta
}
