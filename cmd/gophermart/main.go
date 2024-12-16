package main

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	api "github.com/ulixes-bloom/ya-gophermart/api/gophermart"
	"github.com/ulixes-bloom/ya-gophermart/internal/app"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
	"github.com/ulixes-bloom/ya-gophermart/internal/storage/pg"
)

func main() {
	ctx := context.Background()

	conf, err := config.Parse()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	db, err := sql.Open("pgx", conf.DatabaseURI)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	storage, err := pg.NewStorage(ctx, db)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	app := app.New(storage, conf)

	api := api.New(conf, app)
	err = api.Run(ctx)
	if err != nil {
		log.Panic().Err(err)
	}
}
