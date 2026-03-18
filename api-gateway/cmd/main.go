package main

import (
	"api-gateway/config"
	"api-gateway/controller"
	"api-gateway/lib/sl"
	"api-gateway/server_grpc"
	deliverygrpc "api-gateway/server_grpc/proto/v1/delivery"
	notifygrpc "api-gateway/server_grpc/proto/v1/notify"
	productgrpc "api-gateway/server_grpc/proto/v1/product"
	"api-gateway/service/use_cases"
	"api-gateway/store"
	"api-gateway/store/pg"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
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
	forceShutdownTime = 5
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

	les, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		slog.Error("Failed to lister", err)
	}

	s := grpc.NewServer()
	notifygrpc.RegisterNotifyServiceServer(s, &server_grpc.NotifyGRPCServer{})

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("не удалось подключиться", err)
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
	serverHTTPErrCh := make(chan error, 1)
	serverGRPCErrCh := make(chan error, 1)
	go func() {
		slog.Info("Start server http", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", slog.Any("error", err.Error()))
			serverHTTPErrCh <- err
		}
		close(serverHTTPErrCh)
	}()
	go func() {
		slog.Info("Start server grpc", slog.String("tcp", "localhost:50052"))
		if err := s.Serve(les); !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("Failed to grpc Serve", slog.Any("error", err.Error()))
			serverGRPCErrCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err = <-serverHTTPErrCh:
		slog.Info("Server HTTP error, initiating shutdown", sl.Err(err))
	case err = <-serverGRPCErrCh:
		slog.Info("Server GRPC error, initiating shutdown", sl.Err(err))
	case <-stop:
		slog.Info("Shutdown signal is received")
	case <-ctx.Done():
		slog.Info("Context cancelled")
	}

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	go func() {
		select {
		case <-done:
			slog.Info("GRPC Server stopped gracefully")
		case <-time.After(forceShutdownTime):
			s.Stop()
			slog.Info("GRPC Server Close Force ")
		}
	}()

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
