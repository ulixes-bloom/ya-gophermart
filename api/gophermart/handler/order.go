package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/middleware"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

// @Summary	Загрузка номера заказа
// @ID			RegisterUserOrder
// @Produce	json
// @Success	200	"номер заказа уже был загружен этим пользователем"
// @Success	202	"новый номер заказа принят в обработку"
// @Failure	400	"неверный формат запроса"
// @Failure	401	"пользователь не аутентифицирован"
// @Failure	409	"номер заказа уже был загружен другим пользователем"
// @Failure	422	"неверный формат номера заказа"
// @Failure	500	"внутренняя ошибка сервера"
// @Router		/api/user/orders [post]
// @Param		Authorization	header	string				false	"Bearer"
// @Param		order			body	models.OrderRequest	true	"Новый заказ"
func (h *HTTPHandler) RegisterUserOrder(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

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

	orderReq := &models.OrderRequest{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(orderReq); err != nil {
		h.handleError(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	valid := h.app.ValidateOrderNumber(orderReq.Number)
	if !valid {
		h.handleError(rw,
			appErrors.ErrInvalidOrderNumber,
			appErrors.ErrInvalidOrderNumber.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	err := h.app.RegisterOrder(ctx, userID, orderReq.Number)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrOrderWasUploadedByCurrentUser):
			rw.WriteHeader(http.StatusOK)
		case errors.Is(err, appErrors.ErrOrderWasUploadedByAnotherUser):
			h.handleError(rw,
				appErrors.ErrOrderWasUploadedByAnotherUser,
				appErrors.ErrOrderWasUploadedByAnotherUser.Error(),
				http.StatusConflict)
		default:
			h.handleError(rw, err, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// @Summary	Получение списка загруженных номеров заказов
// @ID			GetUserOrders
// @Produce	json
// @Success	200	"успешная обработка запроса"
// @Success	204	"нет данных для ответа"
// @Failure	401	"пользователь не авторизован"
// @Failure	500	"внутренняя ошибка сервера"
// @Router		/api/user/orders [get]
// @Param		Authorization	header	string	false	"Bearer"
func (h *HTTPHandler) GetUserOrders(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, ok := req.Context().Value(middleware.UserIDContext).(int64)
	if !ok {
		h.handleError(rw,
			appErrors.ErrUserInalidID,
			appErrors.ErrUserInalidID.Error(),
			http.StatusInternalServerError)
		return
	}

	dbOrders, err := h.app.GetOrdersByUser(ctx, userID)
	if err != nil {
		h.handleError(rw, err, "failed to get user orders", http.StatusInternalServerError)
		return
	}

	if len(dbOrders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(rw).Encode(dbOrders); err != nil {
		h.handleError(rw, err, "failed to encode orders", http.StatusInternalServerError)
		return
	}
}
