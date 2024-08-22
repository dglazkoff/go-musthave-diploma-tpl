package db

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (d *dbStorage) GetWithdrawals(ctx context.Context, orderId string) ([]models.Withdrawals, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT order, sum, user_id, processed_at from withdrawals WHERE order = $1", orderId)
	var withdrawals []models.Withdrawals

	if err != nil {
		logger.Log.Debug("error while reading withdrawals: ", err)
		return nil, err
	}

	for rows.Next() {
		var withdrawal models.Withdrawals
		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.UserId, &withdrawal.ProcessedAt)

		if err != nil {
			logger.Log.Debug("error while scan withdrawal: ", err)
			continue
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}

func (d *dbStorage) AddWithdrawal(ctx context.Context, withdrawal models.Withdrawals) (models.Withdrawals, error) {
	_, err := d.db.ExecContext(ctx, "INSERT INTO withdrawals (order, sum, user_id, processed_at) VALUES ($1, $2, $3, $4) ON CONFLICT (order) DO NOTHING", withdrawal.Order, withdrawal.Sum, withdrawal.UserId, withdrawal.ProcessedAt)

	if err != nil {
		logger.Log.Error("error while add withdrawal: ", err)
		return models.Withdrawals{}, err
	}

	return withdrawal, nil
}
