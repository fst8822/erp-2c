package workers

import (
	"context"
	"erp-2c/lib/collection"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
	"log/slog"
	"sync"
	"time"
)

type Worker interface {
	Run(ctx context.Context)
}
type WorkerPoolQueue struct {
	store        *store.Store
	queue        *collection.Queue
	wg           *sync.WaitGroup
	countWorkers int
	cron         time.Duration
}

func NewWorkerPool(
	store *store.Store,
	queue *collection.Queue,
	countWorkers int,
	sec time.Duration) *WorkerPoolQueue {
	return &WorkerPoolQueue{
		store:        store,
		queue:        queue,
		wg:           &sync.WaitGroup{},
		countWorkers: countWorkers,
		cron:         sec}
}

func (w *WorkerPoolQueue) Run(ctx context.Context) {
	const op = "lib.workers.worker_pool.Run"
	logger := slog.With("op", op)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	defer cancelWorker()

	for i := 0; i < w.countWorkers; i++ {
		w.wg.Add(1)
		go w.worker(workerCtx, i)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for deliveryDB := range w.queue.Out() {
			logger.Info("Start processing out channel", slog.Int64("DeliveryID", deliveryDB.ID))
			err := w.store.Delivery.ChangeStatusById(nil, deliveryDB.ID, deliveryDB.Status)
			if err != nil {
				logger.Error("Failed send update deliveries", sl.Err(err))
			}
		}
	}()

	ticker := time.NewTicker(w.cron * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutdown worker poll")
			w.queue.CloseIn()
			w.wg.Wait()
			w.queue.CloseOut()
			logger.Info("Worker pool stopped")
			return
		case <-ticker.C:
			err := w.loadDeliveries(workerCtx)
			if err != nil {
				logger.Error("Failed to load deliveries", sl.Err(err))
			}
		}
	}
}

func (w *WorkerPoolQueue) worker(ctx context.Context, workerId int) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker stopped by context", slog.Int("id", workerId))
			return
		case delivery, ok := <-w.queue.In():
			if !ok {
				slog.Info("Worker stopped, IN channel is close", slog.Int("id", workerId))
				return
			}
			slog.Info("Start processing IN channel",
				slog.Int("id", workerId),
				slog.Int64("DeliveryID", delivery.ID))

			delivery.Status = model.SHIPPED
			w.queue.Out() <- delivery
			slog.Info("Worker completed",
				slog.Int("id", workerId),
				slog.Int64("DeliveryID", delivery.ID))
		}
	}
}
func (w *WorkerPoolQueue) loadDeliveries(ctx context.Context) error {
	const op = "lib.workers.worker_pool.loadDeliveries"
	logger := slog.With("op", op)

	deliveriesDB, err := w.store.Delivery.GetAllByStatus(nil, model.CREATED)
	if err != nil {
		return err
	}

	logger.Info("Load deliveries", slog.Int("count", len(deliveriesDB)))
	for i := range deliveriesDB {
		select {
		case <-ctx.Done():
			logger.Info("Load deliveries stopped by context")
			return ctx.Err()
		case w.queue.In() <- deliveriesDB[i]:
		}
	}
	return nil
}
