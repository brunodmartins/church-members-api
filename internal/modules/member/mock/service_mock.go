// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_member is a generated GoMock package.
package mock_member

import (
	reflect "reflect"
	time "time"

	domain "github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	member "github.com/BrunoDM2943/church-members-api/internal/modules/member"
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

// ChangeStatus mocks base method.
func (m *MockService) ChangeStatus(id string, status bool, reason string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", id, status, reason, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus.
func (mr *MockServiceMockRecorder) ChangeStatus(id, status, reason, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockService)(nil).ChangeStatus), id, status, reason, date)
}

// FindMembers mocks base method.
func (m *MockService) SearchMembers(specification member.Specification) ([]*domain.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchMembers", specification)
	ret0, _ := ret[0].([]*domain.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMembers indicates an expected call of FindMembers.
func (mr *MockServiceMockRecorder) FindMembers(specification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchMembers", reflect.TypeOf((*MockService)(nil).SearchMembers), specification)
}

// FindMembersByID mocks base method.
func (m *MockService) GetMember(id string) (*domain.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMember", id)
	ret0, _ := ret[0].(*domain.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMembersByID indicates an expected call of FindMembersByID.
func (mr *MockServiceMockRecorder) FindMembersByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMember", reflect.TypeOf((*MockService)(nil).GetMember), id)
}

// SaveMember mocks base method.
func (m *MockService) SaveMember(member *domain.Member) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMember", member)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMember indicates an expected call of SaveMember.
func (mr *MockServiceMockRecorder) SaveMember(member interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMember", reflect.TypeOf((*MockService)(nil).SaveMember), member)
}
