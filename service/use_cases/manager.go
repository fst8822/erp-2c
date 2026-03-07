package use_cases

import (
	"erp-2c/cache"
	"erp-2c/service"
	"erp-2c/store/pg"
)

type Manager struct {
	UserService     service.UserService
	ProductService  service.ProductService
	AuthService     service.AuthService
	DeliveryService service.DeliveryService
	NotifyService   service.NotifyService
}

func NewManager(
	userRepo userRepositoryInt,
	productRepo productRepoInt,
	deliveryRepo *pg.DeliveryRepository,
	cacheRepo *cache.RepositoryCache) (*Manager, error) {

	userService := NewUserService(userRepo)
	productService := NewProductService(productRepo)
	authService := NewAuthService(userService)
	deliveryService := NewDeliveryService(deliveryRepo, productRepo, cacheRepo)
	notifyService := NewNotifyService()

	return &Manager{
		UserService:     userService,
		ProductService:  productService,
		AuthService:     authService,
		DeliveryService: deliveryService,
		NotifyService:   notifyService,
	}, nil
}
