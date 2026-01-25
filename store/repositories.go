package store

import (
	"database/sql"
	"erp-2c/model"
)

type ProductRepository interface {
	Save(productToSave model.ProductDB) (*model.ProductDB, error)
	GetById(productId int64) (*model.ProductDB, error)
	GetExistIds(productIds []int64) ([]int64, error)
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

type DeliveryRepository interface {
	SaveDelivery(tx *sql.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error)
	SaveDeliveryProducts(tx *sql.Tx, deliveryProductsDB []model.DeliveryProductDB) error
	GetById(tx *sql.Tx, deliveryId int64) (*model.DeliveryDB, error)
	GetAll() (*[]model.ProductDomain, error)
	GetByStatus(tx *sql.Tx, status string) (*model.DeliveryDB, error)
	UpdateById(tx *sql.Tx, deliveryId int64, status model.UpdateStatus) error
	DeleteById(tx *sql.Tx, deliveryId int64) error
}
