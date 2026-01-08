package controller

import (
	"erp-2c/dto/response"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type ProductController struct {
	services *use_cases.Manager
}

func NewProductController(services *use_cases.Manager) *ProductController {
	return &ProductController{services: services}
}

func (p *ProductController) Save(w http.ResponseWriter, r *http.Request) {
	slog.Info("Post request Save")

	var productToSave model.Product
	if err := render.DecodeJSON(r.Body, &productToSave); err != nil {
		response.BadRequest("Invalid decode json", r.Body)
		return
	}
	//saved, err := p.services.ProductService.Save(productToSave)
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
