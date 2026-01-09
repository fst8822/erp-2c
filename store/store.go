package store

import (
	"erp-2c/lib/sl"
	"erp-2c/store/pg"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	UserRepo    UserRepository
	ProductRepo ProductRepository
}

func NewStore(db *sqlx.DB) (*Store, error) {
	if err := RunPgMigrations(db); err != nil {
		slog.Error("failed to run migration script", sl.Err(err))
		return nil, fmt.Errorf("failed to run migration %w", err)
	}
	store := &Store{
		UserRepo:    pg.NewUserRepository(db),
		ProductRepo: pg.NewProductRepository(db),
	}
	return store, nil
}
