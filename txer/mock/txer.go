// Code generated by MockGen. DO NOT EDIT.
// Source: txer.go

// Package txer_mock is a generated GoMock package.
package txer_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTXer is a mock of TXer interface.
type MockTXer struct {
	ctrl     *gomock.Controller
	recorder *MockTXerMockRecorder
}

// MockTXerMockRecorder is the mock recorder for MockTXer.
type MockTXerMockRecorder struct {
	mock *MockTXer
}

// NewMockTXer creates a new mock instance.
func NewMockTXer(ctrl *gomock.Controller) *MockTXer {
	mock := &MockTXer{ctrl: ctrl}
	mock.recorder = &MockTXerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTXer) EXPECT() *MockTXerMockRecorder {
	return m.recorder
}

// WithTX mocks base method.
func (m *MockTXer) WithTX(arg0 context.Context, arg1 func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTX", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTX indicates an expected call of WithTX.
func (mr *MockTXerMockRecorder) WithTX(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTX", reflect.TypeOf((*MockTXer)(nil).WithTX), arg0, arg1)
}
