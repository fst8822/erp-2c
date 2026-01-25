package pg

import (
	"erp-2c/model"
	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func (d *DeliveryRepository) Save(DeliveryToSave model.DeliveryDBProductDB) (*model.DeliveryDBProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeliveryRepository) GetById(deliveryId int64) (*model.DeliveryDBProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeliveryRepository) GetAll() (*[]model.ProductDomain, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeliveryRepository) GetByStatus(status string) (*model.DeliveryDBProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeliveryRepository) UpdateById(deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}
