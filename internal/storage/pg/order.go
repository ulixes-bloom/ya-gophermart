package pg

import (
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (pg *pgstorage) RegisterOrder(userID int64, orderNumber string) error {
	_, err := pg.db.Exec(`
		INSERT INTO orders (user_id, number, status)
		VALUES ($1, $2, $3)
		RETURNING id;`, userID, orderNumber, "NEW")

	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok &&
			pgerrcode.IsIntegrityConstraintViolation(pgError.Code) &&
			pgError.ConstraintName == "orders_number_key" {

			existingOrder, err := pg.getOrderByNumber(orderNumber)
			if err != nil {
				return fmt.Errorf(`pg.registerOrder: %w`, err)
			}

			if existingOrder.UserID == userID {
				return appErrors.ErrOrderWasUploadedByCurrentUser
			} else {
				return appErrors.ErrOrderWasUploadedByAnotherUser
			}
		}
		return fmt.Errorf("pg.registerOrder: %w", err)
	}

	return nil
}

func (pg *pgstorage) GetOrdersByUser(userID int64) ([]models.Order, error) {
	rows, err := pg.db.Query(`
		SELECT number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id=$1
		ORDER BY uploaded_at;`, userID)
	if err != nil {
		return nil, fmt.Errorf("pg.getOrdersByUser: %w", err)
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("pg.getOrdersByUser: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pg.getOrdersByUser: %w", err)
	}

	return orders, nil
}

func (pg *pgstorage) GetOrdersByStatus(statuses []models.OrderStatus) ([]models.Order, error) {
	orders := []models.Order{}
	for _, status := range statuses {
		rows, err := pg.db.Query(`
			SELECT number, user_id, status, accrual, uploaded_at
			FROM orders
			WHERE status=$1;`, status)
		if err != nil {
			return nil, fmt.Errorf("pg.getOrdersByStatus.query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			order := models.Order{}
			err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
			if err != nil {
				return nil, fmt.Errorf("pg.getOrdersByStatus.scanRow: %w", err)
			}
			orders = append(orders, order)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("pg.getOrdersByStatus.rowsErr: %w", err)
		}
	}

	return orders, nil
}

func (pg *pgstorage) UpdateOrders(orders []models.Order) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, order := range orders {
		_, err := tx.Exec(`
			UPDATE orders
			SET accrual=$1, status=$2
			WHERE number=$3;`, order.Accrual, order.Status, order.Number)
		if err != nil {
			return fmt.Errorf("pg.updateOrders: %w", err)
		}

		_, err = tx.Exec(`
			UPDATE balances
			SET current=balances.current+$1
			WHERE user_id=$2;`, order.Accrual, order.UserID)
		if err != nil {
			return fmt.Errorf("pg.updateOrders: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("pg.updateOrders: %w", err)
	}

	return nil
}

func (pg *pgstorage) getOrderByNumber(orderNumber string) (*models.Order, error) {
	row := pg.db.QueryRow(`
		SELECT id, number, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE number=$1;`, orderNumber)
	order := models.Order{}
	err := row.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("pg.getOrderByNumber: %w", err)
	}
	return &order, nil
}
