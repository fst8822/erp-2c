package controller

import (
	"erp-2c/service"
	"log/slog"
	"net/http"
)

type AuthController struct {
	services *service.Manager
}

func NewAuthController(services *service.Manager) *AuthController {
	return &AuthController{services: services}
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	slog.Info("Post request SignUp")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SignUp"))
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	slog.Info("Post request SignIn")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SignUp"))
}
