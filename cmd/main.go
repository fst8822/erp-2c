package main

import (
	"erp-2c/config"
	"erp-2c/store/pg"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadENV()
	cfg := config.Get()
	_ = cfg
	db, err := pg.Dial()
	if err != nil {
		log.Fatalf(err.Error())
	}
	_ = db
}

func loadENV() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatal("No .env file found")
	}
	fmt.Printf("RUN APP: env=%s\n", env)
}
