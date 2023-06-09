// Code generated by MockGen. DO NOT EDIT.
// Source: gateway.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

// MockMemDBClient is a mock of MemDBClient interface.
type MockMemDBClient struct {
	ctrl     *gomock.Controller
	recorder *MockMemDBClientMockRecorder
}

// MockMemDBClientMockRecorder is the mock recorder for MockMemDBClient.
type MockMemDBClientMockRecorder struct {
	mock *MockMemDBClient
}

// NewMockMemDBClient creates a new mock instance.
func NewMockMemDBClient(ctrl *gomock.Controller) *MockMemDBClient {
	mock := &MockMemDBClient{ctrl: ctrl}
	mock.recorder = &MockMemDBClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemDBClient) EXPECT() *MockMemDBClientMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockMemDBClient) Del(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockMemDBClientMockRecorder) Del(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockMemDBClient)(nil).Del), ctx, key)
}

// Expire mocks base method.
func (m *MockMemDBClient) Expire(ctx context.Context, key string, duration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expire", ctx, key, duration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Expire indicates an expected call of Expire.
func (mr *MockMemDBClientMockRecorder) Expire(ctx, key, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expire", reflect.TypeOf((*MockMemDBClient)(nil).Expire), ctx, key, duration)
}

// Get mocks base method.
func (m *MockMemDBClient) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMemDBClientMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMemDBClient)(nil).Get), ctx, key)
}

// Ping mocks base method.
func (m *MockMemDBClient) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockMemDBClientMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMemDBClient)(nil).Ping), ctx)
}

// Set mocks base method.
func (m *MockMemDBClient) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, duration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockMemDBClientMockRecorder) Set(ctx, key, value, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockMemDBClient)(nil).Set), ctx, key, value, duration)
}

// SetNX mocks base method.
func (m *MockMemDBClient) SetNX(ctx context.Context, key string, value any, duration time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNX", ctx, key, value, duration)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetNX indicates an expected call of SetNX.
func (mr *MockMemDBClientMockRecorder) SetNX(ctx, key, value, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNX", reflect.TypeOf((*MockMemDBClient)(nil).SetNX), ctx, key, value, duration)
}

// SetXX mocks base method.
func (m *MockMemDBClient) SetXX(ctx context.Context, key string, value any, duration time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetXX", ctx, key, value, duration)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetXX indicates an expected call of SetXX.
func (mr *MockMemDBClientMockRecorder) SetXX(ctx, key, value, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetXX", reflect.TypeOf((*MockMemDBClient)(nil).SetXX), ctx, key, value, duration)
}

// MockAIClient is a mock of AIClient interface.
type MockAIClient struct {
	ctrl     *gomock.Controller
	recorder *MockAIClientMockRecorder
}

// MockAIClientMockRecorder is the mock recorder for MockAIClient.
type MockAIClientMockRecorder struct {
	mock *MockAIClient
}

// NewMockAIClient creates a new mock instance.
func NewMockAIClient(ctrl *gomock.Controller) *MockAIClient {
	mock := &MockAIClient{ctrl: ctrl}
	mock.recorder = &MockAIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAIClient) EXPECT() *MockAIClientMockRecorder {
	return m.recorder
}

// ChatCompletion mocks base method.
func (m *MockAIClient) ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatCompletion", ctx, messages)
	ret0, _ := ret[0].(domain.ChatMessages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChatCompletion indicates an expected call of ChatCompletion.
func (mr *MockAIClientMockRecorder) ChatCompletion(ctx, messages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatCompletion", reflect.TypeOf((*MockAIClient)(nil).ChatCompletion), ctx, messages)
}
