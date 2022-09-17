package wallet

import (
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
	Create(Wallet) (Wallet, error)
	Read(id string) (Wallet, error)
	Delete(id string) error
}

type service struct {
	repo Repository
	conf config.Conf
}

func NewService(repo Repository, conf config.Conf) *service {
	return &service{repo, conf}
}

func (s *service) Create(info *WalletCreationInfo) (Wallet, error) {
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

	wallet, err := s.repo.Create(wallet)
	if err != nil {
		return Wallet{}, err
	}

	return wallet, nil
}

func (s *service) Get(id string) (Wallet, error) {
	wallet, err := s.repo.Read(id)
	if err != nil {
		return Wallet{}, ErrWalletNotFound
	}

	return wallet, nil
}

func (s *service) Delete(id string) error {
	return nil
}
