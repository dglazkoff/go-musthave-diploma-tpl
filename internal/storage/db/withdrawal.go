package db

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (s *dbStorage) GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawals, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, sum, user_id, processed_at from withdrawals WHERE user_id = $1", userID)
	var withdrawals []models.Withdrawals

	if err != nil {
		logger.Log.Error("error while reading withdrawals: ", err)
		return nil, err
	}

	for rows.Next() {
		var withdrawal models.Withdrawals
		err = rows.Scan(&withdrawal.ID, &withdrawal.Sum, &withdrawal.UserID, &withdrawal.ProcessedAt)

		if err != nil {
			logger.Log.Error("error while scan withdrawal: ", err)
			continue
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if rows.Err() != nil {
		logger.Log.Debug("error from rows ", err)
	}

	return withdrawals, nil
}

// большое дублирование
func (s *dbStorage) GetWithdrawalsTx(ctx context.Context, tx *sql.Tx, userID string) ([]models.Withdrawals, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id, sum, user_id, processed_at from withdrawals WHERE user_id = $1", userID)
	var withdrawals []models.Withdrawals

	if err != nil {
		logger.Log.Error("error while reading withdrawals: ", err)
		return nil, err
	}

	for rows.Next() {
		var withdrawal models.Withdrawals
		err = rows.Scan(&withdrawal.ID, &withdrawal.Sum, &withdrawal.UserID, &withdrawal.ProcessedAt)

		if err != nil {
			logger.Log.Error("error while scan withdrawal: ", err)
			continue
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if rows.Err() != nil {
		logger.Log.Debug("error from rows ", err)
	}

	return withdrawals, nil
}

func (s *dbStorage) AddWithdrawal(ctx context.Context, withdrawal models.Withdrawals) (models.Withdrawals, error) {
	_, err := s.db.ExecContext(ctx, "INSERT INTO withdrawals (id, sum, user_id, processed_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING", withdrawal.ID, withdrawal.Sum, withdrawal.UserID, withdrawal.ProcessedAt)

	if err != nil {
		logger.Log.Error("error while add withdrawal: ", err)
		return models.Withdrawals{}, err
	}

	return withdrawal, nil
}
