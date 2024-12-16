package models

import (
	"fmt"
	"time"
)

type (
	//easyjson:json
	Order struct {
		ID         int64       `json:"-"`
		UserID     int64       `json:"-"`
		Number     string      `json:"number"`
		Status     OrderStatus `json:"status"`
		Accrual    Money       `json:"accrual,omitempty"`
		UploadedAt time.Time   `json:"uploaded_at"`
	}

	OrderStatus string
)

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type OrderRequest struct {
	Number string `json:"number"`
}

func NewOrder(userID int64, number string, status OrderStatus, accrual Money) *Order {
	return &Order{
		UserID:  userID,
		Number:  number,
		Status:  status,
		Accrual: accrual,
	}
}

func NewOrderRequest(number string) *OrderRequest {
	return &OrderRequest{
		Number: number,
	}
}

func (o *Order) String() string {
	if o == nil {
		return "order is nil pointer"
	}

	return fmt.Sprintf("Number: %s, Status: %s, Accrual: %f", o.Number, o.Status, o.Accrual)
}
