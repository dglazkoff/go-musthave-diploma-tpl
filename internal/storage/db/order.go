package db

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (s *dbStorage) GetOrder(ctx context.Context, orderNumber string) (order models.Order, err error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, user_id, status, uploaded_at, accrual from orders WHERE id = $1", orderNumber)
	err = row.Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt, &order.Accrual)

	return
}

func (s *dbStorage) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	_, err := s.db.ExecContext(
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

func (s *dbStorage) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, user_id, status, uploaded_at, accrual from orders WHERE user_id = $1", userID)
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

	if rows.Err() != nil {
		logger.Log.Debug("error from rows ", err)
	}

	return orders, nil
}

func (s *dbStorage) UpdateOrderTx(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, uploaded_at, accrual) VALUES($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET status = $3, accrual = $5",
		order.ID, order.UserID, order.Status, order.UploadedAt, order.Accrual,
	)

	if err != nil {
		logger.Log.Debug("error while updating order ", err)
		return order, err
	}

	return order, nil
}
