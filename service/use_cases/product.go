package use_cases

import (
	"erp-2c/model"
	"erp-2c/store"
)

type ProductService struct {
	store *store.Store
}

func NewProductService(store *store.Store) *ProductService {
	return &ProductService{store: store}
}

func (p *ProductService) Save(productToSave model.Product) (model.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) GetById(productId int) (model.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) GetByName(productName string) (model.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) GetAll() ([]model.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) UpdateById(product int) error {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) DeleteById(product int) error {
	//TODO implement me
	panic("implement me")
}
