package main

import (
	"context"
	"delivery-service/cache"
	"delivery-service/config"
	_ "delivery-service/config"
	"delivery-service/lib/collection"
	"delivery-service/lib/observability/app_metrics"
	"delivery-service/lib/sl"
	"delivery-service/lib/workers"
	"delivery-service/server_grpc"
	"delivery-service/server_grpc/interceptors_grpc"
	deliverygrpc "delivery-service/server_grpc/proto/v1/delivery"
	clientNotifyrgrpc "delivery-service/server_grpc/proto/v1/notify"
	clientProductrgrpc "delivery-service/server_grpc/proto/v1/product"
	"delivery-service/service/use_cases"
	"delivery-service/store/pg"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	capacity          = 10
	countWorkers      = 5
	cron              = 5
	tTLCache          = time.Minute
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
	cfg := config.Get()

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	db, err := pg.Dial()
	if err != nil {
		slog.Error("Error connect to DB", err)
		return
	}
	defer db.Pg.Close()

	mapCache := cache.NewMapCache(ctx, tTLCache)

	deliveryRepository := pg.NewDeliveryRepository(db.Pg)
	deliveryCache := pg.NewDeliveryCacheCache(ctx, deliveryRepository, mapCache)

	go app_metrics.StartMetricsSync(ctx, deliveryRepository)

	wg := sync.WaitGroup{}
	les, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		slog.Error("Failed to lister", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors_grpc.AuthServerInterceptorUnary),
		grpc.StreamInterceptor(interceptors_grpc.AuthServerInterceptorStream),
	)

	go func() {
		slog.Info("delivery-service: Start server grpc", slog.String("tcp", "localhost:50051"))
		if err := s.Serve(les); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("delivery-service: Failed to grpc Serve", slog.Any("err", err.Error()))
		}
	}()

	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("delivery-service: Connection to grpc server failed", err)
	}
	defer conn.Close()

	clientNotifyGRPC := clientNotifyrgrpc.NewNotifyServiceClient(conn)
	clientProductGRPC := clientProductrgrpc.NewProductServiceClient(conn)
	deliveryService := use_cases.NewDeliveryService(deliveryCache, deliveryRepository, clientProductGRPC)
	deliveryGRPCServer := server_grpc.NewDeliveryGRPCServer(deliveryService)
	deliverygrpc.RegisterDeliveryServiceServer(s, deliveryGRPCServer)

	queue := collection.NewQueue(capacity)
	workPoll := workers.NewWorkerPool(
		clientNotifyGRPC,
		deliveryRepository,
		queue,
		countWorkers,
		cron)

	wg.Add(1)
	go func() {
		defer wg.Done()
		go workPoll.Run(ctx)
	}()

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
		slog.Info("delivery-service: Start server http", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("delivery-service: Failed to start server", slog.Any("error", err.Error()))
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
	slog.Info("Start wait group")
	wg.Wait()
	slog.Info("End wait group")
	slog.Info("Server shutdown gracefully")
}
