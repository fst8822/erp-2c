package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	services *use_cases.Manager
	validate *validator.Validate
}

func NewAuthController(services *use_cases.Manager, validate *validator.Validate) *AuthController {
	return &AuthController{
		services: services,
		validate: validate,
	}
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	const op = "control.auth.SignUp"

	var singUp model.SignUp

	err := render.DecodeJSON(r.Body, &singUp)
	if err != nil {
		slog.Error("failed to decode request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	if err = a.validate.Struct(&singUp); err != nil {
		slog.Error("failed validate request fields", sl.ErrWithOP(err, op))
		response.ValidationError(err).SendResponse(w, r)
		return
	}

	saved, err := a.services.AuthService.SignUp(singUp)
	if err != nil {
		slog.Error("failed to save user", sl.ErrWithOP(err, op))
		response.InternalServerError().SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	const op = "control.auth.SignIn"

	var signIn model.SignIn

	err := render.DecodeJSON(r.Body, &signIn)
	if err != nil {
		slog.Error("failed to decode request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	if err := a.validate.Struct(signIn); err != nil {
		slog.Error("failed validate request fields", sl.ErrWithOP(err, op))
		response.ValidationError(err).SendResponse(w, r)
		return
	}

	token, err := a.services.AuthService.SignIn(signIn)
	if err != nil {
		slog.Error("failed Unauthorized", sl.ErrWithOP(err, op))
		response.Unauthorized("Failed Unauthorized").SendResponse(w, r)
		return
	}

	response.OK(token).SendResponse(w, r)
}
