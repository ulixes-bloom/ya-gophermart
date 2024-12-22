package pg

import (
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

func (pg *pgstorage) AddUser(login, password string) (int64, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("pg.AddUser: %w", err)
	}
	defer tx.Rollback()

	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return -1, fmt.Errorf("pg.AddUser: %w", err)
	}

	var id int64
	err = tx.QueryRow(`
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id;`, login, hashPassword).Scan(&id)
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok &&
			pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == "users_login_key" {
			return -1, fmt.Errorf(`pg.AddUser: %w: %s`, appErrors.ErrUserLoginAlreadyExists, err)
		}
		return -1, fmt.Errorf("pg.AddUser: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO balances (user_id)
		VALUES ($1)
		RETURNING id;`, id)
	if err != nil {
		return -1, fmt.Errorf("pg.AddUser: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("pg.AddUser: %w", err)
	}

	return id, nil
}

func (pg *pgstorage) GetUserByLogin(login string) (*models.User, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("pg.GetUserByLogin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Msg("pg.GetUserByLogin: tx failed to rollback")
		}
	}()

	dbUser := &models.User{}
	row := tx.QueryRow(`
		SELECT id, login, password
		FROM users
		WHERE login=$1;`, login)
	if err := row.Scan(&dbUser.ID, &dbUser.Login, &dbUser.Password); err != nil {
		return nil, fmt.Errorf("pg.GetUserByLogin: %w", err)
	}
	return dbUser, nil
}
