package controller

import (
	"erp-2c/lib/response"
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
	const op = "control.product.Save"
	slog.Info("Post request Save", slog.String("op", op))

	var productToSave model.ProductToSave
	if err := render.DecodeJSON(r.Body, &productToSave); err != nil {
		resp := response.BadRequest("Invalid json decode body", r.Body)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	saved, err := p.services.ProductService.Save(productToSave)
	if err != nil {
		resp := response.InternalServerError("InternalServerError")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.JSON(w, r, saved)
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
