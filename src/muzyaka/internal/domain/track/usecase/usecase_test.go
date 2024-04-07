package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/track/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_UpdatedTrack(t *testing.T) {
	type mock func(r *mock_repository.MockTrackRepository, track *models.Track)

	testTable := []struct {
		name        string
		inputTrack  *models.Track
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			inputTrack: &models.Track{
				Id:     1,
				Name:   "Updated Track Name",
				Source: "updated_source.mp3",
				Genre:  "Pop",
			},
			mock: func(r *mock_repository.MockTrackRepository, track *models.Track) {
				r.EXPECT().UpdateTrack(track).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputTrack: &models.Track{
				Id:     2,
				Name:   "Invalid Track",
				Source: "invalid_source.mp3",
				Genre:  "Rock",
			},
			mock: func(r *mock_repository.MockTrackRepository, track *models.Track) {
				r.EXPECT().UpdateTrack(track).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "track.usecase.UpdatedTrack error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockTrackRepository(ctrl)
			tc.mock(repo, tc.inputTrack)

			u := NewTrackUseCase(repo)
			err := u.UpdatedTrack(tc.inputTrack)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetTrack(t *testing.T) {
	type mock func(r *mock_repository.MockTrackRepository, id uint64)

	testTable := []struct {
		name          string
		id            uint64
		mock          mock
		expectedTrack *models.Track
		expectedErr   error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockTrackRepository, id uint64) {
				expectedTrack := &models.Track{
					Id:     1,
					Name:   "Test Track",
					Source: "test_source.mp3",
					Genre:  "Pop",
				}
				r.EXPECT().GetTrack(id).Return(expectedTrack, nil)
			},
			expectedTrack: &models.Track{
				Id:     1,
				Name:   "Test Track",
				Source: "test_source.mp3",
				Genre:  "Pop",
			},
			expectedErr: nil,
		},
		{
			name: "Track not found test",
			id:   2,
			mock: func(r *mock_repository.MockTrackRepository, id uint64) {
				r.EXPECT().GetTrack(id).Return(nil, errors.New("track not found"))
			},
			expectedTrack: nil,
			expectedErr:   errors.Wrap(errors.New("track not found"), "track.usecase.GetTrack error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockTrackRepository(ctrl)
			tc.mock(repo, tc.id)

			u := NewTrackUseCase(repo)
			track, err := u.GetTrack(tc.id)

			assert.Equal(t, tc.expectedTrack, track)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
