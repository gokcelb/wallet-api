package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var badRequestErrors = []error{
	ErrAboveMaximumBalanceLimit,
	ErrAboveMaximumTransactionLimit,
	ErrBelowMinimumTransactionLimit,
}

type Service interface {
	Create(info *WalletCreationInfo) (Wallet, error)
	Get(id string) (Wallet, error)
	Delete(id string) error
}

type handler struct {
	svc Service
}

type WalletCreationInfo struct {
	UserId                string  `json:"userId"`
	BalanceUpperLimit     float64 `json:"balanceUpperLimit"`
	TransactionUpperLimit float64 `json:"transactionUpperLimit"`
}

type TransactionCreationInfo struct {
	WalletId        string  `json:"walletId"`
	TransactionType string  `json:"type"`
	Amount          float64 `json:"amount"`
}

func NewHandler(svc Service) *handler {
	return &handler{svc}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/wallets", h.CreateWallet)
	e.GET("/wallets/:id", h.GetWallet)
	e.DELETE("/wallets/:id", h.DeleteWallet)

	e.POST("/wallets/:id/transactions", h.CreateTransaction)
	e.GET("/wallets/:id/transactions", h.GetTransactions)
	e.GET("/wallets/:id/transactions/:transactionId", h.GetTransaction)
}

// 201 => successfully created
// 400 => balanceUpperLimit and transactionUpperLimit may not be convertible to float
// 400 => balanceUpperLimit and transactionUpperLimit may not meet config requirements
// 500 => any other error
func (h *handler) CreateWallet(c echo.Context) error {
	var info WalletCreationInfo
	if err := c.Bind(&info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	wallet, err := h.svc.Create(&info)
	if err != nil && isBadRequest(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}

// 200 => successfully read
// 404 => wallet with given id may not exist
// 500 => any other error
func (h *handler) GetWallet(c echo.Context) error {
	wallet, err := h.svc.Get(c.Param("id"))
	if err != nil && err == ErrWalletNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, wallet)
}

// 204 => successfully deleted
// 404 => wallet with given id may not exist
// 500 => any other error
func (h *handler) DeleteWallet(c echo.Context) error {
	err := h.svc.Delete(c.Param("id"))
	if err != nil && err == ErrWalletNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return nil
}

// 201 => successfully created
// 404 => wallet not found
// 400 => type not valid
// 400 => transaction amount does not meet config requirements
// 500 => any other error
func (h *handler) CreateTransaction(c echo.Context) error {
	return nil
}

// 200 => successfully read
// 404 => wallet not found
// 500 => any other error
// * they should be able to filter by type
// * they should be able to paginate
// 400 => type not valid
// 400 => invalid pagination paramaters
func (h *handler) GetTransactions(c echo.Context) error {
	return nil
}

// 200 => successfully read
// 404 => wallet or transaction not found
// 500 => any other error
func (h *handler) GetTransaction(c echo.Context) error {
	return nil
}

func isBadRequest(err error) bool {
	return ContainsError(err, badRequestErrors)
}

func ContainsError(err error, errList []error) bool {
	for _, e := range errList {
		if err == e {
			return true
		}
	}

	return false
}
