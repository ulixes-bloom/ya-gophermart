package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/handler/mocks"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/middleware"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func TestHandler_RegisterUserOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		mockService        func() *mocks.MockApp
		ctx                context.Context
		reqBody            *bytes.Buffer
		expectedStatusCode int
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().RegisterOrder(gomock.Any(), int64(1), "2377225624").Return(nil)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"2377225624\"}\n")),
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name: "Order already restired by current user Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				err := appErrors.ErrOrderWasUploadedByCurrentUser

				mockService.EXPECT().RegisterOrder(gomock.Any(), int64(1), "2377225624").Return(err)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"2377225624\"}\n")),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Order already restired by another user Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				err := appErrors.ErrOrderWasUploadedByAnotherUser

				mockService.EXPECT().RegisterOrder(gomock.Any(), int64(1), "2377225624").Return(err)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"2377225624\"}\n")),
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "Not valid order number Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(false)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"2377225624\"}\n")),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Bad request body Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number:\"2377225624\"}\n")),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Data base error Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				err := errors.New("Table orders does not exist")

				mockService.EXPECT().RegisterOrder(gomock.Any(), int64(1), "2377225624").Return(err)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"2377225624\"}\n")),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPHandler{app: tt.mockService()}

			req := httptest.NewRequest("POST", "/api/user/orders", tt.reqBody)
			req = req.WithContext(tt.ctx)
			rw := httptest.NewRecorder()

			handler.RegisterUserOrder(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
		})
	}
}

func TestHandler_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		mockService        func() *mocks.MockApp
		ctx                context.Context
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				orders := []models.Order{
					{
						Number:     "2377225624",
						Status:     models.OrderStatusNew,
						Accrual:    200,
						UploadedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetOrdersByUser(gomock.Any(), int64(1)).Return(orders, nil)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedBody:       "[{\"number\":\"2377225624\",\"status\":\"NEW\",\"accrual\":200,\"uploaded_at\":\"2024-01-02T00:00:00Z\"}]\n",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "No orders Case",
			mockService: func() *mocks.MockApp {
				orders := []models.Order{}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetOrdersByUser(gomock.Any(), int64(1)).Return(orders, nil)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedBody:       "",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Data base error Case",
			mockService: func() *mocks.MockApp {
				err := errors.New("Table orders does not exist")
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetOrdersByUser(gomock.Any(), int64(1)).Return(nil, err)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedBody:       "failed to get user orders\n",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPHandler{app: tt.mockService()}

			req := httptest.NewRequest("GET", "/api/user/orders", nil)
			req = req.WithContext(tt.ctx)
			rw := httptest.NewRecorder()

			handler.GetUserOrders(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
			assert.Equal(t, tt.expectedBody, rw.Body.String())
		})
	}
}
