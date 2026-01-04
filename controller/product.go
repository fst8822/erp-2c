package controller

import (
	"erp-2c/service"
	"net/http"
)

type ProductController struct {
	services *service.Manager
}

func NewProductController(services *service.Manager) *ProductController {
	return &ProductController{services: services}
}

func (p *ProductController) Save(w http.ResponseWriter, r *http.Request) {

}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {

}

func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (p *ProductController) UpdateById(w http.ResponseWriter, r *http.Request) {

}

func (p *ProductController) DeleteById(w http.ResponseWriter, r *http.Request) {

}
