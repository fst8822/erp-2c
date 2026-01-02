package controller

import "erp-2c/service"

type UserController struct {
	services *service.Manager
}

func NewUserController(services *service.Manager) *UserController {
	return &UserController{services: services}
}

func (c *UserController) register() {
}

func (c *UserController) GetById() {

}

func (c *UserController) GetByName() {

}
