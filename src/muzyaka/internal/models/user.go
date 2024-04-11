package models

type User struct {
	Id       uint64
	Name     string
	Password string
	Role     string
	Email    string
}
