package store

import (
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Save(productToSave model.ProductDB) (*model.ProductDB, error)
	GetById(productId int64) (*model.ProductDB, error)
	GetExistIds(tx *sqlx.Tx, productIds []int64) ([]int64, error)
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
	SaveDelivery(tx *sqlx.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error)
	SaveDeliveryProducts(tx *sqlx.Tx, deliveryProductsDB []model.DeliveryProductDB) error
	GetById(tx *sqlx.Tx, deliveryId int64) (*model.DeliveryDB, error)
	GetAll(tx *sqlx.Tx) (*[]model.ProductDomain, error)
	GetByStatus(tx *sqlx.Tx, status string) (*model.DeliveryDB, error)
	UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error
	DeleteById(tx *sqlx.Tx, deliveryId int64) error
}
