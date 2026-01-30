package main

import (
	"context"
	"erp-2c/config"
	"erp-2c/controller/routers"
	"erp-2c/lib/collection"
	"erp-2c/lib/sl"
	workers "erp-2c/lib/workers"
	"erp-2c/service/use_cases"
	"erp-2c/store"
	"erp-2c/store/pg"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
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

func run(ctx context.Context) error {

	loadENV()

	cfg := config.Get()

	db, err := pg.Dial()
	if err != nil {
		return err
	}
	defer db.Pg.Close()

	if err := store.RunPgMigrations(db.Pg); err != nil {
		return err
	}

	storeRepo := store.NewStore(db.Pg)

	serviceManager, err := use_cases.NewManager(storeRepo)
	if err != nil {
		return err
	}
	queue := collection.NewQueue(10)
	workPoll := workers.NewWorkerPoolQueue(storeRepo, queue, 5, 1)
	workPoll.Run(ctx)
	r := routers.New(serviceManager)

	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	serverErrCh := make(chan error)
	go func() {
		slog.Info("Start server", slog.String("address", cfg.HTTPAddress))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", slog.String("error", err.Error()))
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err = <-serverErrCh:
		return err
	case <-stop:
		slog.Info("Shutdown signal is received")
	case <-ctx.Done():
		slog.Info("Context cancelled")
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown failed, forcing close", sl.Err(err))
		if closeErr := srv.Close(); closeErr != nil {
			slog.Error("Forcing close failed", sl.Err(closeErr))
			return errors.Join(err, closeErr)
		}
		return err
	}
	slog.Info("Server shutdown gracefully")
	return nil
}
