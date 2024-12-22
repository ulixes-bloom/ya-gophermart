package pg

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
	"gotest.tools/v3/assert"
)

func TestStorage_CreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	storage, err := newPostgresStorage(ctx)
	require.NoError(t, err)

	user := models.User{
		Login:    "login",
		Password: "password",
	}

	// Создание пользователя в БД
	userID, err := storage.AddUser(user.Login, user.Password)
	require.NoError(t, err)

	// Поиск созданного пользователя в БД
	dbUser, err := storage.GetUserByLogin(user.Login)
	require.NoError(t, err)
	require.NotNil(t, dbUser)

	// Проверка результата
	assert.Equal(t, user.Login, dbUser.Login)
	assert.Equal(t, userID, dbUser.ID)
	assert.NilError(t, security.CheckPassword(user.Password, dbUser.Password))

	// Создание пользователя с уже существующим логином
	_, err = storage.AddUser(user.Login, user.Password)
	assert.ErrorIs(t, err, appErrors.ErrUserLoginAlreadyExists)

	// Поиск несуществующего пользователя в БД
	dbUser, err = storage.GetUserByLogin("not_existing_login")
	require.Error(t, err)
	require.Nil(t, dbUser)
}

func TestStorage_WithdrawFromBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	storage, err := newPostgresStorage(ctx)
	require.NoError(t, err)

	user := models.User{
		Login:    "login",
		Password: "password",
	}

	// Создание пользователя в БД
	userID, err := storage.AddUser(user.Login, user.Password)
	require.NoError(t, err)

	// Поиск баланса созданного пользователя в БД
	dbBalance, err := storage.GetBalanceByUser(userID)
	require.NoError(t, err)
	require.NotNil(t, dbBalance)

	assert.Equal(t, dbBalance.Current, models.Money(0))
	assert.Equal(t, dbBalance.Withdrawn, models.Money(0))

	withdrawalReq := models.WithdrawalRequest{
		Order: "2377225624",
		Sum:   models.Money(200),
	}

	// Попытка списания со счета при недостаточном значении баланса
	err = storage.WithdrawFromUserBalance(withdrawalReq.Order, withdrawalReq.Sum, userID)
	assert.ErrorIs(t, err, appErrors.ErrNegativeBalance)

	order := models.Order{
		Number:  "12345678903",
		UserID:  userID,
		Status:  models.OrderStatusProcessed,
		Accrual: 300,
	}

	// Регистрация заказа с начислнием бонусов
	err = storage.RegisterOrder(order.Number, userID)
	require.NoError(t, err)
	err = storage.UpdateOrders([]models.Order{order})
	require.NoError(t, err)

	// Получениие баланса из БД
	dbBalance, err = storage.GetBalanceByUser(userID)
	require.NoError(t, err)
	require.NotNil(t, dbBalance)

	// Проверка пополнения баланса после загруски заказа с начислением бонусов
	assert.Equal(t, dbBalance.Current, models.Money(300))
	assert.Equal(t, dbBalance.Withdrawn, models.Money(0))

	// Попытка списания со счета при достаточном значении баланса
	err = storage.WithdrawFromUserBalance(withdrawalReq.Order, withdrawalReq.Sum, userID)
	require.NoError(t, err)

	// Проверка результата списания средств со счета
	dbBalance, err = storage.GetBalanceByUser(userID)
	require.NoError(t, err)
	require.NotNil(t, dbBalance)

	assert.Equal(t, dbBalance.Current, models.Money(100))
	assert.Equal(t, dbBalance.Withdrawn, models.Money(200))
}

func newPostgresStorage(ctx context.Context) (*pgstorage, error) {
	dbName := "gophermart"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	storage, err := NewStorage(ctx, db)
	if err != nil {
		return nil, err
	}

	return storage, err
}
