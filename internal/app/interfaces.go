package app

import (
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

type (
	// Интерфейс хранилища
	Storage interface {
		GetUserByLogin(login string) (*models.User, error)
		AddUser(login, password string) (int64, error)

		RegisterOrder(orderNumber string, userID int64) error
		GetOrdersByUser(userID int64) ([]models.Order, error)
		GetOrdersByStatus(statuses []models.OrderStatus) ([]models.Order, error)
		UpdateOrders(orders []models.Order) error

		GetBalanceByUser(userID int64) (*models.Balance, error)
		GetWithdrawalsByUser(userID int64) ([]models.Withdrawal, error)
		WithdrawFromUserBalance(orderNumber string, sum models.Money, userID int64) error

		Close() error
	}
)
