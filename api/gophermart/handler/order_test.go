package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/handler/mocks"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/middleware"
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
				mockService.EXPECT().RegisterOrder("1234", int64(1)).Return(nil)
				mockService.EXPECT().ValidateOrderNumber("1234").Return(true)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			reqBody:            bytes.NewBuffer([]byte("{\"number\":\"1234\"}\n")),
			expectedStatusCode: http.StatusAccepted,
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
				mockService.EXPECT().GetOrdersByUser(int64(1)).Return(orders, nil).AnyTimes()
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedBody:       "[{\"number\":\"2377225624\",\"status\":\"NEW\",\"accrual\":200,\"uploaded_at\":\"2024-01-02T00:00:00Z\"}]\n",
			expectedStatusCode: http.StatusOK,
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
