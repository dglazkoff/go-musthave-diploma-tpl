package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"net/http"
	"strconv"
	"time"
)

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

func (s *service) GetAccrual(order models.Order, userLogin string) {
	for {
		response, err := client.R().Get(s.cfg.AccrualSystemAddress + "/api/orders/" + order.ID)

		fmt.Println(response.StatusCode())

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
}
