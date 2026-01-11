package pg

import (
	"database/sql"
	"erp-2c/model"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) Save(productToSave model.ProductDB) (*model.ProductDB, error) {
	query := `INSERT INTO products (product_name, product_group, image,  stock, price) 
			  VALUES (:product_name, :product_group, :image,  :stock, :price) RETURNING *`

	rows, err := p.db.NamedQuery(query, productToSave)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}
	defer rows.Close()

	rows.Next()
	productDB := &model.ProductDB{}
	err = rows.StructScan(productDB)
	if err != nil {
		return nil, fmt.Errorf("failed to scan product %w", err)
	}
	return productDB, nil
}

func (p *ProductRepository) GetById(productId int64) (*model.ProductDB, error) {
	productDB := &model.ProductDB{}

	query := `SELECT * FROM products WHERE id = $1`
	if err := p.db.Get(productDB, query, productId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with id %d not found", productId)
		}
		return nil, fmt.Errorf("failed to get product %w", err)
	}
	return productDB, nil
}

func (p *ProductRepository) GetByName(productName string) (*model.ProductDB, error) {
	productDB := &model.ProductDB{}

	query := `SELECT * FROM products WHERE product_name = $1`
	if err := p.db.Get(productDB, query, productName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with name %s not found", productName)
		}
		return nil, fmt.Errorf("failed ot get product %w", err)
	}
	return productDB, nil
}

func (p *ProductRepository) GetAll() ([]model.ProductDB, error) {
	var products []model.ProductDB

	query := ` SELECT * FROM products`

	if err := p.db.Select(&products, query); err != nil {
		return nil, fmt.Errorf("failed to get list product %w", err)
	}
	return products, nil
}

func (p *ProductRepository) UpdateById(productId int64, productToUpdate model.ProductDB) error {
	return nil
}

func (p *ProductRepository) DeleteById(productId int64) error {
	query := `DELETE FROM products WHERE id = $1 RETURNING id`

	var deletedId int
	err := p.db.QueryRow(query, productId).Scan(&deletedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no find to delete product with id %d", productId)
		}
		return fmt.Errorf("failed to delete product %w", err)
	}
	return nil
}

func (p *ProductRepository) GetByGroupName(groupName string) ([]model.ProductDB, error) {
	var products []model.ProductDB

	query := `SELECT * FROM products WHERE product_group = $1`
	if err := p.db.Select(&products, query, groupName); err != nil {
		return nil, fmt.Errorf("failed to list product from group %w", err)
	}
	return products, nil
}
