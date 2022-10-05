package wallet_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gokcelb/wallet-api/internal/wallet"
	"github.com/gokcelb/wallet-api/internal/wallet/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const contentType = "application/json"

type httpErr struct {
	Message string `json:"message"`
}

func (e *httpErr) Error() string {
	return e.Message
}

func createMockWalletService(t *testing.T) *mock.MockWalletService {
	return mock.NewMockWalletService(gomock.NewController(t))
}

func TestHandlerPostWallet(t *testing.T) {
	mockWalletService := createMockWalletService(t)
	h := wallet.NewHandler(mockWalletService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                       string
		givenUserID                string
		givenBalanceUpperLimit     float64
		givenTransactionUpperLimit float64
		mockWSWalletID             string
		mockWSError                error
		expectedStatusCode         int
		expectedResponseBody       interface{}
	}{
		{
			desc:                       "wallet creation info is valid, return new wallet",
			givenUserID:                "1",
			givenBalanceUpperLimit:     1000,
			givenTransactionUpperLimit: 500,
			mockWSWalletID:             "1",
			mockWSError:                nil,
			expectedStatusCode:         201,
			expectedResponseBody:       wallet.PostResponse{"1"},
		},
		{
			desc:                       "balance upper limit is not valid, return error",
			givenUserID:                "2",
			givenBalanceUpperLimit:     30000,
			givenTransactionUpperLimit: 100,
			mockWSWalletID:             "",
			mockWSError:                wallet.ErrAboveMaximumBalanceLimit,
			expectedStatusCode:         422,
			expectedResponseBody:       httpErr{wallet.ErrAboveMaximumBalanceLimit.Error()},
		},
		{
			desc:                       "transaction upper limit is not valid, return error",
			givenUserID:                "3",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 10000,
			mockWSWalletID:             "",
			mockWSError:                wallet.ErrAboveMaximumTransactionLimit,
			expectedStatusCode:         422,
			expectedResponseBody:       httpErr{wallet.ErrAboveMaximumTransactionLimit.Error()},
		},
		{
			desc:                       "wallet with user id already exists, return error",
			givenUserID:                "1",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 1000,
			mockWSWalletID:             "",
			mockWSError:                wallet.ErrWalletWithUserIDExists,
			expectedStatusCode:         422,
			expectedResponseBody:       httpErr{wallet.ErrWalletWithUserIDExists.Error()},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			walletCreationInfo := wallet.WalletCreationInfo{
				UserID:                tC.givenUserID,
				BalanceUpperLimit:     tC.givenBalanceUpperLimit,
				TransactionUpperLimit: tC.givenTransactionUpperLimit,
			}

			mockWalletService.EXPECT().
				CreateWallet(gomock.Any(), &walletCreationInfo).
				Return(tC.mockWSWalletID, tC.mockWSError)

			walletCreationInfoBytes, _ := json.Marshal(walletCreationInfo)
			res, err := http.DefaultClient.Post(
				fmt.Sprintf("%s/wallets", testServer.URL),
				contentType,
				bytes.NewReader(walletCreationInfoBytes),
			)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer res.Body.Close()

			resBodyBytes, _ := io.ReadAll(res.Body)
			expectedResBodyBytes, _ := json.Marshal(tC.expectedResponseBody)

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
			assert.JSONEq(t, string(expectedResBodyBytes), string(resBodyBytes))
		})
	}
}

func TestHandlerGetWallet(t *testing.T) {
	mockWalletService := createMockWalletService(t)
	h := wallet.NewHandler(mockWalletService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                 string
		givenID              string
		mockWSWallet         *wallet.Wallet
		mockWSErr            error
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		{
			desc:    "wallet exists, return wallet",
			givenID: "1",
			mockWSWallet: &wallet.Wallet{
				ID:                    "1",
				UserID:                "1",
				Balance:               0,
				BalanceUpperLimit:     1000,
				TransactionUpperLimit: 100,
			},
			mockWSErr:          nil,
			expectedStatusCode: 200,
			expectedResponseBody: &wallet.Wallet{
				ID:                    "1",
				UserID:                "1",
				Balance:               0,
				BalanceUpperLimit:     1000,
				TransactionUpperLimit: 100,
			},
		},
		{
			desc:                 "wallet does not exist, return error",
			givenID:              "2",
			mockWSWallet:         nil,
			mockWSErr:            wallet.ErrWalletNotFound,
			expectedStatusCode:   404,
			expectedResponseBody: httpErr{wallet.ErrWalletNotFound.Error()},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockWalletService.EXPECT().GetWallet(gomock.Any(), tC.givenID).Return(tC.mockWSWallet, tC.mockWSErr)

			res, err := http.DefaultClient.Get(fmt.Sprintf("%s/wallets/%s", testServer.URL, tC.givenID))
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer res.Body.Close()

			resBodyBytes, _ := io.ReadAll(res.Body)
			expectedResBodyBytes, _ := json.Marshal(tC.expectedResponseBody)

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
			assert.JSONEq(t, string(expectedResBodyBytes), string(resBodyBytes))
		})
	}
}

func TestHandlerDeleteWallet(t *testing.T) {
	mockService := createMockWalletService(t)
	h := wallet.NewHandler(mockService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                   string
		givenWalletID          string
		mockWSErr              error
		expectedResponseStatus int
	}{
		{
			desc:                   "wallet id exists, return success",
			givenWalletID:          "1",
			mockWSErr:              nil,
			expectedResponseStatus: 204,
		},
		{
			desc:                   "wallet id does not exist, return error",
			givenWalletID:          "2",
			mockWSErr:              wallet.ErrWalletNotFound,
			expectedResponseStatus: 404,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockService.EXPECT().DeleteWallet(gomock.Any(), tC.givenWalletID).Return(tC.mockWSErr)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/wallets/%s", testServer.URL, tC.givenWalletID), nil)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			res, err := testServer.Client().Do(req)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			assert.Equal(t, res.StatusCode, tC.expectedResponseStatus)
		})
	}
}

func TestHandlerCreateTransaction(t *testing.T) {
	mockWalletService := createMockWalletService(t)
	h := wallet.NewHandler(mockWalletService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                       string
		givenWalletID              string
		givenTransactionType       string
		givenAmount                float64
		mockWSTransactionID        string
		mockWSErr                  error
		expectedResponseStatusCode int
		expectedResponseBody       interface{}
	}{
		{
			desc:                       "transaction creation info is valid, return transaction",
			givenWalletID:              "1",
			givenTransactionType:       "deposit",
			givenAmount:                300,
			mockWSTransactionID:        "1",
			mockWSErr:                  nil,
			expectedResponseStatusCode: 201,
			expectedResponseBody:       wallet.PostResponse{"1"},
		},
		{
			desc:                       "wallet id does not exist, return error",
			givenWalletID:              "2",
			givenTransactionType:       "deposit",
			givenAmount:                300,
			mockWSTransactionID:        "",
			mockWSErr:                  wallet.ErrWalletNotFound,
			expectedResponseStatusCode: 404,
			expectedResponseBody:       httpErr{wallet.ErrWalletNotFound.Error()},
		},
		{
			desc:                       "transaction type is not valid, return error",
			givenWalletID:              "1",
			givenTransactionType:       "some type",
			givenAmount:                300,
			mockWSTransactionID:        "",
			mockWSErr:                  wallet.ErrInvalidTransactionType,
			expectedResponseStatusCode: 400,
			expectedResponseBody:       httpErr{wallet.ErrInvalidTransactionType.Error()},
		},
		{
			desc:                       "transaction amount above requirements, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                20000,
			mockWSTransactionID:        "",
			mockWSErr:                  wallet.ErrAboveMaximumTransactionLimit,
			expectedResponseStatusCode: 422,
			expectedResponseBody:       httpErr{wallet.ErrAboveMaximumTransactionLimit.Error()},
		},
		{
			desc:                       "transaction amount above requirements, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                1,
			mockWSTransactionID:        "",
			mockWSErr:                  wallet.ErrBelowMinimumTransactionLimit,
			expectedResponseStatusCode: 422,
			expectedResponseBody:       httpErr{wallet.ErrBelowMinimumTransactionLimit.Error()},
		},
		{
			desc:                       "insufficient balance, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                500,
			mockWSTransactionID:        "",
			mockWSErr:                  wallet.ErrInsufficientBalance,
			expectedResponseStatusCode: 422,
			expectedResponseBody:       httpErr{wallet.ErrInsufficientBalance.Error()},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			transactionCreationInfo := wallet.TransactionCreationInfo{
				WalletID:        tC.givenWalletID,
				TransactionType: tC.givenTransactionType,
				Amount:          tC.givenAmount,
			}

			mockWalletService.EXPECT().
				CreateTransaction(gomock.Any(), &transactionCreationInfo).
				Return(tC.mockWSTransactionID, tC.mockWSErr)

			transactionCreationInfoBytes, _ := json.Marshal(transactionCreationInfo)
			res, err := testServer.Client().Post(
				fmt.Sprintf("%s/wallets/%s/transactions", testServer.URL, tC.givenWalletID),
				contentType,
				bytes.NewReader(transactionCreationInfoBytes),
			)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer res.Body.Close()

			resBodyBytes, _ := io.ReadAll(res.Body)
			expectedResBodyBytes, _ := json.Marshal(tC.expectedResponseBody)

			assert.Equal(t, tC.expectedResponseStatusCode, res.StatusCode)
			assert.JSONEq(t, string(expectedResBodyBytes), string(resBodyBytes))
		})
	}
}
