package db

import (
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
		"balance INT NOT NULL" +
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
		"accrual INT NOT NULL" +
		")")

	if err != nil {
		logger.Log.Error("Error while create table orders: ", err)
		return err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS withdrawals (" +
		"order VARCHAR(250) PRIMARY KEY, " +
		"user_id VARCHAR(250) NOT NULL, " +
		"sum INT NOT NULL, " +
		"processed_at VARCHAR(250) NOT NULL" +
		")")

	if err != nil {
		logger.Log.Error("Error while create table withdrawals: ", err)
		return err
	}

	return nil
}
