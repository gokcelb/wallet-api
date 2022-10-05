package wallet_test

import (
	"context"
	"testing"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/gokcelb/wallet-api/internal/wallet"
	"github.com/gokcelb/wallet-api/internal/wallet/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMockWalletRepository(t *testing.T) *mock.MockWalletRepository {
	return mock.NewMockWalletRepository(gomock.NewController(t))
}

func createMockTransactionService(t *testing.T) *mock.MockTransactionService {
	return mock.NewMockTransactionService(gomock.NewController(t))
}

func getConf() config.Conf {
	conf, err := config.Read("../../.config/dev.json")
	if err != nil {
		panic(err)
	}

	return conf
}

func TestServiceCreateWalletWithValidWalletCreationInfo(t *testing.T) {
	mockWalletRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockWalletRepository, nil, getConf())

	walletCreationInfo := &wallet.WalletCreationInfo{
		UserID:                "1",
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}
	convertedWallet := &wallet.Wallet{
		UserID:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}
	expectedWalletID := "1"

	mockWalletRepository.EXPECT().
		ReadByUserID(context.TODO(), walletCreationInfo.UserID).
		Return(nil, wallet.ErrWalletNotFound)

	mockWalletRepository.EXPECT().Create(context.TODO(), convertedWallet).Return(expectedWalletID, nil)

	id, err := s.CreateWallet(context.TODO(), walletCreationInfo)

	assert.Equal(t, expectedWalletID, id)
	assert.Nil(t, err)
}

func TestServiceCreateWalletWithInvalidLimit(t *testing.T) {
	mockWalletRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockWalletRepository, nil, getConf())

	testCases := []struct {
		desc                    string
		givenWalletCreationInfo wallet.WalletCreationInfo
		expectedWalletID        string
		expectedErr             error
	}{
		{
			desc: "balance upper limit is not valid, return error",
			givenWalletCreationInfo: wallet.WalletCreationInfo{
				UserID:                "1",
				BalanceUpperLimit:     100000,
				TransactionUpperLimit: 500,
			},
			expectedWalletID: "",
			expectedErr:      wallet.ErrAboveMaximumBalanceLimit,
		},
		{
			desc: "transaction upper limit is not valid, return error",
			givenWalletCreationInfo: wallet.WalletCreationInfo{
				UserID:                "1",
				BalanceUpperLimit:     10000,
				TransactionUpperLimit: 10000,
			},
			expectedWalletID: "",
			expectedErr:      wallet.ErrAboveMaximumTransactionLimit,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			id, err := s.CreateWallet(context.TODO(), &tC.givenWalletCreationInfo)

			assert.Equal(t, tC.expectedWalletID, id)
			assert.Equal(t, tC.expectedErr, err)
		})
	}
}

func TestServiceCreateWalletWithExistingUserID(t *testing.T) {
	mockRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockRepository, nil, getConf())

	givenWalletCreationInfo := &wallet.WalletCreationInfo{
		UserID:                "1",
		BalanceUpperLimit:     10000,
		TransactionUpperLimit: 1000,
	}
	existingWallet := &wallet.Wallet{
		ID:                    "1",
		UserID:                "1",
		Balance:               0,
		BalanceUpperLimit:     10000,
		TransactionUpperLimit: 1000,
	}

	mockRepository.EXPECT().ReadByUserID(context.TODO(), givenWalletCreationInfo.UserID).Return(existingWallet, nil)

	id, err := s.CreateWallet(context.TODO(), givenWalletCreationInfo)

	assert.Empty(t, id)
	assert.ErrorIs(t, err, wallet.ErrWalletWithUserIDExists)
}

func TestServiceGetWallet(t *testing.T) {
	validWalletId := "1"
	invalidWalletId := "2"
	mockWallet := &wallet.Wallet{
		ID:                    validWalletId,
		UserID:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}

	mockRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockRepository, nil, getConf())

	testCases := []struct {
		desc           string
		givenId        string
		mockRepoWallet *wallet.Wallet
		mockRepoError  error
		expectedWallet *wallet.Wallet
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
			mockRepoWallet: nil,
			mockRepoError:  wallet.ErrWalletNotFound,
			expectedWallet: nil,
			expectedError:  wallet.ErrWalletNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepository.EXPECT().Read(context.TODO(), tC.givenId).Return(tC.mockRepoWallet, tC.mockRepoError)

			w, err := s.GetWallet(context.TODO(), tC.givenId)

			assert.Equal(t, tC.expectedWallet, w)
			assert.Equal(t, tC.expectedError, err)
		})
	}
}

func TestServiceDeleteWallet(t *testing.T) {
	mockRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockRepository, nil, getConf())

	testCases := []struct {
		desc                     string
		givenWalletID            string
		mockRepoReadWalletWallet *wallet.Wallet
		mockRepoReadWalletErr    error
		mockRepoDeleteWalletErr  error
		expectedErr              error
	}{
		{
			desc:                     "wallet id exists, delete wallet",
			givenWalletID:            "1",
			mockRepoReadWalletWallet: &wallet.Wallet{ID: "1"},
			mockRepoReadWalletErr:    nil,
			mockRepoDeleteWalletErr:  nil,
			expectedErr:              nil,
		},
		{
			desc:                     "wallet id does not exist, return error",
			givenWalletID:            "2",
			mockRepoReadWalletWallet: nil,
			mockRepoReadWalletErr:    wallet.ErrWalletNotFound,
			expectedErr:              wallet.ErrWalletNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepository.EXPECT().
				Read(context.TODO(), tC.givenWalletID).
				Return(tC.mockRepoReadWalletWallet, tC.mockRepoReadWalletErr)

			if tC.mockRepoReadWalletErr == nil {
				mockRepository.EXPECT().Delete(context.TODO(), tC.givenWalletID).Return(tC.mockRepoDeleteWalletErr)
			}

			err := s.DeleteWallet(context.TODO(), tC.givenWalletID)

			assert.Equal(t, tC.expectedErr, err)
		})
	}
}

func TestServiceCreateTransactionWithValidTransactionCreationInfo(t *testing.T) {
	mockRepository := createMockWalletRepository(t)
	mockTransactionService := createMockTransactionService(t)
	s := wallet.NewService(mockRepository, mockTransactionService, getConf())

	givenTransactionCreationInfo := &wallet.TransactionCreationInfo{
		WalletID:        "1",
		TransactionType: "withdrawal",
		Amount:          300,
	}
	mockRepoGetWalletWallet := &wallet.Wallet{
		ID:                    "1",
		UserID:                "1",
		Balance:               500,
		BalanceUpperLimit:     10000,
		TransactionUpperLimit: 1000,
	}
	convertedTxnSvcTransaction := &transaction.Transaction{
		WalletID: "1",
		Type:     "withdrawal",
		Amount:   300,
	}
	mockTxnSvcTransactionID := "1"

	mockRepository.EXPECT().
		Read(context.TODO(), givenTransactionCreationInfo.WalletID).
		Return(mockRepoGetWalletWallet, nil)

	mockRepository.EXPECT().
		UpdateBalance(
			context.TODO(),
			givenTransactionCreationInfo.WalletID,
			mockRepoGetWalletWallet.Balance-givenTransactionCreationInfo.Amount,
		).Return(nil)

	mockTransactionService.EXPECT().
		CreateTransaction(context.TODO(), convertedTxnSvcTransaction).
		Return(mockTxnSvcTransactionID, nil)

	txn, err := s.CreateTransaction(context.TODO(), givenTransactionCreationInfo)

	assert.Equal(t, mockTxnSvcTransactionID, txn)
	assert.Nil(t, err)
}

func TestServiceCreateTransactionWithInvalidTransactionCreationInfo(t *testing.T) {
	mockRepository := createMockWalletRepository(t)
	s := wallet.NewService(mockRepository, nil, getConf())

	testCases := []struct {
		desc                         string
		givenTransactionCreationInfo *wallet.TransactionCreationInfo
		mockRepoGetWalletWallet      *wallet.Wallet
		mockRepoGetWalletErr         error
		expectedTransactionID        string
		expectedErr                  error
	}{
		{
			desc: "wallet id does not exist, return error",
			givenTransactionCreationInfo: &wallet.TransactionCreationInfo{
				WalletID:        "2",
				TransactionType: "deposit",
				Amount:          100,
			},
			mockRepoGetWalletWallet: nil,
			mockRepoGetWalletErr:    wallet.ErrWalletNotFound,
			expectedTransactionID:   "",
			expectedErr:             wallet.ErrWalletNotFound,
		},
		{
			desc: "transaction amount is above limit, return error",
			givenTransactionCreationInfo: &wallet.TransactionCreationInfo{
				WalletID:        "1",
				TransactionType: "deposit",
				Amount:          10000,
			},
			mockRepoGetWalletWallet: &wallet.Wallet{
				ID:                    "1",
				UserID:                "1",
				Balance:               0,
				BalanceUpperLimit:     10000,
				TransactionUpperLimit: 1000,
			},
			mockRepoGetWalletErr:  nil,
			expectedTransactionID: "",
			expectedErr:           wallet.ErrAboveMaximumTransactionLimit,
		},
		{
			desc: "transaction amount is below requirement, return error",
			givenTransactionCreationInfo: &wallet.TransactionCreationInfo{
				WalletID:        "1",
				TransactionType: "deposit",
				Amount:          1,
			},
			mockRepoGetWalletWallet: &wallet.Wallet{
				ID:                    "1",
				UserID:                "1",
				Balance:               0,
				BalanceUpperLimit:     10000,
				TransactionUpperLimit: 1000,
			},
			mockRepoGetWalletErr:  nil,
			expectedTransactionID: "",
			expectedErr:           wallet.ErrBelowMinimumTransactionLimit,
		},
		{
			desc: "balance is insufficient, return error",
			givenTransactionCreationInfo: &wallet.TransactionCreationInfo{
				WalletID:        "1",
				TransactionType: "withdrawal",
				Amount:          1000,
			},
			mockRepoGetWalletWallet: &wallet.Wallet{
				ID:                    "1",
				UserID:                "1",
				Balance:               100,
				BalanceUpperLimit:     10000,
				TransactionUpperLimit: 1000,
			},
			mockRepoGetWalletErr:  nil,
			expectedTransactionID: "",
			expectedErr:           wallet.ErrInsufficientBalance,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepository.EXPECT().
				Read(context.TODO(), tC.givenTransactionCreationInfo.WalletID).
				Return(tC.mockRepoGetWalletWallet, tC.mockRepoGetWalletErr)

			id, err := s.CreateTransaction(context.TODO(), tC.givenTransactionCreationInfo)

			assert.Equal(t, tC.expectedTransactionID, id)
			assert.ErrorIs(t, err, tC.expectedErr)
		})
	}
}
