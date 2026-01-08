package store

import (
	"erp-2c/config"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func RunPgMigrations() error {
	cfg := config.Get()

	if cfg.PGMigrationsPath == "" {
		return errors.New("no cfg.PGMigrationsPath provided")
	}
	pgURL := fmt.Sprintf(
		"host=%s port=%s user=%s database=%s password=%s sslmode=%s",
		cfg.HostDB, cfg.PortDB, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.SSLMode)
	m, err := migrate.New(cfg.PGMigrationsPath, pgURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
