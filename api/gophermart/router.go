package api

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/handler"
	"github.com/ulixes-bloom/ya-gophermart/api/gophermart/middleware"
	_ "github.com/ulixes-bloom/ya-gophermart/docs"
	"github.com/ulixes-bloom/ya-gophermart/internal/app"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
)

func NewRouter(app *app.App, conf *config.Config) *chi.Mux {
	r := chi.NewRouter()
	h := handler.New(app, conf)

	r.Use(middleware.WithLogging)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.AuthUser)

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAuth(conf.TokenSecretKey))
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", h.RegisterUserOrder)
				r.Get("/", h.GetUserOrders)
			})

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", h.GetUserBalance)
				r.Post("/withdraw", h.WithdrawFromUserBalance)
			})

			r.Get("/withdrawals", h.GetUserWithdrawals)
		})
	})

	return r
}
