package service

import (
	usercase "erp-2c/service/use_cases"
	"erp-2c/store"
	"fmt"
)

type Manager struct {
	UserService    UserService
	ProductService ProductService
	AuthService    AuthService
}

func NewManager(storeRepo *store.Store) (*Manager, error) {
	if storeRepo == nil {
		return nil, fmt.Errorf("no store provided")
	}
	return &Manager{
		UserService:    usercase.NewUserService(storeRepo),
		ProductService: usercase.NewProductService(storeRepo),
		AuthService:    usercase.NewAuthService(storeRepo),
	}, nil
}
