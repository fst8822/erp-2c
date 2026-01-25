package service

import (
	"erp-2c/model"
)

type AuthService interface {
	SignUp(signUp model.SignUp) (*model.UserDomain, error)
	SignIn(signIn model.SignIn) (string, error)
}

type UserService interface {
	Save(userToSave model.SignUp) (*model.UserDomain, error)
	GetById(userId int64) (*model.UserDomain, error)
	GetByLogin(userId string) (*model.UserDomain, error)
}

type ProductService interface {
	Save(productToSave model.ProductToSave) (*model.ProductDomain, error)
	GetById(productId int64) (*model.ProductDomain, error)
	GetByName(productName string) (*model.ProductDomain, error)
	GetAll() (*[]model.ProductDomain, error)
	UpdateById(productId int64, productToUpdate model.ProductUpdate) error
	DeleteById(productId int64) error
}

type DeliveryService interface {
	Save(DeliveryToSave model.DeliveryToSave) (*model.DeliveryDomain, error)
	GetById(deliveryId int64) (*model.DeliveryDomain, error)
	GetAll() (*[]model.ProductDomain, error)
	GetByStatus(status string) (*model.DeliveryDomain, error)
	UpdateById(deliveryId int64, status model.UpdateStatus) error
	DeleteById(deliveryId int64) error
}
