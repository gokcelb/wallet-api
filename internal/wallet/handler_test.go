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
	mock_wallet "github.com/gokcelb/wallet-api/mocks/wallet"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const contentType = "application/json"

type testHandlerErr struct {
	Message string `json:"message"`
}

func (e *testHandlerErr) Error() string {
	return e.Message
}

func createMockService(t *testing.T) *mock_wallet.MockService {
	return mock_wallet.NewMockService(gomock.NewController(t))
}

func TestHandlerPostWallet(t *testing.T) {
	newWallet := wallet.Wallet{
		Id:                    "1",
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 500,
	}

	e := echo.New()
	mockService := createMockService(t)
	h := wallet.NewHandler(mockService)
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                       string
		givenUserId                string
		givenBalanceUpperLimit     float64
		givenTransactionUpperLimit float64
		mockSvcWallet              wallet.Wallet
		mockSvcError               error
		expectedStatusCode         int
		expectedResponseBody       interface{}
	}{
		{
			desc:                       "wallet creation info is valid, return new wallet",
			givenUserId:                "1",
			givenBalanceUpperLimit:     1000,
			givenTransactionUpperLimit: 500,
			mockSvcWallet:              newWallet,
			mockSvcError:               nil,
			expectedStatusCode:         201,
			expectedResponseBody:       newWallet,
		},
		{
			desc:                       "balance upper limit is not valid, return error",
			givenUserId:                "2",
			givenBalanceUpperLimit:     30000,
			givenTransactionUpperLimit: 100,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrAboveMaximumBalanceLimit,
			expectedStatusCode:         400,
			expectedResponseBody:       testHandlerErr{wallet.ErrAboveMaximumBalanceLimit.Error()},
		},
		{
			desc:                       "transaction upper limit is not valid, return error",
			givenUserId:                "3",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 10000,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrAboveMaximumTransactionLimit,
			expectedStatusCode:         400,
			expectedResponseBody:       testHandlerErr{wallet.ErrAboveMaximumTransactionLimit.Error()},
		},
		{
			desc:                       "wallet with user id already exists, return error",
			givenUserId:                "1",
			givenBalanceUpperLimit:     3000,
			givenTransactionUpperLimit: 1000,
			mockSvcWallet:              wallet.Wallet{},
			mockSvcError:               wallet.ErrWalletWithUserIdExists,
			expectedStatusCode:         400,
			expectedResponseBody:       testHandlerErr{wallet.ErrWalletWithUserIdExists.Error()},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			walletCreationInfo := wallet.WalletCreationInfo{
				UserId:                tC.givenUserId,
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
		Id:                    "1",
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     1000,
		TransactionUpperLimit: 100,
	}

	e := echo.New()
	mockService := createMockService(t)
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
			expectedResponseBody: testHandlerErr{wallet.ErrWalletNotFound.Error()},
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
