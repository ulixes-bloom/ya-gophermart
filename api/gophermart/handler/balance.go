package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/middleware"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

// @Summary	Получение текущего баланса пользователя
// @ID			GetUserBalance
// @Produce	json
// @Success	200	"успешная обработка запроса"
// @Failure	401	"пользователь не авторизован"
// @Failure	500	"внутренняя ошибка сервера"
// @Router		/api/user/balance [get]
// @Param		Authorization	header	string	false	"Bearer"
func (h *HTTPHandler) GetUserBalance(rw http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(middleware.UserIDContext).(int64)
	if !ok {
		h.handleError(rw,
			appErrors.ErrUserInalidID,
			appErrors.ErrUserInalidID.Error(),
			http.StatusInternalServerError)
		return
	}

	dbBalance, err := h.app.GetUserBalance(userID)
	if err != nil {
		h.handleError(rw, err, "failed to get user balance", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(dbBalance); err != nil {
		h.handleError(rw, err, "failed to encode balance", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// @Summary	Запрос на списание средств
// @ID			WithdrawFromUserBalance
// @Produce	json
// @Success	200	"успешная обработка запроса"
// @Failure	401	"пользователь не авторизован"
// @Failure	402	"на счету недостаточно средств"
// @Failure	422	"неверный номер заказа"
// @Failure	500	"внутренняя ошибка сервера"
// @Router		/api/user/balance/withdraw [post]
// @Param		Authorization	header	string						false	"Bearer"
// @Param		user			body	models.WithdrawalRequest	true	"User Registration Information"
func (h *HTTPHandler) WithdrawFromUserBalance(rw http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(middleware.UserIDContext).(int64)
	if !ok {
		h.handleError(rw,
			appErrors.ErrUserInalidID,
			appErrors.ErrUserInalidID.Error(),
			http.StatusInternalServerError)
		return
	}

	if req.Body == nil {
		h.handleError(rw, nil, "request body is missing", http.StatusBadRequest)
		return
	}

	withdrawalReq := &models.WithdrawalRequest{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(withdrawalReq); err != nil {
		h.handleError(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	valid := h.app.ValidateOrderNumber(withdrawalReq.Order)
	if !valid {
		h.handleError(rw,
			appErrors.ErrInvalidOrderNumber,
			appErrors.ErrInvalidOrderNumber.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	err := h.app.WithdrawFromUserBalance(userID, withdrawalReq)
	if err != nil {
		if errors.Is(err, appErrors.ErrNegativeBalance) {
			h.handleError(rw, err, appErrors.ErrNegativeBalance.Error(), http.StatusPaymentRequired)
		} else {
			h.handleError(rw, err, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// @Summary	Получение информации о выводе средств
// @ID			GetUserWithdrawals
// @Produce	json
// @Success	200	"успешная обработка запроса"
// @Success	204	"нет ни одного списания"
// @Failure	401	"пользователь не авторизован"
// @Failure	500	"внутренняя ошибка сервера"
// @Router		/api/user/withdrawals [get]
// @Param		Authorization	header	string	false	"Bearer"
func (h *HTTPHandler) GetUserWithdrawals(rw http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(middleware.UserIDContext).(int64)
	if !ok {
		h.handleError(rw,
			appErrors.ErrUserInalidID,
			appErrors.ErrUserInalidID.Error(),
			http.StatusInternalServerError)
		return
	}

	dbWithdrawals, err := h.app.GetUserWithdrawals(userID)
	if err != nil {
		h.handleError(rw, err, "failed to get user waithdrawals", http.StatusInternalServerError)
		return
	}

	if len(dbWithdrawals) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(rw).Encode(dbWithdrawals); err != nil {
		h.handleError(rw, err, "failed to encode waithdrawals", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
