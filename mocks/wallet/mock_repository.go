// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/wallet/service.go

// Package mock_wallet is a generated GoMock package.
package mock_wallet

import (
	context "context"
	reflect "reflect"

	wallet "github.com/gokcelb/wallet-api/internal/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, wallet wallet.Wallet) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, wallet)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, wallet)
}

// Read mocks base method.
func (m *MockRepository) Read(ctx context.Context, id string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, id)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockRepositoryMockRecorder) Read(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockRepository)(nil).Read), ctx, id)
}

// ReadByUserId mocks base method.
func (m *MockRepository) ReadByUserId(ctx context.Context, userId string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByUserId", ctx, userId)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByUserId indicates an expected call of ReadByUserId.
func (mr *MockRepositoryMockRecorder) ReadByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByUserId", reflect.TypeOf((*MockRepository)(nil).ReadByUserId), ctx, userId)
}