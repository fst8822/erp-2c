package controller

import (
	"api-gateway/lib/response"
	"api-gateway/lib/sl"
	"api-gateway/lib/types"
	"api-gateway/model"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type userServiceInt interface {
	Save(userToSave model.SignUp) (*model.UserDomain, error)
	GetById(userId int64) (*model.UserDomain, error)
	GetByLogin(userId string) (*model.UserDomain, error)
}

type UserController struct {
	userService userServiceInt
	validate    *validator.Validate
}

func NewUserController(userService userServiceInt, validate *validator.Validate) *UserController {
	return &UserController{
		userService: userService,
		validate:    validate,
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

	found, err := c.userService.GetById(userId)
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

	saved, err := c.userService.Save(userToSave)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}
