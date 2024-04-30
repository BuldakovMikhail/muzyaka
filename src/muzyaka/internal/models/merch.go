package models

type Merch struct {
	Id          uint64
	Name        string
	Photos      [][]byte
	Description string
	OrderUrl    string
}
