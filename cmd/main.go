package main

import (
	"context"
	"erp-2c/cache"
	"erp-2c/config"
	"erp-2c/controller"
	"erp-2c/lib/collection"
	"erp-2c/lib/observability/app_metrics"
	"erp-2c/lib/sl"
	"erp-2c/lib/workers"
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
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	capacity     = 10
	countWorkers = 5
	cron         = 5
	tTLCache     = time.Minute
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
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

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

	productRepository := pg.NewProductRepository(db.Pg)
	userRepository := pg.NewUserRepository(db.Pg)
	deliveryRepository := pg.NewDeliveryRepository(db.Pg)

	mapCache := cache.NewMapCache(tTLCache)
	repositoryCache := pg.NewRepositoryCache(ctx, deliveryRepository, &mapCache)

	userService := use_cases.NewUserService(userRepository)
	productService := use_cases.NewProductService(productRepository)
	authService := use_cases.NewAuthService(userService)
	deliveryService := use_cases.NewDeliveryService(repositoryCache, productRepository)
	notifyService := use_cases.NewNotifyService()

	if err != nil {
		return err
	}

	go app_metrics.StartMetricsSync(ctx, deliveryRepository)

	queue := collection.NewQueue(capacity)
	workPoll := workers.NewWorkerPool(
		notifyService,
		deliveryRepository,
		queue,
		countWorkers,
		cron)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go workPoll.Run(ctx)
	}()

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

	wg.Wait()
	notifyService.Shutdown()
	slog.Info("Server shutdown gracefully")
	return nil
}
