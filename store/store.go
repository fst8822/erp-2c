package store

import (
	"erp-2c/store/pg"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	UserRepo    UserRepository
	ProductRepo ProductRepository
}

func NewStore(db *sqlx.DB) (*Store, error) {
	if err := RunPgMigrations(db); err != nil {
		return nil, fmt.Errorf("runPgMigrations failed %w", err)
	}
	store := &Store{
		UserRepo:    pg.NewUserRepository(db),
		ProductRepo: pg.NewProductRepository(db),
	}
	return store, nil
}
