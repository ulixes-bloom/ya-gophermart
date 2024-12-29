package pg

import (
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

func (pg *pgstorage) AddUser(login, password string) (int64, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("pg.addUser: %w", err)
	}
	defer tx.Rollback()

	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return -1, fmt.Errorf("pg.addUser: %w", err)
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
		return -1, fmt.Errorf("pg.addUser: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO balances (user_id)
		VALUES ($1);`, id)
	if err != nil {
		return -1, fmt.Errorf("pg.addUser: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("pg.addUser: %w", err)
	}

	return id, nil
}

func (pg *pgstorage) GetUserByLogin(login string) (*models.User, error) {
	dbUser := &models.User{}
	row := pg.db.QueryRow(`
		SELECT id, login, password
		FROM users
		WHERE login=$1;`, login)
	if err := row.Scan(&dbUser.ID, &dbUser.Login, &dbUser.Password); err != nil {
		return nil, fmt.Errorf("pg.getUserByLogin: %w", err)
	}
	return dbUser, nil
}
