// Code generated by MockGen. DO NOT EDIT.
// Source: ./pdf.go

// Package mock_file is a generated GoMock package.
package mock_file

import (
	reflect "reflect"

	domain "github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockBuilder is a mock of Builder interface.
type MockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderMockRecorder
}

// MockBuilderMockRecorder is the mock recorder for MockBuilder.
type MockBuilderMockRecorder struct {
	mock *MockBuilder
}

// NewMockBuilder creates a new mock instance.
func NewMockBuilder(ctrl *gomock.Controller) *MockBuilder {
	mock := &MockBuilder{ctrl: ctrl}
	mock.recorder = &MockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilder) EXPECT() *MockBuilderMockRecorder {
	return m.recorder
}

// BuildFile mocks base method.
func (m *MockBuilder) BuildFile(title string, data []*domain.Member) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildFile", title, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildFile indicates an expected call of BuildFile.
func (mr *MockBuilderMockRecorder) BuildFile(title, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildFile", reflect.TypeOf((*MockBuilder)(nil).BuildFile), title, data)
}
