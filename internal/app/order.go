package app

import (
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/luhn"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (a *App) RegisterOrder(userID int64, orderNumber string) error {
	err := a.storage.RegisterOrder(userID, orderNumber)
	if err != nil {
		return fmt.Errorf("app.registerOrder: %w", err)
	}

	return nil
}

func (a *App) GetOrdersByUser(userID int64) ([]models.Order, error) {
	orders, err := a.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("app.getOrdersByUser: %w", err)
	}
	return orders, nil
}

func (a *App) ValidateOrderNumber(orderNumber string) bool {
	return luhn.ValidateNumber(orderNumber)
}
