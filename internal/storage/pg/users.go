package pg

import (
	"context"
	"fmt"

	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

func (pg *pgstorage) AddUser(ctx context.Context, login, password string) (int64, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("pg.addUser.beginTx: %w", err)
	}
	defer tx.Rollback()

	hashPassword, err := security.HashPassword(password)
	if err != nil {
		return -1, fmt.Errorf("pg.addUser.hashPassword: %w", err)
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id;`, login, hashPassword).Scan(&id)
	if err != nil {
		if isUniqueViolation(err, "users_login_key") {
			return -1, fmt.Errorf(`pg.AddUser: %w: %s`, appErrors.ErrUserLoginAlreadyExists, err)
		}
		return -1, fmt.Errorf("pg.addUser: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO balances (user_id)
		VALUES ($1);`, id)
	if err != nil {
		return -1, fmt.Errorf("pg.addUser.createBalance: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("pg.addUser.commit: %w", err)
	}

	return id, nil
}

func (pg *pgstorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	dbUser := &models.User{}
	row := pg.db.QueryRowContext(ctx, `
		SELECT id, login, password
		FROM users
		WHERE login=$1;`, login)
	if err := row.Scan(&dbUser.ID, &dbUser.Login, &dbUser.Password); err != nil {
		return nil, fmt.Errorf("pg.getUserByLogin: %w", err)
	}
	return dbUser, nil
}
