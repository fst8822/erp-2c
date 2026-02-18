package workers

import (
	"context"
	"erp-2c/lib/collection"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/service"
	"erp-2c/store"
	"log/slog"
	"sync"
	"time"
)

const bachSize = 4

type Worker interface {
	Run(ctx context.Context)
}
type WorkerPool struct {
	notify       service.NotifyService
	deliveryRepo store.DeliveryRepository
	queue        *collection.Queue
	wg           *sync.WaitGroup
	countWorkers int
	cron         time.Duration
}

func NewWorkerPool(
	notify service.NotifyService,
	delivery store.DeliveryRepository,
	queue *collection.Queue,
	countWorkers int,
	cron time.Duration) *WorkerPool {
	return &WorkerPool{
		notify:       notify,
		deliveryRepo: delivery,
		queue:        queue,
		wg:           &sync.WaitGroup{},
		countWorkers: countWorkers,
		cron:         cron}
}

func (w *WorkerPool) Run(ctx context.Context) {
	const op = "lib.workers.worker_pool.Run"
	logger := slog.With("op", op)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	defer cancelWorker()

	for i := 0; i < w.countWorkers; i++ {
		w.wg.Add(1)
		go w.worker(workerCtx, i)
	}

	w.wg.Add(1)
	go w.updateDBItem()

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

func (w *WorkerPool) worker(ctx context.Context, workerId int) {
	const op = "lib.workers.worker_pool.worker"
	logger := slog.With(
		"op", op,
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

func (w *WorkerPool) loadDeliveries(ctx context.Context) error {
	const op = "lib.workers.worker_pool.loadDeliveries"
	logger := slog.With("op", op)

	deliveriesDB, err := w.deliveryRepo.LockAndGetDeliveries(model.CREATED)
	if err != nil {
		logger.Error("failed Load deliveries", sl.Err(err))
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

func (w *WorkerPool) updateDBItem() {
	const op = "lib.workers.worker_pool.updateDBItem"
	logger := slog.With("op", op)

	groups := make(map[model.DeliveryStatus][]int64, 100)
	ticket := time.NewTicker(3 * time.Second)
	defer ticket.Stop()
	defer w.wg.Done()

	for deliveryDB := range w.queue.Out() {
		groups[deliveryDB.Status] = append(groups[deliveryDB.Status], deliveryDB.ID)
		logger.Info("Start processing out channel", slog.Int64("DeliveryID", deliveryDB.ID))

		go func(deliveryDB model.DeliveryDB) {
			notification := model.Notification{
				DeliveryId: deliveryDB.ID,
				Status:     deliveryDB.Status,
				CreatedAt:  time.Now(),
				UserID:     deliveryDB.UserID,
			}
			logger.Info("Start Send Notification")
			w.notify.SendNotify(notification)
		}(deliveryDB)

		select {
		case <-ticket.C:
			countRows := w.GetTotalCount(groups)
			logger.Info("Start update to db", slog.Int("bachSize", countRows))
			err := w.deliveryRepo.UpdateStatusByIds(nil, groups)
			if err != nil {
				logger.Error("Failed send update deliveries", sl.Err(err))
				return
			}
			groups = make(map[model.DeliveryStatus][]int64, 100)
		default:
			countRows := w.GetTotalCount(groups)
			if countRows > bachSize {
				logger.Info("Start update to db", slog.Int("bachSize", len(groups)))
				err := w.deliveryRepo.UpdateStatusByIds(nil, groups)
				if err != nil {
					logger.Error("Failed send update deliveries", sl.Err(err))
					return
				}
				groups = make(map[model.DeliveryStatus][]int64, 100)
			}
		}
	}
}

func (w *WorkerPool) GetTotalCount(groups map[model.DeliveryStatus][]int64) int {
	var totalCount int
	for _, ids := range groups {
		totalCount += len(ids)
	}
	return totalCount
}
