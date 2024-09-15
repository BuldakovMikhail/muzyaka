package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock_repository "src/internal/domain/album/repository/mocks"
	mock_repository2 "src/internal/domain/track/repository/mocks"
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
					Id:        1,
					Name:      "test_name",
					CoverFile: []byte("test_cover"),
					Type:      "test_type",
				}, nil)
			},
			expectedValue: &models.Album{
				Id:        1,
				Name:      "test_name",
				CoverFile: []byte("test_cover"),
				Type:      "test_type",
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

			storage := mock_repository2.NewMockTrackStorage(c)

			s := NewAlbumUseCase(repo, storage)
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
				Id:        1,
				Name:      "test_name",
				CoverFile: []byte("test_cover"),
				Type:      "test_type",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album models.Album) {
				r.EXPECT().UpdateAlbum(&album).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Fail in repo",
			input: models.Album{
				Id:        1,
				Name:      "test_name",
				CoverFile: []byte("test_cover"),
				Type:      "test_type",
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

			storage := mock_repository2.NewMockTrackStorage(c)

			s := NewAlbumUseCase(repo, storage)
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
	type mock func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.TrackObject)
	type storageMock func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackObject)

	testTable := []struct {
		name        string
		inputAlbum  *models.Album
		inputTracks []*models.TrackObject
		mock        mock
		storageMock storageMock
		expectedID  uint64
		expectedErr error
	}{
		{
			name: "Usual test",
			inputAlbum: &models.Album{
				Id:        1,
				Name:      "Test Album",
				CoverFile: []byte{1, 2, 3},
				Type:      "LP",
			},
			inputTracks: []*models.TrackObject{
				{
					TrackMeta: models.TrackMeta{Id: 1, Name: "TrackMeta 2"},
					Payload:   []byte{1, 2, 3},
				},
				{
					TrackMeta: models.TrackMeta{Id: 2, Name: "TrackMeta 2"},
					Payload:   []byte{1, 2, 3},
				},
			},
			mock: func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.TrackObject) {
				var metaTracks []*models.TrackMeta

				for _, v := range tracks {
					metaTracks = append(metaTracks, v.ExtractMeta())
				}

				r.EXPECT().AddAlbumWithTracksOutbox(album, metaTracks, uint64(1)).Return(uint64(1), nil)
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackObject) {
				for _, v := range tracks {
					r.EXPECT().UploadObject(v).Return(nil)
				}

			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Repo fail test",
			inputAlbum: &models.Album{
				Id:        0,
				Name:      "Invalid Album",
				CoverFile: []byte("Test CoverFile"),
				Type:      "LP",
			},
			inputTracks: []*models.TrackObject{
				{
					TrackMeta: models.TrackMeta{Id: 1, Name: "TrackMeta 2"},
					Payload:   []byte{1, 2, 3},
				},
				{
					TrackMeta: models.TrackMeta{Id: 2, Name: "TrackMeta 2"},
					Payload:   []byte{1, 2, 3},
				},
			},
			mock: func(r *mock_repository.MockAlbumRepository, album *models.Album, tracks []*models.TrackObject) {
				var metaTracks []*models.TrackMeta

				for _, v := range tracks {
					metaTracks = append(metaTracks, v.ExtractMeta())
				}

				r.EXPECT().AddAlbumWithTracksOutbox(album, metaTracks, uint64(1)).Return(uint64(0), errors.New("error in repo"))
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackObject) {
				for _, v := range tracks {
					r.EXPECT().UploadObject(v).Return(nil)
				}
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

			storage := mock_repository2.NewMockTrackStorage(ctrl)
			tc.storageMock(storage, tc.inputTracks)

			u := NewAlbumUseCase(repo, storage)
			id, err := u.AddAlbumWithTracks(tc.inputAlbum, tc.inputTracks, 1)

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
	type mock func(r *mock_repository.MockAlbumRepository, id uint64, tracks []*models.TrackMeta)
	type storageMock func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackMeta)

	testTable := []struct {
		name        string
		input       uint64
		tracks      []*models.TrackMeta
		mock        mock
		storageMock storageMock
		expectedErr error
	}{
		{
			name:  "Usual test",
			input: uint64(1),
			tracks: []*models.TrackMeta{
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
			mock: func(r *mock_repository.MockAlbumRepository, id uint64, tracks []*models.TrackMeta) {
				r.EXPECT().GetAllTracksForAlbum(id).Return(tracks, nil)
				r.EXPECT().DeleteAlbumOutbox(id).Return(nil)
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackMeta) {
				for _, v := range tracks {
					r.EXPECT().DeleteObject(v).Return(nil)
				}
			},
			expectedErr: nil,
		},
		{
			name:  "Fail in repo test",
			input: uint64(110),
			tracks: []*models.TrackMeta{
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
			mock: func(r *mock_repository.MockAlbumRepository, id uint64, tracks []*models.TrackMeta) {
				r.EXPECT().GetAllTracksForAlbum(id).Return(tracks, nil)
				r.EXPECT().DeleteAlbumOutbox(id).Return(errors.New("error in repo"))
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, tracks []*models.TrackMeta) {

			},
			expectedErr: errors.Wrap(errors.New("error in repo"), "album.usecase.DeleteAlbumOutbox error while delete"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockAlbumRepository(c)
			tc.mock(repo, tc.input, tc.tracks)

			storage := mock_repository2.NewMockTrackStorage(c)
			tc.storageMock(storage, tc.tracks)

			s := NewAlbumUseCase(repo, storage)
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
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackObject)
	type storageMock func(r *mock_repository2.MockTrackStorage, tracks models.TrackObject)

	testTable := []struct {
		name          string
		inputId       uint64
		inputTrack    models.TrackObject
		mock          mock
		storageMock   storageMock
		expectedValue uint64
		expectedErr   error
	}{
		{
			name: "Usual test",
			inputTrack: models.TrackObject{
				TrackMeta: models.TrackMeta{
					Id:     10,
					Source: "test_src",
					Name:   "test_name",
					Genre:  "test_genre",
				},
				Payload: []byte{1, 2, 3},
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackObject) {
				r.EXPECT().AddTrackToAlbumOutbox(album_id, track.ExtractMeta()).Return(uint64(10), nil)
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, track models.TrackObject) {
				r.EXPECT().UploadObject(&track).Return(nil)
			},
			expectedValue: uint64(10),
			expectedErr:   nil,
		},
		{
			name: "Repo fail test",
			inputTrack: models.TrackObject{
				TrackMeta: models.TrackMeta{
					Id:     10,
					Source: "test_src",
					Name:   "test_name",
					Genre:  "test_genre",
				},
				Payload: []byte{1, 2, 3},
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, track models.TrackObject) {
				r.EXPECT().UploadObject(&track).Return(nil)
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackObject) {
				r.EXPECT().AddTrackToAlbumOutbox(album_id, track.ExtractMeta()).Return(uint64(0), errors.New("error in repo"))
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
			storage := mock_repository2.NewMockTrackStorage(c)
			tc.storageMock(storage, tc.inputTrack)

			s := NewAlbumUseCase(repo, storage)
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
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackMeta)
	type storageMock func(r *mock_repository2.MockTrackStorage, tracks models.TrackMeta)

	testTable := []struct {
		name        string
		albumId     uint64
		inputTrack  models.TrackMeta
		mock        mock
		storageMock storageMock
		expectedErr error
	}{
		{
			name:    "Usual test",
			albumId: uint64(1),
			inputTrack: models.TrackMeta{
				Id:     10,
				Source: "test_src",
				Name:   "test_name",
				Genre:  "test_genre",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackMeta) {
				r.EXPECT().DeleteTrackFromAlbumOutbox(album_id, &track).Return(nil)
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, track models.TrackMeta) {
				r.EXPECT().DeleteObject(&track).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Repo fail test",
			albumId: uint64(2),
			inputTrack: models.TrackMeta{
				Id:     10,
				Source: "test_src",
				Name:   "test_name",
				Genre:  "test_genre",
			},
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, track models.TrackMeta) {
				r.EXPECT().DeleteTrackFromAlbumOutbox(album_id, &track).Return(errors.New("error in repo"))
			},
			storageMock: func(r *mock_repository2.MockTrackStorage, track models.TrackMeta) {
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
			tc.mock(repo, tc.albumId, tc.inputTrack)

			storage := mock_repository2.NewMockTrackStorage(ctrl)
			tc.storageMock(storage, tc.inputTrack)

			s := NewAlbumUseCase(repo, storage)
			err := s.DeleteTrack(tc.albumId)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestUsecase_GetAllTracks(t *testing.T) {
	type mock func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.TrackMeta)

	testTable := []struct {
		name           string
		albumId        uint64
		mock           mock
		expectedTracks []*models.TrackMeta
		expectedErr    error
	}{
		{
			name:    "Usual test",
			albumId: 1,
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.TrackMeta) {
				r.EXPECT().GetAllTracksForAlbum(album_id).Return(tracks, nil)
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
			name:    "Repo fail test",
			albumId: 2,
			mock: func(r *mock_repository.MockAlbumRepository, album_id uint64, tracks []*models.TrackMeta) {
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

			storage := mock_repository2.NewMockTrackStorage(ctrl)

			u := NewAlbumUseCase(repo, storage)
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
