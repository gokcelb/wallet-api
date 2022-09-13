package wallet_test

import (
	"encoding/json"
	"errors"
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

type errWalletNotFound struct {
	Message string `json:"message"`
}

func createMockService(t *testing.T) *mock_wallet.MockService {
	return mock_wallet.NewMockService(gomock.NewController(t))
}

func TestGetWallet(t *testing.T) {
	validWalletId := "1"
	invalidWalletId := "2"
	mockWallet := wallet.Wallet{
		Id:                    "1",
		UserId:                "1",
		Balance:               0,
		BalanceUpperLimit:     100,
		TransactionUpperLimit: 100,
	}

	e := echo.New()
	mockService := createMockService(t)
	h := wallet.New(mockService)
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
			mockSvcError:         errors.New(wallet.ErrWalletNotFound.Error()),
			expectedStatusCode:   404,
			expectedResponseBody: errWalletNotFound{wallet.ErrWalletNotFound.Error()},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockService.EXPECT().Get(tC.requestedId).Return(tC.mockSvcWallet, tC.mockSvcError)

			url := fmt.Sprintf("%s/wallets/%s", testServer.URL, tC.requestedId)
			res, err := http.DefaultClient.Get(url)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer res.Body.Close()

			var resBody wallet.Wallet
			resBodyBytes, _ := io.ReadAll(res.Body)
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				assert.Fail(t, err.Error())
			}

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tC.mockSvcWallet, resBody)
		})
	}
}
