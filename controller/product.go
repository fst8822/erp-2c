package controller

import (
	"erp-2c/service"
	"log/slog"
	"net/http"
)

type ProductController struct {
	services *service.Manager
}

func NewProductController(services *service.Manager) *ProductController {
	return &ProductController{services: services}
}

func (p *ProductController) Save(w http.ResponseWriter, r *http.Request) {
	slog.Info("Post request Save")
}

func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetAll")
}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetById")
}

func (p *ProductController) GetByName(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetById")
}

func (p *ProductController) UpdateById(w http.ResponseWriter, r *http.Request) {
	slog.Info("PUT request UpdateById")
}

func (p *ProductController) DeleteById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete request DeleteById")
}
