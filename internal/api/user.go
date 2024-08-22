package api

import (
	"encoding/json"
	"errors"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/auth"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
	"net/http"
)

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BalanceResponse struct {
	Current   uint `json:"current"`
	Withdrawn uint `json:"withdrawn"`
}

func (a *api) Register(writer http.ResponseWriter, request *http.Request) {
	var userRequest UserRequest

	if err := json.NewDecoder(request.Body).Decode(&userRequest); err != nil {
		logger.Log.Debug("Error while decode: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.s.Register(request.Context(), userRequest.Login, userRequest.Password)

	if err != nil {
		if errors.Is(err, service.ErrorLoginExists) {
			writer.WriteHeader(http.StatusConflict)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	generatedJWT, err := auth.BuildJWTString(userRequest.Login)

	if err != nil {
		logger.Log.Error("Error while create token: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Authorization", generatedJWT)
	writer.WriteHeader(http.StatusOK)
}

func (a *api) Login(writer http.ResponseWriter, request *http.Request) {
	var userRequest UserRequest

	if err := json.NewDecoder(request.Body).Decode(&userRequest); err != nil {
		logger.Log.Debug("Error while decode: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.s.Login(request.Context(), userRequest.Login, userRequest.Password)

	if err != nil {
		if errors.Is(err, service.ErrorWrongCredentials) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	generatedJWT, err := auth.BuildJWTString(userRequest.Login)

	if err != nil {
		logger.Log.Error("Error while create token: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Authorization", generatedJWT)
	writer.WriteHeader(http.StatusOK)
}

func (a *api) GetBalance(writer http.ResponseWriter, request *http.Request) {
	userId := auth.GetUserIDFromRequest(request)

	balance, err := a.s.GetBalance(request.Context(), userId)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceResponse := BalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(balanceResponse)
	writer.WriteHeader(http.StatusOK)
}
