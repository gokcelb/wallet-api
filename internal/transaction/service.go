package transaction

import (
	"context"
	"errors"
)

var ErrTransactionNotFound = errors.New("no transaction with the given id exists")

type TransactionRepository interface {
	Create(ctx context.Context, txn *Transaction) (string, error)
	Read(ctx context.Context, id string) (*Transaction, error)
}

type service struct {
	tr TransactionRepository
}

func NewService(tr TransactionRepository) *service {
	return &service{tr}
}

func (s *service) CreateTransaction(ctx context.Context, txn *Transaction) (string, error) {
	return s.tr.Create(ctx, txn)
}

func (s *service) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	return s.tr.Read(ctx, id)
}
