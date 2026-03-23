package app_metrics

import (
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
	WebsocketSessionDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "websocket_session_duration_seconds",
			Help:    "How long the websocket connection stayed open",
			Buckets: []float64{60, 180, 360, 600, 1800, 3600},
		},
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
