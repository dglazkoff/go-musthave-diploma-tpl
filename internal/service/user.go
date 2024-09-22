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
	tx, err := s.storage.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.Error("Error while begin transaction: ", err)
		return err
	}

	defer tx.Rollback()

	// получается мы блокируем запись в таблицу на время выполнения регистрации
	// не уверен, что это хорошо, но не вижу других вариантов, если надо предотвращать регистрацию параллельным потоком
	_, err = s.storage.GetUserByLoginForUpdate(ctx, tx, login)

	if err == nil {
		return ErrorLoginExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("Error while get user by login: ", err)
		return err
	}

	err = s.storage.CreateUser(ctx, tx, login, getHashPassword(password))

	if err != nil {
		logger.Log.Error("Error while create user: ", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("Error while commit transaction: ", err)
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

func (s *service) GetBalance(ctx context.Context, userID string) (UserBalance, error) {
	// единственное чем тут транзакция полезна, что может придти withdrawal уже после того, как бы вычитали баланс
	// а поможет ли нам такой уровень изоляции? ведь это другая таблица и снепшота ее не будет
	tx, err := s.storage.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	if err != nil {
		logger.Log.Error("Error while begin transaction: ", err)
		return UserBalance{}, err
	}

	defer tx.Rollback()

	user, err := s.storage.GetUserByLoginTx(ctx, tx, userID)

	if err != nil {
		logger.Log.Error("Error while get user by login during get balance: ", err)
		return UserBalance{}, err
	}

	withdrawals, err := s.storage.GetWithdrawalsTx(ctx, tx, userID)

	if err != nil {
		logger.Log.Error("Error while get withdrawals during get balance: ", err)
		return UserBalance{}, err
	}

	var sum float64

	for _, w := range withdrawals {
		sum += w.Sum
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error("Error while commit transaction: ", err)
		return UserBalance{}, err
	}

	return UserBalance{
		Current:   user.Balance,
		Withdrawn: sum,
	}, nil
}

func (s *service) UpdateBalanceTx(ctx context.Context, tx *sql.Tx, sum float64, userID string) error {
	logger.Log.Debug("Update balance: ", sum)
	user, err := s.storage.GetUserByLoginForUpdate(ctx, tx, userID)

	if err != nil {
		logger.Log.Error("Error while get user by login during update balance: ", err)
		return err
	}

	if sum < 0 && user.Balance < math.Abs(sum) {
		logger.Log.Error("Not enough balance")
		return ErrorNotEnoughBalance
	}

	user.Balance += sum
	_, err = s.storage.UpdateUser(ctx, tx, user)

	if err != nil {
		logger.Log.Error("Error while update user: ", err)
		return err
	}

	return nil
}

func (s *service) UpdateBalance(ctx context.Context, sum float64, userID string) error {
	tx, err := s.storage.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.Error("Error while begin transaction: ", err)
		return err
	}

	defer tx.Rollback()

	err = s.UpdateBalanceTx(ctx, tx, sum, userID)

	if err == nil {
		if err = tx.Commit(); err != nil {
			logger.Log.Error("Error while commit transaction: ", err)
			return err
		}
	}

	return nil
}
