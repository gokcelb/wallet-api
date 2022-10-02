package wallet

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

var badRequestErrors = []error{
	ErrInvalidTransactionType,
}

var notFoundErrors = []error{
	ErrWalletNotFound,
}

var unprocessableEntityErrors = []error{
	ErrWalletWithUserIDExists,
	ErrAboveMaximumBalanceLimit,
	ErrAboveMaximumTransactionLimit,
	ErrBelowMinimumTransactionLimit,
	ErrInsufficientBalance,
}

type WalletService interface {
	CreateWallet(ctx context.Context, info *WalletCreationInfo) (Wallet, error)
	GetWallet(ctx context.Context, id string) (Wallet, error)
	DeleteWallet(ctx context.Context, id string) error
	CreateTransaction(ctx context.Context, info *TransactionCreationInfo) (string, error)
}

type handler struct {
	ws WalletService
}

type WalletCreationInfo struct {
	UserID                string  `json:"userId"`
	BalanceUpperLimit     float64 `json:"balanceUpperLimit"`
	TransactionUpperLimit float64 `json:"transactionUpperLimit"`
}

type TransactionCreationInfo struct {
	WalletID        string  `param:"walletId"`
	TransactionType string  `json:"type"`
	Amount          float64 `json:"amount"`
}

type PostResponse struct {
	ID string `json:"id"`
}

func NewHandler(ws WalletService) *handler {
	return &handler{ws}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/wallets", h.CreateWallet)
	e.GET("/wallets/:id", h.GetWallet)
	e.DELETE("/wallets/:id", h.DeleteWallet)

	e.POST("/wallets/:walletId/transactions", h.CreateTransaction)
	e.GET("/wallets/:walletId/transactions", h.GetTransactions)
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

	wallet, err := h.ws.CreateWallet(c.Request().Context(), &info)
	if err != nil && isBadRequest(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && isUnprocessableEntity(err) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}

// 200 => successfully read
// 404 => wallet with given id may not exist
// 500 => any other error
func (h *handler) GetWallet(c echo.Context) error {
	wallet, err := h.ws.GetWallet(c.Request().Context(), c.Param("id"))
	if err != nil && err == ErrWalletNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, wallet)
}

// 204 => successfully deleted
// 404 => wallet with given id may not exist
// 500 => any other error
func (h *handler) DeleteWallet(c echo.Context) error {
	err := h.ws.DeleteWallet(c.Request().Context(), c.Param("id"))
	if err != nil && err == ErrWalletNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// 201 => successfully created
// 404 => wallet not found
// 400 => type not valid
// 422 => transaction amount does not meet config requirements
// 500 => any other error
func (h *handler) CreateTransaction(c echo.Context) error {
	var info TransactionCreationInfo
	if err := c.Bind(&info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txnID, err := h.ws.CreateTransaction(c.Request().Context(), &info)
	if err != nil && isNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	} else if err != nil && isBadRequest(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && isUnprocessableEntity(err) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, PostResponse{txnID})
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

func isBadRequest(err error) bool {
	return ContainsError(err, badRequestErrors)
}

func isNotFound(err error) bool {
	return ContainsError(err, notFoundErrors)
}

func isUnprocessableEntity(err error) bool {
	return ContainsError(err, unprocessableEntityErrors)
}

func ContainsError(err error, errList []error) bool {
	for _, e := range errList {
		if err == e {
			return true
		}
	}

	return false
}
