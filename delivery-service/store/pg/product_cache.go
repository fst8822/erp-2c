package pg

import (
	"delivery-service/model"

	"github.com/jmoiron/sqlx"
)

// ProductCache Save todo в разработке
type ProductCache struct {
	productRepo *ProductRepository
	cache       CacheInt
}

func NewProductCache(productRepo *ProductRepository, cache CacheInt) *ProductCache {
	return &ProductCache{productRepo: productRepo, cache: cache}
}

func (c *ProductCache) Save(productToSave model.ProductDB) (*model.ProductDB, error) {
	return nil, nil
}
func (c *ProductCache) GetById(productId int64) (*model.ProductDB, error) {
	return nil, nil
}
func (c *ProductCache) GetExistIds(tx *sqlx.Tx, productIds []int64) ([]int64, error) {
	return []int64{}, nil
}
func (c *ProductCache) GetByName(productName string) (*model.ProductDB, error) {
	return nil, nil
}
func (c *ProductCache) GetAll() ([]model.ProductDB, error) {
	return []model.ProductDB{}, nil
}
func (c *ProductCache) UpdateById(productId int64, productToUpdate model.ProductUpdate) error {
	return nil
}
func (c *ProductCache) DeleteById(productId int64) error {
	return nil
}
