package pg

import (
	"fmt"

	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (pg *pgstorage) GetBalanceByUser(userID int64) (*models.Balance, error) {
	row := pg.db.QueryRow(`
	SELECT withdrawn, current
	FROM balances
	WHERE user_id=$1;`, userID)

	balance := models.Balance{}
	if err := row.Scan(&balance.Withdrawn, &balance.Current); err != nil {
		return nil, fmt.Errorf("pg.getOrdersByUser: %w", err)
	}

	return &balance, nil
}

func (pg *pgstorage) GetWithdrawalsByUser(userID int64) ([]models.Withdrawal, error) {
	rows, err := pg.db.Query(`
		SELECT order_number, processed_at, sum
		FROM withdrawals
		WHERE user_id=$1
		ORDER BY processed_at;`, userID)
	if err != nil {
		return nil, fmt.Errorf("pg.getWithdrawalsByUser: %w", err)
	}
	defer rows.Close()

	dbWithdrawal := []models.Withdrawal{}
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		if err := rows.Scan(&withdrawal.Order, &withdrawal.ProcessedAt, &withdrawal.Sum); err != nil {
			return nil, fmt.Errorf("pg.getWithdrawalsByUser: %w", err)
		}
		dbWithdrawal = append(dbWithdrawal, withdrawal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pg.getWithdrawalsByUser: %w", err)
	}

	return dbWithdrawal, nil
}

func (pg *pgstorage) WithdrawFromUserBalance(orderNumber string, sum models.Money, userID int64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.beginTx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO withdrawals (user_id, order_number, sum)
		VALUES ($1, $2, $3)
		RETURNING id;`, userID, orderNumber, sum)
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance: %w", err)
	}

	var newBalance models.Money
	err = tx.QueryRow(`
		UPDATE balances
		SET withdrawn=balances.withdrawn+$1, current=balances.current-$1
		WHERE user_id=$2
		RETURNING current;`, sum, userID).Scan(&newBalance)
	if err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance: %w", err)
	}
	if newBalance < 0 {
		return appErrors.ErrNegativeBalance
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("pg.withdrawFromUserBalance.commit: %w", err)
	}

	return nil
}
