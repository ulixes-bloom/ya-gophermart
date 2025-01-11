package accrual

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-gophermart/internal/accrual/mocks"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func TestAccrualClient_GetOrderInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		mockClient    func() *mocks.MockHTTPClient
		conf          *config.Config
		expectedOrder *models.Order
		expectedErr   error
	}{
		{
			name: "Success Case",
			mockClient: func() *mocks.MockHTTPClient {
				mockService := mocks.NewMockHTTPClient(ctrl)
				mockBody := `{"order":"2377225624","status":"REGISTERED","accrual":200}`
				mockResponse := &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(mockBody)),
				}
				mockService.EXPECT().Do(gomock.Any()).Return(mockResponse, nil)
				return mockService
			},
			conf: config.GetDefault(),
			expectedOrder: &models.Order{
				Number:  "2377225624",
				Status:  models.OrderStatusNew,
				Accrual: 200,
			},
			expectedErr: nil,
		},
		{
			name: "Order not regisered Case",
			mockClient: func() *mocks.MockHTTPClient {
				mockService := mocks.NewMockHTTPClient(ctrl)
				mockResponse := &http.Response{
					StatusCode: 204,
					Body:       io.NopCloser(nil),
				}
				mockService.EXPECT().Do(gomock.Any()).Return(mockResponse, nil)
				return mockService
			},
			conf:          config.GetDefault(),
			expectedOrder: nil,
			expectedErr:   appErrors.ErrAccrualOrderNotRegistered,
		},
		{
			name: "Too many requests Case",
			mockClient: func() *mocks.MockHTTPClient {
				mockService := mocks.NewMockHTTPClient(ctrl)
				mockResponse := &http.Response{
					StatusCode: 429,
					Body:       io.NopCloser(nil),
				}
				mockService.EXPECT().Do(gomock.Any()).Return(mockResponse, nil)
				return mockService
			},
			conf:          config.GetDefault(),
			expectedOrder: nil,
			expectedErr:   appErrors.ErrAccrualTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := Client{conf: tt.conf, http: tt.mockClient()}
			ctx := context.Background()

			order, err := ac.GetOrderInfo(ctx, &models.Order{})

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, order)
			} else {
				assert.NotNil(t, order)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedOrder.Number, order.Number)
				assert.Equal(t, tt.expectedOrder.Accrual, order.Accrual)
				assert.Equal(t, tt.expectedOrder.Status, order.Status)
			}
		})
	}
}
