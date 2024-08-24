package db

import (
	"context"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (d *dbStorage) GetWithdrawals(ctx context.Context, userId string) ([]models.Withdrawals, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, sum, user_id, processed_at from withdrawals WHERE user_id = $1", userId)
	var withdrawals []models.Withdrawals

	if err != nil {
		logger.Log.Debug("error while reading withdrawals: ", err)
		return nil, err
	}

	for rows.Next() {
		var withdrawal models.Withdrawals
		err = rows.Scan(&withdrawal.ID, &withdrawal.Sum, &withdrawal.UserId, &withdrawal.ProcessedAt)

		if err != nil {
			logger.Log.Debug("error while scan withdrawal: ", err)
			continue
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}

func (d *dbStorage) AddWithdrawal(ctx context.Context, withdrawal models.Withdrawals) (models.Withdrawals, error) {
	fmt.Println(withdrawal)
	_, err := d.db.ExecContext(ctx, "INSERT INTO withdrawals (id, sum, user_id, processed_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING", withdrawal.ID, withdrawal.Sum, withdrawal.UserId, withdrawal.ProcessedAt)

	if err != nil {
		logger.Log.Error("error while add withdrawal: ", err)
		return models.Withdrawals{}, err
	}

	return withdrawal, nil
}
