package pg

import (
	"erp-2c/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	Pg *sqlx.DB
}

func Dial() (*DB, error) {
	const op = "store.pg.Dial"

	cfg := config.Get()
	pgURL, err := checkDBFieldsReturnPgUrl(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(cfg.DriverName, pgURL)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db %w %s", err, op)
	}
	return &DB{db}, nil
}

func checkDBFieldsReturnPgUrl(cfg *config.Config) (string, error) {
	const op = "store.pg.checkDBFieldsReturnPgUrl"

	if cfg.DriverName == "" {
		return "", fmt.Errorf("wrong config field %s, driver name is empty, op=%s", cfg.DriverName, op)
	}
	if cfg.DBName == "" {
		return "", fmt.Errorf("wrong config field %s, driver name is empty, op=%s", cfg.DBName, op)
	}
	if cfg.HostDB == "" {
		return "", fmt.Errorf("wrong config field %s, db name is empty, op=%s", cfg.HostDB, op)
	}
	if cfg.PortDB == "" {
		return "", fmt.Errorf("wrong config field %s, host is empty, op=%s", cfg.PortDB, op)
	}
	if cfg.DBUser == "" {
		return "", fmt.Errorf("wrong config field %s,user name is empty, op=%s", cfg.DBUser, op)
	}
	if cfg.DBPassword == "" {
		return "", fmt.Errorf("wrong config field %s, password name is empty, op=%s", cfg.DBPassword, op)
	}
	if cfg.SSLMode == "" {
		return "", fmt.Errorf("wrong config field %s, sslmode name is empty, op=%s", cfg.SSLMode, op)
	}
	pgURL := fmt.Sprintf(
		"host=%s port=%s user=%s database=%s password=%s sslmode=%s",
		cfg.HostDB, cfg.PortDB, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.SSLMode)
	return pgURL, nil
}
