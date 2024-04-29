package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_remote "src/internal/domain/recsys/recsys_client/mocks"
	mock_repository "src/internal/domain/track/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestRecSysUseCase_GetSameTracks(t *testing.T) {
	type mock func(r *mock_remote.MockRecSysProvider, r2 *mock_repository.MockTrackRepository, id uint64)

	testTable := []struct {
		name           string
		id             uint64
		mock           mock
		expectedTracks []*models.TrackMeta
		expectedErr    error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_remote.MockRecSysProvider, r2 *mock_repository.MockTrackRepository, id uint64) {
				r.EXPECT().GetRecs(id).Return([]uint64{1, 2, 3}, nil)

				track1 := &models.TrackMeta{Id: 1, Name: "TrackMeta 1"}
				track2 := &models.TrackMeta{Id: 2, Name: "TrackMeta 2"}
				track3 := &models.TrackMeta{Id: 3, Name: "TrackMeta 3"}

				r2.EXPECT().GetTrack(uint64(1)).Return(track1, nil)
				r2.EXPECT().GetTrack(uint64(2)).Return(track2, nil)
				r2.EXPECT().GetTrack(uint64(3)).Return(track3, nil)
			},
			expectedTracks: []*models.TrackMeta{
				{Id: 1, Name: "TrackMeta 1"},
				{Id: 2, Name: "TrackMeta 2"},
				{Id: 3, Name: "TrackMeta 3"},
			},
			expectedErr: nil,
		},
		{
			name: "GetRecs fails",
			id:   2,
			mock: func(r *mock_remote.MockRecSysProvider, r2 *mock_repository.MockTrackRepository, id uint64) {
				r.EXPECT().GetRecs(id).Return(nil, errors.New("error in GetRecs call"))
			},
			expectedTracks: nil,
			expectedErr:    errors.Wrap(errors.New("error in GetRecs call"), "recsys.usecase.GetSameTracks error while GetRecs call"),
		},
		{
			name: "GetTrack fails",
			id:   3,
			mock: func(r *mock_remote.MockRecSysProvider, r2 *mock_repository.MockTrackRepository, id uint64) {
				r.EXPECT().GetRecs(id).Return([]uint64{1, 2, 3}, nil)

				track1 := &models.TrackMeta{Id: 1, Name: "TrackMeta 1"}
				track2 := &models.TrackMeta{Id: 2, Name: "TrackMeta 2"}

				r2.EXPECT().GetTrack(uint64(1)).Return(track1, nil)
				r2.EXPECT().GetTrack(uint64(2)).Return(track2, nil)
				r2.EXPECT().GetTrack(uint64(3)).Return(nil, errors.New("error in GetTrack call"))
			},
			expectedTracks: nil,
			expectedErr:    errors.Wrap(errors.New("error in GetTrack call"), "recsys.usecase.GetSameTracks error while trackRep call"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recsProvider := mock_remote.NewMockRecSysProvider(ctrl)
			trackRep := mock_repository.NewMockTrackRepository(ctrl)
			tc.mock(recsProvider, trackRep, tc.id)

			u := NewRecSysUseCase(recsProvider, trackRep)
			tracks, err := u.GetSameTracks(tc.id)

			assert.Equal(t, tc.expectedTracks, tracks)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
