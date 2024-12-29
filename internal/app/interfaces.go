package app

import (
	"context"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

type (
	// Интерфейс хранилища
	Storage interface {
		GetUserByLogin(ctx context.Context, login string) (*models.User, error)
		AddUser(ctx context.Context, login, password string) (int64, error)

		RegisterOrder(ctx context.Context, userID int64, orderNumber string) error
		GetOrdersByUser(ctx context.Context, userID int64) ([]models.Order, error)
		GetOrdersByStatus(ctx context.Context, statuses []models.OrderStatus) ([]models.Order, error)
		SetOrdersAccrualAndUpdateBalance(ctx context.Context, orders []models.Order) error

		GetBalanceByUser(ctx context.Context, userID int64) (*models.Balance, error)
		GetWithdrawalsByUser(ctx context.Context, userID int64) ([]models.Withdrawal, error)
		WithdrawFromUserBalance(ctx context.Context, userID int64, orderNumber string, sum models.Money) error

		Close() error
	}
)
