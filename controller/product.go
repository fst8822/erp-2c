package controller

import (
	"erp-2c/lib/response"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	const op = "control.product.GetAll"

	products, err := p.services.ProductService.GetAll()
	if err != nil {
		resp := response.NotFound("NotFound")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("Products", products))
}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetById"

	productIdParam := chi.URLParam(r, "id")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		resp := response.BadRequest("failed to parse path variable", productId)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	productDomain, err := p.services.ProductService.GetById(int64(productId))
	if err != nil {
		resp := response.NotFound("NotFound")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("OK", productDomain))
}

func (p *ProductController) GetByName(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetByName"

	productName := chi.URLParam(r, "name")
	productDomain, err := p.services.ProductService.GetByName(productName)
	if err != nil {
		resp := response.NotFound("NotFound")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("OK", productDomain))
}

func (p *ProductController) UpdateById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.UpdateById"

	idParam := chi.URLParam(r, "id")
	productId, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.BadRequest("failed to parse path variable", productId)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	var productToUpdate model.ProductUpdate

	if err := render.DecodeJSON(r.Body, &productToUpdate); err != nil {
		resp := response.BadRequest("Invalid json decode body", r.Body)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	if err := p.services.ProductService.UpdateById(int64(productId), productToUpdate); err != nil {
		resp := response.NotFound("NotFound")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("ok", nil))
}

func (p *ProductController) DeleteById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.DeleteById"

	productIdParam := chi.URLParam(r, "id")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		resp := response.BadRequest("failed to parse path veriable", productId)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}

	if err := p.services.ProductService.DeleteById(int64(productId)); err != nil {
		resp := response.NotFound("NotFound")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("ok", nil))
}
