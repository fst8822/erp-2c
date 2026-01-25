package pg

import (
	"database/sql"
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (d *DeliveryRepository) SaveDelivery(tx *sql.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetById(tx *sql.Tx, deliveryId int64) (*model.DeliveryDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetAll() (*[]model.ProductDomain, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetByStatus(tx *sql.Tx, status string) (*model.DeliveryDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) UpdateById(tx *sql.Tx, deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func (d *DeliveryRepository) DeleteById(tx *sql.Tx, deliveryId int64) error { return nil }

func (d *DeliveryRepository) SaveDeliveryProducts(tx *sql.Tx, deliveryProductsDB []model.DeliveryProductDB) error {
	return nil
}
