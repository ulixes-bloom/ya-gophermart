package handler

import "github.com/ulixes-bloom/ya-gophermart/internal/models"

type App interface {
	GetUserBalance(userID int64) (*models.Balance, error)
	GetUserWithdrawals(userID int64) ([]models.Withdrawal, error)
	WithdrawFromUserBalance(withdrawalReq *models.WithdrawalRequest, userID int64) error

	RegisterOrder(orderNumber string, userID int64) error
	GetOrdersByUser(userID int64) ([]models.Order, error)
	ValidateOrderNumber(orderNumber string) bool

	ValidateUser(user *models.User) (*models.User, error)
	RegisterUser(user *models.User) (int64, error)
}
