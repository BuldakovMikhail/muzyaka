package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/user/repository"
	"src/internal/lib/jwt"
	"src/internal/models"
)

type AuthUseCase interface {
	SignUp(user *models.User) (*models.AuthToken, error)
	SignIn(user *models.User) (*models.AuthToken, error)
	Authorization(token *models.AuthToken, role string) (bool, error)
}

type usecase struct {
	userRep       repository.UserRepository
	tokenProvider jwt.TokenProvider
	encryptor     Encryptor
}

func NewAuthUseCase(tokenProvider jwt.TokenProvider,
	userRep repository.UserRepository,
	enc Encryptor) AuthUseCase {
	return &usecase{
		tokenProvider: tokenProvider,
		userRep:       userRep,
		encryptor:     enc,
	}
}

func (u *usecase) SignUp(user *models.User) (*models.AuthToken, error) {
	encPassword, err := u.encryptor.EncodePassword([]byte(user.Password))
	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp encode error")
	}

	temp := user
	temp.Password = string(encPassword)
	_, err = u.userRep.AddUser(temp)

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp AddUser error")
	}

	jwtToken, err := u.tokenProvider.GenerateToken(user)

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp token generation error")
	}

	user.Password = ""
	return jwtToken, nil
}

func (u *usecase) SignIn(user *models.User) (*models.AuthToken, error) {
	repUser, err := u.userRep.GetUser(user.Id)

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn user get error")
	}

	err = u.encryptor.CompareHashAndPassword([]byte(repUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn compare error")
	}
	user.Password = ""

	jwtToken, err := u.tokenProvider.GenerateToken(user)
	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn token generation error")
	}

	return jwtToken, nil
}

func (u *usecase) Authorization(token *models.AuthToken, role string) (bool, error) {
	tokenRole, err := u.tokenProvider.GetRole(token)
	if err != nil {
		return false, errors.Wrap(err, "auth.usecase.Authorization token parse error")
	}

	return tokenRole == role, nil
}
