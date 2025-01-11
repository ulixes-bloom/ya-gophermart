package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
)

type HTTPHandler struct {
	app  App
	conf *config.Config
}

func New(app App, conf *config.Config) *HTTPHandler {
	return &HTTPHandler{
		app:  app,
		conf: conf,
	}
}

func (h *HTTPHandler) handleError(rw http.ResponseWriter, err error, errMsg string, statusCode int) {
	log.Error().Err(err).Msg(errMsg)
	http.Error(rw, errMsg, statusCode)
}
