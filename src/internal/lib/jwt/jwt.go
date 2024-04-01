package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"src/internal/models"
	"time"
)

func NewToken(user *models.User, key string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["sub"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "lib.jwt.NewToken error in sign")
	}

	return tokenString, nil
}
