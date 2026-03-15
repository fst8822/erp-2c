package main

import (
	"context"
	"erp-2c/config"
	"erp-2c/controller"
	"erp-2c/lib/sl"
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
	"time"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatal("No .env file found")
	}
	fmt.Printf("RUN APP: env=%s\n", env)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	cfg := config.Get()

	db, err := pg.Dial()
	if err != nil {
		slog.Error("Error connect to DB", err)
		return
	}
	defer db.Pg.Close()

	if err := store.RunPgMigrations(db.Pg); err != nil {
		slog.Error("Error run migrations DB", err)
		return
	}

	userRepository := pg.NewUserRepository(db.Pg)
	userService := use_cases.NewUserService(userRepository)
	authService := use_cases.NewAuthService(userService)
	notifyService := use_cases.NewNotifyService()

	r := controller.NewRouters(
		authService,
		deliveryService,
		notifyService,
		productService,
		userService,
	)
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
		slog.Info("Server error, initiating shutdown", sl.Err(err))
	case <-stop:
		slog.Info("Shutdown signal is received")
	case <-ctx.Done():
		slog.Info("Context cancelled")
	}

	slog.Info("Starting graceful shutdown")
	ctxCancel()

	slog.Info("Send signal done graceful is successful")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown failed, forcing close", sl.Err(err))
		if closeErr := srv.Close(); closeErr != nil {
			slog.Error("Forcing close failed", sl.Err(closeErr))
		}
	}

	notifyService.Shutdown()
	slog.Info("Server shutdown gracefully")
}
