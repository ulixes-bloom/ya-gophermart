package models

import (
	"fmt"
	"time"
)

type Money float64

//easyjson:json
type (
	Balance struct {
		ID        int64 `json:"-"`
		UserID    int64 `json:"-"`
		Withdrawn Money `json:"withdrawn"`
		Current   Money `json:"current"`
	}

	//easyjson:json
	Withdrawal struct {
		ID          int64     `json:"-"`
		UserID      int64     `json:"-"`
		Order       string    `json:"order"`
		ProcessedAt time.Time `json:"processed_at"`
		Sum         Money     `json:"sum"`
	}

	//easyjson:json//easyjson:json
	WithdrawalRequest struct {
		Order string `json:"order"`
		Sum   Money  `json:"sum"`
	}
)

func (b *Balance) String() string {
	if b == nil {
		return "balance is nil pointer"
	}

	return fmt.Sprintf("UserID: %d, Withdrawn: %f, Current: %f", b.UserID, b.Withdrawn, b.Current)
}

func (w *Withdrawal) String() string {
	if w == nil {
		return "balance is nil pointer"
	}

	return fmt.Sprintf("UserID: %d, Order: %s, Sum: %f", w.UserID, w.Order, w.Sum)
}
