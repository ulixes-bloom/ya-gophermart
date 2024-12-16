package pg

import (
	"context"
	"database/sql"
	"fmt"
)

type pgstorage struct {
	db *sql.DB
}

func NewStorage(ctx context.Context, db *sql.DB) (*pgstorage, error) {
	newPg := &pgstorage{db: db}

	if err := newPg.createTables(ctx); err != nil {
		return nil, fmt.Errorf("pg.newStorage: %w", err)
	}

	return newPg, nil
}

func (pg *pgstorage) Close() error {
	if err := pg.db.Close(); err != nil {
		return fmt.Errorf("pg.close: %w", err)
	}

	return nil
}

func (ps *pgstorage) createTables(ctx context.Context) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("pg.createTables.beginTx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DO $$ BEGIN
			CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)
	if err != nil {
		return fmt.Errorf("pg.createTables.orderStatusType: %w", err)
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users
		(
			id       bigint  PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			login    varchar NOT NULL UNIQUE,
			password varchar NOT NULL
		);`)
	if err != nil {
		return fmt.Errorf("pg.createTables.usersTable: %w", err)
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS orders
		(
			id          bigint  PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			number      varchar NOT NULL UNIQUE,
			user_id     bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status      order_status NOT NULL,
			accrual     double precision NOT NULL DEFAULT 0,
			uploaded_at timestamp NOT NULL DEFAULT NOW()
		);`)
	if err != nil {
		return fmt.Errorf("pg.createTables.ordersTable: %w", err)
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS withdrawals
		(
			id           bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			user_id      bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			order_number varchar NOT NULL UNIQUE,
			processed_at timestamp NOT NULL DEFAULT NOW(),
			sum          double precision NOT NULL
		);`)
	if err != nil {
		return fmt.Errorf("pg.createTables.withdrawalsTable: %w", err)
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS balances
		(
			id        bigint  PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			user_id   bigint NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			withdrawn double precision NOT NULL DEFAULT 0,
			current   double precision NOT NULL DEFAULT 0
		);`)
	if err != nil {
		return fmt.Errorf("pg.createTables.balancesTable: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("pg.createTables.commit: %w", err)
	}

	return nil
}
