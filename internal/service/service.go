package service

import (
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage"
)

type service struct {
	storage storage.Gophermart
	cfg     *config.Config
}

func New(storage storage.Gophermart, cfg *config.Config) *service {
	return &service{storage: storage, cfg: cfg}
}
