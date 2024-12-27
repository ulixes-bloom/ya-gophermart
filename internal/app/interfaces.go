package app

import (
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

type (
	// Интерфейс хранилища
	Storage interface {
		GetUserByLogin(login string) (*models.User, error)
		AddUser(login, password string) (int64, error)

		RegisterOrder(userID int64, orderNumber string) error
		GetOrdersByUser(userID int64) ([]models.Order, error)
		GetOrdersByStatus(statuses []models.OrderStatus) ([]models.Order, error)
		UpdateOrders(orders []models.Order) error

		GetBalanceByUser(userID int64) (*models.Balance, error)
		GetWithdrawalsByUser(userID int64) ([]models.Withdrawal, error)
		WithdrawFromUserBalance(userID int64, orderNumber string, sum models.Money) error

		Close() error
	}
)
