// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	models "src/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockMerchRepository is a mock of MerchRepository interface.
type MockMerchRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMerchRepositoryMockRecorder
}

// MockMerchRepositoryMockRecorder is the mock recorder for MockMerchRepository.
type MockMerchRepositoryMockRecorder struct {
	mock *MockMerchRepository
}

// NewMockMerchRepository creates a new mock instance.
func NewMockMerchRepository(ctrl *gomock.Controller) *MockMerchRepository {
	mock := &MockMerchRepository{ctrl: ctrl}
	mock.recorder = &MockMerchRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerchRepository) EXPECT() *MockMerchRepositoryMockRecorder {
	return m.recorder
}

// AddMerch mocks base method.
func (m *MockMerchRepository) AddMerch(album *models.Merch) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMerch", album)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMerch indicates an expected call of AddMerch.
func (mr *MockMerchRepositoryMockRecorder) AddMerch(album interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMerch", reflect.TypeOf((*MockMerchRepository)(nil).AddMerch), album)
}

// DeleteMerch mocks base method.
func (m *MockMerchRepository) DeleteMerch(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMerch", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMerch indicates an expected call of DeleteMerch.
func (mr *MockMerchRepositoryMockRecorder) DeleteMerch(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMerch", reflect.TypeOf((*MockMerchRepository)(nil).DeleteMerch), id)
}

// GetMerch mocks base method.
func (m *MockMerchRepository) GetMerch(id uint64) (*models.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerch", id)
	ret0, _ := ret[0].(*models.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerch indicates an expected call of GetMerch.
func (mr *MockMerchRepositoryMockRecorder) GetMerch(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerch", reflect.TypeOf((*MockMerchRepository)(nil).GetMerch), id)
}

// UpdateMerch mocks base method.
func (m *MockMerchRepository) UpdateMerch(album *models.Merch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMerch", album)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMerch indicates an expected call of UpdateMerch.
func (mr *MockMerchRepositoryMockRecorder) UpdateMerch(album interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMerch", reflect.TypeOf((*MockMerchRepository)(nil).UpdateMerch), album)
}