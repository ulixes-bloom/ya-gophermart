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

func TestHandler_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		mockService        func() *mocks.MockApp
		reqBody            *bytes.Buffer
		ctx                context.Context
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				mockBalance := &models.Balance{
					Current:   0,
					Withdrawn: 0,
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetUserBalance(int64(1)).Return(mockBalance, nil)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "{\"withdrawn\":0,\"current\":0}\n",
		},
		{
			name: "Balance does not exist Case",
			mockService: func() *mocks.MockApp {
				err := errors.New("User balance does not exist")

				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetUserBalance(int64(1)).Return(nil, err)
				return mockService
			},
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "failed to get user balance\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPHandler{app: tt.mockService()}

			req := httptest.NewRequest("GET", "/api/user/balance", nil)
			req = req.WithContext(tt.ctx)
			rw := httptest.NewRecorder()

			handler.GetUserBalance(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
			assert.Equal(t, tt.expectedBody, rw.Body.String())
		})
	}
}

func TestHandler_WithdrawFromUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		mockService        func() *mocks.MockApp
		reqBody            *bytes.Buffer
		ctx                context.Context
		expectedStatusCode int
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				withdrawalReq := &models.WithdrawalRequest{
					Order: "2377225624",
					Sum:   200,
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().WithdrawFromUserBalance(int64(1), withdrawalReq).Return(nil)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			reqBody:            bytes.NewBuffer([]byte("{\"order\":\"2377225624\",\"sum\":200.0}\n")),
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Negative balance Case",
			mockService: func() *mocks.MockApp {
				withdrawalReq := &models.WithdrawalRequest{
					Order: "2377225624",
					Sum:   200,
				}
				err := appErrors.ErrNegativeBalance

				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().WithdrawFromUserBalance(int64(1), withdrawalReq).Return(err)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			reqBody:            bytes.NewBuffer([]byte("{\"order\":\"2377225624\",\"sum\":200.0}\n")),
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name: "Not valid order number Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(false)
				return mockService
			},
			reqBody:            bytes.NewBuffer([]byte("{\"order\":\"2377225624\",\"sum\":200.0}\n")),
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Data base error Case",
			mockService: func() *mocks.MockApp {
				withdrawalReq := &models.WithdrawalRequest{
					Order: "2377225624",
					Sum:   200,
				}
				err := errors.New("Table witdrawals does not exist")

				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().WithdrawFromUserBalance(int64(1), withdrawalReq).Return(err)
				mockService.EXPECT().ValidateOrderNumber("2377225624").Return(true)
				return mockService
			},
			reqBody:            bytes.NewBuffer([]byte("{\"order\":\"2377225624\",\"sum\":200.0}\n")),
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPHandler{app: tt.mockService()}

			req := httptest.NewRequest("POST", "/api/user/balance/withdraw", tt.reqBody)
			req = req.WithContext(tt.ctx)
			rw := httptest.NewRecorder()

			handler.WithdrawFromUserBalance(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
		})
	}
}

func TestHandler_GetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		mockService        func() *mocks.MockApp
		expectedBody       string
		ctx                context.Context
		expectedStatusCode int
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				withdrawals := []models.Withdrawal{
					{
						Order:       "2377225624",
						Sum:         200,
						ProcessedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetUserWithdrawals(int64(1)).Return(withdrawals, nil)
				return mockService
			},
			expectedBody:       "[{\"order\":\"2377225624\",\"processed_at\":\"2024-01-02T00:00:00Z\",\"sum\":200}]\n",
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "No wihdrawals Case",
			mockService: func() *mocks.MockApp {
				withdrawals := []models.Withdrawal{}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetUserWithdrawals(int64(1)).Return(withdrawals, nil)
				return mockService
			},
			expectedBody:       "",
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Data base error Case",
			mockService: func() *mocks.MockApp {
				err := errors.New("Table witdrawals does not exist")
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().GetUserWithdrawals(int64(1)).Return(nil, err)
				return mockService
			},
			expectedBody:       "failed to get user waithdrawals\n",
			ctx:                context.WithValue(context.Background(), middleware.UserIDContext, int64(1)),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPHandler{app: tt.mockService()}

			req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			req = req.WithContext(tt.ctx)
			rw := httptest.NewRecorder()

			handler.GetUserWithdrawals(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
			assert.Equal(t, tt.expectedBody, rw.Body.String())
		})
	}
}
