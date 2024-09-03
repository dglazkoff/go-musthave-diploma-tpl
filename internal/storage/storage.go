package storage

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

// стор один просто в нем будут различные таблицы и методы для работы с этими таблицами?
type Gophermart interface {
	CreateUser(ctx context.Context, tx *sql.Tx, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByLoginTx(ctx context.Context, tx *sql.Tx, login string) (models.User, error)
	GetUserByLoginForUpdate(ctx context.Context, tx *sql.Tx, login string) (models.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, user models.User) (models.User, error)

	GetNotAccrualOrders(ctx context.Context) ([]models.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]models.Order, error)
	GetOrder(ctx context.Context, orderID string) (models.Order, error)
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	UpdateOrderTx(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error)

	GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawals, error)
	GetWithdrawalsTx(ctx context.Context, tx *sql.Tx, userID string) ([]models.Withdrawals, error)
	AddWithdrawal(ctx context.Context, withdrawal models.Withdrawals) (models.Withdrawals, error)

	BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error)
}
