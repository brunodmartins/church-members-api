// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_church is a generated GoMock package.
package mock_church

import (
	reflect "reflect"

	domain "github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetChurch mocks base method.
func (m *MockService) GetChurch(id string) (*domain.Church, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChurch", id)
	ret0, _ := ret[0].(*domain.Church)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChurch indicates an expected call of GetChurch.
func (mr *MockServiceMockRecorder) GetChurch(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChurch", reflect.TypeOf((*MockService)(nil).GetChurch), id)
}

// List mocks base method.
func (m *MockService) List() ([]*domain.Church, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*domain.Church)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List))
}
