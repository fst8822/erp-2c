package controller

import (
	"erp-2c/dto/response"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type AuthController struct {
	services *use_cases.Manager
}

func NewAuthController(services *use_cases.Manager) *AuthController {
	return &AuthController{services: services}
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	const op = "control.routers.auth.SignUp"
	slog.With("op", op)

	var request model.User

	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		slog.Error("Failed to decode request body", sl.Err(err))

		resp := response.BadRequest("Invalid json decode", r.Body)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}

	saved, err := a.services.AuthService.SignUp(request)
	if err != nil {
		resp := response.InternalServerError("Failed to save user")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response.Created("Success", saved))
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SignUp"))
}
