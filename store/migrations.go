package store

import (
	"erp-2c/config"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func RunPgMigrations(db *sqlx.DB) error {
	slog.Info("Running PgMigrations")
	const op = "store.RunPgMigrations"

	cfg := config.Get()

	if cfg.PGMigrationsPath == "" {
		return errors.New("no cfg.PGMigrationsPath provided")
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres instance failed %w, op = %s", err, op)
	}

	m, err := migrate.NewWithDatabaseInstance(
		cfg.PGMigrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("postgres instance failed %w, op = %s", err, op)
	}
	return nil
}
