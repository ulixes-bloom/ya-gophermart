package pg

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"gotest.tools/v3/assert"
)

func TestStorage_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		expectedError error
	}{
		{
			name: "Success Case",
			user: &models.User{
				Login:    "user",
				Password: "password",
			},
			expectedError: nil,
		},
		{
			name: "Success Case",
			user: &models.User{
				Login:    "user",
				Password: "password",
			},
			expectedError: appErrors.ErrUserLoginAlreadyExists,
		},
		{
			name: "Success Case",
			user: &models.User{
				Login:    "user1",
				Password: "password",
			},
			expectedError: nil,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctr, err := NewPostgresContainer(ctx)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	storage, err := NewStorage(ctx, db)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUserID, err := storage.AddUser(tt.user.Login, tt.user.Password)
			if tt.expectedError != nil {
				assert.Equal(t, createdUserID, int64(-1))
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)

				dbUser, err := storage.GetUserByLogin(tt.user.Login)
				require.NoError(t, err)
				assert.Equal(t, dbUser.ID, createdUserID)
				assert.Equal(t, dbUser.Login, tt.user.Login)
			}
		})
	}
}