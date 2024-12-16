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
			authHeaderSplit := strings.Split(authHeader, " ")
			if len(authHeaderSplit) != 2 || authHeaderSplit[0] != "Bearer" {
				log.Error().Msg(appErrors.ErrUserUnauthorized.Error())
				http.Error(rw, appErrors.ErrUserUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			jwtToken := authHeaderSplit[1]
			userID, err := security.GetUserID(jwtToken, secretKey)
			if err != nil {
				log.Error().Msg(appErrors.ErrUserUnauthorized.Error())
				http.Error(rw, appErrors.ErrUserUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			ctx := req.Context()
			ctx = context.WithValue(ctx, UserIDContext, userID)

			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
