package use_cases

import (
	"erp-2c/model"
	"erp-2c/store"
)

type MetricsService struct {
	repo *store.Store
}

func NewMetricsService(repo *store.Store) *MetricsService {
	return &MetricsService{repo: repo}
}

func (m *MetricsService) GetMetrics() model.MetricsDomain {
	return model.MetricsDomain{}
}
