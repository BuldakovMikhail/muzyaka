package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/musician/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_UpdateMusician(t *testing.T) {
	type mock func(r *mock_repository.MockMusicianRepository, musician *models.Musician)

	testTable := []struct {
		name          string
		inputMusician *models.Musician
		mock          mock
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputMusician: &models.Musician{
				Id:          1,
				Name:        "Updated Musician",
				PhotoFiles:  [][]byte{[]byte("updated_photo.jpg")},
				Description: "Updated Description of Musician",
			},
			mock: func(r *mock_repository.MockMusicianRepository, musician *models.Musician) {
				r.EXPECT().UpdateMusician(musician).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputMusician: &models.Musician{
				Id:          2,
				Name:        "Invalid Musician",
				PhotoFiles:  nil,
				Description: "Invalid Description",
			},
			mock: func(r *mock_repository.MockMusicianRepository, musician *models.Musician) {
				r.EXPECT().UpdateMusician(musician).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "musician.usecase.UpdatedMusician error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMusicianRepository(ctrl)
			tc.mock(repo, tc.inputMusician)

			u := NewMusicianUseCase(repo)
			err := u.UpdatedMusician(tc.inputMusician)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_AddMusician(t *testing.T) {
	type mock func(r *mock_repository.MockMusicianRepository, musician *models.Musician)

	testTable := []struct {
		name          string
		inputMusician *models.Musician
		mock          mock
		expectedID    uint64
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputMusician: &models.Musician{
				Name:        "John Doe",
				PhotoFiles:  [][]byte{[]byte("photo1.jpg"), []byte("photo2.jpg")},
				Description: "Description of John Doe",
			},
			mock: func(r *mock_repository.MockMusicianRepository, musician *models.Musician) {
				r.EXPECT().AddMusician(musician).Return(uint64(1), nil)
			},
			expectedID:  uint64(1),
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputMusician: &models.Musician{
				Name:        "Invalid Musician",
				PhotoFiles:  nil,
				Description: "Invalid Description",
			},
			mock: func(r *mock_repository.MockMusicianRepository, musician *models.Musician) {
				r.EXPECT().AddMusician(musician).Return(uint64(0), errors.New("error in repo"))
			},
			expectedID: uint64(0),
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"musician.usecase.AddMusician error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMusicianRepository(ctrl)
			tc.mock(repo, tc.inputMusician)

			u := NewMusicianUseCase(repo)
			id, err := u.AddMusician(tc.inputMusician)

			assert.Equal(t, tc.expectedID, id)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeleteMusician(t *testing.T) {
	type mock func(r *mock_repository.MockMusicianRepository, id uint64)

	testTable := []struct {
		name        string
		id          uint64
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockMusicianRepository, id uint64) {
				r.EXPECT().DeleteMusician(id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   2,
			mock: func(r *mock_repository.MockMusicianRepository, id uint64) {
				r.EXPECT().DeleteMusician(id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"musician.usecase.DeleteMusician error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMusicianRepository(ctrl)
			tc.mock(repo, tc.id)

			u := NewMusicianUseCase(repo)
			err := u.DeleteMusician(tc.id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetMusician(t *testing.T) {
	type mock func(r *mock_repository.MockMusicianRepository, id uint64, musician *models.Musician)

	testTable := []struct {
		name             string
		id               uint64
		mock             mock
		expectedMusician *models.Musician
		expectedErr      error
	}{
		{
			name: "Usual test",
			id:   uint64(1),
			mock: func(r *mock_repository.MockMusicianRepository, id uint64, musician *models.Musician) {
				r.EXPECT().GetMusician(id).Return(musician, nil)
			},
			expectedMusician: &models.Musician{
				Id:          uint64(1),
				Name:        "John Doe",
				PhotoFiles:  [][]byte{[]byte("photo1.jpg"), []byte("photo2.jpg")},
				Description: "Description of John Doe",
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   uint64(2),
			mock: func(r *mock_repository.MockMusicianRepository, id uint64, musician *models.Musician) {
				r.EXPECT().GetMusician(id).Return(nil, errors.New("error in repo"))
			},
			expectedMusician: nil,
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"musician.usecase.GetMusician error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMusicianRepository(ctrl)
			tc.mock(repo, tc.id, tc.expectedMusician)

			u := NewMusicianUseCase(repo)
			res, err := u.GetMusician(tc.id)

			assert.Equal(t, tc.expectedMusician, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
