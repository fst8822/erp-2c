package use_cases

import "erp-2c/store"

type ProductService struct {
	store *store.Store
}

func NewProductService(store *store.Store) *ProductService {
	return &ProductService{store: store}
}

func (p ProductService) Save() {
	//TODO implement me
	panic("implement me")
}

func (p ProductService) GetById() {
	//TODO implement me
	panic("implement me")
}

func (p ProductService) GetAll() {
	//TODO implement me
	panic("implement me")
}

func (p ProductService) UpdateById() {
	//TODO implement me
	panic("implement me")
}

func (p ProductService) DeleteById() {
	//TODO implement me
	panic("implement me")
}
