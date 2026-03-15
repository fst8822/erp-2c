package controller

import (
	"api-gateway/lib/response"
	"api-gateway/lib/sl"
	"api-gateway/lib/types"
	"api-gateway/model"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type productServiceInt interface {
	Save(productToSave model.ProductToSave) (*model.ProductDomain, error)
	GetById(productId int64) (*model.ProductDomain, error)
	GetByName(productName string) (*model.ProductDomain, error)
	GetAll() ([]model.ProductDomain, error)
	UpdateById(productId int64, productToUpdate model.ProductUpdate) error
	DeleteById(productId int64) error
}

type ProductController struct {
	productService productServiceInt
	validate       *validator.Validate
}

func NewProductController(productService productServiceInt, validate *validator.Validate) *ProductController {
	return &ProductController{
		productService: productService,
		validate:       validate,
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

	if err := p.validate.Struct(productToSave); err != nil {
		slog.Error("failed validate request fields", sl.ErrWithOP(err, op))
		response.ValidationError(err).SendResponse(w, r)
		return
	}

	saved, err := p.productService.Save(productToSave)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}

func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetAll"

	products, err := p.productService.GetAll()
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(products).SendResponse(w, r)
}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetWithItemsById"

	param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	found, err := p.productService.GetById(productId)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(found).SendResponse(w, r)
}

func (p *ProductController) GetByName(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetByName"

	productName := chi.URLParam(r, "name")
	if strings.TrimSpace(productName) == "" {
		slog.Error("Path variable is empty", slog.StringValue(productName))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	productDomain, err := p.productService.GetByName(productName)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
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

	if err := p.productService.UpdateById(productId, productToUpdate); err != nil {
		types.HandleError(err).SendResponse(w, r)
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

	if err := p.productService.DeleteById(productId); err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.NoContent()
}
