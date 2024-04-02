package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/album/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_SignUp(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, id uint64)

	testTable := []struct {
		name          string
		input         uint64
		mock          mock
		expectedValue *models.Album
		expectedErr   error
	}{
		{
			name:  "Usual test",
			input: uint64(1),
			mock: func(r *mock_repository.MockAlbumRepository, id uint64) {
				r.EXPECT().GetAlbum(id).Return(&models.Album{
					Id:    1,
					Name:  "test_name",
					Cover: "test_cover",
					Type:  "test_type",
				}, nil)
			},
			expectedValue: &models.Album{
				Id:    1,
				Name:  "test_name",
				Cover: "test_cover",
				Type:  "test_type",
			},
			expectedErr: nil,
		},
		{
			name:  "Fail in repo test",
			input: uint64(110),
			mock: func(r *mock_repository.MockAlbumRepository, id uint64) {
				r.EXPECT().GetAlbum(id).Return(nil, errors.New("error in repo"))
			},
			expectedValue: nil,
			expectedErr:   errors.Wrap(errors.New("error in repo"), "album.usecase.GetAlbum error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockAlbumRepository(c)
			tc.mock(repo, tc.input)

			s := NewAlbumUseCase(repo)
			res, err := s.GetAlbum(tc.input)

			assert.Equal(t, tc.expectedValue, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
