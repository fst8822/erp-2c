package use_cases

import (
	"erp-2c/cache"
	"erp-2c/service"
	"erp-2c/store"
	"fmt"
)

type Manager struct {
	UserService     service.UserService
	ProductService  service.ProductService
	AuthService     service.AuthService
	DeliveryService service.DeliveryService
	NotifyService   service.NotifyService
}

func NewManager(
	productRepo productRepoInt,
	storeRepo *store.Store,
	cacheRepo *cache.RepositoryCache,
) (*Manager, error) {
	if storeRepo == nil {
		return nil, fmt.Errorf("no store provided")
	}

	userService := NewUserService(storeRepo)
	productService := NewProductService(productRepo)
	authService := NewAuthService(storeRepo, userService)
	deliveryService := NewDeliveryService(storeRepo, productRepo, cacheRepo)
	notifyService := NewNotifyService()

	return &Manager{
		UserService:     userService,
		ProductService:  productService,
		AuthService:     authService,
		DeliveryService: deliveryService,
		NotifyService:   notifyService,
	}, nil
}
