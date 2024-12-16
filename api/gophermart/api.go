package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-gophermart/internal/app"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
)

type API struct {
	router *chi.Mux
	app    *app.App
	conf   *config.Config
}

func New(conf *config.Config, app *app.App) *API {
	return &API{
		router: NewRouter(app, conf),
		app:    app,
		conf:   conf,
	}
}

func (a *API) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- http.ListenAndServe(a.conf.RunAddr, a.router)
	}()

	go a.updateNotProcessedOrders(ctx)

	for {
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return a.app.Shutdown()
		}
	}
}

func (a *API) updateNotProcessedOrders(ctx context.Context) {
	updateNotProcessedOrdersTicker := time.NewTicker(a.conf.OrderInfoUpdateInterval)
	defer updateNotProcessedOrdersTicker.Stop()

	for {
		select {
		case <-updateNotProcessedOrdersTicker.C:
			err := a.app.UpdateNotProcessedOrders()
			if err != nil {
				log.Error().Err(err).Msg(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}
