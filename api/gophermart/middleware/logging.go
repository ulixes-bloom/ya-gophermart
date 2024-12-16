package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		uri := req.RequestURI
		method := req.Method

		next.ServeHTTP(rw, req)

		duration := time.Since(start)

		log.Debug().
			Str("uri", uri).
			Str("method", method).
			Str("duration", duration.String()).
			Msg("got incomming HTTP request")
	})
}
