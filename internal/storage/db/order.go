package db

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (d *dbStorage) GetOrder(ctx context.Context, orderNumber string) (order models.Order, err error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, user_id, status, uploaded_at, accrual from orders WHERE id = $1", orderNumber)
	err = row.Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt, &order.Accrual)

	return
}

func (d *dbStorage) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	_, err := d.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, uploaded_at, accrual) VALUES($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
		order.ID, order.UserID, order.Status, order.UploadedAt, order.Accrual,
	)

	if err != nil {
		logger.Log.Debug("error while creating order ", err)
		return order, err
	}

	return order, nil
}

func (d *dbStorage) GetOrders(ctx context.Context, userId string) ([]models.Order, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, user_id, status, uploaded_at, accrual from orders WHERE user_id = $1", userId)
	var orders []models.Order

	if err != nil {
		logger.Log.Debug("error while reading orders ", err)
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt, &order.Accrual)

		if err != nil {
			logger.Log.Debug("error while scan order ", err)
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}
