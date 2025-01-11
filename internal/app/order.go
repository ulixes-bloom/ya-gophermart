package app

import (
	"context"
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/luhn"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (a *App) RegisterOrder(ctx context.Context, userID int64, orderNumber string) error {
	err := a.storage.RegisterOrder(ctx, userID, orderNumber)
	if err != nil {
		return fmt.Errorf("app.registerOrder: %w", err)
	}

	return nil
}

func (a *App) GetOrdersByUser(ctx context.Context, userID int64) ([]models.Order, error) {
	orders, err := a.storage.GetOrdersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.getOrdersByUser: %w", err)
	}
	return orders, nil
}

func (a *App) ValidateOrderNumber(orderNumber string) bool {
	return luhn.ValidateNumber(orderNumber)
}
