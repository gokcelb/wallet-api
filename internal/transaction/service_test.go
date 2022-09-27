package transaction_test

import (
	"context"
	"testing"

	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/gokcelb/wallet-api/internal/transaction/mock"
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
