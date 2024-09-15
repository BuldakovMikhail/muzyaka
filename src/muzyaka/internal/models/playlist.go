package models

type Playlist struct {
	Id          uint64
	Name        string
	CoverFile   []byte
	Description string
}
