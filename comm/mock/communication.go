// Code generated by MockGen. DO NOT EDIT.
// Source: ./comm/communication.go

// Package mock_comm is a generated GoMock package.
package mock_comm

import (
	reflect "reflect"

	comm "github.com/ChainSafe/sygma-relayer/comm"
	gomock "github.com/golang/mock/gomock"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

// MockCommunication is a mock of Communication interface.
type MockCommunication struct {
	ctrl     *gomock.Controller
	recorder *MockCommunicationMockRecorder
}

// MockCommunicationMockRecorder is the mock recorder for MockCommunication.
type MockCommunicationMockRecorder struct {
	mock *MockCommunication
}

// NewMockCommunication creates a new mock instance.
func NewMockCommunication(ctrl *gomock.Controller) *MockCommunication {
	mock := &MockCommunication{ctrl: ctrl}
	mock.recorder = &MockCommunicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunication) EXPECT() *MockCommunicationMockRecorder {
	return m.recorder
}

// Broadcast mocks base method.
func (m *MockCommunication) Broadcast(peers peer.IDSlice, msg []byte, msgType comm.MessageType, sessionID string, errChan chan error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Broadcast", peers, msg, msgType, sessionID, errChan)
}

// Broadcast indicates an expected call of Broadcast.
func (mr *MockCommunicationMockRecorder) Broadcast(peers, msg, msgType, sessionID, errChan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockCommunication)(nil).Broadcast), peers, msg, msgType, sessionID, errChan)
}

// CloseSession mocks base method.
func (m *MockCommunication) CloseSession(sessionID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseSession", sessionID)
}

// CloseSession indicates an expected call of CloseSession.
func (mr *MockCommunicationMockRecorder) CloseSession(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSession", reflect.TypeOf((*MockCommunication)(nil).CloseSession), sessionID)
}

// Subscribe mocks base method.
func (m *MockCommunication) Subscribe(sessionID string, msgType comm.MessageType, channel chan *comm.WrappedMessage) comm.SubscriptionID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", sessionID, msgType, channel)
	ret0, _ := ret[0].(comm.SubscriptionID)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockCommunicationMockRecorder) Subscribe(sessionID, msgType, channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockCommunication)(nil).Subscribe), sessionID, msgType, channel)
}

// UnSubscribe mocks base method.
func (m *MockCommunication) UnSubscribe(subID comm.SubscriptionID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnSubscribe", subID)
}

// UnSubscribe indicates an expected call of UnSubscribe.
func (mr *MockCommunicationMockRecorder) UnSubscribe(subID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnSubscribe", reflect.TypeOf((*MockCommunication)(nil).UnSubscribe), subID)
}
