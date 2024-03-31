package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/user/repository"
	"src/internal/models"
)

type UserUseCase interface {
	UpdateUser(user *models.User) error
	GetUser(id uint64) (*models.User, error)
	AddUser(user *models.User) (uint64, error)
	DeleteUser(id uint64) error
}

type usecase struct {
	userRep repository.UserRepository
}

func (u *usecase) UpdateUser(user *models.User) error {
	err := u.userRep.UpdateUser(user)

	if err != nil {
		return errors.Wrap(err, "user.usecase.UpdateUser error while update")
	}

	return nil
}

func (u *usecase) GetUser(id uint64) (*models.User, error) {
	res, err := u.userRep.GetUser(id)

	if err != nil {
		return nil, errors.Wrap(err, "user.usecase.GetUser error while get")
	}

	return res, nil
}

func (u *usecase) AddUser(user *models.User) (uint64, error) {
	id, err := u.userRep.AddUser(user)

	if err != nil {
		return 0, errors.Wrap(err, "user.usecase.AddUser error while add")
	}

	return id, nil
}

func (u usecase) DeleteUser(id uint64) error {
	err := u.userRep.DeleteUser(id)

	if err != nil {
		return errors.Wrap(err, "user.usecase.DeleteUser error while delete")
	}

	return nil
}
