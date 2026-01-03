package main

import (
	"erp-2c/config"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadENV()
	cfg := config.Get()
	_ = cfg

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
