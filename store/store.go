package store

import (
	"context"
	"database/sql"
	"erp-2c/lib/sl"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	Db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Db: db,
	}
}

func (s *Store) BeginTxx(ctx context.Context) (*sqlx.Tx, error) {
	const op = "store.store.BeginTxx"

	tx, err := s.Db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		slog.Error("failed to begin transaction", sl.ErrWithOP(err, op))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
