// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	models "src/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockPlaylistRepository is a mock of PlaylistRepository interface.
type MockPlaylistRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPlaylistRepositoryMockRecorder
}

// MockPlaylistRepositoryMockRecorder is the mock recorder for MockPlaylistRepository.
type MockPlaylistRepositoryMockRecorder struct {
	mock *MockPlaylistRepository
}

// NewMockPlaylistRepository creates a new mock instance.
func NewMockPlaylistRepository(ctrl *gomock.Controller) *MockPlaylistRepository {
	mock := &MockPlaylistRepository{ctrl: ctrl}
	mock.recorder = &MockPlaylistRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlaylistRepository) EXPECT() *MockPlaylistRepositoryMockRecorder {
	return m.recorder
}

// AddPlaylist mocks base method.
func (m *MockPlaylistRepository) AddPlaylist(playlist *models.Playlist, userId uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlaylist", playlist, userId)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPlaylist indicates an expected call of AddPlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) AddPlaylist(playlist, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).AddPlaylist), playlist, userId)
}

// AddTrackToPlaylist mocks base method.
func (m *MockPlaylistRepository) AddTrackToPlaylist(playlistId, trackId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrackToPlaylist", playlistId, trackId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTrackToPlaylist indicates an expected call of AddTrackToPlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) AddTrackToPlaylist(playlistId, trackId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrackToPlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).AddTrackToPlaylist), playlistId, trackId)
}

// DeletePlaylist mocks base method.
func (m *MockPlaylistRepository) DeletePlaylist(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlaylist", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlaylist indicates an expected call of DeletePlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) DeletePlaylist(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).DeletePlaylist), id)
}

// DeleteTrackFromPlaylist mocks base method.
func (m *MockPlaylistRepository) DeleteTrackFromPlaylist(playlistId, trackId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrackFromPlaylist", playlistId, trackId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrackFromPlaylist indicates an expected call of DeleteTrackFromPlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) DeleteTrackFromPlaylist(playlistId, trackId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrackFromPlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).DeleteTrackFromPlaylist), playlistId, trackId)
}

// GetAllPlaylistsForUser mocks base method.
func (m *MockPlaylistRepository) GetAllPlaylistsForUser(userId uint64) ([]*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPlaylistsForUser", userId)
	ret0, _ := ret[0].([]*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPlaylistsForUser indicates an expected call of GetAllPlaylistsForUser.
func (mr *MockPlaylistRepositoryMockRecorder) GetAllPlaylistsForUser(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPlaylistsForUser", reflect.TypeOf((*MockPlaylistRepository)(nil).GetAllPlaylistsForUser), userId)
}

// GetAllTracks mocks base method.
func (m *MockPlaylistRepository) GetAllTracks(playlistId uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTracks", playlistId)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTracks indicates an expected call of GetAllTracks.
func (mr *MockPlaylistRepositoryMockRecorder) GetAllTracks(playlistId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTracks", reflect.TypeOf((*MockPlaylistRepository)(nil).GetAllTracks), playlistId)
}

// GetPlaylist mocks base method.
func (m *MockPlaylistRepository) GetPlaylist(id uint64) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaylist", id)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaylist indicates an expected call of GetPlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) GetPlaylist(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).GetPlaylist), id)
}

// GetUserForPlaylist mocks base method.
func (m *MockPlaylistRepository) GetUserForPlaylist(playlistId uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserForPlaylist", playlistId)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserForPlaylist indicates an expected call of GetUserForPlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) GetUserForPlaylist(playlistId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserForPlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).GetUserForPlaylist), playlistId)
}

// IsPlaylistOwned mocks base method.
func (m *MockPlaylistRepository) IsPlaylistOwned(playlistId, userId uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPlaylistOwned", playlistId, userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPlaylistOwned indicates an expected call of IsPlaylistOwned.
func (mr *MockPlaylistRepositoryMockRecorder) IsPlaylistOwned(playlistId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlaylistOwned", reflect.TypeOf((*MockPlaylistRepository)(nil).IsPlaylistOwned), playlistId, userId)
}

// UpdatePlaylist mocks base method.
func (m *MockPlaylistRepository) UpdatePlaylist(playlist *models.Playlist) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePlaylist", playlist)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePlaylist indicates an expected call of UpdatePlaylist.
func (mr *MockPlaylistRepositoryMockRecorder) UpdatePlaylist(playlist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePlaylist", reflect.TypeOf((*MockPlaylistRepository)(nil).UpdatePlaylist), playlist)
}
