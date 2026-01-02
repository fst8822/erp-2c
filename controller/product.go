package controller

import "erp-2c/service"

type ProductController struct {
	services *service.Manager
}

func NewProductController(services *service.Manager) *ProductController {
	return &ProductController{services: services}
}

func (p *ProductController) save() {

}

func (p *ProductController) getById() {

}

func (p *ProductController) getAll() {

}

func (p *ProductController) updateById() {

}

func (p *ProductController) deleteById() {

}
