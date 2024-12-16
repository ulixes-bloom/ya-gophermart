package app

import (
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

var NotProcessedOrderStatuses = []models.OrderStatus{
	models.OrderStatusNew,
	models.OrderStatusProcessing,
}

func (a *App) UpdateNotProcessedOrders() error {
	notProcessedOrders, err := a.storage.GetOrdersByStatus(NotProcessedOrderStatuses)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.getOrdersFromStorage: %w", err)
	}

	updatedOrders, err := a.ac.GetOrdersInfo(notProcessedOrders)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.getOrdersInfo: %w", err)
	}

	err = a.storage.UpdateOrders(updatedOrders)
	if err != nil {
		return fmt.Errorf("app.updateNotProcessedOrders.updateOrders: %w", err)
	}

	return nil
}
