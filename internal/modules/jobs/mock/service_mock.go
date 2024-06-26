// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go
//
// Generated by this command:
//
//	mockgen -source=./service.go -destination=./mock/service_mock.go
//

// Package mock_jobs is a generated GoMock package.
package mock_jobs

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockJob is a mock of Job interface.
type MockJob struct {
	ctrl     *gomock.Controller
	recorder *MockJobMockRecorder
}

// MockJobMockRecorder is the mock recorder for MockJob.
type MockJobMockRecorder struct {
	mock *MockJob
}

// NewMockJob creates a new mock instance.
func NewMockJob(ctrl *gomock.Controller) *MockJob {
	mock := &MockJob{ctrl: ctrl}
	mock.recorder = &MockJobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJob) EXPECT() *MockJobMockRecorder {
	return m.recorder
}

// RunJob mocks base method.
func (m *MockJob) RunJob(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunJob", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunJob indicates an expected call of RunJob.
func (mr *MockJobMockRecorder) RunJob(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunJob", reflect.TypeOf((*MockJob)(nil).RunJob), ctx)
}
