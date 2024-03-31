package usecase

import (
	"src/internal/domain/user/repository"
	"src/internal/models"
)

type AuthUseCase interface {
	SignUp(user *models.User) (*models.AuthToken, error)
	SignIn(user *models.User) (*models.AuthToken, error)
}

type usecase struct {
	userRep repository.UserRepository
}

func (u *usecase) SignUp(user *models.User) (*models.AuthToken, error) {
	u.userRep.
}

func (u *usecase) SignIn(user *models.User) (*models.AuthToken, error) {
	//TODO implement me
	panic("implement me")
}
