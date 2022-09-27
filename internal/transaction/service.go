package transaction

import (
	"context"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *Transaction) (string, error)
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
