package models

type Merch struct {
	Id          uint64
	Name        string
	Photos      []string
	Description string
	OrderUrl    string
}
