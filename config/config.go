package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPAddress      string        `envconfig:"HTTP_ADDRESS"`
	ReadTimeout      time.Duration `envconfig:"READ_TIMEOUT"`
	WriteTimeout     time.Duration `envconfig:"WRITE_TIMEOUT"`
	IdleTimeout      time.Duration `envconfig:"IDLE_TIMEOUT"`
	DriverName       string        `envconfig:"DRIVER_NAME"`
	HostDB           string        `envconfig:"HOST_DB"`
	PortDB           string        `envconfig:"PORT_DB"`
	DBUser           string        `envconfig:"DB_USER"`
	DBPassword       string        `envconfig:"DB_PASSWORD"`
	DBName           string        `envconfig:"DB_NAME"`
	SSLMode          string        `envconfig:"SSLMODE"`
	PGMigrationsPath string        `envconfig:"PG_MIGRATIONS_PATH"`
}

var (
	conf Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &conf)
		if err != nil {
			log.Fatalf(err.Error())
		}
		configBytes, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("Configuration", string(configBytes))

	})
	return &conf
}
