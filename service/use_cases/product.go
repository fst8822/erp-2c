package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
	"fmt"
	"log/slog"
)

type ProductService struct {
	store *store.Store
}

func NewProductService(store *store.Store) *ProductService {
	return &ProductService{store: store}
}

func (p *ProductService) Save(productToSave model.ProductToSave) (*model.ProductDomain, error) {
	const op = "service.usescases.product.Save"

	productDB := model.ProductDB{
		ProductName:  productToSave.ProductName,
		ProductGroup: productToSave.ProductGroup,
		Image:        productToSave.Image,
		Stock:        productToSave.Stock,
		Price:        productToSave.Price,
	}
	saved, err := p.store.ProductRepo.Save(productDB)
	if err != nil {
		slog.Error("failed save product",
			slog.String("product id", productToSave.ProductName), sl.ErrWithOP(err, op))

		return nil, fmt.Errorf("failed to save product %w", err)
	}

	productDomain := model.ProductDomain{
		Id:           saved.Id,
		ProductName:  saved.ProductName,
		ProductGroup: saved.ProductGroup,
		Image:        saved.Image,
		Stock:        saved.Stock,
		Price:        saved.Price,
	}
	return &productDomain, nil
}

func (p *ProductService) GetById(productId int) (*model.ProductDomain, error) {
	return &model.ProductDomain{}, nil
}

func (p *ProductService) GetByName(productName string) (*model.ProductDomain, error) {
	return &model.ProductDomain{}, nil
}

func (p *ProductService) GetAll() (*[]model.ProductDomain, error) {
	return &[]model.ProductDomain{}, nil
}

func (p *ProductService) UpdateById(product int) error {
	return nil
}

func (p *ProductService) DeleteById(product int) error {
	return nil
}
