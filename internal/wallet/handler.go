package wallet

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/labstack/echo/v4"
)

var badRequestErrors = []error{
	ErrInvalidTransactionType,
}

var notFoundErrors = []error{
	ErrWalletNotFound,
	transaction.ErrTransactionNotFound,
}

var unprocessableEntityErrors = []error{
	ErrWalletWithUserIDExists,
	ErrAboveMaximumBalanceLimit,
	ErrAboveMaximumTransactionLimit,
	ErrBelowMinimumTransactionLimit,
	ErrInsufficientBalance,
}

var (
	DefaultPageNo   = 0
	DefaultPageSize = 10

	ErrInvalidPageNo   = errors.New("pageNo cannot be converted to integer")
	ErrInvalidPageSize = errors.New("pageSize cannot be converted to integer")
)

type WalletService interface {
	CreateWallet(ctx context.Context, info *WalletCreationInfo) (string, error)
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	DeleteWallet(ctx context.Context, id string) error
	CreateTransaction(ctx context.Context, info *TransactionCreationInfo) (string, error)
	GetTransactions(ctx context.Context, walletID, typeFilter string, pageNo, pageSize int) ([]*transaction.Transaction, error)
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
	WalletID        string  `param:"id"`
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

	e.POST("/wallets/:id/transactions", h.CreateTransaction)
	e.GET("/wallets/:id/transactions", h.GetTransactions)
}

func (h *handler) CreateWallet(c echo.Context) error {
	var info WalletCreationInfo
	if err := c.Bind(&info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.ws.CreateWallet(c.Request().Context(), &info)
	if err != nil && isBadRequest(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && isUnprocessableEntity(err) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, PostResponse{id})
}

func (h *handler) GetWallet(c echo.Context) error {
	w, err := h.ws.GetWallet(c.Request().Context(), c.Param("id"))
	if err != nil && isNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, w)
}

func (h *handler) DeleteWallet(c echo.Context) error {
	err := h.ws.DeleteWallet(c.Request().Context(), c.Param("id"))
	if err != nil && isNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

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

func (h *handler) GetTransactions(c echo.Context) error {
	pageNo, pageSize, err := h.getPaginationParamsOrDefault(c.QueryParam("pageNo"), c.QueryParam("pageSize"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txns, err := h.ws.GetTransactions(c.Request().Context(), c.Param("id"), c.QueryParam("type"), pageNo, pageSize)
	if err != nil && isBadRequest(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && isNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, txns)
}

func (h *handler) getPaginationParamsOrDefault(pageNoQuery string, pageSizeQuery string) (int, int, error) {
	if pageNoQuery == "" || pageSizeQuery == "" {
		return DefaultPageNo, DefaultPageSize, nil
	}

	pageNo, err := strconv.Atoi(pageNoQuery)
	if err != nil {
		return 0, 0, ErrInvalidPageNo
	}

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil {
		return 0, 0, ErrInvalidPageSize
	}

	return pageNo, pageSize, nil
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
		if errors.Is(e, err) {
			return true
		}
	}

	return false
}
