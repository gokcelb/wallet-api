package transaction_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/gokcelb/wallet-api/internal/transaction/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type httpErr struct {
	Message string `json:"message"`
}

func createMockTransactionService(t *testing.T) *mock.MockTransactionService {
	return mock.NewMockTransactionService(gomock.NewController(t))
}

func TestHandlerGetTransaction(t *testing.T) {
	mockTransactionService := createMockTransactionService(t)
	h := transaction.NewHandler(mockTransactionService)

	e := echo.New()
	h.RegisterRoutes(e)
	testServer := httptest.NewServer(e.Server.Handler)
	defer testServer.Close()

	testCases := []struct {
		desc                       string
		givenID                    string
		mockTxnSvcTxn              *transaction.Transaction
		mockTxnSvcErr              error
		expectedResponseStatusCode int
		expectedResponseBody       interface{}
	}{
		{
			desc:    "transaction id exists, return transaction",
			givenID: "1",
			mockTxnSvcTxn: &transaction.Transaction{
				ID:       "1",
				WalletID: "1",
				Type:     "deposit",
				Amount:   200,
			},
			mockTxnSvcErr:              nil,
			expectedResponseStatusCode: 200,
			expectedResponseBody: &transaction.Transaction{
				ID:       "1",
				WalletID: "1",
				Type:     "deposit",
				Amount:   200,
			},
		},
		{
			desc:                       "transaction id does not exist, return error",
			givenID:                    "2",
			mockTxnSvcTxn:              nil,
			mockTxnSvcErr:              transaction.ErrTransactionNotFound,
			expectedResponseStatusCode: 404,
			expectedResponseBody:       httpErr{transaction.ErrTransactionNotFound.Error()},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockTransactionService.EXPECT().
				GetTransaction(gomock.Any(), tC.givenID).
				Return(tC.mockTxnSvcTxn, tC.mockTxnSvcErr)

			res, err := testServer.Client().Get(fmt.Sprintf("%s/transactions/%s", testServer.URL, tC.givenID))
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
