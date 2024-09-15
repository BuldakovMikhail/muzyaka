package usecase

import "golang.org/x/crypto/bcrypt"

//go:generate mockgen -source=encryptor.go -destination=mocks/mock.go

type Encryptor interface {
	EncodePassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type encryptor struct{}

func NewEncryptor() Encryptor {
	return &encryptor{}
}

func (e *encryptor) EncodePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (e *encryptor) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
