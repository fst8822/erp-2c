package main

import (
	"api-gateway/config"
	"api-gateway/controller"
	"api-gateway/lib/sl"
	"api-gateway/server_grpc/interceptors_grpc"
	deliverygrpc "api-gateway/server_grpc/proto/v1/delivery"
	productgrpc "api-gateway/server_grpc/proto/v1/product"
	"api-gateway/service/use_cases"
	"api-gateway/store"
	"api-gateway/store/pg"
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	forceShutdownTime = 5 * time.Second
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

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors_grpc.TokenInterceptorUnary()),
		grpc.WithStreamInterceptor(interceptors_grpc.TokenInterceptorStream()),
	)
	if err != nil {
		slog.Error("Service api-gateway: Unable to connect grpc client", err)
	}
	defer conn.Close()

	clientDeliveryGRPC := deliverygrpc.NewDeliveryServiceClient(conn)
	clientProductGRPC := productgrpc.NewProductServiceClient(conn)

	userRepository := pg.NewUserRepository(db.Pg)
	userService := use_cases.NewUserService(userRepository)
	authService := use_cases.NewAuthService(userService)
	notifyService := use_cases.NewNotifyService()

	r := controller.NewRouters(
		authService,
		clientDeliveryGRPC,
		notifyService,
		clientProductGRPC,
		userService,
	)
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		slog.Info("Service api-gateway: Start server http", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Service api-gateway: Failed to start server", slog.Any("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		slog.Info("Shutdown signal is received")
	case <-ctx.Done():
		slog.Info("Context cancelled")
	}
	slog.Info("Starting graceful shutdown")

	ctxCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown failed, forcing close", sl.Err(err))
		if closeErr := srv.Close(); closeErr != nil {
			slog.Error("Forcing close failed", sl.Err(closeErr))
		}
	}

	slog.Info("Server shutdown gracefully")
}
