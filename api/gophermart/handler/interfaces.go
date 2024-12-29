package handler

import (
	"context"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

type App interface {
	GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error)
	GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error)
	WithdrawFromUserBalance(ctx context.Context, userID int64, withdrawalReq *models.WithdrawalRequest) error

	RegisterOrder(ctx context.Context, userID int64, orderNumber string) error
	GetOrdersByUser(ctx context.Context, userID int64) ([]models.Order, error)
	ValidateOrderNumber(orderNumber string) bool

	ValidateUser(ctx context.Context, user *models.User) (*models.User, error)
	RegisterUser(ctx context.Context, user *models.User) (int64, error)
}
