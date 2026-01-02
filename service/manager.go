package service

import (
	usercase "erp-2c/service/user-case"
	"erp-2c/store"
)

type Manager struct {
	UserService    UserService
	ProductService ProductService
}

func NewManager(store *store.Store) *Manager {
	return &Manager{
		UserService:    usercase.NewUserService(store),
		ProductService: usercase.NewProductService(store),
	}
}
