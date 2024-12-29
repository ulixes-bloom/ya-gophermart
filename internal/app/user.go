package app

import (
	"context"
	"fmt"

	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

func (a *App) ValidateUser(ctx context.Context, user *models.User) (*models.User, error) {
	dbUser, err := a.storage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("app.validateUser: %w", err)
	}

	if err := security.CheckPassword(user.Password, dbUser.Password); err != nil {
		return nil, fmt.Errorf("app.validateUser: %w", err)
	}

	return dbUser, nil
}

func (a *App) RegisterUser(ctx context.Context, user *models.User) (int64, error) {
	createdUserID, err := a.storage.AddUser(ctx, user.Login, user.Password)
	if err != nil {
		return -1, fmt.Errorf("app.registerUser: %w", err)
	}

	return createdUserID, nil
}
