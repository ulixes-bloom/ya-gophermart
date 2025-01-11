package pg

import (
	"context"
	"fmt"

	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (pg *pgstorage) GetBalanceByUser(ctx context.Context, userID int64) (*models.Balance, error) {
	row := pg.db.QueryRowContext(ctx, `
		SELECT withdrawn, current
		FROM balances
		WHERE user_id=$1;`, userID)

	balance := models.Balance{}
	if err := row.Scan(&balance.Withdrawn, &balance.Current); err != nil {
		return nil, fmt.Errorf("pg.getOrdersByUser: %w", err)
	}

	return &balance, nil
}

func (pg *pgstorage) GetWithdrawalsByUser(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT order_number, processed_at, sum
		FROM withdrawals
		WHERE user_id=$1
		ORDER BY processed_at;`, userID)
	if err != nil {
		return nil, fmt.Errorf("pg.getWithdrawalsByUser.selectWithdrawal: %w", err)
	}
	defer rows.Close()

	dbWithdrawal := []models.Withdrawal{}
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		if err := rows.Scan(&withdrawal.Order, &withdrawal.ProcessedAt, &withdrawal.Sum); err != nil {
			return nil, fmt.Errorf("pg.getWithdrawalsByUser.scanWithdrawal: %w", err)
		}
		dbWithdrawal = append(dbWithdrawal, withdrawal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pg.getWithdrawalsByUser.err: %w", err)
	}

	return dbWithdrawal, nil
}

func (pg *pgstorage) WithdrawFromUserBalance(ctx context.Context, userID int64, orderNumber string, sum models.Money) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.beginTx: %w", err)
	}
	defer tx.Rollback()

	var newBalance models.Money
	err = tx.QueryRowContext(ctx, `
		UPDATE balances
		SET withdrawn=balances.withdrawn+$1, current=balances.current-$1
		WHERE user_id=$2
		RETURNING current;`, sum, userID).Scan(&newBalance)
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.updateBalance: %w", err)
	}
	if newBalance < 0 {
		return appErrors.ErrNegativeBalance
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO withdrawals (user_id, order_number, sum)
		VALUES ($1, $2, $3);`, userID, orderNumber, sum)
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.insertWithdrawal: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.commit: %w", err)
	}

	return nil
}
