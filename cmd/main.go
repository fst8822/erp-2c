package main

import (
	"erp-2c/config"
	"erp-2c/controller/routers"
	"erp-2c/service/use_cases"
	"erp-2c/store"
	"erp-2c/store/pg"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
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

func run() error {
	loadENV()

	cfg := config.Get()

	db, err := pg.Dial()
	if err != nil {
		return err
	}

	storeRepo, err := store.NewStore(db.Pg)
	if err != nil {
		return err
	}

	serviceManager, err := use_cases.NewManager(storeRepo)
	if err != nil {
		return err
	}

	r := routers.New(serviceManager)

	slog.Info("Start server", slog.String("address", cfg.HTTPAddress))
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	//todo add method for grace full shutdown
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		return err
	}
	return nil
}
