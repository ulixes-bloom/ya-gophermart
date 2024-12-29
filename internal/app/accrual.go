package app

import (
	"context"
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

var NotProcessedOrderStatuses = []models.OrderStatus{
	models.OrderStatusNew,
	models.OrderStatusProcessing,
}

func (a *App) UpdateNotProcessedOrders(ctx context.Context) error {
	notProcessedOrders, err := a.storage.GetOrdersByStatus(ctx, NotProcessedOrderStatuses)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.getOrdersFromStorage: %w", err)
	}

	updatedOrders, err := a.ac.GetOrdersInfo(ctx, notProcessedOrders)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.getOrdersInfo: %w", err)
	}

	err = a.storage.SetOrdersAccrualAndUpdateBalance(ctx, updatedOrders)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.setOrdersAccrualAndUpdateBalance: %w", err)
	}

	return nil
}
