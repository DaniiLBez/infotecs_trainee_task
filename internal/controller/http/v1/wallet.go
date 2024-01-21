package v1

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo/repoerrors"
	"infotecs_trainee_task/internal/service"
	"net/http"
)

type walletRoutes struct {
	walletService service.Wallet
}

func newWalletRoutes(g *echo.Group, walletService service.Wallet) {
	r := &walletRoutes{walletService: walletService}

	g.POST("/", r.CreateWallet)
	g.POST("/:walletId/send", r.TransferCash)
	g.GET("/:walletId/history", r.GetHistory)
	g.GET("/:walletId", r.GetState)
}

func (r *walletRoutes) CreateWallet(c echo.Context) error {

	wallet, err := r.walletService.CreateWallet(c.Request().Context())

	if err != nil {
		if errors.Is(err, service.ErrWalletAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Uuid uuid.UUID `json:"uuid"`
	}

	return c.JSON(http.StatusOK, response{Uuid: wallet})
}

type transferCashRequest struct {
	To     uuid.UUID `json:"to" validator:"required"`
	Amount float64   `json:"amount" validator:"required,min=0"`
}

func (r *walletRoutes) TransferCash(c echo.Context) error {
	walletId, err := uuid.Parse(c.Param("walletId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request params "+err.Error())
		return err
	}

	var transferInput transferCashRequest
	if err = c.Bind(&transferInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err = c.Validate(transferInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err = r.walletService.MakeTransaction(
		c.Request().Context(),
		walletId,
		transferInput.To,
		transferInput.Amount,
	)

	if err != nil {
		if errors.Is(err, service.ErrTransactionAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}

		if errors.Is(err, service.ErrCannotGetWallet) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error "+err.Error())
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (r *walletRoutes) GetHistory(c echo.Context) error {
	walletId, err := uuid.Parse(c.Param("walletId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request params "+err.Error())
		return err
	}

	transactions, err := r.walletService.GetTransactionsHistory(c.Request().Context(), walletId)
	if err != nil {
		if errors.Is(err, service.ErrCannotGetTransaction) {
			newErrorResponse(c, http.StatusBadRequest, "can not get transaction "+err.Error())
			return err
		}

		if errors.Is(err, service.ErrCannotGetWallet) {
			newErrorResponse(c, http.StatusNotFound, "wallet not found")
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error "+err.Error())
	}

	type response struct {
		Transactions []entity.Transaction `json:"array"`
	}

	return c.JSON(http.StatusOK, response{Transactions: transactions})
}

func (r *walletRoutes) GetState(c echo.Context) error {
	walletId, err := uuid.Parse(c.Param("walletId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request params "+err.Error())
		return err
	}

	walletState, err := r.walletService.GetWalletState(c.Request().Context(), walletId)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, "wallet not found")
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error "+err.Error())
	}

	type response struct {
		State entity.Wallet `json:"state"`
	}

	return c.JSON(http.StatusOK, response{State: walletState})
}
