package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"math"
)

var ErrorLoginExists = fmt.Errorf("login exists")
var ErrorWrongCredentials = fmt.Errorf("wrong pair login/password")
var ErrorNotEnoughBalance = fmt.Errorf("not enough balance")

type UserBalance struct {
	Current   float64
	Withdrawn float64
}

func getHashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hp := h.Sum(nil)
	return hex.EncodeToString(hp)
}

func (s *service) Register(ctx context.Context, login, password string) error {
	// нужна ли транзацкция которая возмет юзера и если такого нет то запишет нового ??
	// может ли за время работы функции другой поток создать юзера с таким же логином?
	_, err := s.storage.GetUserByLogin(ctx, login)

	if err == nil {
		return ErrorLoginExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("Error while get user by login: ", err)
		return err
	}

	err = s.storage.CreateUser(ctx, login, getHashPassword(password))

	if err != nil {
		logger.Log.Error("Error while create user: ", err)
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, login, password string) error {
	user, err := s.storage.GetUserByLogin(ctx, login)

	if err != nil {
		logger.Log.Error("Error while get user by login: ", err)
		return ErrorWrongCredentials
	}

	if user.Password != getHashPassword(password) {
		return ErrorWrongCredentials
	}

	return nil
}

func (s *service) GetBalance(ctx context.Context, userId string) (UserBalance, error) {
	user, err := s.storage.GetUserByLogin(ctx, userId)

	if err != nil {
		logger.Log.Error("Error while get user by login during get balance: ", err)
		return UserBalance{}, err
	}

	withdrawals, err := s.GetWithdrawals(ctx, userId)

	if err != nil {
		logger.Log.Error("Error while get withdrawals during get balance: ", err)
		return UserBalance{}, err
	}

	var sum float64

	for _, w := range withdrawals {
		sum += w.Sum
	}

	return UserBalance{
		Current:   user.Balance,
		Withdrawn: sum,
	}, nil
}

func (s *service) UpdateBalance(ctx context.Context, sum float64, userId string) error {
	logger.Log.Debug("Update balance: ", sum)
	user, err := s.storage.GetUserByLogin(ctx, userId)
	copyUser := user

	if err != nil {
		logger.Log.Error("Error while get user by login during update balance: ", err)
		return err
	}

	if sum < 0 && copyUser.Balance < math.Abs(sum) {
		logger.Log.Error("Not enough balance")
		return ErrorNotEnoughBalance
	}

	copyUser.Balance += sum
	_, err = s.storage.UpdateUser(ctx, copyUser)

	if err != nil {
		logger.Log.Error("Error while update user: ", err)
		return err
	}

	return nil
}
