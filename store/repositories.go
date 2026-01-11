package store

import "erp-2c/model"

type ProductRepository interface {
	Save(productToSave model.ProductDB) (*model.ProductDB, error)
	GetById(productId int64) (*model.ProductDB, error)
	GetByName(productName string) (*model.ProductDB, error)
	GetAll() ([]model.ProductDB, error)
	UpdateById(productId int64, productToUpdate model.ProductUpdate) error
	DeleteById(productId int64) error
	GetByGroupName(groupId string) ([]model.ProductDB, error)
}

type UserRepository interface {
	Save(userToSave model.UserDB) (*model.UserDB, error)
	GetById(userId int64) (*model.UserDB, error)
	GetByLogin(userId string) (*model.UserDB, error)
}
