// Code generated by MockGen. DO NOT EDIT.
// Source: ./calendar.go
//
// Generated by this command:
//
//	mockgen -source=./calendar.go -destination=./mock/calendar_mock.go
//

// Package mock_calendar is a generated GoMock package.
package mock_calendar

import (
	context "context"
	reflect "reflect"

	calendar "github.com/brunodmartins/church-members-api/internal/services/calendar"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetURL mocks base method.
func (m *MockStorage) GetURL(ctx context.Context, name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", ctx, name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockStorageMockRecorder) GetURL(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockStorage)(nil).GetURL), ctx, name)
}

// Store mocks base method.
func (m *MockStorage) Store(ctx context.Context, calendar *calendar.Calendar) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, calendar)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockStorageMockRecorder) Store(ctx, calendar any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockStorage)(nil).Store), ctx, calendar)
}
