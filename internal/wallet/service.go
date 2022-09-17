package wallet

import (
	"context"
	"fmt"

	"github.com/gokcelb/wallet-api/config"
)

var (
	ErrWalletNotFound               = fmt.Errorf("no wallet with the given id exists")
	ErrAboveMaximumBalanceLimit     = fmt.Errorf("wallet balance is above maximum balance limit")
	ErrAboveMaximumTransactionLimit = fmt.Errorf("transaction is above maximum transaction limit")
	ErrBelowMinimumTransactionLimit = fmt.Errorf("transaction is below minimum transaction limit")
)

type Repository interface {
	Create(ctx context.Context, wallet Wallet) (string, error)
	Read(ctx context.Context, id string) (Wallet, error)
	// Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	conf config.Conf
}

func NewService(repo Repository, conf config.Conf) *service {
	return &service{repo, conf}
}

func (s *service) CreateWallet(ctx context.Context, info *WalletCreationInfo) (Wallet, error) {
	if info.BalanceUpperLimit > s.conf.Wallet.MaxBalance {
		return Wallet{}, ErrAboveMaximumBalanceLimit
	}

	if info.TransactionUpperLimit > s.conf.Transaction.MaxAmount {
		return Wallet{}, ErrAboveMaximumTransactionLimit
	}

	wallet := Wallet{
		UserId:                info.UserId,
		Balance:               s.conf.Wallet.InitialBalance,
		BalanceUpperLimit:     info.BalanceUpperLimit,
		TransactionUpperLimit: info.TransactionUpperLimit,
	}

	walletId, err := s.repo.Create(ctx, wallet)
	if err != nil {
		return Wallet{}, err
	}

	wallet.Id = walletId
	return wallet, nil
}

func (s *service) GetWallet(ctx context.Context, id string) (Wallet, error) {
	wallet, err := s.repo.Read(ctx, id)
	if err != nil {
		return Wallet{}, ErrWalletNotFound
	}

	return wallet, nil
}

func (s *service) Delete(id string) error {
	return nil
}
