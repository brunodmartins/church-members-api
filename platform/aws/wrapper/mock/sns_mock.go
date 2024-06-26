// Code generated by MockGen. DO NOT EDIT.
// Source: ./sns.go
//
// Generated by this command:
//
//	mockgen -source=./sns.go -destination=./mock/sns_mock.go
//

// Package mock_wrapper is a generated GoMock package.
package mock_wrapper

import (
	context "context"
	reflect "reflect"

	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	gomock "go.uber.org/mock/gomock"
)

// MockSNSAPI is a mock of SNSAPI interface.
type MockSNSAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSNSAPIMockRecorder
}

// MockSNSAPIMockRecorder is the mock recorder for MockSNSAPI.
type MockSNSAPIMockRecorder struct {
	mock *MockSNSAPI
}

// NewMockSNSAPI creates a new mock instance.
func NewMockSNSAPI(ctrl *gomock.Controller) *MockSNSAPI {
	mock := &MockSNSAPI{ctrl: ctrl}
	mock.recorder = &MockSNSAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNSAPI) EXPECT() *MockSNSAPIMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockSNSAPI) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(*sns.PublishOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockSNSAPIMockRecorder) Publish(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockSNSAPI)(nil).Publish), varargs...)
}
