package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/transaction"
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

type service struct {
	wr   WalletRepository
	conf config.Conf
}

func NewService(wr WalletRepository, conf config.Conf) *service {
	return &service{wr, conf}
}

func (s *service) CreateWallet(ctx context.Context, info *WalletCreationInfo) (Wallet, error) {
	if info.BalanceUpperLimit > s.conf.Wallet.MaxBalance {
		return Wallet{}, ErrAboveMaximumBalanceLimit
	}

	if info.TransactionUpperLimit > s.conf.Transaction.MaxAmount {
		return Wallet{}, ErrAboveMaximumTransactionLimit
	}

	if s.checkIfWalletWithUserIDExists(ctx, info.UserID) {
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

func (s *service) checkIfWalletWithUserIDExists(ctx context.Context, userID string) bool {
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

func (s *service) CreateTransaction(ctx context.Context, info *TransactionCreationInfo) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}
