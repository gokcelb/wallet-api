// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gokcelb/wallet-api/internal/transaction (interfaces: TransactionRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	transaction "github.com/gokcelb/wallet-api/internal/transaction"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTransactionRepository) Create(arg0 context.Context, arg1 *transaction.Transaction) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactionRepository)(nil).Create), arg0, arg1)
}

// Read mocks base method.
func (m *MockTransactionRepository) Read(arg0 context.Context, arg1 string) (*transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(*transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockTransactionRepositoryMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockTransactionRepository)(nil).Read), arg0, arg1)
}

// ReadByType mocks base method.
func (m *MockTransactionRepository) ReadByType(arg0 context.Context, arg1, arg2 string) (*transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByType", arg0, arg1, arg2)
	ret0, _ := ret[0].(*transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByType indicates an expected call of ReadByType.
func (mr *MockTransactionRepositoryMockRecorder) ReadByType(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByType", reflect.TypeOf((*MockTransactionRepository)(nil).ReadByType), arg0, arg1, arg2)
}
