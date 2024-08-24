package router

import (
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/api"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/auth"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(store storage.Gophermart, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	s := service.New(store, cfg)
	newAPI := api.New(s)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", newAPI.Register)
		r.Post("/login", newAPI.Login)
		r.Post("/orders", auth.Auth(newAPI.AddOrder))
		r.Get("/orders", auth.Auth(newAPI.GetOrders))
		r.Get("/balance", auth.Auth(newAPI.GetBalance))
		r.Get("/withdrawals", auth.Auth(newAPI.GetWithdrawals))
		r.Post("/balance/withdraw", auth.Auth(newAPI.CreateWithdrawal))
	})

	return r
}
