package service

import (
	"erp-2c/model"
)

type AuthService interface {
	SignUp(UserToSave model.User) (*model.User, error)
	SignIn(login string, password string) (string, error)
}

type UserService interface {
	Save(userToSave model.User) (*model.User, error)
	GetById(userId int) (*model.User, error)
	GetByName(userName string) (*model.User, error)
}
type ProductService interface {
	Save(productToSave model.Product) (*model.Product, error)
	GetById(productId int) (*model.Product, error)
	GetByName(productName string) (*model.Product, error)
	GetAll() (*[]model.Product, error)
	UpdateById(productId int) error
	DeleteById(productId int) error
}
