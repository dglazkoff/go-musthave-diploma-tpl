package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"time"
)

var ErrorOrderAlreadyAdded = fmt.Errorf("order already added")
var ErrorOrderAlreadyAddedByAnotherUser = fmt.Errorf("order already added by another user")
var ErrorNoOrders = fmt.Errorf("no orders")

// var ErrorWrongOrderNumber = fmt.Errorf("wrong order number")

func (s *service) AddOrder(ctx context.Context, orderNumber string, userLogin string) error {
	// транзакцию тоже на случай если ордер добавят пока мы делаем проверки ??
	order, err := s.storage.GetOrder(ctx, orderNumber)

	if err == nil {
		if order.UserID == userLogin {
			return ErrorOrderAlreadyAdded
		}

		return ErrorOrderAlreadyAddedByAnotherUser
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("Error while get user by login: ", err)
		return err
	}

	// запускать сразу горутину которая опрашивает или через какую то общую очередь задач?
	_, err = s.storage.CreateOrder(ctx, models.Order{ID: orderNumber, UserID: userLogin, Status: models.New, Accrual: 0, UploadedAt: time.Now().Format(time.RFC3339)})

	if err != nil {
		logger.Log.Error("Error while create order: ", err)
		return err
	}

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
