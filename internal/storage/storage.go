package storage

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

// стор один просто в нем будут различные таблицы и методы для работы с этими таблицами?
type Gophermart interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)

	GetOrders(ctx context.Context, userID string) ([]models.Order, error)
	GetOrder(ctx context.Context, orderID string) (models.Order, error)
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order) (models.Order, error)

	GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawals, error)
	AddWithdrawal(ctx context.Context, withdrawal models.Withdrawals) (models.Withdrawals, error)
}
