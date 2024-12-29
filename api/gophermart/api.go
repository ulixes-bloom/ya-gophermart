package api

import (
	"context"
	"net/http"
	"sync"
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
	var wg sync.WaitGroup

	// Запуск HTTP сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- http.ListenAndServe(a.conf.RunAddr, a.router)
	}()

	// Обновление информации по необработанным заказам
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.updateNotProcessedOrders(ctx)
	}()

	// Ожидание завершения работы
	for {
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			wg.Wait()
			return a.app.Shutdown()
		}
	}
}

func (a *API) updateNotProcessedOrders(ctx context.Context) {
	updateNotProcessedOrdersTicker := time.NewTicker(a.conf.OrderInfoUpdateInterval)
	defer updateNotProcessedOrdersTicker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		select {
		case <-updateNotProcessedOrdersTicker.C:
			err := a.app.UpdateNotProcessedOrders(ctx)
			if err != nil {
				log.Error().Err(err).Msg(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}
