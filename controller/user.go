package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	services *use_cases.Manager
	validate *validator.Validate
}

func NewUserController(services *use_cases.Manager, validate *validator.Validate) *UserController {
	return &UserController{
		services: services,
		validate: validate,
	}
}

func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.user.GetById"

	idParam := chi.URLParam(r, "id")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("invalid user id").SendResponse(w, r)
		return
	}

	found, err := c.services.UserService.GetById(userId)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
	}
	response.OK(found).SendResponse(w, r)
}

func (c *UserController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.user.Save"

	var userToSave model.SignUp

	err := render.DecodeJSON(r.Body, &userToSave)
	if err != nil {
		slog.Error("failed parse request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	saved, err := c.services.UserService.Save(userToSave)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}
