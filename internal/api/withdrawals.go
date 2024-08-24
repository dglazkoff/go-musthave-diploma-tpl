package api

import (
	"encoding/json"
	"errors"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/auth"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
	"net/http"
)

type CreateWithdrawalRequest struct {
	Order string `json:"order"`
	Sum   uint   `json:"sum"`
}

func (a *api) GetWithdrawals(writer http.ResponseWriter, request *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(request)
	if !ok {
		logger.Log.Error("Error while get userID from request")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := a.s.GetWithdrawals(request.Context(), userID)

	if err != nil {
		if errors.Is(err, service.ErrorNoOrders) {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(withdrawals)

	if err != nil {
		logger.Log.Error("Error while encode: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (a *api) CreateWithdrawal(writer http.ResponseWriter, request *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(request)
	if !ok {
		logger.Log.Error("Error while get userID from request")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var withdrawalRequest CreateWithdrawalRequest

	err := json.NewDecoder(request.Body).Decode(&withdrawalRequest)

	if err != nil {
		logger.Log.Error("Error while decode: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isValidOrderNumber(withdrawalRequest.Order) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return

	}

	err = a.s.CreateWithdrawal(request.Context(), withdrawalRequest.Order, withdrawalRequest.Sum, userID)

	if err != nil {
		if errors.Is(err, service.ErrorNotEnoughBalance) {
			writer.WriteHeader(http.StatusPaymentRequired)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
