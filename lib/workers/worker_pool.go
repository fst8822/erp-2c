package workers

import (
	"context"
	"erp-2c/lib/collection"
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
	wq           *sync.WaitGroup
	countWorkers int
	cron         time.Duration
}

func NewWorkerPoolQueue(
	store *store.Store,
	queue *collection.Queue,
	countWorkers int,
	cron time.Duration) *WorkerPoolQueue {
	return &WorkerPoolQueue{store: store, queue: queue, wq: &sync.WaitGroup{}, countWorkers: countWorkers, cron: cron}
}

func (qpq *WorkerPoolQueue) Run(ctx context.Context) {

	for i := 0; i < qpq.countWorkers; i++ {
		qpq.wq.Add(1)
		go func(workerId int) {
			defer qpq.wq.Done()

			slog.Info("Worker run GetAllByStatus", slog.Int("id", workerId))
			deliveriesDB, _ := qpq.store.Delivery.GetAllByStatus(nil, model.CREATED)

			for _, it := range deliveriesDB {
				qpq.queue.In() <- it
			}
			slog.Info("Worker end GetAllByStatus", slog.Int("id", workerId))
		}(i)

		qpq.wq.Add(1)
		qpq.worker(ctx, i)
	}

	for deliveryDB := range qpq.queue.Out() {
		qpq.wq.Add(1)
		go func(deliveryDB model.DeliveryDB) {
			defer qpq.wq.Done()

			slog.Info("Worker run ChangeStatusById", slog.Int64("deliveryDB.ID", deliveryDB.ID))
			_ = qpq.store.Delivery.ChangeStatusById(nil, deliveryDB.ID, deliveryDB.Status)
			slog.Info("Worker end GetAllByStatus", slog.Int64("deliveryDB.ID", deliveryDB.ID))

		}(deliveryDB)
	}
	qpq.wq.Wait()
	qpq.queue.CloseIn()
	qpq.queue.CloseOut()
}

func (qpq *WorkerPoolQueue) worker(ctx context.Context, workerId int) {
	defer qpq.wq.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-qpq.queue.In():
			if !ok {
				return
			}
			slog.Info("Worker run", slog.Int("id", workerId))
			delivery.Status = model.SHIPPED
			qpq.queue.Out() <- delivery
			slog.Info("Worker end", slog.Int("id", workerId))
		}
	}
}
