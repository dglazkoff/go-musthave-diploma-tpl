package service

import (
	"context"
	"encoding/json"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"net/http"
	"strconv"
	"time"
)

type AccrualOrderStatus string

const (
	Registered AccrualOrderStatus = "REGISTERED"
	Processing AccrualOrderStatus = "PROCESSING"
	Invalid    AccrualOrderStatus = "INVALID"
	Processed  AccrualOrderStatus = "PROCESSED"
)

type AccrualSystemResponse struct {
	Order   string             `json:"order"`
	Status  AccrualOrderStatus `json:"status"`
	Accrual float64            `json:"accrual"`
}

func (s *service) handleSuccessAccrualResponse(accrual *AccrualSystemResponse, order models.Order, userLogin string) error {
	ctx := context.Background()
	logger.Log.Debug("Handle accrual response: ", accrual.Accrual)
	logger.Log.Debug("Handle accrual response with status: ", accrual.Status)
	logger.Log.Debug("Handle accrual response for Order: ", accrual.Order)

	tx, err := s.storage.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.Error("Error while begin transaction: ", err)
		return err
	}

	defer tx.Rollback()

	if accrual.Status == Invalid {
		_, err := s.storage.UpdateOrderTx(
			ctx,
			tx,
			models.Order{ID: order.ID, UserID: order.UserID, UploadedAt: order.UploadedAt, Status: models.Invalid, Accrual: accrual.Accrual},
		)

		if err != nil {
			logger.Log.Error("Error while update order: ", err)
			return err
		}

		err = s.UpdateBalanceTx(context.Background(), tx, accrual.Accrual, userLogin)

		if err != nil {
			logger.Log.Error("Error while update balance: ", err)
			return err
		}
	}

	if accrual.Status == Processed {
		/*
			что будет если после UpdateOrderTx сервис перезапустят для обновления и UpdateBalance вы выпонится?
			что будет если UpdateOrderTx вернет ошибку, баланс все равно нужно обновить?

			сделал транзакцию на обновление данных
		*/
		_, err := s.storage.UpdateOrderTx(
			context.Background(),
			tx,
			models.Order{ID: order.ID, UserID: order.UserID, UploadedAt: order.UploadedAt, Status: models.Processed, Accrual: accrual.Accrual},
		)

		if err != nil {
			logger.Log.Error("Error while update order: ", err)
			return err
		}

		err = s.UpdateBalanceTx(context.Background(), tx, accrual.Accrual, userLogin)

		if err != nil {
			logger.Log.Error("Error while update balance: ", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("Error while commit transaction: ", err)
		return err
	}

	return nil
}

func (s *service) GetAccrual(order models.Order, userLogin string) {
	for {
		logger.Log.Debug("Get accrual for order: ", order.ID)
		response, err := client.R().Get(s.cfg.AccrualSystemAddress + "/api/orders/" + order.ID)

		if err != nil {
			logger.Log.Error("Error while get order status: ", err)
			continue
		}

		if response.StatusCode() == 200 {
			accrualResponse := AccrualSystemResponse{}
			if err := json.Unmarshal(response.Body(), &accrualResponse); err != nil {
				logger.Log.Debug("Error while decode accrual response: ", err)
				continue
			}

			if err = s.handleSuccessAccrualResponse(&accrualResponse, order, userLogin); err != nil {
				continue
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

		break
	}
}
