package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"time"
)

var ErrorOrderAlreadyAdded = fmt.Errorf("order already added")
var ErrorOrderAlreadyAddedByAnotherUser = fmt.Errorf("order already added by another user")
var ErrorNoOrders = fmt.Errorf("no orders")

var client = resty.New()

func (s *service) AddOrder(ctx context.Context, orderNumber string, userLogin string) error {
	// нужно ли транзакцию на случай если ордер добавят пока мы делаем проверки ??
	order, err := s.storage.GetOrder(ctx, orderNumber)

	if err == nil {
		if order.UserID == userLogin {
			return ErrorOrderAlreadyAdded
		}

		return ErrorOrderAlreadyAddedByAnotherUser
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("Error while get order: ", err)
		return err
	}

	order, err = s.storage.CreateOrder(ctx, models.Order{ID: orderNumber, UserID: userLogin, Status: models.New, Accrual: 0, UploadedAt: time.Now().Format(time.RFC3339)})

	if err != nil {
		logger.Log.Error("Error while create order: ", err)
		return err
	}

	go func() {
		s.GetAccrual(order, userLogin)
	}()

	return nil
}

func (s *service) GetOrders(ctx context.Context, userId string) ([]models.Order, error) {
	orders, err := s.storage.GetOrders(ctx, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orders, ErrorNoOrders
		}

		logger.Log.Error("Error while get orders: ", err)
		return orders, err
	}

	return orders, nil
}
