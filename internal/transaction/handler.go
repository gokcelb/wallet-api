package transaction

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransactionService interface {
	GetTransaction(ctx context.Context, id string) (*Transaction, error)
}

type handler struct {
	ts TransactionService
}

func NewHandler(ts TransactionService) *handler {
	return &handler{ts}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/transactions/:id", h.GetTransaction)
}

func (h *handler) GetTransaction(c echo.Context) error {
	txn, err := h.ts.GetTransaction(c.Request().Context(), c.Param("id"))
	if err != nil && errors.Is(err, ErrTransactionNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, txn)
}
