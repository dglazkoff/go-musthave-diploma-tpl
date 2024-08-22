package service

import "github.com/dglazkoff/go-musthave-diploma-tpl/internal/storage"

// создавать userStorage или использовать storage.Gophermart ??
//type userStorage interface {
//
//}

type service struct {
	storage storage.Gophermart
}

func New(storage storage.Gophermart) *service {
	return &service{storage: storage}
}
