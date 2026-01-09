package store

import (
	"erp-2c/store/pg"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	UserRepo    UserRepository
	ProductRepo ProductRepository
}

func NewStore(db *sqlx.DB) (*Store, error) {
	store := &Store{
		UserRepo:    pg.NewUserRepository(db),
		ProductRepo: pg.NewProductRepository(db),
	}
	return store, nil
}
