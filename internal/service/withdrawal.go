package service

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"time"
)

func (s *service) GetWithdrawals(ctx context.Context, userId string) ([]models.Withdrawals, error) {
	withdrawals, err := s.storage.GetWithdrawals(ctx, userId)

	if err != nil {
		logger.Log.Debug("error while reading withdrawals: ", err)
		return nil, err
	}

	return withdrawals, nil
}

func (s *service) CreateWithdrawal(ctx context.Context, orderNumber string, sum float64, userId string) error {
	err := s.UpdateBalance(ctx, -sum, userId)

	if err != nil {
		logger.Log.Error("Error while update balance: ", err)
		return err
	}

	_, err = s.storage.AddWithdrawal(ctx, models.Withdrawals{ID: orderNumber, UserId: userId, Sum: sum, ProcessedAt: time.Now().Format(time.RFC3339)})

	if err != nil {
		logger.Log.Error("Error while add withdrawal: ", err)
		return err
	}

	return nil
}
