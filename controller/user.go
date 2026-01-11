package controller

import (
	"erp-2c/lib/response"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserController struct {
	services *use_cases.Manager
}

func NewUserController(services *use_cases.Manager) *UserController {
	return &UserController{services: services}
}

func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.user.GetById"

	idParam := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.BadRequest("failed to parse path variable", idParam)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}

	found, _ := c.services.UserService.GetById(int64(userId))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response.OK("successful", found))
}

func (c *UserController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.user.Save"

	var userToSave model.SignUp

	err := render.DecodeJSON(r.Body, &userToSave)
	if err != nil {
		resp := response.BadRequest("Invalid json decode body", r.Body)
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	saved, err := c.services.UserService.Save(userToSave)
	if err != nil {
		resp := response.InternalServerError("InternalServerError")
		render.Status(r, resp.Code)
		render.JSON(w, r, resp)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK("successful", saved))
}
