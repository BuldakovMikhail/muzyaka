package models

type Album struct {
	Id        uint64
	Name      string
	CoverFile []byte
	Type      string
}
