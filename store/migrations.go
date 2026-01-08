package store

import (
	"erp-2c/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"os"
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
	cwd, _ := os.Getwd()
	log.Println("cwd:", cwd)
	m, err := migrate.NewWithDatabaseInstance(
		"file://./store/pg/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
