package postgres

import (
	"gorm.io/gorm"
	"src/internal/models"
)

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) GetUser(id uint64) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) UpdateUser(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) AddUser(user *models.User) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) DeleteUser(id uint64) error {
	//TODO implement me
	panic("implement me")
}
