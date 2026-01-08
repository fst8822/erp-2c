package store

import "erp-2c/model"

type ProductRepository interface {
	Save(productToSave model.ProductDB) (*model.ProductDB, error)
	GetById(productId int) (*model.ProductDB, error)
	GetByName(productName string) (*model.ProductDB, error)
	GetAll() (*[]model.ProductDB, error)
	UpdateById(productId int) error
	DeleteById(productId int) error
}

type UserRepository interface {
	Save(userToSave model.UserDB) (*model.UserDB, error)
	GetById(userId int) (*model.UserDB, error)
	GetByLogin(login string) (*model.UserDB, error)
}
