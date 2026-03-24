package main

import (
	"notification-service/config"
	"notification-service/lib/sl"
	"notification-service/server_grpc/interceptors_grpc"
	notifygrpc "notification-service/server_grpc/proto/v1/notify"
	"notification-service/service/use_cases"

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

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	les, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		slog.Error("notification-service: Failed to lister", err)
		return
	}

	s := grpc.NewServer()
	notifygrpc.RegisterNotifyServiceServer(s, &use_cases.NotifyService{})

	go func() {
		slog.Info("notification-service: Start server grpc", slog.String("tcp", "localhost:50052"))
		if err := s.Serve(les); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("Failed to grpc Serve", slog.Any("error", err.Error()))
		}
	}()
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors_grpc.TokenInterceptorUnary()),
		grpc.WithStreamInterceptor(interceptors_grpc.TokenInterceptorStream()),
	)
	if err != nil {
		slog.Error("notification-service: Unable to connect grpc client", err)
	}
	defer conn.Close()

	notifyService := use_cases.NewNotifyService()

	route := chi.NewRouter()
	route.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      route,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		slog.Info("notification-service: Start server http", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("notification-service: Failed to start server", slog.Any("error", err.Error()))
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

	done := make(chan struct{})
	go func() {
		slog.Info("Start Server GRPC shutdown")
		s.GracefulStop()
		close(done)
		slog.Info("End Server GRPC shutdown")
	}()

	select {
	case <-done:
		slog.Info("GRPC Server stopped gracefully")
	case <-time.After(forceShutdownTime):
		s.Stop()
		slog.Info("GRPC Server Close Force ")
	}
	ctxCancel()

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
