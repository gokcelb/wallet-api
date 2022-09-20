package wallet_test

import (
	"context"
	"errors"
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

func TestServiceCreateWalletWithValidWalletCreationInfo(t *testing.T) {
	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	walletCreationInfo := wallet.WalletCreationInfo{
		UserId:                "1",
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}
	convertedWallet := wallet.Wallet{
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}
	expectedWallet := wallet.Wallet{
		Id:                    "1",
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}

	mockRepository.
		EXPECT().
		ReadByUserId(context.TODO(), walletCreationInfo.UserId).
		Return(wallet.Wallet{}, nil)

	mockRepository.
		EXPECT().
		Create(context.TODO(), convertedWallet).
		Return(expectedWallet.Id, nil)

	result, err := s.CreateWallet(context.TODO(), &walletCreationInfo)

	assert.Equal(t, expectedWallet, result)
	assert.Nil(t, err)
}

func TestServiceCreateWalletWithInvalidLimit(t *testing.T) {
	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	testCases := []struct {
		desc                    string
		givenWalletCreationInfo wallet.WalletCreationInfo
		expectedWallet          wallet.Wallet
		expectedErr             error
	}{
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
			wallet, err := s.CreateWallet(context.TODO(), &tC.givenWalletCreationInfo)

			assert.Equal(t, tC.expectedWallet, wallet)
			assert.Equal(t, tC.expectedErr, err)
		})
	}
}

func TestServiceCreateWalletWithExistingUserId(t *testing.T) {
	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	givenWalletCreationInfo := wallet.WalletCreationInfo{
		UserId:                "1",
		BalanceUpperLimit:     10000,
		TransactionUpperLimit: 1000,
	}
	existingWallet := wallet.Wallet{
		Id:                    "1",
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     10000,
		TransactionUpperLimit: 1000,
	}

	mockRepository.
		EXPECT().
		ReadByUserId(context.TODO(), givenWalletCreationInfo.UserId).
		Return(existingWallet, nil)

	result, err := s.CreateWallet(context.TODO(), &givenWalletCreationInfo)

	assert.Equal(t, wallet.Wallet{}, result)
	assert.Equal(t, wallet.ErrWalletWithUserIdExists, err)
}

func TestServiceGetWallet(t *testing.T) {
	validWalletId := "1"
	invalidWalletId := "2"
	mockWallet := wallet.Wallet{
		Id:                    validWalletId,
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}

	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	testCases := []struct {
		desc           string
		givenId        string
		mockRepoWallet wallet.Wallet
		mockRepoError  error
		expectedWallet wallet.Wallet
		expectedError  error
	}{
		{
			desc:           "wallet exists, return wallet",
			givenId:        validWalletId,
			mockRepoWallet: mockWallet,
			mockRepoError:  nil,
			expectedWallet: mockWallet,
			expectedError:  nil,
		},
		{
			desc:           "wallet does not exist, return error",
			givenId:        invalidWalletId,
			mockRepoWallet: wallet.Wallet{},
			mockRepoError:  errors.New(""),
			expectedWallet: wallet.Wallet{},
			expectedError:  wallet.ErrWalletNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepository.
				EXPECT().
				Read(context.TODO(), tC.givenId).
				Return(tC.mockRepoWallet, tC.mockRepoError)

			wallet, err := s.GetWallet(context.TODO(), tC.givenId)

			assert.Equal(t, tC.expectedWallet, wallet)
			assert.Equal(t, tC.expectedError, err)
		})
	}
}

func TestServiceDeleteWallet(t *testing.T) {
	mockRepository := createMockRepository(t)
	s := wallet.NewService(mockRepository, getConf())

	testCases := []struct {
		desc                     string
		givenWalletID            string
		mockRepoReadWalletWallet wallet.Wallet
		mockRepoReadWalletErr    error
		mockRepoDeleteWalletErr  error
		expectedErr              error
	}{
		{
			desc:                     "wallet id exists, delete wallet",
			givenWalletID:            "1",
			mockRepoReadWalletWallet: wallet.Wallet{Id: "1"},
			mockRepoReadWalletErr:    nil,
			mockRepoDeleteWalletErr:  nil,
			expectedErr:              nil,
		},
		{
			desc:                     "wallet id does not exist, return error",
			givenWalletID:            "2",
			mockRepoReadWalletWallet: wallet.Wallet{},
			mockRepoReadWalletErr:    wallet.ErrWalletNotFound,
			expectedErr:              wallet.ErrWalletNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepository.
				EXPECT().
				Read(context.TODO(), tC.givenWalletID).
				Return(tC.mockRepoReadWalletWallet, tC.mockRepoReadWalletErr)

			if tC.mockRepoReadWalletErr == nil {
				mockRepository.
					EXPECT().
					Delete(context.TODO(), tC.givenWalletID).
					Return(tC.mockRepoDeleteWalletErr)
			}

			err := s.DeleteWallet(context.TODO(), tC.givenWalletID)

			assert.Equal(t, tC.expectedErr, err)
		})
	}
}
