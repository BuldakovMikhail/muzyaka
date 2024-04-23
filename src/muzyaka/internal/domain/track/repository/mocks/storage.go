// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	models "src/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockTrackStorage is a mock of TrackStorage interface.
type MockTrackStorage struct {
	ctrl     *gomock.Controller
	recorder *MockTrackStorageMockRecorder
}

// MockTrackStorageMockRecorder is the mock recorder for MockTrackStorage.
type MockTrackStorageMockRecorder struct {
	mock *MockTrackStorage
}

// NewMockTrackStorage creates a new mock instance.
func NewMockTrackStorage(ctrl *gomock.Controller) *MockTrackStorage {
	mock := &MockTrackStorage{ctrl: ctrl}
	mock.recorder = &MockTrackStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrackStorage) EXPECT() *MockTrackStorageMockRecorder {
	return m.recorder
}

// DeleteObject mocks base method.
func (m *MockTrackStorage) DeleteObject(track *models.TrackMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObject", track)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObject indicates an expected call of DeleteObject.
func (mr *MockTrackStorageMockRecorder) DeleteObject(track interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObject", reflect.TypeOf((*MockTrackStorage)(nil).DeleteObject), track)
}

// LoadObject mocks base method.
func (m *MockTrackStorage) LoadObject(track *models.TrackMeta) (*models.TrackObject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadObject", track)
	ret0, _ := ret[0].(*models.TrackObject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadObject indicates an expected call of LoadObject.
func (mr *MockTrackStorageMockRecorder) LoadObject(track interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadObject", reflect.TypeOf((*MockTrackStorage)(nil).LoadObject), track)
}

// UploadObject mocks base method.
func (m *MockTrackStorage) UploadObject(track *models.TrackObject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadObject", track)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadObject indicates an expected call of UploadObject.
func (mr *MockTrackStorageMockRecorder) UploadObject(track interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadObject", reflect.TypeOf((*MockTrackStorage)(nil).UploadObject), track)
}