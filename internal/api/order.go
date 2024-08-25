package api

import (
	"encoding/json"
	"errors"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/auth"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func isValidOrderNumber(orderNumber string) bool {
	sum := 0
	orderNumberSlice := strings.Split(orderNumber, "")
	parity := len(orderNumberSlice) % 2

	for i := 0; i < len(orderNumberSlice); i++ {
		number, err := strconv.Atoi(orderNumberSlice[i])

		if err != nil {
			logger.Log.Error("Error while convert to int: ", err)
			return false
		}

		if i%2 == parity {
			if number*2 > 9 {
				sum += number*2 - 9
			} else {
				sum += number * 2
			}
		} else {
			sum += number
		}
	}

	return sum%10 == 0
}

func (a *api) AddOrder(writer http.ResponseWriter, request *http.Request) {
	orderNumber, err := io.ReadAll(request.Body)

	if err != nil {
		logger.Log.Error("Error while reading request body: ", err)
		writer.WriteHeader(http.StatusBadRequest)
	}

	if !isValidOrderNumber(string(orderNumber)) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return

	}

	userID, ok := auth.GetUserIDFromRequest(request)
	if !ok {
		logger.Log.Error("Error while get userID from request")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.s.AddOrder(request.Context(), string(orderNumber), userID)

	if err != nil {
		if errors.Is(err, service.ErrorOrderAlreadyAdded) {
			writer.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, service.ErrorOrderAlreadyAddedByAnotherUser) {
			writer.WriteHeader(http.StatusConflict)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

func (a *api) GetOrders(writer http.ResponseWriter, request *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(request)

	if !ok {
		logger.Log.Error("Error while get userID from request")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := a.s.GetOrders(request.Context(), userID)

	if err != nil {
		if errors.Is(err, service.ErrorNoOrders) {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(orders)

	if err != nil {
		logger.Log.Error("Error while encode: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
