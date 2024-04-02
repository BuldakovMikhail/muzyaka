package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type UserRepository interface {
	GetUser(id uint64) (*models.User, error)
	UpdateUser(user *models.User) error
	AddUser(user *models.User) (uint64, error)
	DeleteUser(id uint64) error
}
