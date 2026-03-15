package model

import "time"

type MetricsDomain struct {
	CountDeliveries int
	ActiveConnWS    int
	AverageProcess  time.Time
}
