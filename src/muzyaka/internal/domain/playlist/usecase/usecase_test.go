package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/playlist/repository/mocks"
	mock_repository2 "src/internal/domain/track/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_UpdatedPlaylist(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist)

	testTable := []struct {
		name          string
		inputPlaylist *models.Playlist
		mock          mock
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputPlaylist: &models.Playlist{
				Id:          1,
				Name:        "Updated Playlist",
				CoverFile:   []byte("updated_cover.jpg"),
				Description: "Updated Description of Playlist",
			},
			mock: func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist) {
				r.EXPECT().UpdatePlaylist(playlist).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputPlaylist: &models.Playlist{
				Id:          2,
				Name:        "Invalid Playlist",
				CoverFile:   []byte("invalid_cover.jpg"),
				Description: "Invalid Description",
			},
			mock: func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist) {
				r.EXPECT().UpdatePlaylist(playlist).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.UpdatedPlaylist error while update"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.inputPlaylist)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			err := u.UpdatedPlaylist(tc.inputPlaylist)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_AddPlaylist(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist)

	testTable := []struct {
		name          string
		inputPlaylist *models.Playlist
		mock          mock
		expectedID    uint64
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputPlaylist: &models.Playlist{
				Name:        "New Playlist",
				CoverFile:   []byte("playlist_cover.jpg"),
				Description: "Description of New Playlist",
			},
			mock: func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist) {
				r.EXPECT().AddPlaylist(playlist, uint64(0)).Return(uint64(1), nil)
			},
			expectedID:  uint64(1),
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputPlaylist: &models.Playlist{
				Name:        "Invalid Playlist",
				CoverFile:   []byte("invalid_cover.jpg"),
				Description: "Invalid Description",
			},
			mock: func(r *mock_repository.MockPlaylistRepository, playlist *models.Playlist) {
				r.EXPECT().AddPlaylist(playlist, uint64(0)).Return(uint64(0), errors.New("error in repo"))
			},
			expectedID: uint64(0),
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.AddPlaylist error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.inputPlaylist)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			id, err := u.AddPlaylist(tc.inputPlaylist, 0)

			assert.Equal(t, tc.expectedID, id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeletePlaylist(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, id uint64)

	testTable := []struct {
		name        string
		id          uint64
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockPlaylistRepository, id uint64) {
				r.EXPECT().DeletePlaylist(id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			id:   2,
			mock: func(r *mock_repository.MockPlaylistRepository, id uint64) {
				r.EXPECT().DeletePlaylist(id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.DeletePlaylist error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.id)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			err := u.DeletePlaylist(tc.id)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetPlaylist(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, id uint64)

	testTable := []struct {
		name             string
		id               uint64
		mock             mock
		expectedPlaylist *models.Playlist
		expectedErr      error
	}{
		{
			name: "Usual test",
			id:   1,
			mock: func(r *mock_repository.MockPlaylistRepository, id uint64) {
				expectedPlaylist := &models.Playlist{
					Id:          1,
					Name:        "Test Playlist",
					CoverFile:   []byte("playlist_cover.jpg"),
					Description: "Description of Test Playlist",
				}
				r.EXPECT().GetPlaylist(id).Return(expectedPlaylist, nil)
			},
			expectedPlaylist: &models.Playlist{
				Id:          1,
				Name:        "Test Playlist",
				CoverFile:   []byte("playlist_cover.jpg"),
				Description: "Description of Test Playlist",
			},
			expectedErr: nil,
		},
		{
			name: "Playlist not found test",
			id:   2,
			mock: func(r *mock_repository.MockPlaylistRepository, id uint64) {
				r.EXPECT().GetPlaylist(id).Return(nil, errors.New("playlist not found"))
			},
			expectedPlaylist: nil,
			expectedErr: errors.Wrap(errors.New("playlist not found"),
				"playlist.usecase.GetPlaylist error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.id)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			playlist, err := u.GetPlaylist(tc.id)

			assert.Equal(t, tc.expectedPlaylist, playlist)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_AddTrack(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64)

	testTable := []struct {
		name        string
		playlistId  uint64
		trackId     uint64
		mock        mock
		expectedErr error
	}{
		{
			name:       "Usual test",
			playlistId: 1,
			trackId:    10,
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64) {
				r.EXPECT().AddTrackToPlaylist(playlistId, trackId).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Repo fail test",
			playlistId: 2,
			trackId:    20,
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64) {
				r.EXPECT().AddTrackToPlaylist(playlistId, trackId).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.AddTrack error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.playlistId, tc.trackId)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			err := u.AddTrack(tc.playlistId, tc.trackId)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeleteTrack(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64)

	testTable := []struct {
		name        string
		playlistId  uint64
		trackId     uint64
		mock        mock
		expectedErr error
	}{
		{
			name:       "Usual test",
			playlistId: 1,
			trackId:    10,
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64) {
				r.EXPECT().DeleteTrackFromPlaylist(playlistId, trackId).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Repo fail test",
			playlistId: 2,
			trackId:    20,
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, trackId uint64) {
				r.EXPECT().DeleteTrackFromPlaylist(playlistId, trackId).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.DeleteTrack error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.playlistId, tc.trackId)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)

			u := NewPlaylistUseCase(repo, trackRepo)
			err := u.DeleteTrack(tc.playlistId, tc.trackId)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetAllTracks(t *testing.T) {
	type mock func(r *mock_repository.MockPlaylistRepository, playlistId uint64, tracks []uint64)
	type tracksMock func(r *mock_repository2.MockTrackRepository, tracks []*models.TrackMeta)

	testTable := []struct {
		name           string
		playlistId     uint64
		returnIds      []uint64
		mock           mock
		tracksMock     tracksMock
		expectedTracks []*models.TrackMeta
		expectedErr    error
	}{
		{
			name:       "Usual test",
			playlistId: 1,
			returnIds:  []uint64{1, 2},
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, tracks []uint64) {
				r.EXPECT().GetAllTracks(playlistId).Return(tracks, nil)
			},
			tracksMock: func(r *mock_repository2.MockTrackRepository, tracks []*models.TrackMeta) {
				for _, v := range tracks {
					r.EXPECT().GetTrack(v.Id).Return(v, nil)
				}
			},
			expectedTracks: []*models.TrackMeta{
				{
					Id:     1,
					Source: "track_src_1",
					Name:   "track_name_1",
					Genre:  "track_genre_1",
				},
				{
					Id:     2,
					Source: "track_src_2",
					Name:   "track_name_2",
					Genre:  "track_genre_2",
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Repo fail test",
			playlistId: 2,
			mock: func(r *mock_repository.MockPlaylistRepository, playlistId uint64, tracks []uint64) {
				r.EXPECT().GetAllTracks(playlistId).Return(nil, errors.New("error in repo"))
			},
			tracksMock: func(r *mock_repository2.MockTrackRepository, tracks []*models.TrackMeta) {
			},
			expectedTracks: nil,
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"playlist.usecase.GetAllTracks error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockPlaylistRepository(ctrl)
			tc.mock(repo, tc.playlistId, tc.returnIds)

			trackRepo := mock_repository2.NewMockTrackRepository(ctrl)
			tc.tracksMock(trackRepo, tc.expectedTracks)

			u := NewPlaylistUseCase(repo, trackRepo)
			tracks, err := u.GetAllTracks(tc.playlistId)

			assert.Equal(t, tc.expectedTracks, tracks)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
