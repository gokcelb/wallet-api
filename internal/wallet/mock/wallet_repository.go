// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gokcelb/wallet-api/internal/wallet (interfaces: WalletRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	wallet "github.com/gokcelb/wallet-api/internal/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
	mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
	mock := &MockWalletRepository{ctrl: ctrl}
	mock.recorder = &MockWalletRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWalletRepository) Create(arg0 context.Context, arg1 wallet.Wallet) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockWalletRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWalletRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockWalletRepository) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockWalletRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWalletRepository)(nil).Delete), arg0, arg1)
}

// Read mocks base method.
func (m *MockWalletRepository) Read(arg0 context.Context, arg1 string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockWalletRepositoryMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockWalletRepository)(nil).Read), arg0, arg1)
}

// ReadByUserID mocks base method.
func (m *MockWalletRepository) ReadByUserID(arg0 context.Context, arg1 string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByUserID", arg0, arg1)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByUserID indicates an expected call of ReadByUserID.
func (mr *MockWalletRepositoryMockRecorder) ReadByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByUserID", reflect.TypeOf((*MockWalletRepository)(nil).ReadByUserID), arg0, arg1)
}

// UpdateBalance mocks base method.
func (m *MockWalletRepository) UpdateBalance(arg0 context.Context, arg1 string, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockWalletRepositoryMockRecorder) UpdateBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockWalletRepository)(nil).UpdateBalance), arg0, arg1, arg2)
}