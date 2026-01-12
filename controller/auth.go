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
		slog.Error("Failed to decode request body", sl.ErrWithOP(err, op))

		resp := response.BadRequest("Invalid json decode body", r.Body)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}

	validate := validator.New()
	err = validate.Struct(&singUp)
	if err != nil {
		slog.Error("Failed validate resuest fields", sl.ErrWithOP(err, op))
		resp := response.BadRequest("Failed validate resuest fields", err.Error())
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}

	saved, err := a.services.AuthService.SignUp(singUp)
	if err != nil {
		slog.Error("Failed to save user", sl.ErrWithOP(err, op))

		resp := response.InternalServerError("Failed to save user")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response.Created("Success", saved))
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	const op = "control.auth.SignIn"

	var signIn model.SignIn

	err := render.DecodeJSON(r.Body, &signIn)
	if err != nil {
		slog.Error("Invalid json decode body", sl.ErrWithOP(err, op))

		res := response.BadRequest("Invalid json decode body", r.Body)
		render.Status(r, res.Code)
		render.JSON(w, r, res)
		return
	}

	token, err := a.services.AuthService.SignIn(signIn)
	if err != nil {
		slog.Error("failed to get jwt token", sl.ErrWithOP(err, op))

		res := response.Unauthorized("Unauthorized")
		render.Status(r, res.Code)
		render.JSON(w, r, res)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("OK", token))
}
