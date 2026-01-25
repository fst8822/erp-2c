package pg

import (
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (d *DeliveryRepository) Save(DeliveryToSave model.DeliveryProductDB) (*model.DeliveryProductDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetById(deliveryId int64) (*model.DeliveryProductDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetAll() (*[]model.ProductDomain, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetByStatus(status string) (*model.DeliveryProductDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) UpdateById(deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func (d *DeliveryRepository) DeleteById(deliveryId int64) error { return nil }
