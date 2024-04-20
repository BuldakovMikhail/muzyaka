package models

// TODO: replace photos with array of bytes

type Musician struct {
	Id          uint64
	Name        string
	Photos      []string
	Description string
}
