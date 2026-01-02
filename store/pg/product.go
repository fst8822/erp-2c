package pg

import "github.com/jmoiron/sqlx"

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p ProductRepository) Save() {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) GetById() {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) GetAll() {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) UpdateById() {
	//TODO implement me
	panic("implement me")
}

func (p ProductRepository) DeleteById() {
	//TODO implement me
	panic("implement me")
}
