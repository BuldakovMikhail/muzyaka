package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"src/internal/models"
	"time"
)

//go:generate mockgen -source=jwt.go -destination=mocks/mock.go

type TokenProvider interface {
	GenerateToken(user *models.User) (*models.AuthToken, error)
	IsTokenValid(token *models.AuthToken) (bool, error)
	GetRole(token *models.AuthToken) (string, error)
	GetId(token *models.AuthToken) (uint64, error)
}

type Payload struct {
	uid  uint64
	role string
	exp  time.Duration
}

type tokenProvider struct {
	asyncKey string
	duration time.Duration
}

func NewTokenProvider(key string, dur time.Duration) TokenProvider {
	return &tokenProvider{asyncKey: key, duration: dur}
}

func (t *tokenProvider) GenerateToken(user *models.User) (*models.AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(t.duration).Unix()

	tokenString, err := token.SignedString([]byte(t.asyncKey))
	if err != nil {
		return nil, errors.Wrap(err, "auth.tokenhelper.GenerateToken error in sign")
	}

	return &models.AuthToken{Secret: []byte(tokenString)}, nil
}

func (t *tokenProvider) IsTokenValid(token *models.AuthToken) (bool, error) {
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(string(token.Secret), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.asyncKey), nil
	})
	if err != nil {
		return false, errors.Wrap(err, "auth.tokenhelper.GetRole error in parse")
	}

	return jwtToken.Valid, nil
}

func (t *tokenProvider) GetRole(token *models.AuthToken) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(string(token.Secret), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.asyncKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "auth.tokenhelper.GetRole error in parse")
	}

	return claims["role"].(string), nil
}

func (t *tokenProvider) GetId(token *models.AuthToken) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(string(token.Secret), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.asyncKey), nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "auth.tokenhelper.GetId error in parse")
	}

	return uint64(claims["uid"].(float64)), nil
}
