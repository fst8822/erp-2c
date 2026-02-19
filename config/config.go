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
	HTTPAddress      string        `envconfig:"HTTP_ADDRESS" default:"localhost:8080"`
	ReadTimeout      time.Duration `envconfig:"READ_TIMEOUT" default:"4s"`
	WriteTimeout     time.Duration `envconfig:"WRITE_TIMEOUT" default:"4s"`
	IdleTimeout      time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
	ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
	DriverName       string        `envconfig:"DRIVER_NAME"`
	HostDB           string        `envconfig:"HOST_DB"`
	PortDB           string        `envconfig:"PORT_DB"`
	DBUser           string        `envconfig:"DB_USER"`
	DBPassword       string        `envconfig:"DB_PASSWORD" json:"-"`
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
