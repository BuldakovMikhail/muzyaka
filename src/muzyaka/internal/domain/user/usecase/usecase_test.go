package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/user/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_UpdateUser(t *testing.T) {
	type mock func(r *mock_repository.MockUserRepository, user *models.User)

	testTable := []struct {
		name        string
		inputUser   *models.User
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			inputUser: &models.User{
				Id:       1,
				Name:     "Updated User Name",
				Password: "newpassword",
				Role:     "admin",
				Email:    "updatedemail@example.com",
			},
			mock: func(r *mock_repository.MockUserRepository, user *models.User) {
				r.EXPECT().UpdateUser(user).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputUser: &models.User{
				Id:       2,
				Name:     "Invalid User",
				Password: "invalidpassword",
				Role:     "user",
				Email:    "invalidemail@example.com",
			},
			mock: func(r *mock_repository.MockUserRepository, user *models.User) {
				r.EXPECT().UpdateUser(user).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"user.usecase.UpdateUser error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			tc.mock(repo, tc.inputUser)

			u := NewUserUseCase(repo)
			err := u.UpdateUser(tc.inputUser)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetUser(t *testing.T) {
	type mock func(r *mock_repository.MockUserRepository, id uint64)

	testTable := []struct {
		name         string
		id           uint64
		mock         mock
		expectedUser *models.User
		expectedErr  error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockUserRepository, id uint64) {
				expectedUser := &models.User{
					Id:       1,
					Name:     "Test User",
					Password: "testpassword",
					Role:     "user",
					Email:    "testemail@example.com",
				}
				r.EXPECT().GetUser(id).Return(expectedUser, nil)
			},
			expectedUser: &models.User{
				Id:       1,
				Name:     "Test User",
				Password: "testpassword",
				Role:     "user",
				Email:    "testemail@example.com",
			},
			expectedErr: nil,
		},
		{
			name: "User not found test",
			id:   2,
			mock: func(r *mock_repository.MockUserRepository, id uint64) {
				r.EXPECT().GetUser(id).Return(nil, errors.New("user not found"))
			},
			expectedUser: nil,
			expectedErr:  errors.Wrap(errors.New("user not found"), "user.usecase.GetUser error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			tc.mock(repo, tc.id)

			u := NewUserUseCase(repo)
			user, err := u.GetUser(tc.id)

			assert.Equal(t, tc.expectedUser, user)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
func TestUsecase_AddUser(t *testing.T) {
	type mock func(r *mock_repository.MockUserRepository, user *models.User)

	testTable := []struct {
		name        string
		inputUser   *models.User
		mock        mock
		expectedID  uint64
		expectedErr error
	}{
		{
			name: "Usual test",
			inputUser: &models.User{
				Id:       1,
				Name:     "New User",
				Password: "newuserpassword",
				Role:     "user",
				Email:    "newuser@example.com",
			},
			mock: func(r *mock_repository.MockUserRepository, user *models.User) {
				r.EXPECT().AddUser(user).Return(uint64(1), nil)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputUser: &models.User{
				Id:       0, // Zero ID indicates new user
				Name:     "Invalid User",
				Password: "invalidpassword",
				Role:     "user",
				Email:    "invalidemail@example.com",
			},
			mock: func(r *mock_repository.MockUserRepository, user *models.User) {
				r.EXPECT().AddUser(user).Return(uint64(0), errors.New("error in repo"))
			},
			expectedID:  0,
			expectedErr: errors.Wrap(errors.New("error in repo"), "user.usecase.AddUser error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			tc.mock(repo, tc.inputUser)

			u := NewUserUseCase(repo)
			id, err := u.AddUser(tc.inputUser)

			assert.Equal(t, tc.expectedID, id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeleteUser(t *testing.T) {
	type mock func(r *mock_repository.MockUserRepository, id uint64)

	testTable := []struct {
		name        string
		id          uint64
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockUserRepository, id uint64) {
				r.EXPECT().DeleteUser(id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   2,
			mock: func(r *mock_repository.MockUserRepository, id uint64) {
				r.EXPECT().DeleteUser(id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "user.usecase.DeleteUser error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			tc.mock(repo, tc.id)

			u := NewUserUseCase(repo)
			err := u.DeleteUser(tc.id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
