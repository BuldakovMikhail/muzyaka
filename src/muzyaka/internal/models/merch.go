package models

type Merch struct {
	Id          uint64
	Name        string
	PhotoFiles  [][]byte
	Description string
	OrderUrl    string
}
