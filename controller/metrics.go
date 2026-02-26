package controller

import (
	"erp-2c/service/use_cases"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type MetricsController struct {
	services *use_cases.Manager
}

func NewMetricsController(services *use_cases.Manager) *MetricsController {
	return &MetricsController{services: services}
}

func (m *MetricsController) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "control.metrics.GetAll"
	sLogger := slog.With("OP", op)
	sLogger.Info("get request all Metrics")
	render.JSON(w, r, "metrics")
}
