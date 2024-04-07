package models

type Track struct {
	Id      uint64
	Source  string
	Authors []string
	Name    string
	Genre   string
}
