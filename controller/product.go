package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	services *use_cases.Manager
	validate *validator.Validate
}

func NewProductController(services *use_cases.Manager, validate *validator.Validate) *ProductController {
	return &ProductController{
		services: services,
		validate: validate,
	}
}

func (p *ProductController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.Save"

	var productToSave model.ProductToSave
	if err := render.DecodeJSON(r.Body, &productToSave); err != nil {
		slog.Error("failed parse request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}
	saved, err := p.services.ProductService.Save(productToSave)
	if err != nil {
		slog.Error("failed save product", sl.ErrWithOP(err, op))
		//todo need discern, maybe constrain or err sql
		response.InternalServerError().SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}

func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetAll"

	products, err := p.services.ProductService.GetAll()
	if err != nil {
		//todo need discern, maybe constrain or err sql
		slog.Error("InternalServerError", sl.ErrWithOP(err, op))
		response.InternalServerError().SendResponse(w, r)
		return
	}
	response.OK(products).SendResponse(w, r)
}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetById"

	param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}
	found, err := p.services.ProductService.GetById(productId)
	if err != nil {
		//todo need discern, maybe constrain or err sql
		slog.Error("failed convert str to int", sl.ErrWithOP(err, op))
		response.NotFound("Product not found").SendResponse(w, r)
		return
	}
	response.OK(found).SendResponse(w, r)
}

func (p *ProductController) GetByName(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetByName"

	productName := chi.URLParam(r, "name")
	productDomain, err := p.services.ProductService.GetByName(productName)
	if err != nil {
		slog.Error("Product not found", sl.ErrWithOP(err, op))
		response.NotFound("Product not found").SendResponse(w, r)
		return
	}
	response.OK(productDomain).SendResponse(w, r)
}

func (p *ProductController) UpdateById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.UpdateById"

	Param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(Param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}
	var productToUpdate model.ProductUpdate

	if err := render.DecodeJSON(r.Body, &productToUpdate); err != nil {
		slog.Error("failed parse request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}
	if err := p.services.ProductService.UpdateById(productId, productToUpdate); err != nil {
		slog.Error("failed update product", sl.ErrWithOP(err, op))
		response.NotFound("product not found").SendResponse(w, r)
		return
	}
	response.OK(nil).SendResponse(w, r)
}

func (p *ProductController) DeleteById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.DeleteById"

	param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	if err := p.services.ProductService.DeleteById(productId); err != nil {
		slog.Error("failed delete product", sl.ErrWithOP(err, op))
		response.NotFound("NotFound").SendResponse(w, r)
		return
	}
	response.NoContent()
}
