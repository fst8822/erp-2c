package pg

import (
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p ProductRepository) Save(productToSave model.ProductDB) (model.ProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) GetById(productId int) (model.ProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) GetByName(productName string) (model.ProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) GetAll() ([]model.ProductDB, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) UpdateById(productId int) error {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) DeleteById(productId int) error {
	//TODO implement me
	panic("implement me")
}
