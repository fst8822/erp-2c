package workers

import (
	"context"
	"database/sql"
	"erp-2c/lib/collection"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

const bachSize = 4

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
	instanceID := uuid.New().String()
	logger := slog.With("op", op, "instanceID", instanceID)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	defer cancelWorker()

	for i := 0; i < w.countWorkers; i++ {
		w.wg.Add(1)
		go w.worker(workerCtx, i, instanceID)
	}

	w.wg.Add(1)
	go w.updateDBItem(instanceID)

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
			err := w.loadDeliveries(workerCtx, instanceID)
			if err != nil {
				logger.Error("Failed to load deliveries", sl.Err(err))
			}
		}
	}
}

func (w *WorkerPoolQueue) worker(ctx context.Context, workerId int, instanceID string) {
	const op = "lib.workers.worker_pool.worker"
	logger := slog.With(
		"op", op,
		"instanceID", instanceID,
		"workerId", workerId)

	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker stopped by context")
			return
		case delivery, ok := <-w.queue.In():
			if !ok {
				logger.Info("Worker stopped, IN channel is close")
				return
			}
			logger.Info("Start processing IN channel",
				slog.Int64("DeliveryID", delivery.ID))

			delivery.Status = model.SHIPPED
			w.queue.Out() <- delivery
			logger.Info("Worker completed",
				slog.Int64("DeliveryID", delivery.ID))
		}
	}
}

func (w *WorkerPoolQueue) loadDeliveries(ctx context.Context, instanceID string) error {
	const op = "lib.workers.worker_pool.loadDeliveries"
	logger := slog.With("op", op, "instanceID", instanceID)

	tx, err := w.store.BeginTxx(ctx)
	if err != nil {
		logger.Error("Failed to open transaction", sl.Err(err))
		return err
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error("Failed to Rollback transaction", sl.Err(err))
		} else {
			logger.Info("Transaction Rollback is successful")
		}
	}()

	deliveriesDB, err := w.store.Delivery.LockAndGetDeliveries(tx, model.CREATED, instanceID)
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
	err = tx.Commit()
	if err != nil {
		logger.Error("Failed to commit transaction", sl.Err(err))
		return err
	}
	committed = true
	return nil
}

func (w *WorkerPoolQueue) updateDBItem(instanceID string) {
	const op = "lib.workers.worker_pool.updateDBItem"
	logger := slog.With("op", op, "instanceID", instanceID)

	groups := make(map[model.DeliveryStatus][]int64, 100)
	ticket := time.NewTicker(3 * time.Second)
	defer w.wg.Done()
	defer ticket.Stop()

	for deliveryDB := range w.queue.Out() {
		groups[deliveryDB.Status] = append(groups[deliveryDB.Status], deliveryDB.ID)
		logger.Info("Start processing out channel", slog.Int64("DeliveryID", deliveryDB.ID))
		select {
		case <-ticket.C:
			countRows := w.GetTotalCount(groups)
			logger.Info("Start update to db", slog.Int("bachSize", countRows))
			err := w.store.Delivery.UpdateStatusByIds(nil, groups)
			if err != nil {
				logger.Error("Failed send update deliveries", sl.Err(err))
				return
			}
			groups = make(map[model.DeliveryStatus][]int64, 100)
		default:
			countRows := w.GetTotalCount(groups)
			if countRows > bachSize {
				logger.Info("Start update to db", slog.Int("bachSize", len(groups)))
				err := w.store.Delivery.UpdateStatusByIds(nil, groups)
				if err != nil {
					logger.Error("Failed send update deliveries", sl.Err(err))
					return
				}
				groups = make(map[model.DeliveryStatus][]int64, 100)
			}
		}
	}
}

func (w *WorkerPoolQueue) GetTotalCount(groups map[model.DeliveryStatus][]int64) int {
	var totalCount int
	for _, ids := range groups {
		totalCount += len(ids)
	}
	return totalCount
}
