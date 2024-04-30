package models

// TODO: replace photos with array of bytes

type Merch struct {
	Id          uint64
	Name        string
	Photos      [][]byte
	Description string
	OrderUrl    string
}
