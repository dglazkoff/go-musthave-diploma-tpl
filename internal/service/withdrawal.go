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
	/*
		что будет если между UpdateBalance и AddWithdrawal другой запрос спишет баланс до 0?

		пока не понятно, как будто UpdateBalance при выполнении должен делать это в транзакции и если какая то транзакция уже успела
		раньше выполниться и обновить баланс до 0, то эта транзакция должна вернуть ошибку и мы не создаем withdrawal
	*/
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
