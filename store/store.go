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
	}
}

func (s Store) BeginTx(ctx context.Context) (*sql.Tx, error) {
	const op = "store.store.BeginTx"

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		slog.Error("failed save delivery", sl.ErrWithOP(err, op))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
