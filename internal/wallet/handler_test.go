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

type httpPostRes struct {
	ID string `json:"id"`
}

type httpErr struct {
	Message string `json:"message"`
}

func (e *httpErr) Error() string {
	return e.Message
}

func createMockWalletService(t *testing.T) *mock.MockWalletService {
	return mock.NewMockWalletService(gomock.NewController(t))
}

func TestHandlerCreateWallet(t *testing.T) {
	newWallet := wallet.Wallet{
		ID:                    "1",
		UserID:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 500,
	}

	mockService := createMockWalletService(t)
	h := wallet.NewHandler(mockService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                       string
		givenUserID                string
		givenBalanceUpperLimit     float64
		givenTransactionUpperLimit float64
		mockSvcWallet              wallet.Wallet
		mockSvcError               error
		expectedStatusCode         int
		expectedResponseBody       interface{}
	}{
		{
			desc:                       "wallet creation info is valid, return new wallet",
			givenUserID:                "1",
			givenBalanceUpperLimit:     1000,
			givenTransactionUpperLimit: 500,
			mockSvcWallet:              newWallet,
			mockSvcError:               nil,
			expectedStatusCode:         201,
			expectedResponseBody:       newWallet,
		},
		{
			desc:                       "balance upper limit is not valid, return error",
			givenUserID:                "2",
			givenBalanceUpperLimit:     30000,
			givenTransactionUpperLimit: 100,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrAboveMaximumBalanceLimit,
			expectedStatusCode:         422,
			expectedResponseBody:       httpErr{wallet.ErrAboveMaximumBalanceLimit.Error()},
		},
		{
			desc:                       "transaction upper limit is not valid, return error",
			givenUserID:                "3",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 10000,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrAboveMaximumTransactionLimit,
			expectedStatusCode:         422,
			expectedResponseBody:       httpErr{wallet.ErrAboveMaximumTransactionLimit.Error()},
		},
		{
			desc:                       "wallet with user id already exists, return error",
			givenUserID:                "1",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 1000,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrWalletWithUserIDExists,
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
			mockService.
				EXPECT().
				CreateWallet(gomock.Any(), &walletCreationInfo).
				Return(tC.mockSvcWallet, tC.mockSvcError)

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
	validWalletId := "1"
	invalidWalletId := "2"
	mockWallet := wallet.Wallet{
		ID:                    "1",
		UserID:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}

	e := echo.New()
	mockService := createMockWalletService(t)
	h := wallet.NewHandler(mockService)
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                 string
		requestedId          string
		mockSvcWallet        wallet.Wallet
		mockSvcError         error
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		{
			desc:                 "wallet exists, return wallet",
			requestedId:          validWalletId,
			mockSvcWallet:        mockWallet,
			mockSvcError:         nil,
			expectedStatusCode:   200,
			expectedResponseBody: mockWallet,
		},
		{
			desc:                 "wallet does not exist, return error",
			requestedId:          invalidWalletId,
			mockSvcWallet:        wallet.Wallet{},
			mockSvcError:         wallet.ErrWalletNotFound,
			expectedStatusCode:   404,
			expectedResponseBody: httpErr{wallet.ErrWalletNotFound.Error()},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockService.
				EXPECT().
				GetWallet(gomock.Any(), tC.requestedId).
				Return(tC.mockSvcWallet, tC.mockSvcError)

			url := fmt.Sprintf("%s/wallets/%s", testServer.URL, tC.requestedId)
			res, err := http.DefaultClient.Get(url)
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
	e := echo.New()
	mockService := createMockWalletService(t)
	h := wallet.NewHandler(mockService)
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                   string
		givenWalletID          string
		mockSvcError           error
		expectedResponseStatus int
	}{
		{
			desc:                   "wallet id exists, return success",
			givenWalletID:          "1",
			mockSvcError:           nil,
			expectedResponseStatus: 204,
		},
		{
			desc:                   "wallet id does not exist, return error",
			givenWalletID:          "2",
			mockSvcError:           wallet.ErrWalletNotFound,
			expectedResponseStatus: 404,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockService.
				EXPECT().
				DeleteWallet(gomock.Any(), tC.givenWalletID).
				Return(tC.mockSvcError)

			url := fmt.Sprintf("%s/wallets/%s", testServer.URL, tC.givenWalletID)
			req, err := http.NewRequest("DELETE", url, nil)
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
		mockSvcTransactionID       string
		mockSvcError               error
		expectedResponseStatusCode int
		expectedResponseBody       interface{}
	}{
		{
			desc:                       "transaction creation info is valid, return transaction",
			givenWalletID:              "1",
			givenTransactionType:       "deposit",
			givenAmount:                300,
			mockSvcTransactionID:       "1",
			mockSvcError:               nil,
			expectedResponseStatusCode: 201,
			expectedResponseBody:       httpPostRes{"1"},
		},
		{
			desc:                       "wallet id does not exist, return error",
			givenWalletID:              "2",
			givenTransactionType:       "deposit",
			givenAmount:                300,
			mockSvcTransactionID:       "",
			mockSvcError:               wallet.ErrWalletNotFound,
			expectedResponseStatusCode: 404,
			expectedResponseBody:       httpErr{wallet.ErrWalletNotFound.Error()},
		},
		{
			desc:                       "transaction type is not valid, return error",
			givenWalletID:              "1",
			givenTransactionType:       "some type",
			givenAmount:                300,
			mockSvcTransactionID:       "",
			mockSvcError:               wallet.ErrInvalidTransactionType,
			expectedResponseStatusCode: 400,
			expectedResponseBody:       httpErr{wallet.ErrInvalidTransactionType.Error()},
		},
		{
			desc:                       "transaction amount above requirements, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                20000,
			mockSvcTransactionID:       "",
			mockSvcError:               wallet.ErrAboveMaximumTransactionLimit,
			expectedResponseStatusCode: 422,
			expectedResponseBody: httpErr{
				Message: wallet.ErrAboveMaximumTransactionLimit.Error(),
			},
		},
		{
			desc:                       "transaction amount above requirements, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                1,
			mockSvcTransactionID:       "",
			mockSvcError:               wallet.ErrBelowMinimumTransactionLimit,
			expectedResponseStatusCode: 422,
			expectedResponseBody: httpErr{
				Message: wallet.ErrBelowMinimumTransactionLimit.Error(),
			},
		},
		{
			desc:                       "insufficient balance, return error",
			givenWalletID:              "1",
			givenTransactionType:       "withdrawal",
			givenAmount:                500,
			mockSvcTransactionID:       "",
			mockSvcError:               wallet.ErrInsufficientBalance,
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

			mockWalletService.
				EXPECT().
				CreateTransaction(gomock.Any(), &transactionCreationInfo).
				Return(tC.mockSvcTransactionID, tC.mockSvcError)

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
