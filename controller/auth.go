package controller

import "erp-2c/service"

type AuthController struct {
	services *service.Manager
}

func NewAuthController(services *service.Manager) *AuthController {
	return &AuthController{services: services}
}

func (a *AuthController) signUp() {

}

func (a *AuthController) signIn() {

}
