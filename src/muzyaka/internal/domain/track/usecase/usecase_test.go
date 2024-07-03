package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/track/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_DeleteTrack(t *testing.T) {
	type mock func(r *mock_repository.MockTrackRepository, trackId uint64, track models.TrackMeta)
	type storageMock func(r *mock_repository.MockTrackStorage, tracks models.TrackMeta)

	testTable := []struct {
		name        string
		trackId     uint64
		inputTrack  models.TrackMeta
		mock        mock
		storageMock storageMock
		expectedErr error
	}{
		{
			name:    "Usual test",
			trackId: uint64(1),
			inputTrack: models.TrackMeta{
				Id:     10,
				Source: "test_src",
				Name:   "test_name",
				Genre:  "test_genre",
			},
			mock: func(r *mock_repository.MockTrackRepository, trackId uint64, track models.TrackMeta) {
				r.EXPECT().GetTrack(trackId).Return(&track, nil)
				r.EXPECT().DeleteTrack(trackId).Return(nil)
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackMeta) {
				r.EXPECT().DeleteObject(&track).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Repo fail test",
			trackId: uint64(2),
			inputTrack: models.TrackMeta{
				Id:     10,
				Source: "test_src",
				Name:   "test_name",
				Genre:  "test_genre",
			},
			mock: func(r *mock_repository.MockTrackRepository, trackId uint64, track models.TrackMeta) {
				r.EXPECT().GetTrack(trackId).Return(&track, nil)
				r.EXPECT().DeleteTrack(trackId).Return(errors.New("error in repo"))
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackMeta) {
			},
			expectedErr: errors.Wrap(
				errors.New("error in repo"),
				"track.usecase.DeleteTrack error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockTrackRepository(ctrl)
			tc.mock(repo, tc.trackId, tc.inputTrack)

			storage := mock_repository.NewMockTrackStorage(ctrl)
			tc.storageMock(storage, tc.inputTrack)

			s := NewTrackUseCase(repo, storage)
			err := s.DeleteTrack(tc.trackId)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_UpdatedTrack(t *testing.T) {
	type mock func(r *mock_repository.MockTrackRepository, track models.TrackObject)
	type storageMock func(r *mock_repository.MockTrackStorage, track models.TrackObject)

	testTable := []struct {
		name        string
		inputTrack  models.TrackObject
		mock        mock
		storageMock storageMock
		expectedErr error
	}{
		{
			name: "Usual test",
			inputTrack: models.TrackObject{
				TrackMeta: models.TrackMeta{
					Id:     1,
					Name:   "Updated TrackMeta Name",
					Source: "updated_source.mp3",
					Genre:  "Pop",
				},
				Payload: []byte{1, 2, 3},
			},
			mock: func(r *mock_repository.MockTrackRepository, track models.TrackObject) {
				r.EXPECT().UpdateTrack(track.ExtractMeta()).Return(nil)
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackObject) {
				r.EXPECT().UploadObject(&track).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputTrack: models.TrackObject{
				TrackMeta: models.TrackMeta{
					Id:     1,
					Name:   "Updated TrackMeta Name",
					Source: "updated_source.mp3",
					Genre:  "Pop",
				},
				Payload: []byte{1, 2, 3},
			},
			mock: func(r *mock_repository.MockTrackRepository, track models.TrackObject) {
				r.EXPECT().UpdateTrack(track.ExtractMeta()).Return(errors.New("error in repo"))
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackObject) {
				r.EXPECT().UploadObject(&track).Return(nil)
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "track.usecase.UpdateTrack error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockTrackRepository(ctrl)
			tc.mock(repo, tc.inputTrack)

			storage := mock_repository.NewMockTrackStorage(ctrl)
			tc.storageMock(storage, tc.inputTrack)

			u := NewTrackUseCase(repo, storage)
			err := u.UpdateTrack(&tc.inputTrack)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetTrack(t *testing.T) {
	type mock func(r *mock_repository.MockTrackRepository, id uint64, track models.TrackMeta)
	type storageMock func(r *mock_repository.MockTrackStorage, track models.TrackMeta)

	testTable := []struct {
		name          string
		id            uint64
		mock          mock
		storageMock   storageMock
		returnTrack   models.TrackMeta
		expectedTrack *models.TrackObject
		expectedErr   error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockTrackRepository, id uint64, track models.TrackMeta) {
				r.EXPECT().GetTrack(id).Return(&track, nil)
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackMeta) {
				ret := &models.TrackObject{
					TrackMeta: models.TrackMeta{
						Id:     1,
						Name:   "Test TrackMeta",
						Source: "test_source.mp3",
						Genre:  "Pop",
					},
					Payload: []byte{1, 2, 3},
				}

				r.EXPECT().LoadObject(&track).Return(ret, nil)
			},
			returnTrack: models.TrackMeta{
				Id:     1,
				Name:   "Test TrackMeta",
				Source: "test_source.mp3",
				Genre:  "Pop",
			},
			expectedTrack: &models.TrackObject{
				TrackMeta: models.TrackMeta{
					Id:     1,
					Name:   "Test TrackMeta",
					Source: "test_source.mp3",
					Genre:  "Pop",
				},
				Payload: []byte{1, 2, 3},
			},
			expectedErr: nil,
		},
		{
			name: "TrackMeta not found test",
			id:   2,
			mock: func(r *mock_repository.MockTrackRepository, id uint64, track models.TrackMeta) {
				r.EXPECT().GetTrack(id).Return(nil, errors.New("track not found"))
			},
			storageMock: func(r *mock_repository.MockTrackStorage, track models.TrackMeta) {

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
			tc.mock(repo, tc.id, tc.returnTrack)

			storage := mock_repository.NewMockTrackStorage(ctrl)
			tc.storageMock(storage, tc.returnTrack)

			u := NewTrackUseCase(repo, storage)
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
