package service

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"time"
)

func (s *service) GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawals, error) {
	withdrawals, err := s.storage.GetWithdrawals(ctx, userID)

	if err != nil {
		logger.Log.Debug("error while reading withdrawals: ", err)
		return nil, err
	}

	return withdrawals, nil
}

func (s *service) CreateWithdrawal(ctx context.Context, orderNumber string, sum float64, userID string) error {
	logger.Log.Debug("Create withdrawal: ", sum)
	err := s.UpdateBalance(ctx, -sum, userID)

	if err != nil {
		logger.Log.Error("Error while update balance: ", err)
		return err
	}

	_, err = s.storage.AddWithdrawal(ctx, models.Withdrawals{ID: orderNumber, UserID: userID, Sum: sum, ProcessedAt: time.Now().Format(time.RFC3339)})

	if err != nil {
		logger.Log.Error("Error while add withdrawal: ", err)
		return err
	}

	return nil
}
