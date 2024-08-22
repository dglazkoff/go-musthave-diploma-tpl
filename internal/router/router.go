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

	s := service.New(store)
	newAPI := api.New(s)

	r.Post("/user/register", newAPI.Register)
	r.Post("/user/login", newAPI.Login)
	r.Post("/user/orders", auth.Auth(newAPI.AddOrder))
	r.Get("/user/orders", auth.Auth(newAPI.GetOrders))
	r.Get("/user/balance", auth.Auth(newAPI.GetBalance))
	r.Get("/user/withdrawals", auth.Auth(newAPI.GetWithdrawals))
	r.Post("/user/balance/withdraw", auth.Auth(newAPI.CreateWithdrawal))

	return r
}
