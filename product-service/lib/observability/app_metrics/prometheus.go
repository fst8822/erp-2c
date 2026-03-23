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
