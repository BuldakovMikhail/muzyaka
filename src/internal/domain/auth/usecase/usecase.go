package usecase

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"src/internal/domain/user/repository"
	"src/internal/lib/jwt"
	"src/internal/models"
	"time"
)

type AuthUseCase interface {
	SignUp(user *models.User) (*models.AuthToken, error)
	SignIn(user *models.User) (*models.AuthToken, error)
	Authorization(token *models.AuthToken, role string) (bool, error)
}

type usecase struct {
	asyncKey string
	duration time.Duration
	userRep  repository.UserRepository
}

func NewAuthUseCase(key string,
	duration time.Duration,
	userRep repository.UserRepository) AuthUseCase {
	return &usecase{
		asyncKey: key,
		duration: duration,
		userRep:  userRep,
	}
}

func (u *usecase) SignUp(user *models.User) (*models.AuthToken, error) {
	_, err := u.userRep.GetUser(user.Id)
	if err != models.ErrNotFound && err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp error while search for user")
	} else if err == nil {
		return nil, models.ErrAlredyExists
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp bcrypt error")
	}

	user.Password = string(encPassword)
	_, err = u.userRep.AddUser(user)
	user.Password = ""

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp AddUser error")
	}

	jwtToken, err := jwt.NewToken(user, u.asyncKey, u.duration)

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignUp token generation error")
	}

	token := models.AuthToken{jwtToken}
	return &token, nil
}

func (u *usecase) SignIn(user *models.User) (*models.AuthToken, error) {
	repUser, err := u.userRep.GetUser(user.Id)

	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn user get error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(repUser.Password), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, models.ErrInvalidPassword
	} else if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn bcrypt error")
	}
	user.Password = ""

	jwtToken, err := jwt.NewToken(user, u.asyncKey, u.duration)
	if err != nil {
		return nil, errors.Wrap(err, "auth.usecase.SignIn token generation error")
	}

	token := models.AuthToken{jwtToken}
	return &token, nil
}

func (u *usecase) Authorization(token *models.AuthToken, role string) (bool, error) {
	panic("implement")
}
