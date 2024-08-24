package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"time"
)

var ErrorOrderAlreadyAdded = fmt.Errorf("order already added")
var ErrorOrderAlreadyAddedByAnotherUser = fmt.Errorf("order already added by another user")
var ErrorNoOrders = fmt.Errorf("no orders")

// var ErrorWrongOrderNumber = fmt.Errorf("wrong order number")

var client = resty.New()

type AccrualOrderStatus string

const (
	Registered AccrualOrderStatus = "REGISTERED"
	Processing                    = "PROCESSING"
	Invalid                       = "INVALID"
	Processed                     = "PROCESSED"
)

type AccrualSystemResponse struct {
	Order   string             `json:"order"`
	Status  AccrualOrderStatus `json:"status"`
	Accrual float64            `json:"accrual"`
}

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

	order, err = s.storage.CreateOrder(ctx, models.Order{ID: orderNumber, UserID: userLogin, Status: models.New, Accrual: 0, UploadedAt: time.Now().Format(time.RFC3339)})

	if err != nil {
		logger.Log.Error("Error while create order: ", err)
		return err
	}

	go func() {
		for {
			response, err := client.R().Get(s.cfg.AccrualSystemAddress + "/api/orders/" + orderNumber)

			fmt.Println(response.StatusCode())

			if err != nil {
				logger.Log.Error("Error while get order status: ", err)
				continue
			}

			if response.StatusCode() == 200 {
				accrualResponse := AccrualSystemResponse{}
				if err := json.Unmarshal(response.Body(), &accrualResponse); err != nil {
					logger.Log.Debug("Error while decode accrual response: ", err)
					// return err
				}

				if accrualResponse.Status == Invalid {
					_, err := s.storage.UpdateOrder(
						context.Background(),
						models.Order{ID: order.ID, UserID: order.UserID, UploadedAt: order.UploadedAt, Status: models.Invalid, Accrual: accrualResponse.Accrual},
					)

					if err != nil {
						logger.Log.Error("Error while update order: ", err)
					}

					err = s.UpdateBalance(context.Background(), accrualResponse.Accrual, userLogin)

					if err != nil {
						logger.Log.Error("Error while update balance: ", err)
					}

					return
				}

				if accrualResponse.Status == Processed {
					_, err := s.storage.UpdateOrder(
						context.Background(),
						models.Order{ID: order.ID, UserID: order.UserID, UploadedAt: order.UploadedAt, Status: models.Processed, Accrual: accrualResponse.Accrual},
					)

					if err != nil {
						logger.Log.Error("Error while update order: ", err)
					}

					err = s.UpdateBalance(context.Background(), accrualResponse.Accrual, userLogin)

					if err != nil {
						logger.Log.Error("Error while update balance: ", err)
					}

					return
				}
			}

			if response.StatusCode() == http.StatusNoContent {
				time.Sleep(1 * time.Second)
				continue
			}

			if response.StatusCode() == http.StatusTooManyRequests {
				timeToRetry, err := strconv.Atoi(response.Header().Get("Retry-After"))

				if err != nil {
					logger.Log.Error("Error while parse Retry-After header: ", err)
					continue
				}

				time.Sleep(time.Duration(timeToRetry) * time.Second)
			}
		}

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
