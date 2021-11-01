// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mock_member is a generated GoMock package.
package mock_member

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	gomock "github.com/golang/mock/gomock"
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

// FindAll mocks base method.
func (m *MockRepository) FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, specification)
	ret0, _ := ret[0].([]*domain.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockRepositoryMockRecorder) FindAll(ctx, specification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockRepository)(nil).FindAll), ctx, specification)
}

// FindByID mocks base method.
func (m *MockRepository) FindByID(ctx context.Context, id string) (*domain.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRepository)(nil).FindByID), ctx, id)
}

// GenerateStatusHistory mocks base method.
func (m *MockRepository) GenerateStatusHistory(id string, status bool, reason string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateStatusHistory", id, status, reason, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateStatusHistory indicates an expected call of GenerateStatusHistory.
func (mr *MockRepositoryMockRecorder) GenerateStatusHistory(id, status, reason, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateStatusHistory", reflect.TypeOf((*MockRepository)(nil).GenerateStatusHistory), id, status, reason, date)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, member *domain.Member) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, member)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, member interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, member)
}

// UpdateStatus mocks base method.
func (m *MockRepository) UpdateStatus(ctx context.Context, member *domain.Member) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, member)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockRepositoryMockRecorder) UpdateStatus(ctx, member interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockRepository)(nil).UpdateStatus), ctx, member)
}
