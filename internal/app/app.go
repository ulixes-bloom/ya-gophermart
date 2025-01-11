package app

import (
	"github.com/ulixes-bloom/ya-gophermart/internal/accrual"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
)

type App struct {
	storage Storage
	ac      *accrual.Client
	conf    *config.Config
}

func New(storage Storage, conf *config.Config) *App {
	return &App{
		storage: storage,
		conf:    conf,
		ac:      accrual.NewClient(conf),
	}
}

func (a *App) Shutdown() error {
	err := a.storage.Close()
	if err != nil {
		return err
	}
	return nil
}
