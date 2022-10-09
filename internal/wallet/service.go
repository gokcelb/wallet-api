package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/transaction"
)

const (
	Deposit    string = "deposit"
	Withdrawal string = "withdrawal"
)

var (
	ErrWalletNotFound               = fmt.Errorf("no wallet with the given id exists")
	ErrWalletWithUserIDExists       = fmt.Errorf("wallet with user id already exists")
	ErrAboveMaximumBalanceLimit     = fmt.Errorf("wallet balance is above maximum balance limit")
	ErrAboveMaximumTransactionLimit = fmt.Errorf("transaction is above maximum transaction limit")
	ErrBelowMinimumTransactionLimit = fmt.Errorf("transaction is below minimum transaction limit")
	ErrInvalidTransactionType       = fmt.Errorf("transaction type is invalid")
	ErrInsufficientBalance          = fmt.Errorf("balance is insufficient")
	ErrWalletBalanceUpdateFailed    = fmt.Errorf("wallet balance could not be updated")
)

type WalletRepository interface {
	Create(ctx context.Context, w *Wallet) (string, error)
	Read(ctx context.Context, id string) (*Wallet, error)
	ReadByUserID(ctx context.Context, userID string) (*Wallet, error)
	Delete(ctx context.Context, id string) error
	UpdateBalance(ctx context.Context, id string, newBalance float64) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, txn *transaction.Transaction) (string, error)
	GetTransactionsByWalletID(ctx context.Context, walletID, typeFilter string, pageNo, pageSize int) ([]*transaction.Transaction, error)
}

type service struct {
	wr   WalletRepository
	ts   TransactionService
	conf config.Conf
}

func NewService(wr WalletRepository, ts TransactionService, conf config.Conf) *service {
	return &service{wr, ts, conf}
}

func (s *service) CreateWallet(ctx context.Context, info *WalletCreationInfo) (string, error) {
	if info.BalanceUpperLimit > s.conf.Wallet.MaxBalance {
		return "", ErrAboveMaximumBalanceLimit
	}

	if info.TransactionUpperLimit > s.conf.Transaction.MaxAmount {
		return "", ErrAboveMaximumTransactionLimit
	}

	if s.checkWalletWithUserIDExists(ctx, info.UserID) {
		return "", ErrWalletWithUserIDExists
	}

	wallet := &Wallet{
		UserID:                info.UserID,
		Balance:               s.conf.Wallet.InitialBalance,
		BalanceUpperLimit:     info.BalanceUpperLimit,
		TransactionUpperLimit: info.TransactionUpperLimit,
	}

	return s.wr.Create(ctx, wallet)
}

func (s *service) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	return s.wr.Read(ctx, id)
}

func (s *service) DeleteWallet(ctx context.Context, id string) error {
	_, err := s.wr.Read(ctx, id)
	if err != nil {
		return err
	}

	return s.wr.Delete(ctx, id)
}

func (s *service) CreateTransaction(ctx context.Context, info *TransactionCreationInfo) (string, error) {
	if info.TransactionType != Deposit && info.TransactionType != Withdrawal {
		return "", ErrInvalidTransactionType
	}

	w, err := s.wr.Read(ctx, info.WalletID)
	if err != nil && errors.Is(err, ErrWalletNotFound) {
		return "", err
	}

	err = s.processTransaction(ctx, w, info.Amount, info.TransactionType)
	if err != nil {
		return "", err
	}

	return s.ts.CreateTransaction(ctx, s.transactionFromTransactionCreationInfo(info))
}

func (s *service) GetTransactions(ctx context.Context, walletID, typeFilter string, pageNo, pageSize int) ([]*transaction.Transaction, error) {
	if typeFilter != "" && typeFilter != Deposit && typeFilter != Withdrawal {
		return nil, ErrInvalidTransactionType
	}

	_, err := s.wr.Read(ctx, walletID)
	if err != nil {
		return nil, err
	}

	return s.ts.GetTransactionsByWalletID(ctx, walletID, typeFilter, pageNo, pageSize)
}

func (s *service) checkWalletWithUserIDExists(ctx context.Context, userID string) bool {
	w, err := s.wr.ReadByUserID(ctx, userID)
	return w != nil && err == nil
}

func (s *service) processTransaction(ctx context.Context, w *Wallet, txnAmount float64, txnType string) error {
	if txnAmount > w.TransactionUpperLimit {
		return ErrAboveMaximumTransactionLimit
	}

	if txnAmount < s.conf.Transaction.MinAmount {
		return ErrBelowMinimumTransactionLimit
	}

	if txnType == Deposit && w.Balance+txnAmount > w.BalanceUpperLimit {
		return ErrAboveMaximumBalanceLimit
	}

	if txnType == Withdrawal && w.Balance-txnAmount < s.conf.Wallet.MinBalance {
		return ErrInsufficientBalance
	}

	var newBalance float64
	if txnType == Deposit {
		newBalance = w.Balance + txnAmount
	} else if txnType == Withdrawal {
		newBalance = w.Balance - txnAmount
	}

	return s.wr.UpdateBalance(ctx, w.ID, newBalance)
}

func (s *service) transactionFromTransactionCreationInfo(info *TransactionCreationInfo) *transaction.Transaction {
	return &transaction.Transaction{
		WalletID: info.WalletID,
		Type:     info.TransactionType,
		Amount:   info.Amount,
	}
}
