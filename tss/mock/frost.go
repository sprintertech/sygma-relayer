// Code generated by MockGen. DO NOT EDIT.
// Source: ./tss/frost/keygen/keygen.go

// Package mock_tss is a generated GoMock package.
package mock_tss

import (
	reflect "reflect"

	keyshare "github.com/ChainSafe/sygma-relayer/keyshare"
	gomock "github.com/golang/mock/gomock"
)

// MockFrostKeyshareStorer is a mock of FrostKeyshareStorer interface.
type MockFrostKeyshareStorer struct {
	ctrl     *gomock.Controller
	recorder *MockFrostKeyshareStorerMockRecorder
}

// MockFrostKeyshareStorerMockRecorder is the mock recorder for MockFrostKeyshareStorer.
type MockFrostKeyshareStorerMockRecorder struct {
	mock *MockFrostKeyshareStorer
}

// NewMockFrostKeyshareStorer creates a new mock instance.
func NewMockFrostKeyshareStorer(ctrl *gomock.Controller) *MockFrostKeyshareStorer {
	mock := &MockFrostKeyshareStorer{ctrl: ctrl}
	mock.recorder = &MockFrostKeyshareStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFrostKeyshareStorer) EXPECT() *MockFrostKeyshareStorerMockRecorder {
	return m.recorder
}

// GetKeyshare mocks base method.
func (m *MockFrostKeyshareStorer) GetKeyshare(publicKey string) (keyshare.FrostKeyshare, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeyshare", publicKey)
	ret0, _ := ret[0].(keyshare.FrostKeyshare)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeyshare indicates an expected call of GetKeyshare.
func (mr *MockFrostKeyshareStorerMockRecorder) GetKeyshare(publicKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeyshare", reflect.TypeOf((*MockFrostKeyshareStorer)(nil).GetKeyshare), publicKey)
}

// LockKeyshare mocks base method.
func (m *MockFrostKeyshareStorer) LockKeyshare() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LockKeyshare")
}

// LockKeyshare indicates an expected call of LockKeyshare.
func (mr *MockFrostKeyshareStorerMockRecorder) LockKeyshare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockKeyshare", reflect.TypeOf((*MockFrostKeyshareStorer)(nil).LockKeyshare))
}

// StoreKeyshare mocks base method.
func (m *MockFrostKeyshareStorer) StoreKeyshare(keyshare keyshare.FrostKeyshare) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreKeyshare", keyshare)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreKeyshare indicates an expected call of StoreKeyshare.
func (mr *MockFrostKeyshareStorerMockRecorder) StoreKeyshare(keyshare interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreKeyshare", reflect.TypeOf((*MockFrostKeyshareStorer)(nil).StoreKeyshare), keyshare)
}

// UnlockKeyshare mocks base method.
func (m *MockFrostKeyshareStorer) UnlockKeyshare() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnlockKeyshare")
}

// UnlockKeyshare indicates an expected call of UnlockKeyshare.
func (mr *MockFrostKeyshareStorerMockRecorder) UnlockKeyshare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockKeyshare", reflect.TypeOf((*MockFrostKeyshareStorer)(nil).UnlockKeyshare))
}
