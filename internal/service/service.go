package service

import (
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage"
)

type accrualService interface {
	GetAccrual(order models.Order, userLogin string)
}

type service struct {
	storage        storage.Gophermart
	cfg            *config.Config
	accrualService accrualService
}

func New(storage storage.Gophermart, cfg *config.Config) *service {
	return &service{storage: storage, cfg: cfg}
}

func (s *service) SetAccrualService(accrualService accrualService) {
	s.accrualService = accrualService
}
