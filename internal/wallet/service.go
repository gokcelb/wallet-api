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
	Withdrawal        = "withdrawal"
)

var (
	ErrWalletNotFound               = fmt.Errorf("no wallet with the given id exists")
	ErrWalletWithUserIDExists       = fmt.Errorf("wallet with user id already exists")
	ErrAboveMaximumBalanceLimit     = fmt.Errorf("wallet balance is above maximum balance limit")
	ErrAboveMaximumTransactionLimit = fmt.Errorf("transaction is above maximum transaction limit")
	ErrBelowMinimumTransactionLimit = fmt.Errorf("transaction is below minimum transaction limit")
	ErrInvalidTransactionType       = fmt.Errorf("transaction type is invalid")
	ErrInsufficientBalance          = fmt.Errorf("balance is insufficient")
)

type WalletRepository interface {
	Create(ctx context.Context, wallet Wallet) (string, error)
	Read(ctx context.Context, id string) (Wallet, error)
	ReadByUserID(ctx context.Context, userId string) (Wallet, error)
	Delete(ctx context.Context, id string) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, txn *transaction.Transaction) (string, error)
}

type service struct {
	wr   WalletRepository
	ts   TransactionService
	conf config.Conf
}

func NewService(wr WalletRepository, ts TransactionService, conf config.Conf) *service {
	return &service{wr, ts, conf}
}

func (s *service) CreateWallet(ctx context.Context, info *WalletCreationInfo) (Wallet, error) {
	if info.BalanceUpperLimit > s.conf.Wallet.MaxBalance {
		return Wallet{}, ErrAboveMaximumBalanceLimit
	}

	if info.TransactionUpperLimit > s.conf.Transaction.MaxAmount {
		return Wallet{}, ErrAboveMaximumTransactionLimit
	}

	if s.checkWalletWithUserIDExists(ctx, info.UserID) {
		return Wallet{}, ErrWalletWithUserIDExists
	}

	wallet := Wallet{
		UserID:                info.UserID,
		Balance:               s.conf.Wallet.InitialBalance,
		BalanceUpperLimit:     info.BalanceUpperLimit,
		TransactionUpperLimit: info.TransactionUpperLimit,
	}

	walletID, err := s.wr.Create(ctx, wallet)
	if err != nil {
		return Wallet{}, err
	}

	wallet.ID = walletID
	return wallet, nil
}

func (s *service) GetWallet(ctx context.Context, id string) (Wallet, error) {
	wallet, err := s.wr.Read(ctx, id)
	if err != nil {
		return Wallet{}, ErrWalletNotFound
	}

	return wallet, nil
}

func (s *service) checkWalletWithUserIDExists(ctx context.Context, userID string) bool {
	wallet, err := s.wr.ReadByUserID(ctx, userID)
	return wallet != (Wallet{}) && err == nil
}

func (s *service) DeleteWallet(ctx context.Context, id string) error {
	_, err := s.wr.Read(ctx, id)
	if err != nil && errors.Is(err, ErrWalletNotFound) {
		return ErrWalletNotFound
	}

	return s.wr.Delete(ctx, id)
}

func (s *service) CreateTransaction(ctx context.Context, info *TransactionCreationInfo) (string, error) {
	w, err := s.wr.Read(ctx, info.WalletID)
	if err != nil && errors.Is(err, ErrWalletNotFound) {
		return "", ErrWalletNotFound
	} else if err != nil {
		return "", err
	}

	err = s.checkTransactionIsProcessable(w, info.Amount, info.TransactionType)
	if err != nil {
		return "", err
	}

	return s.ts.CreateTransaction(ctx, s.transactionFromTransactionCreationInfo(info))
}

func (s *service) checkTransactionIsProcessable(w Wallet, txnAmount float64, txnType string) error {
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

	return nil
}

func (s *service) transactionFromTransactionCreationInfo(info *TransactionCreationInfo) *transaction.Transaction {
	return &transaction.Transaction{
		WalletID: info.WalletID,
		Type:     info.TransactionType,
		Amount:   info.Amount,
	}
}
