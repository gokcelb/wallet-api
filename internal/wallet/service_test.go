package wallet_test

import (
	"context"
	"testing"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/wallet"
	mock_wallet "github.com/gokcelb/wallet-api/mocks/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMockRepository(t *testing.T) *mock_wallet.MockRepository {
	return mock_wallet.NewMockRepository(gomock.NewController(t))
}

func getConf() config.Conf {
	conf, err := config.Read("../../.config/dev.json")
	if err != nil {
		panic(err)
	}

	return conf
}

func TestServiceCreateWallet(t *testing.T) {
	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	testCases := []struct {
		desc                    string
		givenWalletCreationInfo wallet.WalletCreationInfo
		convertedWallet         wallet.Wallet
		mockRepoWalletId        string
		mockRepoErr             error
		expectedWallet          wallet.Wallet
		expectedErr             error
	}{
		{
			desc: "wallet creation info is valid, return wallet",
			givenWalletCreationInfo: wallet.WalletCreationInfo{
				UserId:                "1",
				BalanceUpperLimit:     1000,
				TransactionUpperLimit: 100,
			},
			convertedWallet: wallet.Wallet{
				UserId:                "1",
				Balance:               0,
				BalanceUpperLimit:     1000,
				TransactionUpperLimit: 100,
			},
			mockRepoWalletId: "1",
			mockRepoErr:      nil,
			expectedWallet: wallet.Wallet{
				Id:                    "1",
				UserId:                "1",
				Balance:               0,
				BalanceUpperLimit:     1000,
				TransactionUpperLimit: 100,
			},
			expectedErr: nil,
		},
		{
			desc: "balance upper limit is not valid, return error",
			givenWalletCreationInfo: wallet.WalletCreationInfo{
				UserId:                "1",
				BalanceUpperLimit:     100000,
				TransactionUpperLimit: 500,
			},
			expectedWallet: wallet.Wallet{},
			expectedErr:    wallet.ErrAboveMaximumBalanceLimit,
		},
		{
			desc: "transaction upper limit is not valid, return error",
			givenWalletCreationInfo: wallet.WalletCreationInfo{
				UserId:                "1",
				BalanceUpperLimit:     10000,
				TransactionUpperLimit: 10000,
			},
			expectedWallet: wallet.Wallet{},
			expectedErr:    wallet.ErrAboveMaximumTransactionLimit,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if !returnsNativeServiceError(tC.expectedErr) {
				mockRepository.
					EXPECT().
					Create(context.TODO(), tC.convertedWallet).
					Return(tC.mockRepoWalletId, tC.mockRepoErr)
			}

			wallet, err := s.CreateWallet(context.TODO(), &tC.givenWalletCreationInfo)

			assert.Equal(t, tC.expectedWallet, wallet)
			assert.Equal(t, tC.expectedErr, err)
		})
	}
}

// func TestServiceGet(t *testing.T) {
// 	validWalletId := "1"
// 	invalidWalletId := "2"
// 	mockWallet := wallet.Wallet{
// 		Id:                    validWalletId,
// 		UserId:                "1",
// 		Balance:               0,
// 		BalanceUpperLimit:     100,
// 		TransactionUpperLimit: 100,
// 	}

// 	mockRepository := createMockRepository(t)
// 	s := wallet.NewService(mockRepository, getConf())

// 	testCases := []struct {
// 		desc           string
// 		givenId        string
// 		mockRepoWallet wallet.Wallet
// 		mockRepoError  error
// 		expectedWallet wallet.Wallet
// 		expectedError  error
// 	}{
// 		{
// 			desc:           "wallet exists, return wallet",
// 			givenId:        validWalletId,
// 			mockRepoWallet: mockWallet,
// 			mockRepoError:  nil,
// 			expectedWallet: mockWallet,
// 			expectedError:  nil,
// 		},
// 		{
// 			desc:           "wallet does not exist, return error",
// 			givenId:        invalidWalletId,
// 			mockRepoWallet: wallet.Wallet{},
// 			mockRepoError:  errors.New(""),
// 			expectedWallet: wallet.Wallet{},
// 			expectedError:  wallet.ErrWalletNotFound,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			mockRepository.
// 				EXPECT().
// 				Read(context.TODO(), tC.givenId).
// 				Return(tC.mockRepoWallet, tC.mockRepoError)

// 			wallet, err := s.GetWallet(context.TODO(), tC.givenId)

// 			assert.Equal(t, tC.expectedWallet, wallet)
// 			assert.Equal(t, tC.expectedError, err)
// 		})
// 	}
// }

func returnsNativeServiceError(err error) bool {
	return wallet.ContainsError(err, []error{
		wallet.ErrAboveMaximumBalanceLimit,
		wallet.ErrAboveMaximumTransactionLimit,
		wallet.ErrBelowMinimumTransactionLimit,
	})
}
