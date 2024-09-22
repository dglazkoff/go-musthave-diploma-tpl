package db

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
)

type dbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *dbStorage {
	return &dbStorage{db: db}
}

func (s *dbStorage) Bootstrap() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"login VARCHAR(250) PRIMARY KEY, " +
		"password VARCHAR(250) NOT NULL," +
		"balance DOUBLE PRECISION NOT NULL" +
		")")

	if err != nil {
		logger.Log.Error("Error while create table users: ", err)
		return err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS orders (" +
		"id VARCHAR(250) PRIMARY KEY, " +
		"user_id VARCHAR(250) NOT NULL, " +
		"status VARCHAR(250) NOT NULL, " +
		"uploaded_at VARCHAR(250) NOT NULL, " +
		"accrual DOUBLE PRECISION NOT NULL" +
		")")

	if err != nil {
		logger.Log.Error("Error while create table orders: ", err)
		return err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS withdrawals (" +
		"id VARCHAR(250) PRIMARY KEY, " +
		"user_id VARCHAR(250) NOT NULL, " +
		"processed_at VARCHAR(250) NOT NULL, " +
		"sum DOUBLE PRECISION NOT NULL" +
		")")

	if err != nil {
		logger.Log.Error("Error while create table withdrawals: ", err)
		return err
	}

	return nil
}

func (s *dbStorage) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, options)
}
