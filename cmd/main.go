package main

import (
	"erp-2c/config"
	"erp-2c/controller/router"
	"erp-2c/service"
	store "erp-2c/store"
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
	loadENV()

	cfg := config.Get()

	db, err := pg.Dial()
	if err != nil {
		log.Fatalf(err.Error())
	}

	storeRepo := store.NewStore(db.Pg)
	serviceManager, err := service.NewManager(storeRepo)
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(serviceManager)

	slog.Info("Start server", slog.String("address", cfg.HTTPAddress))
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
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
