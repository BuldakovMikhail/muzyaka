package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/merch/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_GetMerch(t *testing.T) {
	type mock func(r *mock_repository.MockMerchRepository, id uint64, merch *models.Merch)

	testTable := []struct {
		name          string
		id            uint64
		mock          mock
		expectedMerch *models.Merch
		expectedErr   error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockMerchRepository, id uint64, merch *models.Merch) {
				r.EXPECT().GetMerch(id).Return(merch, nil)
			},
			expectedMerch: &models.Merch{
				Id:          1,
				Name:        "Test Merch",
				PhotoFiles:  [][]byte{[]byte("photo1.jpg"), []byte("photo2.jpg")},
				Description: "Description of Test Merch",
				OrderUrl:    "http://example.com/order",
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   2,
			mock: func(r *mock_repository.MockMerchRepository, id uint64, merch *models.Merch) {
				r.EXPECT().GetMerch(id).Return(nil, errors.New("error in repo"))
			},
			expectedMerch: nil,
			expectedErr:   errors.Wrap(errors.New("error in repo"), "merch.usecase.GetMerch error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMerchRepository(ctrl)
			tc.mock(repo, tc.id, tc.expectedMerch)

			u := NewMerchUseCase(repo)
			merch, err := u.GetMerch(tc.id)

			assert.Equal(t, tc.expectedMerch, merch)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_UpdateMerch(t *testing.T) {
	type mock func(r *mock_repository.MockMerchRepository, merch *models.Merch)

	testTable := []struct {
		name        string
		inputMerch  *models.Merch
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			inputMerch: &models.Merch{
				Id:          1,
				Name:        "Updated Merch",
				PhotoFiles:  [][]byte{[]byte("updated_photo.jpg")},
				Description: "Updated Description of Merch",
				OrderUrl:    "http://example.com/update",
			},
			mock: func(r *mock_repository.MockMerchRepository, merch *models.Merch) {
				r.EXPECT().UpdateMerch(merch).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputMerch: &models.Merch{
				Id:          2,
				Name:        "Invalid Merch",
				PhotoFiles:  nil,
				Description: "Invalid Description",
				OrderUrl:    "http://example.com/update",
			},
			mock: func(r *mock_repository.MockMerchRepository, merch *models.Merch) {
				r.EXPECT().UpdateMerch(merch).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"merch.usecase.UpdateMerch error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMerchRepository(ctrl)
			tc.mock(repo, tc.inputMerch)

			u := NewMerchUseCase(repo)
			err := u.UpdateMerch(tc.inputMerch)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_AddMerch(t *testing.T) {
	type mock func(r *mock_repository.MockMerchRepository, merch *models.Merch)

	testTable := []struct {
		name          string
		inputMerch    *models.Merch
		mock          mock
		expectedValue uint64
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputMerch: &models.Merch{
				Name:        "New Merch",
				PhotoFiles:  [][]byte{[]byte("photo1.jpg"), []byte("photo2.jpg")},
				Description: "Description of New Merch",
				OrderUrl:    "http://example.com/order",
			},
			mock: func(r *mock_repository.MockMerchRepository, merch *models.Merch) {
				r.EXPECT().AddMerch(merch).Return(uint64(1), nil)
			},
			expectedValue: uint64(1),
			expectedErr:   nil,
		},
		{
			name: "Repo fail test",
			inputMerch: &models.Merch{
				Name:        "Invalid Merch",
				PhotoFiles:  nil,
				Description: "Invalid Description",
				OrderUrl:    "http://example.com/order",
			},
			mock: func(r *mock_repository.MockMerchRepository, merch *models.Merch) {
				r.EXPECT().AddMerch(merch).Return(uint64(0), errors.New("error in repo"))
			},
			expectedValue: uint64(0),
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"merch.usecase.AddMerch error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMerchRepository(ctrl)
			tc.mock(repo, tc.inputMerch)

			u := NewMerchUseCase(repo)
			res, err := u.AddMerch(tc.inputMerch)

			assert.Equal(t, tc.expectedValue, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeleteMerch(t *testing.T) {
	type mock func(r *mock_repository.MockMerchRepository, id uint64)

	testTable := []struct {
		name        string
		id          uint64
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockMerchRepository, id uint64) {
				r.EXPECT().DeleteMerch(id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   2,
			mock: func(r *mock_repository.MockMerchRepository, id uint64) {
				r.EXPECT().DeleteMerch(id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"merch.usecase.DeleteMerch error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockMerchRepository(ctrl)
			tc.mock(repo, tc.id)

			u := NewMerchUseCase(repo)
			err := u.DeleteMerch(tc.id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
