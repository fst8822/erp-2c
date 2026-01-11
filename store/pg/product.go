package pg

import (
	"database/sql"
	"erp-2c/model"
	"errors"
	"fmt"
	"strings"

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

func (p *ProductRepository) UpdateById(productId int64, productToUpdate model.ProductUpdate) error {
	params, fields := buildUpdateParams(productId, productToUpdate)
	if len(fields) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = :id", strings.Join(fields, ", "))
	res, err := p.db.NamedExec(query, params)
	if err != nil {
		return fmt.Errorf("failed to update product with id %d: %w", productId, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows count: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", productId)
	}
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

func buildUpdateParams(productId int64, productToUpdate model.ProductUpdate) (map[string]any, []string) {
	params := make(map[string]any)
	var setFields []string

	if productToUpdate.ProductName != nil {
		params["name"] = *productToUpdate.ProductName
		setFields = append(setFields, "product_name = :name")
	}
	if productToUpdate.ProductGroup != nil {
		params["pGroup"] = *productToUpdate.ProductGroup
		setFields = append(setFields, "product_group = :pGroup")
	}
	if productToUpdate.Image != nil {
		params["pImage"] = *productToUpdate.Image
		setFields = append(setFields, "image = :pImage")
	}
	if productToUpdate.Stock != nil {
		params["stock"] = *productToUpdate.Stock
		setFields = append(setFields, "stock = :stock")
	}
	if productToUpdate.Price != nil {
		params["price"] = *productToUpdate.Price
		setFields = append(setFields, "price = :price")
	}
	params["id"] = productId
	return params, setFields
}
