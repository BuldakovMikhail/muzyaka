package models

type Musician struct {
	Id          uint64
	Name        string
	PhotoFiles  [][]byte
	Description string
}
