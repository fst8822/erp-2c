package controller

import (
	"erp-2c/service"
	"net/http"
)

type AuthController struct {
	services *service.Manager
}

func NewAuthController(services *service.Manager) *AuthController {
	return &AuthController{services: services}
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {}
