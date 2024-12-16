package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/handler/mocks"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func TestHandler_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                      string
		mockService               func() *mocks.MockApp
		conf                      config.Config
		reqBody                   *bytes.Buffer
		expectedStatusCode        int
		expectedHeader            string
		expectedHeaderValContains string
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				user := &models.User{
					Login:    "login",
					Password: "password",
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().RegisterUser(user).Return(int64(1), nil)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusOK,
			expectedHeader:            "Authorization",
			expectedHeaderValContains: "Bearer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.mockService(), &tt.conf)

			req := httptest.NewRequest("POST", "/api/user/register", tt.reqBody)
			rw := httptest.NewRecorder()

			handler.RegisterUser(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
			assert.Contains(t, rw.Header().Get(tt.expectedHeader), tt.expectedHeaderValContains)
		})
	}
}

func TestHandler_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                      string
		mockService               func() *mocks.MockApp
		conf                      config.Config
		reqBody                   *bytes.Buffer
		expectedStatusCode        int
		expectedHeader            string
		expectedHeaderValContains string
	}{
		{
			name: "Success Case",
			mockService: func() *mocks.MockApp {
				user := &models.User{
					Login:    "login",
					Password: "password",
				}
				dbUser := &models.User{
					ID:       1,
					Login:    "login",
					Password: "password",
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().ValidateUser(user).Return(dbUser, nil)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusOK,
			expectedHeader:            "Authorization",
			expectedHeaderValContains: "Bearer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.mockService(), &tt.conf)

			req := httptest.NewRequest("POST", "/api/user/login", tt.reqBody)
			rw := httptest.NewRecorder()

			handler.AuthUser(rw, req)

			assert.Equal(t, tt.expectedStatusCode, rw.Code)
			assert.Contains(t, rw.Header().Get(tt.expectedHeader), tt.expectedHeaderValContains)
		})
	}
}
