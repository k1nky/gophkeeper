// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	user "github.com/k1nky/gophkeeper/internal/entity/user"
	vault "github.com/k1nky/gophkeeper/internal/entity/vault"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// GetSecretData mocks base method.
func (m *Mockstorage) GetSecretData(ctx context.Context, uk vault.UniqueKey) (*vault.DataReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretData", ctx, uk)
	ret0, _ := ret[0].(*vault.DataReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretData indicates an expected call of GetSecretData.
func (mr *MockstorageMockRecorder) GetSecretData(ctx, uk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretData", reflect.TypeOf((*Mockstorage)(nil).GetSecretData), ctx, uk)
}

// GetSecretMeta mocks base method.
func (m *Mockstorage) GetSecretMeta(ctx context.Context, uk vault.UniqueKey) (*vault.Meta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretMeta", ctx, uk)
	ret0, _ := ret[0].(*vault.Meta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretMeta indicates an expected call of GetSecretMeta.
func (mr *MockstorageMockRecorder) GetSecretMeta(ctx, uk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretMeta", reflect.TypeOf((*Mockstorage)(nil).GetSecretMeta), ctx, uk)
}

// ListSecretsByUser mocks base method.
func (m *Mockstorage) ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecretsByUser", ctx, userID)
	ret0, _ := ret[0].(vault.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecretsByUser indicates an expected call of ListSecretsByUser.
func (mr *MockstorageMockRecorder) ListSecretsByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecretsByUser", reflect.TypeOf((*Mockstorage)(nil).ListSecretsByUser), ctx, userID)
}

// PutSecret mocks base method.
func (m *Mockstorage) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutSecret", ctx, meta, data)
	ret0, _ := ret[0].(*vault.Meta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutSecret indicates an expected call of PutSecret.
func (mr *MockstorageMockRecorder) PutSecret(ctx, meta, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutSecret", reflect.TypeOf((*Mockstorage)(nil).PutSecret), ctx, meta, data)
}

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Debugf mocks base method.
func (m *Mocklogger) Debugf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockloggerMockRecorder) Debugf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*Mocklogger)(nil).Debugf), varargs...)
}

// Errorf mocks base method.
func (m *Mocklogger) Errorf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockloggerMockRecorder) Errorf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*Mocklogger)(nil).Errorf), varargs...)
}