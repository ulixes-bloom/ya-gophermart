package app

import (
	"context"
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (a *App) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	dbBalance, err := a.storage.GetBalanceByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.getUserBalance: %w", err)
	}
	return dbBalance, nil
}

func (a *App) GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	dbWithdrawals, err := a.storage.GetWithdrawalsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.getUserBalance: %w", err)
	}
	return dbWithdrawals, nil
}

func (a *App) WithdrawFromUserBalance(ctx context.Context, userID int64, withdrawalReq *models.WithdrawalRequest) error {
	err := a.storage.WithdrawFromUserBalance(ctx, userID, withdrawalReq.Order, withdrawalReq.Sum)
	if err != nil {
		return fmt.Errorf("app.withdrawFromUserBalance: %w", err)
	}
	return nil
}
