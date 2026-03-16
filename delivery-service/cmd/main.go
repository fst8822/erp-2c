package main

import (
	"context"
	"delivery-service/cache"
	_ "delivery-service/config"
	"delivery-service/lib/collection"
	"delivery-service/lib/observability/app_metrics"
	"delivery-service/lib/workers"
	"delivery-service/service/use_cases"
	"delivery-service/store/pg"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

const (
	capacity     = 10
	countWorkers = 5
	cron         = 5
	tTLCache     = time.Minute
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

	db, err := pg.Dial()
	if err != nil {
		slog.Error("Error connect to DB", err)
		return
	}
	defer db.Pg.Close()

	productRepository := pg.NewProductRepository(db.Pg)
	productService := use_cases.NewProductService(productRepository)
	_ = productService

	deliveryRepository := pg.NewDeliveryRepository(db.Pg)
	mapCache := cache.NewMapCache(ctx, tTLCache)
	deliveryCache := pg.NewDeliveryCacheCache(ctx, deliveryRepository, mapCache)
	deliveryService := use_cases.NewDeliveryService(deliveryCache, deliveryRepository, productRepository)
	_ = deliveryService

	go app_metrics.StartMetricsSync(ctx, deliveryRepository)
	queue := collection.NewQueue(capacity)
	workPoll := workers.NewWorkerPool(
		nil, //notifyService,
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
	wg.Wait()

	//todo add graceful shutdown
}
