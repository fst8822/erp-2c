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
	SaveWithItems(tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error)
	Save(tx *sqlx.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error)
	GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (*model.DeliveryWithItemsDB, error)
	GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error)
	GetAllWithItemsByStatus(tx *sqlx.Tx, status model.DeliveryStatus) (*model.DeliverListDB, error)
	GetAllByStatus(tx *sqlx.Tx, status model.DeliveryStatus) ([]model.DeliveryDB, error)
	UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error
	DeleteById(tx *sqlx.Tx, deliveryId int64) error
	UpdateStatusById(tx *sqlx.Tx, id int64, status model.DeliveryStatus) error
	UpdateStatusByIds(tx *sqlx.Tx, groups map[model.DeliveryStatus][]int64) error
}
