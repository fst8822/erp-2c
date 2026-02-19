package store

import (
	"context"
	"database/sql"
	"erp-2c/lib/sl"
	"erp-2c/store/pg"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db          *sqlx.DB
	UserRepo    UserRepository
	ProductRepo ProductRepository
	Delivery    DeliveryRepository
}

func NewStore(db *sqlx.DB) *Store {

	return &Store{
		UserRepo:    pg.NewUserRepository(db),
		ProductRepo: pg.NewProductRepository(db),
		Delivery:    pg.NewDeliveryRepository(db),
		db:          db,
	}
}

func (s *Store) BeginTxx(ctx context.Context) (*sqlx.Tx, error) {
	const op = "store.store.BeginTxx"

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		slog.Error("failed to begin transaction", sl.ErrWithOP(err, op))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
