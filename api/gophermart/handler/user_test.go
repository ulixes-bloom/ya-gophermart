package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/handler/mocks"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
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
		{
			name: "Login alrady exists Case",
			mockService: func() *mocks.MockApp {
				user := &models.User{
					Login:    "login",
					Password: "password",
				}
				err := appErrors.ErrUserLoginAlreadyExists

				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().RegisterUser(user).Return(int64(-1), err)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusConflict,
			expectedHeader:            "",
			expectedHeaderValContains: "",
		},
		{
			name: "Login is empty Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusBadRequest,
			expectedHeader:            "",
			expectedHeaderValContains: "",
		},
		{
			name: "Password is empty Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"\"}\n")),
			expectedStatusCode:        http.StatusBadRequest,
			expectedHeader:            "",
			expectedHeaderValContains: "",
		},
		{
			name: "Data base error Case",
			mockService: func() *mocks.MockApp {
				user := &models.User{
					Login:    "login",
					Password: "password",
				}
				err := errors.New("Table users does not exist")

				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().RegisterUser(user).Return(int64(-1), err)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusInternalServerError,
			expectedHeader:            "",
			expectedHeaderValContains: "",
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
		{
			name: "Invalid login or password Case",
			mockService: func() *mocks.MockApp {
				user := &models.User{
					Login:    "login",
					Password: "password",
				}
				mockService := mocks.NewMockApp(ctrl)
				mockService.EXPECT().ValidateUser(user).Return(nil, appErrors.ErrInvalidUserLoginOrPassword)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"login\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusUnauthorized,
			expectedHeader:            "",
			expectedHeaderValContains: "",
		},
		{
			name: "Empty login Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login\":\"\",\"password\":\"password\"}\n")),
			expectedStatusCode:        http.StatusBadRequest,
			expectedHeader:            "",
			expectedHeaderValContains: "",
		},
		{
			name: "Bad request body Case",
			mockService: func() *mocks.MockApp {
				mockService := mocks.NewMockApp(ctrl)
				return mockService
			},
			conf:                      *config.GetDefault(),
			reqBody:                   bytes.NewBuffer([]byte("{\"login:\"login\",\"password\":\"\"}\n")),
			expectedStatusCode:        http.StatusBadRequest,
			expectedHeader:            "",
			expectedHeaderValContains: "",
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
