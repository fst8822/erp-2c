package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
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
		return nil, err
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

func (p *ProductService) GetById(productId int64) (*model.ProductDomain, error) {
	const op = "service.usescases.product.GetById"

	product, err := p.store.ProductRepo.GetById(productId)
	if err != nil {
		slog.Error("failed to find product", sl.ErrWithOP(err, op))
		return nil, err
	}

	productDB := model.ProductDomain{
		Id:           product.Id,
		ProductName:  product.ProductName,
		ProductGroup: product.ProductGroup,
		Image:        product.Image,
		Stock:        product.Stock,
		Price:        product.Price,
	}
	return &productDB, nil
}

func (p *ProductService) GetByName(productName string) (*model.ProductDomain, error) {
	const op = "service.usescases.product.GetByName"

	product, err := p.store.ProductRepo.GetByName(productName)
	if err != nil {
		slog.Error("failed to find user", sl.ErrWithOP(err, op))
		return nil, err
	}

	productDB := model.ProductDomain{
		Id:           product.Id,
		ProductName:  product.ProductName,
		ProductGroup: product.ProductGroup,
		Image:        product.Image,
		Stock:        product.Stock,
		Price:        product.Price,
	}
	return &productDB, nil
}

func (p *ProductService) GetAll() (*[]model.ProductDomain, error) {
	const op = "service.usescases.product.GetAll"

	productsDB, err := p.store.ProductRepo.GetAll()
	if err != nil {
		slog.Error("failed to find products", sl.ErrWithOP(err, op))
		return nil, err
	}

	var productsDomain []model.ProductDomain

	for _, productDB := range productsDB {
		productDomain := model.ProductDomain{
			Id:           productDB.Id,
			ProductName:  productDB.ProductName,
			ProductGroup: productDB.ProductGroup,
			Image:        productDB.Image,
			Stock:        productDB.Stock,
			Price:        productDB.Price,
		}
		productsDomain = append(productsDomain, productDomain)
	}
	return &productsDomain, nil
}

func (p *ProductService) UpdateById(productId int64, productToUpdate model.ProductUpdate) error {
	const op = "service.usescases.product.UpdateById"

	if err := p.store.ProductRepo.UpdateById(productId, productToUpdate); err != nil {
		slog.Error("failed to update product", sl.ErrWithOP(err, op))
		return err
	}
	return nil
}

func (p *ProductService) DeleteById(productId int64) error {
	const op = "service.usescases.product.DeleteById"

	if err := p.store.ProductRepo.DeleteById(productId); err != nil {
		slog.Error("failed to delete product", sl.ErrWithOP(err, op))
		return err
	}
	return nil
}
