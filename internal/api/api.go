package api

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/service"
)

type apiService interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login string, password string) error
	GetBalance(ctx context.Context, userID string) (service.UserBalance, error)

	GetUserOrders(ctx context.Context, userID string) ([]models.Order, error)
	AddOrder(ctx context.Context, orderNumber string, userLogin string) error

	GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawals, error)
	CreateWithdrawal(ctx context.Context, orderNumber string, sum float64, userID string) error
}

type api struct {
	s apiService
}

func New(s apiService) *api {
	return &api{s}
}
