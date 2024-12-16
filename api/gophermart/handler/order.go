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

//	@Summary	Загрузка номера заказа
//	@ID			RegisterUserOrder
//	@Produce	json
//	@Success	200	"номер заказа уже был загружен этим пользователем"
//	@Success	202	"новый номер заказа принят в обработку"
//	@Failure	400	"неверный формат запроса"
//	@Failure	401	"пользователь не аутентифицирован"
//	@Failure	409	"номер заказа уже был загружен другим пользователем"
//	@Failure	422	"неверный формат номера заказа"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Router		/api/user/orders [post]
//	@Param		Authorization	header	string				false	"Bearer"
//	@Param		user			body	models.OrderRequest	true	"User Registration Information"
func (h *HTTPHandler) RegisterUserOrder(rw http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(middleware.UserIDContext).(int64)

	orderReq := &models.OrderRequest{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(orderReq); err != nil {
		h.Error(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	valid := h.app.ValidateOrderNumber(orderReq.Number)
	if !valid {
		h.Error(rw,
			appErrors.ErrInvalidOrderNumber,
			appErrors.ErrInvalidOrderNumber.Error(),
			http.StatusUnprocessableEntity)
	}

	err := h.app.RegisterOrder(orderReq.Number, userID)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrOrderWasUploadedByCurrentUser):
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(appErrors.ErrOrderWasUploadedByCurrentUser.Error()))
		case errors.Is(err, appErrors.ErrOrderWasUploadedByAnotherUser):
			h.Error(rw,
				appErrors.ErrOrderWasUploadedByCurrentUser,
				appErrors.ErrOrderWasUploadedByCurrentUser.Error(),
				http.StatusConflict)
		default:
			h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

//	@Summary	Получение списка загруженных номеров заказов
//	@ID			GetUserOrders
//	@Produce	json
//	@Success	200	"успешная обработка запроса"
//	@Success	204	"нет данных для ответа"
//	@Failure	401	"пользователь не авторизован"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Router		/api/user/orders [get]
//	@Param		Authorization	header	string	false	"Bearer"
func (h *HTTPHandler) GetUserOrders(rw http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(middleware.UserIDContext).(int64)
	dbOrders, err := h.app.GetOrdersByUser(userID)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(dbOrders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	var bufOrders bytes.Buffer
	jsonEncoder := json.NewEncoder(&bufOrders)
	jsonEncoder.Encode(dbOrders)

	rw.WriteHeader(http.StatusOK)
	rw.Write(bufOrders.Bytes())
}
