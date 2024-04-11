package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_usecase "src/internal/domain/auth/usecase/mocks"
	mock_repository3 "src/internal/domain/user/repository/mocks"
	mock_jwt "src/internal/lib/jwt/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_SignUp(t *testing.T) {
	type mockUser func(r *mock_repository3.MockUserRepository, user *models.User)
	type mockToken func(r *mock_jwt.MockTokenProvider, user *models.User)
	type mockEnc func(r *mock_usecase.MockEncryptor, password []byte)

	testTable := []struct {
		name          string
		input         *models.User
		mockUser      mockUser
		mockToken     mockToken
		mockEnc       mockEnc
		expectedValue *models.AuthToken
		expectedErr   error
	}{
		{
			name: "Usual test",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, user *models.User) {
				r.EXPECT().AddUser(user).Return(uint64(1), nil)
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
				r.EXPECT().GenerateToken(user).Return(&models.AuthToken{Secret: []byte("aboba")}, nil)
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, password []byte) {
				r.EXPECT().EncodePassword(password).Return(password, nil)
			},
			expectedValue: &models.AuthToken{Secret: []byte("aboba")},
			expectedErr:   nil,
		},
		{
			name: "Add Error",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, user *models.User) {
				r.EXPECT().AddUser(user).Return(uint64(0), errors.New("repo error"))
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, password []byte) {
				r.EXPECT().EncodePassword(password).Return(password, nil)
			},
			expectedValue: nil,
			expectedErr: errors.Wrap(errors.New("repo error"),
				"auth.usecase.SignUp AddUser error"),
		},
		{
			name: "Token Provider error",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, user *models.User) {
				r.EXPECT().AddUser(user).Return(uint64(1), nil)
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
				r.EXPECT().GenerateToken(user).Return(nil, errors.New("error in token provider"))
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, password []byte) {
				r.EXPECT().EncodePassword(password).Return(password, nil)
			},
			expectedValue: nil,
			expectedErr: errors.Wrap(errors.New("error in token provider"),
				"auth.usecase.SignUp token generation error"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			c2 := gomock.NewController(t)
			defer c2.Finish()

			c3 := gomock.NewController(t)
			defer c3.Finish()

			repoUser := mock_repository3.NewMockUserRepository(c)
			tokenMock := mock_jwt.NewMockTokenProvider(c2)
			dummyEnc := mock_usecase.NewMockEncryptor(c3)

			tc.mockUser(repoUser, tc.input)
			tc.mockToken(tokenMock, tc.input)
			tc.mockEnc(dummyEnc, []byte(tc.input.Password))

			s := NewAuthUseCase(tokenMock, repoUser, dummyEnc)

			res, err := s.SignUp(tc.input)

			assert.Equal(t, tc.expectedValue, res)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.input.Password, "")
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_SignIn(t *testing.T) {
	type mockUser func(r *mock_repository3.MockUserRepository, id uint64)
	type mockToken func(r *mock_jwt.MockTokenProvider, user *models.User)
	type mockEnc func(r *mock_usecase.MockEncryptor, hashedPass []byte, password []byte, retVal error)

	testTable := []struct {
		name          string
		input         *models.User
		mockUser      mockUser
		mockToken     mockToken
		mockEnc       mockEnc
		compRes       error
		expectedValue *models.AuthToken
		expectedErr   error
	}{
		{
			name: "Usual test",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, id uint64) {
				r.EXPECT().GetUser(id).Return(&models.User{
					Id:       uint64(10),
					Name:     "test",
					Password: "test",
					Role:     "test",
					Email:    "test",
				}, nil)
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
				r.EXPECT().GenerateToken(user).Return(&models.AuthToken{Secret: []byte("aboba")}, nil)
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, hashedPass []byte, password []byte, retVal error) {
				r.EXPECT().CompareHashAndPassword(hashedPass, password).Return(retVal)
			},
			compRes:       nil,
			expectedValue: &models.AuthToken{Secret: []byte("aboba")},
			expectedErr:   nil,
		},
		{
			name: "Get user fault",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, id uint64) {
				r.EXPECT().GetUser(id).Return(nil, errors.New("repo error"))
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, hashedPass []byte, password []byte, retVal error) {
			},
			compRes:       nil,
			expectedValue: nil,
			expectedErr: errors.Wrap(errors.New("repo error"),
				"auth.usecase.SignIn user get error"),
		},
		{
			name: "Compare fault",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, id uint64) {
				r.EXPECT().GetUser(id).Return(&models.User{
					Id:       uint64(10),
					Name:     "test",
					Password: "test",
					Role:     "test",
					Email:    "test",
				}, nil)
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, hashedPass []byte, password []byte, retVal error) {
				r.EXPECT().CompareHashAndPassword(hashedPass, password).Return(retVal)
			},
			compRes:       errors.New("comp error"),
			expectedValue: nil,
			expectedErr: errors.Wrap(errors.New("comp error"),
				"auth.usecase.SignIn compare error"),
		},
		{
			name: "Generate token fault",
			input: &models.User{
				Id:       uint64(10),
				Name:     "test",
				Password: "test",
				Role:     "test",
				Email:    "test",
			},
			mockUser: func(r *mock_repository3.MockUserRepository, id uint64) {
				r.EXPECT().GetUser(id).Return(&models.User{
					Id:       uint64(10),
					Name:     "test",
					Password: "test",
					Role:     "test",
					Email:    "test",
				}, nil)
			},
			mockToken: func(r *mock_jwt.MockTokenProvider, user *models.User) {
				r.EXPECT().GenerateToken(user).Return(nil, errors.New("token error"))
			},
			mockEnc: func(r *mock_usecase.MockEncryptor, hashedPass []byte, password []byte, retVal error) {
				r.EXPECT().CompareHashAndPassword(hashedPass, password).Return(retVal)
			},
			compRes:       nil,
			expectedValue: nil,
			expectedErr:   errors.Wrap(errors.New("token error"), "auth.usecase.SignIn token generation error"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			c2 := gomock.NewController(t)
			defer c2.Finish()

			c3 := gomock.NewController(t)
			defer c3.Finish()

			repoUser := mock_repository3.NewMockUserRepository(c)
			tokenMock := mock_jwt.NewMockTokenProvider(c2)
			dummyEnc := mock_usecase.NewMockEncryptor(c3)

			tc.mockUser(repoUser, tc.input.Id)
			tc.mockToken(tokenMock, tc.input)
			tc.mockEnc(dummyEnc, []byte(tc.input.Password), []byte(tc.input.Password), tc.compRes)

			s := NewAuthUseCase(tokenMock, repoUser, dummyEnc)

			res, err := s.SignIn(tc.input)

			assert.Equal(t, tc.expectedValue, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.input.Password, "")
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_Authorization(t *testing.T) {
	type mockToken func(r *mock_jwt.MockTokenProvider, user *models.AuthToken)

	testTable := []struct {
		name          string
		input         *models.AuthToken
		role          string
		mockToken     mockToken
		expectedValue bool
		expectedErr   error
	}{
		{
			name:  "Usual test",
			input: &models.AuthToken{Secret: []byte("aboba")},
			role:  "aboba",
			mockToken: func(r *mock_jwt.MockTokenProvider, token *models.AuthToken) {
				r.EXPECT().GetRole(token).Return("aboba", nil)
			},
			expectedValue: true,
			expectedErr:   nil,
		},
		{
			name:  "Access denied",
			input: &models.AuthToken{Secret: []byte("aboba")},
			role:  "ne aboba",
			mockToken: func(r *mock_jwt.MockTokenProvider, token *models.AuthToken) {
				r.EXPECT().GetRole(token).Return("aboba", nil)
			},
			expectedValue: false,
			expectedErr:   nil,
		},
		{
			name:  "Token fault",
			input: &models.AuthToken{Secret: []byte("aboba")},
			role:  "aboba",
			mockToken: func(r *mock_jwt.MockTokenProvider, token *models.AuthToken) {
				r.EXPECT().GetRole(token).Return("", errors.New("token error"))
			},
			expectedValue: false,
			expectedErr:   errors.Wrap(errors.New("token error"), "auth.usecase.Authorization token parse error"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			c2 := gomock.NewController(t)
			defer c2.Finish()

			c3 := gomock.NewController(t)
			defer c3.Finish()

			repoUser := mock_repository3.NewMockUserRepository(c)
			tokenMock := mock_jwt.NewMockTokenProvider(c2)
			dummyEnc := mock_usecase.NewMockEncryptor(c3)

			tc.mockToken(tokenMock, tc.input)

			s := NewAuthUseCase(tokenMock, repoUser, dummyEnc)

			res, err := s.Authorization(tc.input, tc.role)

			assert.Equal(t, tc.expectedValue, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
