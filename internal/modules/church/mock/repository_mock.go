// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go
//
// Generated by this command:
//
//	mockgen -source=./repository.go -destination=./mock/repository_mock.go
//

// Package mock_church is a generated GoMock package.
package mock_church

import (
	reflect "reflect"

	domain "github.com/brunodmartins/church-members-api/internal/constants/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ID string) (*domain.Church, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(*domain.Church)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ID)
}

// List mocks base method.
func (m *MockRepository) List() ([]*domain.Church, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*domain.Church)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List))
}
