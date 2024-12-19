package handler

import (
	"bytes"
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
	userID := req.Context().Value(middleware.UserIDContext).(int64)
	dbBalance, err := h.app.GetUserBalance(userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var bufBalance bytes.Buffer
	jsonEncoder := json.NewEncoder(&bufBalance)
	err = jsonEncoder.Encode(dbBalance)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bufBalance.Bytes())
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
	userID := req.Context().Value(middleware.UserIDContext).(int64)

	withdrawalReq := &models.WithdrawalRequest{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(withdrawalReq); err != nil {
		h.Error(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	valid := h.app.ValidateOrderNumber(withdrawalReq.Order)
	if !valid {
		h.Error(rw,
			appErrors.ErrInvalidOrderNumber,
			appErrors.ErrInvalidOrderNumber.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	err := h.app.WithdrawFromUserBalance(withdrawalReq, userID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNegativeBalance) {
			h.Error(rw, err, appErrors.ErrNegativeBalance.Error(), http.StatusPaymentRequired)
		} else {
			h.Error(rw, err, err.Error(), http.StatusInternalServerError)
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
	userID := req.Context().Value(middleware.UserIDContext).(int64)
	dbWithdrawals, err := h.app.GetUserWithdrawals(userID)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(dbWithdrawals) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	var bufWithdrawals bytes.Buffer
	jsonEncoder := json.NewEncoder(&bufWithdrawals)
	err = jsonEncoder.Encode(dbWithdrawals)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bufWithdrawals.Bytes())
}
