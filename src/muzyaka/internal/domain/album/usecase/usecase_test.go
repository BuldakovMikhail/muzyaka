package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/album/repository/mocks"
	"src/internal/models"
	"testing"
)

func TestUsecase_GetAlbum(t *testing.T) {
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

func TestUsecase_UpdateAlbum(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album models.Album)

	testTable := []struct {
		name        string
		input       models.Album
		mock        mock
		expectedErr error
	}{
		{
			name: "Usual test",
			input: models.Album{
				Id:    1,
				Name:  "test_name",
				Cover: "test_cover",
				Type:  "test_type",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album models.Album) {
				r.EXPECT().UpdateAlbum(&album).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Fail in repo",
			input: models.Album{
				Id:    1,
				Name:  "test_name",
				Cover: "test_cover",
				Type:  "test_type",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album models.Album) {
				r.EXPECT().UpdateAlbum(&album).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "album.usecase.UpdateAlbum error while update"),
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
			err := s.UpdateAlbum(&tc.input)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUseCase_AddAlbumWithTracks(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.Track)

	testTable := []struct {
		name        string
		inputAlbum  *models.Album
		inputTracks []*models.Track
		mock        mock
		expectedID  uint64
		expectedErr error
	}{
		{
			name: "Usual test",
			inputAlbum: &models.Album{
				Id:    1,
				Name:  "Test Album",
				Cover: "Test Cover",
				Type:  "LP",
			},
			inputTracks: []*models.Track{
				{Id: 1, Name: "Track 1"},
				{Id: 2, Name: "Track 2"},
			},
			mock: func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.Track) {
				r.EXPECT().AddAlbumWithTracks(album, tracks).Return(uint64(1), nil)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputAlbum: &models.Album{
				Id:    0,
				Name:  "Invalid Album",
				Cover: "Test Cover",
				Type:  "LP",
			},
			inputTracks: []*models.Track{
				{Id: 1, Name: "Track 1"},
				{Id: 2, Name: "Track 2"},
			},
			mock: func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.Track) {
				r.EXPECT().AddAlbumWithTracks(album, tracks).Return(uint64(0), errors.New("error in repo"))
			},
			expectedID:  0,
			expectedErr: errors.Wrap(errors.New("error in repo"), "album.usecase.AddAlbum error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockAlbumRepository(ctrl)
			tc.mock(repo, tc.inputAlbum, tc.inputTracks)

			u := NewAlbumUseCase(repo)
			id, err := u.AddAlbumWithTracks(tc.inputAlbum, tc.inputTracks)

			assert.Equal(t, tc.expectedID, id)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsecase_DeleteAlbum(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, id uint64)

	testTable := []struct {
		name        string
		input       uint64
		mock        mock
		expectedErr error
	}{
		{
			name:  "Usual test",
			input: uint64(1),
			mock: func(r *mock_repository.MockAlbumRepository, id uint64) {
				r.EXPECT().DeleteAlbum(id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Fail in repo test",
			input: uint64(110),
			mock: func(r *mock_repository.MockAlbumRepository, id uint64) {
				r.EXPECT().DeleteAlbum(id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "album.usecase.DeleteAlbum error while delete"),
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
			err := s.DeleteAlbum(tc.input)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_AddTrack(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.Track)

	testTable := []struct {
		name          string
		inputId       uint64
		inputTrack    models.Track
		mock          mock
		expectedValue uint64
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputTrack: models.Track{
				Id:         10,
				Source:     "test_src",
				Producers:  nil,
				Authors:    nil,
				Performers: nil,
				Name:       "test_name",
				Genre:      "test_genre",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.Track) {
				r.EXPECT().AddTrackToAlbum(album_id, &track).Return(uint64(10), nil)
			},
			expectedValue: uint64(10),
			expectedErr:   nil,
		},
		{
			name: "Repo fail test",
			inputTrack: models.Track{
				Id:         10,
				Source:     "test_src",
				Producers:  nil,
				Authors:    nil,
				Performers: nil,
				Name:       "test_name",
				Genre:      "test_genre",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.Track) {
				r.EXPECT().AddTrackToAlbum(album_id, &track).Return(uint64(0), errors.New("error in repo"))
			},
			expectedValue: uint64(0),
			expectedErr:   errors.Wrap(errors.New("error in repo"), "album.usecase.AddTrack error while add"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockAlbumRepository(c)
			tc.mock(repo, tc.inputId, tc.inputTrack)

			s := NewAlbumUseCase(repo)
			res, err := s.AddTrack(tc.inputId, &tc.inputTrack)

			assert.Equal(t, tc.expectedValue, res)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_DeleteTrack(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, track_id uint64)

	testTable := []struct {
		name        string
		albumId     uint64
		trackId     uint64
		mock        mock
		expectedErr error
	}{
		{
			name:    "Usual test",
			albumId: uint64(1),
			trackId: uint64(10),
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track_id uint64) {
				r.EXPECT().DeleteTrackFromAlbum(album_id, track_id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Repo fail test",
			albumId: uint64(2),
			trackId: uint64(20),
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track_id uint64) {
				r.EXPECT().DeleteTrackFromAlbum(album_id, track_id).Return(errors.New("error in repo"))
			},
			expectedErr: errors.Wrap(
				errors.New("error in repo"),
				"album.usecase.DeleteTrack error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockAlbumRepository(ctrl)
			tc.mock(repo, tc.albumId, tc.trackId)

			s := NewAlbumUseCase(repo)
			err := s.DeleteTrack(tc.albumId, tc.trackId)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetAllTracks(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.Track)

	testTable := []struct {
		name           string
		albumId        uint64
		mock           mock
		expectedTracks []*models.Track
		expectedErr    error
	}{
		{
			name:    "Usual test",
			albumId: 1,
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.Track) {
				r.EXPECT().GetAllTracksForAlbum(album_id).Return(tracks, nil)
			},
			expectedTracks: []*models.Track{
				{
					Id:         1,
					Source:     "track_src_1",
					Producers:  []string{"producer_1"},
					Authors:    []string{"author_1"},
					Performers: []string{"performer_1"},
					Name:       "track_name_1",
					Genre:      "track_genre_1",
				},
				{
					Id:         2,
					Source:     "track_src_2",
					Producers:  []string{"producer_2"},
					Authors:    []string{"author_2"},
					Performers: []string{"performer_2"},
					Name:       "track_name_2",
					Genre:      "track_genre_2",
				},
			},
			expectedErr: nil,
		},
		{
			name:    "Repo fail test",
			albumId: 2,
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.Track) {
				r.EXPECT().GetAllTracksForAlbum(album_id).Return(nil, errors.New("error in repo"))
			},
			expectedTracks: nil,
			expectedErr: errors.Wrap(errors.New("error in repo"),
				"album.usecase.GetAllTracks error while get"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockAlbumRepository(ctrl)
			tc.mock(repo, tc.albumId, tc.expectedTracks)

			u := NewAlbumUseCase(repo)
			tracks, err := u.GetAllTracks(tc.albumId)

			assert.Equal(t, tc.expectedTracks, tracks)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
