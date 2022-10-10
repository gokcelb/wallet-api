package transaction_test

import (
	"context"
	"testing"

	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/gokcelb/wallet-api/internal/transaction/mock"
	"github.com/gokcelb/wallet-api/internal/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMockTransactionRepository(t *testing.T) *mock.MockTransactionRepository {
	return mock.NewMockTransactionRepository(gomock.NewController(t))
}

func TestServiceCreateTransaction(t *testing.T) {
	mockRepository := createMockTransactionRepository(t)
	s := transaction.NewService(mockRepository)

	givenTxn := &transaction.Transaction{
		WalletID: "1",
		Type:     "deposit",
		Amount:   200,
	}
	mockRepoTxnID := "1"

	mockRepository.EXPECT().Create(context.TODO(), givenTxn).Return(mockRepoTxnID, nil)

	id, err := s.CreateTransaction(context.TODO(), givenTxn)

	assert.Equal(t, mockRepoTxnID, id)
	assert.Nil(t, err)
}

func TestServiceGetTransaction(t *testing.T) {
	mockRepository := createMockTransactionRepository(t)
	s := transaction.NewService(mockRepository)

	mockTxn := &transaction.Transaction{
		ID:       "1",
		WalletID: "1",
		Type:     "deposit",
		Amount:   100,
	}

	mockRepository.EXPECT().Read(context.TODO(), "1").Return(mockTxn, nil)

	txn, err := s.GetTransaction(context.TODO(), "1")

	assert.Equal(t, mockTxn, txn)
	assert.Nil(t, err)
}

func TestServiceGetTransactionsWithoutType(t *testing.T) {
	mockRepository := createMockTransactionRepository(t)
	s := transaction.NewService(mockRepository)

	mockTxns := []*transaction.Transaction{
		{
			ID:       "1",
			WalletID: "1",
			Type:     "deposit",
			Amount:   100,
		},
		{
			ID:       "2",
			WalletID: "1",
			Type:     "withdrawal",
			Amount:   100,
		},
	}

	mockRepository.EXPECT().
		ReadByWalletID(context.TODO(), "1", wallet.DefaultPageNo, wallet.DefaultPageSize).
		Return(mockTxns, nil)

	txns, err := s.GetTransactionsByWalletID(
		context.TODO(),
		"1",
		"",
		wallet.DefaultPageNo,
		wallet.DefaultPageSize,
	)

	assert.Equal(t, mockTxns, txns)
	assert.Nil(t, err)
}

func TestServiceGetTransactionsWithType(t *testing.T) {
	mockRepository := createMockTransactionRepository(t)
	s := transaction.NewService(mockRepository)

	mockTxns := []*transaction.Transaction{
		{
			ID:       "1",
			WalletID: "1",
			Type:     "deposit",
			Amount:   100,
		},
	}

	mockRepository.EXPECT().
		ReadByWalletIDFilterByType(context.TODO(), "1", "deposit", wallet.DefaultPageNo, wallet.DefaultPageSize).
		Return(mockTxns, nil)

	txns, err := s.GetTransactionsByWalletID(
		context.TODO(),
		"1",
		"deposit",
		wallet.DefaultPageNo,
		wallet.DefaultPageSize,
	)

	assert.Equal(t, mockTxns, txns)
	assert.Nil(t, err)
}
