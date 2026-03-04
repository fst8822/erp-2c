package app_metrics

import (
	"context"
	"erp-2c/store"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const cron = time.Second * 10

var (
	HttpReqTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP request received",
		},
		[]string{"path", "method", "status"},
	)
	HttpReqDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_second",
			Help:    "Duration of HTTP request in second",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	WebsocketActiveConn = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_active_connection",
			Help: "Current number of websocket connection",
		},
	)
	TotalCountDelivery = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_count_deliveries",
			Help: "All total count deliveries",
		},
	)
	CurrentCountDeliveriesByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_count_deliveries_by_status",
			Help: "Current total count deliveries group by status",
		},
		[]string{"status"},
	)
)

func StartMetricsSync(ctx context.Context, repo *store.Store) {
	const op = "lib.observability.app_metrics.prometheus.StartMetricsSync"
	logger := slog.With("op", op)

	ticket := time.NewTicker(cron)

	defer ticket.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("Context is done, stopped method StartMetricsSync")
			return
		case <-ticket.C:
			res, err := repo.Delivery.GetStatusCount(nil)
			logger.Info("Start work StartMetricsSync", slog.Int("len", len(res)))
			if err != nil {
				continue
			}
			CurrentCountDeliveriesByStatus.Reset()
			for _, value := range res {
				CurrentCountDeliveriesByStatus.WithLabelValues(value.Status).Set(float64(value.Count))
			}
		}
	}
}
