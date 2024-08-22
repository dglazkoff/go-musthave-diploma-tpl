package server

import (
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/router"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage/db"
	"net/http"
)

func Run(cfg *config.Config) error {
	pgDB, err := sql.Open("pgx", cfg.DatabaseURI)

	if err != nil {
		logger.Log.Debug("Error on open db", "err", err)
		panic(err)
	}
	defer pgDB.Close()

	store := db.New(pgDB)
	err = store.Bootstrap()

	if err != nil {
		logger.Log.Debug("Error on bootstrap db ", err)
		panic(err)
	}

	logger.Log.Infow("Starting Server on ", "addr", cfg.RunAddr)

	return http.ListenAndServe(cfg.RunAddr, router.Router(store, cfg))
}
