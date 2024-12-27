package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

const (
	UserIDContext = "userID"
)

func WithAuth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")

			// Проверяем, что заголовок пристствует в запросе и имеет правильный формат
			if !isValidAuthHeader(authHeader) {
				handleUnauthorized(rw, "invalid or missing Authorization header")
				return
			}

			// Извлекаем JWT токен
			authHeaderSplit := strings.Split(authHeader, " ")
			jwtToken := authHeaderSplit[1]

			// Получаем userID из токена
			userID, err := security.GetUserID(jwtToken, secretKey)
			if err != nil {
				handleUnauthorized(rw, "invalid JWT token")
				return
			}

			// Добавляем userID в контекст
			ctx := context.WithValue(req.Context(), UserIDContext, userID)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}

// Проверяет корректность заголовка Authorization
func isValidAuthHeader(authHeader string) bool {
	if authHeader == "" {
		return false
	}
	authHeaderSplit := strings.Split(authHeader, " ")
	return len(authHeaderSplit) == 2 && authHeaderSplit[0] == "Bearer"
}

// Обрабатывает ошибку Unauthorized (401)
func handleUnauthorized(rw http.ResponseWriter, message string) {
	log.Error().Msg(message)
	http.Error(rw, appErrors.ErrUserUnauthorized.Error(), http.StatusUnauthorized)
}
