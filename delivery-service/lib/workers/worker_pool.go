package workers

import (
	"context"
	"delivery-service/lib/collection"
	"delivery-service/lib/sl"
	"delivery-service/model"
	clientrgrpc "delivery-service/server_grpc/proto/v1/notify"
	"log/slog"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type deliveryRepositoryInt interface {
	LockAndGetDeliveries(status model.DeliveryStatus) ([]model.DeliveryDB, error)
	UpdateStatusByIds(tx *sqlx.Tx, groups map[model.DeliveryStatus][]int64) error
}

const bachSize = 4

type Worker interface {
	Run(ctx context.Context)
}
type WorkerPool struct {
	notify       clientrgrpc.NotifyServiceClient
	deliveryRepo deliveryRepositoryInt
	queue        *collection.Queue
	wg           *sync.WaitGroup
	countWorkers int
	cron         time.Duration
}

func NewWorkerPool(
	notify clientrgrpc.NotifyServiceClient,
	delivery deliveryRepositoryInt,
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
	slogger := slog.With("op", op)
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
			slogger.Info("Shutdown worker poll")
			w.queue.CloseIn()
			w.wg.Wait()
			w.queue.CloseOut()
			slogger.Info("Worker pool stopped")
			return
		case <-ticker.C:
			err := w.loadDeliveries(workerCtx)
			if err != nil {
				slogger.Error("Failed to load deliveries", sl.Err(err))
			}
		}
	}
}

func (w *WorkerPool) worker(ctx context.Context, workerId int) {
	const op = "lib.workers.worker_pool.worker"
	slogger := slog.With(
		"op", op,
		"workerId", workerId)

	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			slogger.Info("Worker stopped by context")
			return
		case delivery, ok := <-w.queue.In():
			if !ok {
				slogger.Info("Worker stopped, IN channel is close")
				return
			}
			slogger.Info("Start processing IN channel",
				slog.Int64("DeliveryID", delivery.ID))

			delivery.Status = model.SHIPPED
			w.queue.Out() <- delivery

			slogger.Info("Worker completed",
				slog.Int64("DeliveryID", delivery.ID))
		}
	}
}

func (w *WorkerPool) loadDeliveries(ctx context.Context) error {
	const op = "lib.workers.worker_pool.loadDeliveries"
	slogger := slog.With("op", op)

	deliveriesDB, err := w.deliveryRepo.LockAndGetDeliveries(model.CREATED)
	if err != nil {
		slogger.Error("failed Load deliveries", sl.Err(err))
		return err
	}

	slogger.Info("Load deliveries", slog.Int("count", len(deliveriesDB)))
	for i := range deliveriesDB {
		select {
		case <-ctx.Done():
			slogger.Info("Load deliveries stopped by context")
			return ctx.Err()
		case w.queue.In() <- deliveriesDB[i]:
		}
	}
	return nil
}

func (w *WorkerPool) updateDBItem() {
	const op = "lib.workers.worker_pool.updateDBItem"
	slogger := slog.With("op", op)

	groups := make(map[model.DeliveryStatus][]int64, 100)
	ticket := time.NewTicker(3 * time.Second)
	notificationCH := make(chan model.DeliveryDB)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	defer close(notificationCH)
	defer ticket.Stop()
	defer w.wg.Done()

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.sendToNotification(ctx, notificationCH)
	}()

	for deliveryDB := range w.queue.Out() {
		groups[deliveryDB.Status] = append(groups[deliveryDB.Status], deliveryDB.ID)
		slogger.Info("Start processing out channel", slog.Int64("DeliveryID", deliveryDB.ID))

		select {
		case <-ticket.C:
			countRows := w.GetTotalCount(groups)
			slogger.Info("Start update to db", slog.Int("bachSize", countRows))
			err := w.deliveryRepo.UpdateStatusByIds(nil, groups)
			if err != nil {
				slogger.Error("Failed send update deliveries", sl.Err(err))
				continue
			}
			groups = make(map[model.DeliveryStatus][]int64, 100)
		default:
			countRows := w.GetTotalCount(groups)
			if countRows > bachSize {
				slogger.Info("Start update to db", slog.Int("bachSize", len(groups)))
				err := w.deliveryRepo.UpdateStatusByIds(nil, groups)
				if err != nil {
					slogger.Error("Failed send update deliveries", sl.Err(err))
					continue
				}
				groups = make(map[model.DeliveryStatus][]int64, 100)
			}
		}
		select {
		case notificationCH <- deliveryDB:
			slogger.Info("Send notification to channel")
		case <-ctx.Done():
			slog.Info("context canceled, stopping worker")
			return
		default:
			slogger.Info("notification channel full, dropping", "delivery_id", deliveryDB.ID)
		}
	}
}

func (w *WorkerPool) sendToNotification(ctx context.Context, ch <-chan model.DeliveryDB) {
	stream, err := w.notify.SendNotify(ctx)
	if err != nil {
		slog.Error("Failed get client grpc stream", slog.Any("err", err))
		return
	}
	defer stream.CloseSend()

	//todo when is close, it read all message
	for it := range ch {
		notification := clientrgrpc.Notification{
			DeliveryId: it.ID,
			Status:     string(it.Status),
			CreatedAt:  timestamppb.New(it.CreatedAt),
			UserId:     it.UserID,
		}
		slog.Info("Start Send Notification")
		if err := stream.Send(&notification); err != nil {
			slog.Error("failed to send notification",
				slog.Int64("deliveryId", it.ID),
				slog.String("status", string(it.Status)),
				err)
		} else {
			slog.Info("Notification sent successfully", slog.Int64("deliveryId", it.ID))
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
