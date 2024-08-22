package main

import (
	server "github.com/dglazkoff/go-musthave-diploma-tpl/internal"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.ParseConfig()
	err := logger.Initialize()

	if err != nil {
		panic(err)
	}

	if err := server.Run(&cfg); err != nil {
		logger.Log.Debug("Error on starting server", "err", err)
	}
}
