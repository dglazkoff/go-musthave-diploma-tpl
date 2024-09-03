package router

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/api"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/auth"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(store storage.Gophermart, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	s := service.New(store, cfg)
	newAPI := api.New(s)

	err := s.UpdateOrdersAccrual(context.Background())

	if err != nil {
		logger.Log.Error("Error while update orders accrual: ", err)
	}

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", logger.Log.Request(newAPI.Register))
		r.Post("/login", logger.Log.Request(newAPI.Login))
		r.Post("/orders", logger.Log.Request(auth.Auth(newAPI.AddOrder)))
		r.Get("/orders", logger.Log.Request(auth.Auth(newAPI.GetOrders)))
		r.Get("/balance", logger.Log.Request(auth.Auth(newAPI.GetBalance)))
		r.Get("/withdrawals", logger.Log.Request(auth.Auth(newAPI.GetWithdrawals)))
		r.Post("/balance/withdraw", logger.Log.Request(auth.Auth(newAPI.CreateWithdrawal)))
	})

	return r
}
