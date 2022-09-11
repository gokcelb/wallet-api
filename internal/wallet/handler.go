package wallet

import "github.com/labstack/echo/v4"

// should these functions return the value or the reference of Wallet?
type Service interface {
	Create(info *WalletCreationInfo) (Wallet, error)
	Get(id string) (Wallet, error)
	Delete(id string) error
}

type Handler struct {
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

func New(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/wallets", h.CreateWallet)
	e.GET("/wallets/:id", h.GetWallet)
	e.DELETE("/wallets/:id", h.DeleteWallet)

	e.POST("/wallets/:id/transactions", h.CreateTransaction)
	e.GET("/wallets/:id/transactions", h.GetTransactions)
	e.GET("/wallets/:id/transactions/:transactionId", h.GetTransaction)
}

func (h *Handler) CreateWallet(c echo.Context) error {
	return nil
}

func (h *Handler) GetWallet(c echo.Context) error {
	return nil
}

func (h *Handler) DeleteWallet(c echo.Context) error {
	return nil
}

func (h *Handler) CreateTransaction(c echo.Context) error {
	return nil
}

func (h *Handler) GetTransactions(c echo.Context) error {
	return nil
}

func (h *Handler) GetTransaction(c echo.Context) error {
	return nil
}
