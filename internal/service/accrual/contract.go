package accrual

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

type storage interface {
	UpdateOrderTx(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error)
	BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error)
	GetNotAccrualOrders(ctx context.Context) ([]models.Order, error)
}

type externalService interface {
	UpdateBalanceTx(ctx context.Context, tx *sql.Tx, accrual float64, userLogin string) error
}
