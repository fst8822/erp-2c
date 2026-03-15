package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	DriverName       string `envconfig:"DRIVER_NAME"`
	HostDB           string `envconfig:"HOST_DB"`
	PortDB           string `envconfig:"PORT_DB"`
	DBUser           string `envconfig:"DB_USER"`
	DBPassword       string `envconfig:"DB_PASSWORD" json:"-"`
	DBName           string `envconfig:"DB_NAME"`
	SSLMode          string `envconfig:"SSLMODE"`
	PGMigrationsPath string `envconfig:"PG_MIGRATIONS_PATH"`
}

var (
	conf DBConfig
	once sync.Once
)

func Get() *DBConfig {
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
