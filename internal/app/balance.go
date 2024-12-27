package app

import (
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
)

func (a *App) GetUserBalance(userID int64) (*models.Balance, error) {
	dbBalance, err := a.storage.GetBalanceByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("app.getUserBalance: %w", err)
	}
	return dbBalance, nil
}

func (a *App) GetUserWithdrawals(userID int64) ([]models.Withdrawal, error) {
	dbWithdrawals, err := a.storage.GetWithdrawalsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("app.getUserBalance: %w", err)
	}
	return dbWithdrawals, nil
}

func (a *App) WithdrawFromUserBalance(userID int64, withdrawalReq *models.WithdrawalRequest) error {
	err := a.storage.WithdrawFromUserBalance(userID, withdrawalReq.Order, withdrawalReq.Sum)
	if err != nil {
		return fmt.Errorf("app.withdrawFromUserBalance: %w", err)
	}
	return nil
}
