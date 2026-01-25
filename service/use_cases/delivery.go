package use_cases

import (
	"erp-2c/model"
	"erp-2c/store"
)

type DeliveryService struct {
	repo *store.Store
}

func NewDeliveryService(repo *store.Store) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (d *DeliveryService) Save(DeliveryToSave model.DeliveryToSave) (*model.DeliveryDomain, error) {
	return nil, nil
}
func (d *DeliveryService) GetById(deliveryId int64) (*model.DeliveryDomain, error) {
	return nil, nil
}
func (d *DeliveryService) GetAll() (*[]model.ProductDomain, error) {
	return nil, nil

}
func (d *DeliveryService) GetByStatus(status string) (*model.DeliveryDomain, error) {
	return nil, nil
}
func (d *DeliveryService) UpdateById(deliveryId int64, status model.UpdateStatus) error {
	return nil
}
