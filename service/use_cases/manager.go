package use_cases

import (
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
	MetricsService  service.MetricsService
}

func NewManager(storeRepo *store.Store) (*Manager, error) {
	if storeRepo == nil {
		return nil, fmt.Errorf("no store provided")
	}

	userService := NewUserService(storeRepo)
	productService := NewProductService(storeRepo)
	authService := NewAuthService(storeRepo, userService)
	deliveryService := NewDeliveryService(storeRepo)
	notifyService := NewNotifyService()
	metricsService := NewMetricsService(storeRepo)

	return &Manager{
		UserService:     userService,
		ProductService:  productService,
		AuthService:     authService,
		DeliveryService: deliveryService,
		NotifyService:   notifyService,
		MetricsService:  metricsService,
	}, nil
}
