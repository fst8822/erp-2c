package controller

import (
	"erp-2c/service"
	"net/http"
)

type UserController struct {
	services *service.Manager
}

func NewUserController(services *service.Manager) *UserController {
	return &UserController{services: services}
}
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
}

func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {

}

func (c *UserController) GetByName(w http.ResponseWriter, r *http.Request) {

}
